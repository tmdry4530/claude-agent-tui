package bridge

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
)

// TrackingFile represents the OMC subagent-tracking.json format.
type TrackingFile struct {
	Agents         []TrackedAgent `json:"agents"`
	TotalSpawned   int            `json:"total_spawned"`
	TotalCompleted int            `json:"total_completed"`
	TotalFailed    int            `json:"total_failed"`
	LastUpdated    string         `json:"last_updated"`
}

// TrackedAgent represents a single agent entry in the tracking file.
type TrackedAgent struct {
	AgentID     string `json:"agent_id"`
	AgentType   string `json:"agent_type"`
	StartedAt   string `json:"started_at"`
	ParentMode  string `json:"parent_mode"`
	Status      string `json:"status"`
	CompletedAt string `json:"completed_at,omitempty"`
	DurationMs  int64  `json:"duration_ms,omitempty"`
}

// ConvertTracking reads a subagent-tracking.json file and returns CanonicalEvents.
// Each tracked agent produces 2 events: spawn + terminal (done/error/failed).
// Agents with status "running" produce only a spawn event.
func ConvertTracking(path string) ([]schema.CanonicalEvent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read tracking file: %w", err)
	}

	var tf TrackingFile
	if err := json.Unmarshal(data, &tf); err != nil {
		return nil, fmt.Errorf("parse tracking file: %w", err)
	}

	return ConvertAgents(tf.Agents)
}

// ConvertAgents converts a slice of TrackedAgent into CanonicalEvents.
func ConvertAgents(agents []TrackedAgent) ([]schema.CanonicalEvent, error) {
	var events []schema.CanonicalEvent

	for _, a := range agents {
		evts, err := agentToEvents(a)
		if err != nil {
			return nil, fmt.Errorf("convert agent %s: %w", a.AgentID, err)
		}
		events = append(events, evts...)
	}

	return events, nil
}

func agentToEvents(a TrackedAgent) ([]schema.CanonicalEvent, error) {
	startedAt, err := time.Parse(time.RFC3339Nano, a.StartedAt)
	if err != nil {
		return nil, fmt.Errorf("parse started_at: %w", err)
	}

	role := mapAgentTypeToRole(a.AgentType)
	mode := mapParentMode(a.ParentMode)
	runID := "omc-" + a.AgentID

	// Spawn event
	spawn := schema.CanonicalEvent{
		Ts:       startedAt,
		RunID:    runID,
		Provider: schema.ProviderClaude,
		Mode:     mode,
		AgentID:  a.AgentID,
		Role:     role,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
	}
	events := []schema.CanonicalEvent{spawn}

	// Terminal event (if not still running)
	if a.Status != "running" && a.CompletedAt != "" {
		completedAt, err := time.Parse(time.RFC3339Nano, a.CompletedAt)
		if err != nil {
			return nil, fmt.Errorf("parse completed_at: %w", err)
		}

		state, eventType := mapStatus(a.Status)
		var payload json.RawMessage
		if a.DurationMs > 0 {
			p, err := json.Marshal(map[string]int64{"duration_ms": a.DurationMs})
			if err != nil {
				return nil, fmt.Errorf("marshal duration payload: %w", err)
			}
			payload = p
		}

		terminal := schema.CanonicalEvent{
			Ts:       completedAt,
			RunID:    runID,
			Provider: schema.ProviderClaude,
			Mode:     mode,
			AgentID:  a.AgentID,
			Role:     role,
			State:    state,
			Type:     eventType,
			Payload:  payload,
		}
		events = append(events, terminal)
	}

	return events, nil
}

// mapAgentTypeToRole strips the "oh-my-claudecode:" prefix and looks up the role.
func mapAgentTypeToRole(agentType string) schema.Role {
	name := strings.TrimPrefix(agentType, "oh-my-claudecode:")
	if r, ok := schema.LookupRole(name); ok {
		return r
	}
	return schema.RoleCustom
}

// mapParentMode converts OMC parent_mode to schema.Mode.
func mapParentMode(parentMode string) schema.Mode {
	if parentMode == "" || parentMode == "none" {
		return ""
	}
	mode := schema.Mode(parentMode)
	if mode.IsValid() {
		return mode
	}
	return schema.ModeUnknown
}

// mapStatus converts OMC tracking status to (AgentState, EventType).
func mapStatus(status string) (schema.AgentState, schema.EventType) {
	switch status {
	case "completed":
		return schema.StateDone, schema.TypeTaskDone
	case "failed":
		return schema.StateFailed, schema.TypeError
	case "cancelled":
		return schema.StateCancelled, schema.TypeStateChange
	default:
		return schema.StateDone, schema.TypeTaskDone
	}
}

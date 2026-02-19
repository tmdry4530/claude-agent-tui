package bridge

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
)

// EmitEvent appends a single CanonicalEvent as JSONL to the session file.
// File path: <eventDir>/<sessionID>.jsonl
func EmitEvent(eventDir, sessionID string, event schema.CanonicalEvent) error {
	if err := os.MkdirAll(eventDir, 0755); err != nil {
		return fmt.Errorf("create event dir: %w", err)
	}

	path := filepath.Join(eventDir, sessionID+".jsonl")

	line, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	line = append(line, '\n')

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return fmt.Errorf("open event file: %w", err)
	}
	defer func() { _ = f.Close() }()

	if _, err := f.Write(line); err != nil {
		return fmt.Errorf("write event: %w", err)
	}

	return nil
}

// WriteEventsFile writes a slice of CanonicalEvents as JSONL to the given path.
func WriteEventsFile(path string, events []schema.CanonicalEvent) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() { _ = f.Close() }()

	enc := json.NewEncoder(f)
	for _, evt := range events {
		if err := enc.Encode(evt); err != nil {
			return fmt.Errorf("encode event: %w", err)
		}
	}

	return nil
}

// NewSpawnEvent creates a task_spawn CanonicalEvent for hook usage.
func NewSpawnEvent(agentID, agentType, parentMode string) schema.CanonicalEvent {
	role := mapAgentTypeToRole(agentType)
	mode := mapParentMode(parentMode)
	return schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "omc-" + agentID,
		Provider: schema.ProviderClaude,
		Mode:     mode,
		AgentID:  agentID,
		Role:     role,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
	}
}

// NewUpdateEvent creates a task_update CanonicalEvent.
func NewUpdateEvent(agentID, agentType, parentMode string, state schema.AgentState) schema.CanonicalEvent {
	role := mapAgentTypeToRole(agentType)
	mode := mapParentMode(parentMode)
	return schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "omc-" + agentID,
		Provider: schema.ProviderClaude,
		Mode:     mode,
		AgentID:  agentID,
		Role:     role,
		State:    state,
		Type:     schema.TypeTaskUpdate,
	}
}

// NewDoneEvent creates a task_done CanonicalEvent.
func NewDoneEvent(agentID, agentType, parentMode string) schema.CanonicalEvent {
	role := mapAgentTypeToRole(agentType)
	mode := mapParentMode(parentMode)
	return schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "omc-" + agentID,
		Provider: schema.ProviderClaude,
		Mode:     mode,
		AgentID:  agentID,
		Role:     role,
		State:    schema.StateDone,
		Type:     schema.TypeTaskDone,
	}
}

// NewErrorEvent creates an error CanonicalEvent.
func NewErrorEvent(agentID, agentType, parentMode, errMsg string) schema.CanonicalEvent {
	role := mapAgentTypeToRole(agentType)
	mode := mapParentMode(parentMode)
	var payload json.RawMessage
	if errMsg != "" {
		p, err := json.Marshal(map[string]string{"error": errMsg})
		if err == nil {
			payload = p
		}
	}
	return schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "omc-" + agentID,
		Provider: schema.ProviderClaude,
		Mode:     mode,
		AgentID:  agentID,
		Role:     role,
		State:    schema.StateError,
		Type:     schema.TypeError,
		Payload:  payload,
	}
}

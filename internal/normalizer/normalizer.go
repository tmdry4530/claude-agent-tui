package normalizer

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
)

// Normalizer transforms raw events into canonical events.
type Normalizer struct {
	redactor *Redactor
}

// New creates a new Normalizer instance.
func New() *Normalizer {
	return &Normalizer{
		redactor: NewRedactor(),
	}
}

// Normalize transforms a RawEvent into a CanonicalEvent.
// It performs:
// - JSON parsing and field mapping
// - Enum validation (unknown values downgraded to fallback enums with warning)
// - Role mapping via schema.LookupRole
// - Payload redaction
func (n *Normalizer) Normalize(raw schema.RawEvent) (*schema.CanonicalEvent, error) {
	// Parse raw data into a map
	var data map[string]interface{}
	if err := json.Unmarshal(raw.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to parse raw event data: %w", err)
	}

	// Extract and validate timestamp
	ts := raw.Received
	if tsStr, ok := data["ts"].(string); ok {
		if parsed, err := time.Parse(time.RFC3339, tsStr); err == nil {
			ts = parsed
		}
	}

	// Extract required fields
	runID := extractString(data, "run_id")
	agentID := extractString(data, "agent_id")
	if runID == "" {
		return nil, fmt.Errorf("missing required field: run_id")
	}
	if agentID == "" {
		return nil, fmt.Errorf("missing required field: agent_id")
	}

	// Normalize provider
	provider := n.normalizeProvider(extractString(data, "provider"))

	// Normalize mode (optional)
	mode := n.normalizeMode(extractString(data, "mode"))

	// Normalize role
	roleStr := extractString(data, "role")
	agentType := extractString(data, "agent_type")
	role := n.normalizeRole(roleStr, agentType)

	// Normalize state
	state := n.normalizeState(extractString(data, "state"))

	// Normalize event type
	eventType := n.normalizeType(extractString(data, "type"))

	// Extract optional fields
	parentAgentID := extractString(data, "parent_agent_id")
	taskID := extractString(data, "task_id")
	intentRef := extractString(data, "intent_ref")
	rawRef := extractString(data, "raw_ref")

	// Extract and redact payload
	var payload json.RawMessage
	if payloadData, ok := data["payload"]; ok {
		payloadBytes, err := json.Marshal(payloadData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		payload = n.redactor.Redact(payloadBytes)
	}

	// Extract metrics
	var metrics *schema.EventMetrics
	if metricsData, ok := data["metrics"].(map[string]interface{}); ok {
		metrics = extractMetrics(metricsData)
	}

	event := &schema.CanonicalEvent{
		Ts:            ts,
		RunID:         runID,
		Provider:      provider,
		Mode:          mode,
		AgentID:       agentID,
		ParentAgentID: parentAgentID,
		Role:          role,
		State:         state,
		Type:          eventType,
		TaskID:        taskID,
		IntentRef:     intentRef,
		Payload:       payload,
		Metrics:       metrics,
		RawRef:        rawRef,
	}

	return event, nil
}

func (n *Normalizer) normalizeProvider(p string) schema.Provider {
	provider := schema.Provider(p)
	if !provider.IsValid() {
		log.Printf("warning: unknown provider %q, using 'system'", p)
		return schema.ProviderSystem
	}
	return provider
}

func (n *Normalizer) normalizeMode(m string) schema.Mode {
	if m == "" {
		return ""
	}
	mode := schema.Mode(m)
	if !mode.IsValid() {
		log.Printf("warning: unknown mode %q, using 'unknown'", m)
		return schema.ModeUnknown
	}
	return mode
}

func (n *Normalizer) normalizeRole(roleStr, agentType string) schema.Role {
	// Try direct role mapping first
	if roleStr != "" {
		role := schema.Role(roleStr)
		if role.IsValid() {
			return role
		}
	}

	// Try agent type lookup
	if agentType != "" {
		if role, ok := schema.LookupRole(agentType); ok {
			return role
		}
		log.Printf("warning: unknown agent_type %q, using 'custom'", agentType)
		return schema.RoleCustom
	}

	// Fallback
	if roleStr != "" {
		log.Printf("warning: unknown role %q, using 'custom'", roleStr)
	}
	return schema.RoleCustom
}

func (n *Normalizer) normalizeState(s string) schema.AgentState {
	state := schema.AgentState(s)
	if !state.IsValid() {
		log.Printf("warning: unknown state %q, using 'idle'", s)
		return schema.StateIdle
	}
	return state
}

func (n *Normalizer) normalizeType(t string) schema.EventType {
	eventType := schema.EventType(t)
	if !eventType.IsValid() {
		log.Printf("warning: unknown event type %q, using 'state_change'", t)
		return schema.TypeStateChange
	}
	return eventType
}

func extractString(data map[string]interface{}, key string) string {
	if v, ok := data[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func extractMetrics(data map[string]interface{}) *schema.EventMetrics {
	metrics := &schema.EventMetrics{}

	if v, ok := data["latency_ms"].(float64); ok {
		metrics.LatencyMs = &v
	}
	if v, ok := data["tokens_in"].(float64); ok {
		tokensIn := int(v)
		metrics.TokensIn = &tokensIn
	}
	if v, ok := data["tokens_out"].(float64); ok {
		tokensOut := int(v)
		metrics.TokensOut = &tokensOut
	}
	if v, ok := data["cost_usd"].(float64); ok {
		metrics.CostUSD = &v
	}

	return metrics
}

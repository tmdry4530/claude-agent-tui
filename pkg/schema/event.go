package schema

import (
	"encoding/json"
	"fmt"
	"time"
)

// CanonicalEvent is the unified event format consumed by the TUI.
// See references/event-schema.md for the full specification.
type CanonicalEvent struct {
	Ts            time.Time        `json:"ts"`
	RunID         string           `json:"run_id"`
	Provider      Provider         `json:"provider"`
	Mode          Mode             `json:"mode,omitempty"`
	AgentID       string           `json:"agent_id"`
	ParentAgentID string           `json:"parent_agent_id,omitempty"`
	Role          Role             `json:"role"`
	State         AgentState       `json:"state"`
	Type          EventType        `json:"type"`
	TaskID        string           `json:"task_id,omitempty"`
	IntentRef     string           `json:"intent_ref,omitempty"`
	Payload       json.RawMessage  `json:"payload,omitempty"`
	Metrics       *EventMetrics    `json:"metrics,omitempty"`
	RawRef        string           `json:"raw_ref,omitempty"`
}

// EventMetrics holds performance and cost metadata.
type EventMetrics struct {
	LatencyMs *float64 `json:"latency_ms,omitempty"`
	TokensIn  *int     `json:"tokens_in,omitempty"`
	TokensOut *int     `json:"tokens_out,omitempty"`
	CostUSD   *float64 `json:"cost_usd,omitempty"`
}

// Validate checks required fields and enum validity.
// Returns nil if valid, or a descriptive error.
func (e *CanonicalEvent) Validate() error {
	if e.Ts.IsZero() {
		return fmt.Errorf("ts is required")
	}
	if e.RunID == "" {
		return fmt.Errorf("run_id is required")
	}
	if e.AgentID == "" {
		return fmt.Errorf("agent_id is required")
	}
	if !e.Provider.IsValid() {
		return fmt.Errorf("invalid provider: %q", e.Provider)
	}
	if e.Mode != "" && !e.Mode.IsValid() {
		return fmt.Errorf("invalid mode: %q", e.Mode)
	}
	if !e.Role.IsValid() {
		return fmt.Errorf("invalid role: %q", e.Role)
	}
	if !e.State.IsValid() {
		return fmt.Errorf("invalid state: %q", e.State)
	}
	if !e.Type.IsValid() {
		return fmt.Errorf("invalid type: %q", e.Type)
	}
	return nil
}

// RawEvent is a raw event before normalization.
type RawEvent struct {
	Source    string          `json:"source"`
	Data     json.RawMessage `json:"data"`
	Received time.Time       `json:"received"`
}

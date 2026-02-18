package inspector

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
)

func TestNewModel(t *testing.T) {
	m := NewModel()
	if m.viewport.Width != 0 || m.viewport.Height != 0 {
		t.Error("expected zero-sized viewport on init")
	}
	if m.event != nil {
		t.Error("expected nil event on init")
	}
}

func TestSetSize(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 24)

	if m.width != 80 {
		t.Errorf("expected width 80, got %d", m.width)
	}
	if m.height != 24 {
		t.Errorf("expected height 24, got %d", m.height)
	}
	// Viewport should be smaller to account for border
	if m.viewport.Width != 78 {
		t.Errorf("expected viewport width 78, got %d", m.viewport.Width)
	}
	if m.viewport.Height != 22 {
		t.Errorf("expected viewport height 22, got %d", m.viewport.Height)
	}
}

func TestClearEvent_ShowsNoEventSelected(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 24)

	// Set and then clear
	event := createTestEvent()
	m.SetEvent(&event)
	m.ClearEvent()

	view := m.View()
	if !strings.Contains(view, "No event selected") {
		t.Error("expected 'No event selected' message after ClearEvent")
	}
}

func TestSetEvent_RendersEventDetails(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 40)

	event := createTestEvent()
	m.SetEvent(&event)

	view := m.View()

	// Check core fields are present
	expectations := []string{
		"=== Event Detail ===",
		"Time:",
		"Run ID:    run-123",
		"Provider:  claude",
		"Mode:      ralph",
		"Agent:     agent-executor-1",
		"Parent:    agent-planner-main",
		"Role:      executor",
		"State:     running",
		"Type:      tool_call",
		"Task:      task-42",
		"Intent:    plan-7",
	}

	for _, exp := range expectations {
		if !strings.Contains(view, exp) {
			t.Errorf("expected view to contain %q", exp)
		}
	}
}

func TestSetEvent_RendersPayloadJSON(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 40)

	payload := map[string]interface{}{
		"tool_name": "Edit",
		"args": map[string]interface{}{
			"file": "auth.go",
		},
	}
	payloadBytes, _ := json.Marshal(payload)

	event := schema.CanonicalEvent{
		Ts:            time.Now(),
		RunID:         "run-1",
		Provider:      schema.ProviderClaude,
		AgentID:       "agent-1",
		Role:          schema.RoleExecutor,
		State:         schema.StateRunning,
		Type:          schema.TypeToolCall,
		Payload:       payloadBytes,
	}

	m.SetEvent(&event)
	view := m.View()

	if !strings.Contains(view, "--- Payload ---") {
		t.Error("expected payload section header")
	}
	if !strings.Contains(view, "Edit") {
		t.Error("expected payload to contain tool_name value")
	}
	if !strings.Contains(view, "auth.go") {
		t.Error("expected payload to contain file argument")
	}
}

func TestSetEvent_RendersMetrics(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 40)

	latency := 420.5
	tokensIn := 210
	tokensOut := 95
	cost := 0.0021

	event := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeToolCall,
		Metrics: &schema.EventMetrics{
			LatencyMs: &latency,
			TokensIn:  &tokensIn,
			TokensOut: &tokensOut,
			CostUSD:   &cost,
		},
	}

	m.SetEvent(&event)
	view := m.View()

	if !strings.Contains(view, "--- Metrics ---") {
		t.Error("expected metrics section header")
	}
	if !strings.Contains(view, "Latency:") {
		t.Error("expected latency label")
	}
	if !strings.Contains(view, "420ms") {
		t.Error("expected latency value")
	}
	if !strings.Contains(view, "Tokens In:") {
		t.Error("expected tokens in label")
	}
	if !strings.Contains(view, "210") {
		t.Error("expected tokens in value")
	}
	if !strings.Contains(view, "Tokens Out:") {
		t.Error("expected tokens out label")
	}
	if !strings.Contains(view, "95") {
		t.Error("expected tokens out value")
	}
	if !strings.Contains(view, "Cost:") {
		t.Error("expected cost label")
	}
	if !strings.Contains(view, "$0.0021") {
		t.Error("expected cost value")
	}
}

func TestSetEvent_NilEvent_ShowsNoEventSelected(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 24)
	m.SetEvent(nil)

	view := m.View()
	if !strings.Contains(view, "No event selected") {
		t.Error("expected 'No event selected' for nil event")
	}
}

func TestSetEvent_EmptyPayload_NoPayloadSection(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 40)

	event := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		Payload:  nil,
	}

	m.SetEvent(&event)
	view := m.View()

	if strings.Contains(view, "--- Payload ---") {
		t.Error("expected no payload section when payload is empty")
	}
}

func TestSetEvent_NoMetrics_NoMetricsSection(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 40)

	event := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		Metrics:  nil,
	}

	m.SetEvent(&event)
	view := m.View()

	if strings.Contains(view, "--- Metrics ---") {
		t.Error("expected no metrics section when metrics is nil")
	}
}

// createTestEvent returns a fully populated test event.
func createTestEvent() schema.CanonicalEvent {
	latency := 420.0
	tokensIn := 210
	tokensOut := 95
	cost := 0.0021

	payload := map[string]interface{}{
		"tool_name": "Edit",
		"args": map[string]interface{}{
			"file": "auth.go",
		},
	}
	payloadBytes, _ := json.Marshal(payload)

	return schema.CanonicalEvent{
		Ts:            time.Date(2026, 2, 17, 22, 27, 0, 0, time.UTC),
		RunID:         "run-123",
		Provider:      schema.ProviderClaude,
		Mode:          schema.ModeRalph,
		AgentID:       "agent-executor-1",
		ParentAgentID: "agent-planner-main",
		Role:          schema.RoleExecutor,
		State:         schema.StateRunning,
		Type:          schema.TypeToolCall,
		TaskID:        "task-42",
		IntentRef:     "plan-7",
		Payload:       payloadBytes,
		Metrics: &schema.EventMetrics{
			LatencyMs: &latency,
			TokensIn:  &tokensIn,
			TokensOut: &tokensOut,
			CostUSD:   &cost,
		},
	}
}

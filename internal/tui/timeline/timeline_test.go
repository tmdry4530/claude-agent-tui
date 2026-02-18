package timeline

import (
	"strings"
	"testing"
	"time"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	if m.events == nil {
		t.Error("events slice should be initialized")
	}
	if len(m.events) != 0 {
		t.Errorf("expected 0 events, got %d", len(m.events))
	}
	if cap(m.events) != maxEvents {
		t.Errorf("expected events capacity %d, got %d", maxEvents, cap(m.events))
	}
}

func TestAddEvent(t *testing.T) {
	m := NewModel()
	event := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-001",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-001",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		TaskID:   "task-001",
	}

	m.AddEvent(event)

	if len(m.events) != 1 {
		t.Errorf("expected 1 event, got %d", len(m.events))
	}
	if m.events[0].AgentID != "agent-001" {
		t.Errorf("expected AgentID agent-001, got %s", m.events[0].AgentID)
	}
}

func TestAddEventNewestFirst(t *testing.T) {
	m := NewModel()
	event1 := schema.CanonicalEvent{
		Ts:       time.Now().Add(-1 * time.Hour),
		RunID:    "run-001",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-001",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		TaskID:   "task-001",
	}
	event2 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-001",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-002",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskDone,
		TaskID:   "task-002",
	}

	m.AddEvent(event1)
	m.AddEvent(event2)

	if len(m.events) != 2 {
		t.Errorf("expected 2 events, got %d", len(m.events))
	}
	// Newest should be first
	if m.events[0].AgentID != "agent-002" {
		t.Errorf("expected newest event first (agent-002), got %s", m.events[0].AgentID)
	}
	if m.events[1].AgentID != "agent-001" {
		t.Errorf("expected oldest event second (agent-001), got %s", m.events[1].AgentID)
	}
}

func TestAddEventTrimToMaxEvents(t *testing.T) {
	m := NewModel()

	// Add maxEvents + 10 events
	for i := 0; i < maxEvents+10; i++ {
		event := schema.CanonicalEvent{
			Ts:       time.Now(),
			RunID:    "run-001",
			Provider: schema.ProviderClaude,
			AgentID:  "agent-001",
			Role:     schema.RoleExecutor,
			State:    schema.StateRunning,
			Type:     schema.TypeTaskSpawn,
		}
		m.AddEvent(event)
	}

	if len(m.events) != maxEvents {
		t.Errorf("expected %d events (trimmed to maxEvents), got %d", maxEvents, len(m.events))
	}
}

func TestSetSize(t *testing.T) {
	m := NewModel()
	m.SetSize(100, 50)

	if m.width != 100 {
		t.Errorf("expected width 100, got %d", m.width)
	}
	if m.height != 50 {
		t.Errorf("expected height 50, got %d", m.height)
	}
	// Viewport should account for border (-2)
	if m.viewport.Width != 98 {
		t.Errorf("expected viewport width 98 (100-2), got %d", m.viewport.Width)
	}
	if m.viewport.Height != 48 {
		t.Errorf("expected viewport height 48 (50-2), got %d", m.viewport.Height)
	}
}

func TestViewZeroSize(t *testing.T) {
	m := NewModel()
	event := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-001",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-001",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
	}
	m.AddEvent(event)

	// Without SetSize, width and height are 0
	view := m.View()
	if view != "" {
		t.Error("view should be empty when size is 0")
	}
}

func TestViewWithSize(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)
	event := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-001",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-001",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		TaskID:   "task-001",
	}
	m.AddEvent(event)

	view := m.View()
	if view == "" {
		t.Error("view should not be empty when size is set and events exist")
	}
	// View should contain agent-001 somewhere in the styled output
	if !strings.Contains(view, "agent-001") {
		t.Error("view should contain agent-001")
	}
}

func TestRenderEvent(t *testing.T) {
	ts := time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)
	event := schema.CanonicalEvent{
		Ts:       ts,
		RunID:    "run-001",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-123",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		TaskID:   "task-456",
	}

	result := renderEvent(event)

	// Check timestamp (should be in HH:MM:SS format)
	if !strings.Contains(result, "14:30:45") {
		t.Errorf("expected timestamp 14:30:45 in result, got: %s", result)
	}

	// Check agentID
	if !strings.Contains(result, "agent-123") {
		t.Errorf("expected agentID agent-123 in result, got: %s", result)
	}

	// Check event type
	if !strings.Contains(result, "task_spawn") {
		t.Errorf("expected event type task_spawn in result, got: %s", result)
	}

	// Check task summary
	if !strings.Contains(result, "task:task-456") {
		t.Errorf("expected task summary task:task-456 in result, got: %s", result)
	}
}

func TestGetEventTypeStyle(t *testing.T) {
	tests := []struct {
		name      string
		eventType schema.EventType
	}{
		{"TaskSpawn", schema.TypeTaskSpawn},
		{"TaskDone", schema.TypeTaskDone},
		{"Error", schema.TypeError},
		{"ToolCall", schema.TypeToolCall},
		{"Message", schema.TypeMessage},
		{"Unknown", schema.EventType("unknown")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := getEventTypeStyle(tt.eventType)
			rendered := style.Render("test")
			if rendered == "" {
				t.Error("style should render non-empty output")
			}
		})
	}
}

func TestGetSummaryEmpty(t *testing.T) {
	event := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-001",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-001",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		// No TaskID or IntentRef
	}

	result := getSummary(event)
	if result != "" {
		t.Errorf("expected empty summary, got %s", result)
	}
}

func TestGetSummaryWithTaskID(t *testing.T) {
	event := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-001",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-001",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		TaskID:   "task-123",
	}

	result := getSummary(event)
	if result != "task:task-123" {
		t.Errorf("expected 'task:task-123', got %s", result)
	}
}

func TestGetSummaryWithIntentRef(t *testing.T) {
	event := schema.CanonicalEvent{
		Ts:        time.Now(),
		RunID:     "run-001",
		Provider:  schema.ProviderClaude,
		AgentID:   "agent-001",
		Role:      schema.RoleExecutor,
		State:     schema.StateRunning,
		Type:      schema.TypeTaskSpawn,
		IntentRef: "intent-456",
	}

	result := getSummary(event)
	if result != "intent:intent-456" {
		t.Errorf("expected 'intent:intent-456', got %s", result)
	}
}

func TestGetSummaryWithBoth(t *testing.T) {
	event := schema.CanonicalEvent{
		Ts:        time.Now(),
		RunID:     "run-001",
		Provider:  schema.ProviderClaude,
		AgentID:   "agent-001",
		Role:      schema.RoleExecutor,
		State:     schema.StateRunning,
		Type:      schema.TypeTaskSpawn,
		TaskID:    "task-123",
		IntentRef: "intent-456",
	}

	result := getSummary(event)
	expected := "task:task-123 intent:intent-456"
	if result != expected {
		t.Errorf("expected '%s', got %s", expected, result)
	}
}

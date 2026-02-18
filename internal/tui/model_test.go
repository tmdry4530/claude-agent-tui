package tui

import (
	"testing"
	"time"

	"github.com/chamdom/omc-agent-tui/internal/store"
	"github.com/chamdom/omc-agent-tui/pkg/schema"
	tea "github.com/charmbracelet/bubbletea"
)

func TestModelInit(t *testing.T) {
	m := NewModel(store.NewStore(100))
	cmd := m.Init()
	if cmd == nil {
		t.Error("Init() should return a command")
	}
}

func TestModelUpdate_Quit(t *testing.T) {
	m := NewModel(store.NewStore(100))

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Error("'q' key should trigger quit command")
	}
	if _, ok := updated.(Model); !ok {
		t.Error("Update should return Model type")
	}
}

func TestModelUpdate_WindowSize(t *testing.T) {
	m := NewModel(store.NewStore(100))

	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updated, _ := m.Update(msg)

	model, ok := updated.(Model)
	if !ok {
		t.Fatal("Update should return Model type")
	}
	if model.width != 100 || model.height != 30 {
		t.Errorf("Expected width=100, height=30; got width=%d, height=%d", model.width, model.height)
	}
}

func TestModelView(t *testing.T) {
	m := NewModel(store.NewStore(100))

	view := m.View()
	if view != "Initializing..." {
		t.Errorf("Expected 'Initializing...', got %q", view)
	}

	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	view = m.View()
	if view == "" {
		t.Error("View should return non-empty string")
	}
}

func TestModelAddEvent(t *testing.T) {
	s := store.NewStore(100)
	m := NewModel(s)

	event := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "test-run",
		Provider: schema.ProviderClaude,
		AgentID:  "test-agent",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		TaskID:   "task-001",
	}

	m.AddEvent(event)
	if s.EventCount() != 1 {
		t.Errorf("Expected 1 event in store, got %d", s.EventCount())
	}

	// Error event
	errorEvent := event
	errorEvent.Type = schema.TypeError
	m.AddEvent(errorEvent)
	if s.EventCount() != 2 {
		t.Errorf("Expected 2 events in store, got %d", s.EventCount())
	}
}

func TestModelFocusSwitch(t *testing.T) {
	m := NewModel(store.NewStore(100))

	if m.focused != 0 {
		t.Errorf("Expected initial focus 0 (arena), got %d", m.focused)
	}

	// Press tab -> focus 1 (timeline)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	model := updated.(Model)
	if model.focused != 1 {
		t.Errorf("Expected focus 1, got %d", model.focused)
	}

	// Cycle through all panels
	for i := 0; i < panelCount; i++ {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
		model = updated.(Model)
	}
	if model.focused != 1 {
		t.Errorf("Expected focus 1 after full cycle, got %d", model.focused)
	}
}

func TestModelEventMsg(t *testing.T) {
	s := store.NewStore(100)
	m := NewModel(s)

	event := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "test-run",
		Provider: schema.ProviderClaude,
		AgentID:  "test-agent",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		TaskID:   "task-001",
	}

	updated, _ := m.Update(EventMsg(event))
	_ = updated.(Model)

	if s.EventCount() != 1 {
		t.Errorf("Expected 1 event in store, got %d", s.EventCount())
	}
}

func TestModelGraphIntegration(t *testing.T) {
	s := store.NewStore(100)
	m := NewModel(s)

	parentEvent := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "test-run",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RolePlanner,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		TaskID:   "task-001",
	}
	m.AddEvent(parentEvent)

	childEvent := schema.CanonicalEvent{
		Ts:            time.Now(),
		RunID:         "test-run",
		Provider:      schema.ProviderClaude,
		AgentID:       "agent-2",
		ParentAgentID: "agent-1",
		Role:          schema.RoleExecutor,
		State:         schema.StateRunning,
		Type:          schema.TypeTaskSpawn,
		TaskID:        "task-002",
	}
	m.AddEvent(childEvent)

	if s.EventCount() != 2 {
		t.Errorf("Expected 2 events in store, got %d", s.EventCount())
	}
	if m.agentTasks["agent-1"] != "task-001" {
		t.Errorf("Expected agent-1 mapped to task-001, got %s", m.agentTasks["agent-1"])
	}
	if m.agentTasks["agent-2"] != "task-002" {
		t.Errorf("Expected agent-2 mapped to task-002, got %s", m.agentTasks["agent-2"])
	}

	doneEvent := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "test-run",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-2",
		Role:     schema.RoleExecutor,
		State:    schema.StateDone,
		Type:     schema.TypeTaskDone,
		TaskID:   "task-002",
	}
	m.AddEvent(doneEvent)

	if _, ok := m.agentTasks["agent-2"]; ok {
		t.Error("agent-2 should be removed from agentTasks after TaskDone")
	}
}

// New integration tests for Arena features

func TestModelArenaKeyNavigation(t *testing.T) {
	m := NewModel(store.NewStore(100))

	// Add agents
	m.AddEvent(schema.CanonicalEvent{
		Ts: time.Now(), RunID: "r", Provider: schema.ProviderClaude,
		AgentID: "a1", Role: schema.RolePlanner, State: schema.StateRunning,
		Type: schema.TypeMessage,
	})
	m.AddEvent(schema.CanonicalEvent{
		Ts: time.Now(), RunID: "r", Provider: schema.ProviderClaude,
		AgentID: "a2", Role: schema.RoleExecutor, State: schema.StateRunning,
		Type: schema.TypeMessage,
	})

	// Arena is focused by default (focused=0)
	// Navigate down with 'j'
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model := updated.(Model)

	selected := model.arena.SelectedAgent()
	if selected == nil {
		t.Fatal("Expected non-nil selected agent after j")
	}
	if selected.AgentID != "a2" {
		t.Errorf("Expected a2 selected after j, got %q", selected.AgentID)
	}
}

func TestModelArenaEnterInspector(t *testing.T) {
	s := store.NewStore(100)
	m := NewModel(s)

	event := schema.CanonicalEvent{
		Ts: time.Now(), RunID: "r", Provider: schema.ProviderClaude,
		AgentID: "a1", Role: schema.RolePlanner, State: schema.StateRunning,
		Type: schema.TypeMessage,
	}
	m.AddEvent(event)

	// Press enter while arena is focused -> should set inspector event
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model := updated.(Model)
	_ = model // inspector event was set internally
}

func TestModelAgentEventTracking(t *testing.T) {
	m := NewModel(store.NewStore(100))

	event1 := schema.CanonicalEvent{
		Ts: time.Now(), RunID: "r", Provider: schema.ProviderClaude,
		AgentID: "a1", Role: schema.RolePlanner, State: schema.StateRunning,
		Type: schema.TypeTaskSpawn, TaskID: "t1",
	}
	m.AddEvent(event1)

	if _, ok := m.agentEvents["a1"]; !ok {
		t.Error("Expected agent event to be tracked")
	}

	// Second event should update the tracked event
	event2 := schema.CanonicalEvent{
		Ts: time.Now(), RunID: "r", Provider: schema.ProviderClaude,
		AgentID: "a1", Role: schema.RolePlanner, State: schema.StateDone,
		Type: schema.TypeTaskDone, TaskID: "t1",
	}
	m.AddEvent(event2)

	tracked := m.agentEvents["a1"]
	if tracked.State != schema.StateDone {
		t.Errorf("Expected tracked event state done, got %q", tracked.State)
	}
}

func TestModelArenaSummary(t *testing.T) {
	m := NewModel(store.NewStore(100))

	event := schema.CanonicalEvent{
		Ts: time.Now(), RunID: "r", Provider: schema.ProviderClaude,
		AgentID: "a1", Role: schema.RoleExecutor, State: schema.StateRunning,
		Type: schema.TypeToolCall,
	}
	m.AddEvent(event)

	agent := m.arena.SelectedAgent()
	if agent == nil {
		t.Fatal("Expected non-nil agent")
	}
	if agent.Summary != "tool call" {
		t.Errorf("Expected summary 'tool call', got %q", agent.Summary)
	}
}

func TestBuildEventSummary(t *testing.T) {
	tests := []struct {
		eventType schema.EventType
		taskID    string
		state     schema.AgentState
		expected  string
	}{
		{schema.TypeTaskSpawn, "task-001", schema.StateRunning, "spawn: task-001"},
		{schema.TypeTaskDone, "", schema.StateDone, "task done"},
		{schema.TypeToolCall, "", schema.StateRunning, "tool call"},
		{schema.TypeToolResult, "", schema.StateRunning, "tool result"},
		{schema.TypeMessage, "", schema.StateRunning, "message"},
		{schema.TypeError, "", schema.StateError, "error"},
		{schema.TypeReplan, "", schema.StateRunning, "replanning"},
		{schema.TypeVerify, "", schema.StateRunning, "verifying"},
		{schema.TypeFix, "", schema.StateRunning, "fixing"},
		{schema.TypeRecover, "", schema.StateRunning, "recovering"},
		{schema.TypeStateChange, "", schema.StateRunning, "-> running"},
	}

	for _, tt := range tests {
		t.Run(string(tt.eventType), func(t *testing.T) {
			event := schema.CanonicalEvent{
				Type:   tt.eventType,
				TaskID: tt.taskID,
				State:  tt.state,
			}
			result := buildEventSummary(event)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestModelPanelCoexistence(t *testing.T) {
	s := store.NewStore(100)
	m := NewModel(s)

	// Set window size
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model := updated.(Model)

	// Add events
	events := []schema.CanonicalEvent{
		{
			Ts: time.Now(), RunID: "r", Provider: schema.ProviderClaude,
			AgentID: "planner", Role: schema.RolePlanner, State: schema.StateRunning,
			Type: schema.TypeTaskSpawn, TaskID: "t1",
		},
		{
			Ts: time.Now(), RunID: "r", Provider: schema.ProviderClaude,
			AgentID: "exec", Role: schema.RoleExecutor, State: schema.StateRunning,
			Type: schema.TypeToolCall,
		},
		{
			Ts: time.Now(), RunID: "r", Provider: schema.ProviderClaude,
			AgentID: "reviewer", Role: schema.RoleReviewer, State: schema.StateWaiting,
			Type: schema.TypeMessage,
		},
	}

	for _, e := range events {
		updated, _ = model.Update(EventMsg(e))
		model = updated.(Model)
	}

	// View should render all panels without panic
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view with all panels")
	}

	// Arena should have 3 agents
	if model.arena.AgentCount() != 3 {
		t.Errorf("Expected 3 agents in arena, got %d", model.arena.AgentCount())
	}
}

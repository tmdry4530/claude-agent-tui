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

	// Test 'q' key
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

	// View should show initializing without size
	view := m.View()
	if view != "Initializing..." {
		t.Errorf("Expected 'Initializing...', got %q", view)
	}

	// Set size via Update
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// View should return non-empty (even without applied size, "Initializing..." is non-empty)
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

	// Should not panic
	m.AddEvent(event)

	// Verify store received the event
	if s.EventCount() != 1 {
		t.Errorf("Expected 1 event in store, got %d", s.EventCount())
	}

	// Error event should increment error counter
	errorEvent := event
	errorEvent.Type = schema.TypeError
	m.AddEvent(errorEvent)

	if s.EventCount() != 2 {
		t.Errorf("Expected 2 events in store, got %d", s.EventCount())
	}
}

func TestModelFocusSwitch(t *testing.T) {
	m := NewModel(store.NewStore(100))

	initialFocus := m.focused

	// Press tab
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	model := updated.(Model)

	if model.focused == initialFocus {
		t.Error("Tab key should change focus")
	}

	// Cycle through all 4 panels
	for i := 0; i < panelCount; i++ {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
		model = updated.(Model)
	}

	// After panelCount tabs from position 1, should be back at 1
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

	// Send event via Update (simulates pipeline delivery)
	updated, _ := m.Update(EventMsg(event))
	_ = updated.(Model)

	// Verify store received the event
	if s.EventCount() != 1 {
		t.Errorf("Expected 1 event in store, got %d", s.EventCount())
	}
}

func TestModelGraphIntegration(t *testing.T) {
	s := store.NewStore(100)
	m := NewModel(s)

	// Spawn parent task
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

	// Spawn child task (parent agent has an active task)
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

	// Verify both events in store
	if s.EventCount() != 2 {
		t.Errorf("Expected 2 events in store, got %d", s.EventCount())
	}

	// Verify agent-task mapping
	if m.agentTasks["agent-1"] != "task-001" {
		t.Errorf("Expected agent-1 mapped to task-001, got %s", m.agentTasks["agent-1"])
	}
	if m.agentTasks["agent-2"] != "task-002" {
		t.Errorf("Expected agent-2 mapped to task-002, got %s", m.agentTasks["agent-2"])
	}

	// Complete task - should remove from agentTasks
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

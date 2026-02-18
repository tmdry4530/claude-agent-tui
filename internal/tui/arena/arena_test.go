package arena

import (
	"strings"
	"testing"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	if m.agents == nil {
		t.Error("Expected agents map to be initialized, got nil")
	}

	if len(m.agents) != 0 {
		t.Errorf("Expected empty agents map, got length %d", len(m.agents))
	}

	if m.width != 0 {
		t.Errorf("Expected width to be 0, got %d", m.width)
	}

	if m.height != 0 {
		t.Errorf("Expected height to be 0, got %d", m.height)
	}
}

func TestUpdateAgent_NewAgent(t *testing.T) {
	m := NewModel()

	agentID := "agent-001"
	role := schema.RoleExecutor
	state := schema.StateRunning

	m.UpdateAgent(agentID, role, state)

	if len(m.agents) != 1 {
		t.Fatalf("Expected 1 agent, got %d", len(m.agents))
	}

	card, exists := m.agents[agentID]
	if !exists {
		t.Fatal("Expected agent to exist in map")
	}

	if card.AgentID != agentID {
		t.Errorf("Expected AgentID %q, got %q", agentID, card.AgentID)
	}

	if card.Role != role {
		t.Errorf("Expected Role %q, got %q", role, card.Role)
	}

	if card.State != state {
		t.Errorf("Expected State %q, got %q", state, card.State)
	}
}

func TestUpdateAgent_ExistingAgent(t *testing.T) {
	m := NewModel()
	agentID := "agent-001"

	m.UpdateAgent(agentID, schema.RoleExecutor, schema.StateRunning)
	m.UpdateAgent(agentID, schema.RolePlanner, schema.StateDone)

	if len(m.agents) != 1 {
		t.Errorf("Expected 1 agent after update, got %d", len(m.agents))
	}

	card := m.agents[agentID]
	if card.Role != schema.RolePlanner {
		t.Errorf("Expected updated role %q, got %q", schema.RolePlanner, card.Role)
	}

	if card.State != schema.StateDone {
		t.Errorf("Expected updated state %q, got %q", schema.StateDone, card.State)
	}
}

func TestUpdateAgent_MultipleAgents(t *testing.T) {
	m := NewModel()

	m.UpdateAgent("agent-001", schema.RoleExecutor, schema.StateRunning)
	m.UpdateAgent("agent-002", schema.RolePlanner, schema.StateIdle)
	m.UpdateAgent("agent-003", schema.RoleReviewer, schema.StateDone)

	if len(m.agents) != 3 {
		t.Errorf("Expected 3 agents, got %d", len(m.agents))
	}
}

func TestSetSize(t *testing.T) {
	m := NewModel()

	m.SetSize(80, 24)

	if m.width != 80 {
		t.Errorf("Expected width 80, got %d", m.width)
	}

	if m.height != 24 {
		t.Errorf("Expected height 24, got %d", m.height)
	}

	m.SetSize(120, 40)

	if m.width != 120 {
		t.Errorf("Expected width 120 after update, got %d", m.width)
	}

	if m.height != 40 {
		t.Errorf("Expected height 40 after update, got %d", m.height)
	}
}

func TestView_NoSize(t *testing.T) {
	m := NewModel()
	m.UpdateAgent("agent-001", schema.RoleExecutor, schema.StateRunning)

	result := m.View()

	if result != "" {
		t.Errorf("Expected empty string when size is not set, got %q", result)
	}
}

func TestView_EmptyWithSize(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 24)

	result := m.View()

	if result == "" {
		t.Error("Expected non-empty output with size set")
	}

	if !strings.Contains(result, "No agents") {
		t.Errorf("Expected 'No agents' message, got %q", result)
	}
}

func TestView_WithAgents(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 24)
	m.UpdateAgent("agent-001", schema.RoleExecutor, schema.StateRunning)
	m.UpdateAgent("agent-002", schema.RolePlanner, schema.StateDone)

	result := m.View()

	if result == "" {
		t.Error("Expected non-empty output")
	}

	if strings.Contains(result, "No agents") {
		t.Error("Should not contain 'No agents' when agents exist")
	}
}

func TestGetRoleColor_KnownRoles(t *testing.T) {
	tests := []struct {
		role     schema.Role
		expected string
	}{
		{schema.RolePlanner, "#7B68EE"},
		{schema.RoleExecutor, "#4FC3F7"},
		{schema.RoleReviewer, "#FFD54F"},
		{schema.RoleGuard, "#EF5350"},
		{schema.RoleTester, "#66BB6A"},
		{schema.RoleWriter, "#A1887F"},
		{schema.RoleExplorer, "#BA68C8"},
		{schema.RoleArchitect, "#FF8A65"},
		{schema.RoleDebugger, "#F06292"},
		{schema.RoleVerifier, "#81C784"},
		{schema.RoleDesigner, "#4DD0E1"},
		{schema.RoleCustom, "#9E9E9E"},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			color := getRoleColor(tt.role)
			if color != tt.expected {
				t.Errorf("Expected color %q for role %q, got %q", tt.expected, tt.role, color)
			}
		})
	}
}

func TestGetRoleColor_UnknownRole(t *testing.T) {
	unknownRole := schema.Role("unknown-role")
	color := getRoleColor(unknownRole)
	expected := "#9E9E9E"

	if color != expected {
		t.Errorf("Expected default color %q for unknown role, got %q", expected, color)
	}
}

func TestGetStateStyle_AllStates(t *testing.T) {
	tests := []struct {
		state schema.AgentState
		name  string
	}{
		{schema.StateIdle, "idle"},
		{schema.StateRunning, "running"},
		{schema.StateError, "error"},
		{schema.StateDone, "done"},
		{schema.StateFailed, "failed"},
		{schema.StateCancelled, "cancelled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := getStateStyle(tt.state)
			rendered := style.Render(string(tt.state))
			if rendered == "" {
				t.Error("Expected style to render non-empty output")
			}
		})
	}
}

func TestGetStateStyle_DefaultCase(t *testing.T) {
	unknownState := schema.AgentState("unknown-state")
	style := getStateStyle(unknownState)
	rendered := style.Render(string(unknownState))

	if rendered == "" {
		t.Error("Expected default style to render non-empty output")
	}
}

func TestRenderCard(t *testing.T) {
	card := &AgentCard{
		AgentID: "agent-001",
		Role:    schema.RoleExecutor,
		State:   schema.StateRunning,
	}

	result := renderCard(card)

	if result == "" {
		t.Error("Expected non-empty card output")
	}

	if !strings.Contains(result, string(schema.RoleExecutor)) {
		t.Errorf("Expected card to contain role %q", schema.RoleExecutor)
	}

	if !strings.Contains(result, "agent-001") {
		t.Error("Expected card to contain agent ID")
	}

	if !strings.Contains(result, string(schema.StateRunning)) {
		t.Errorf("Expected card to contain state %q", schema.StateRunning)
	}
}

func TestRenderCard_AllRoles(t *testing.T) {
	roles := []schema.Role{
		schema.RolePlanner,
		schema.RoleExecutor,
		schema.RoleReviewer,
		schema.RoleGuard,
		schema.RoleTester,
		schema.RoleWriter,
		schema.RoleExplorer,
		schema.RoleArchitect,
		schema.RoleDebugger,
		schema.RoleVerifier,
		schema.RoleDesigner,
	}

	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			card := &AgentCard{
				AgentID: "test-agent",
				Role:    role,
				State:   schema.StateRunning,
			}

			result := renderCard(card)

			if !strings.Contains(result, string(role)) {
				t.Errorf("Expected card to contain role %q", role)
			}
		})
	}
}

func TestRenderCard_AllStates(t *testing.T) {
	states := []schema.AgentState{
		schema.StateIdle,
		schema.StateRunning,
		schema.StateError,
		schema.StateDone,
		schema.StateFailed,
		schema.StateCancelled,
	}

	for _, state := range states {
		t.Run(string(state), func(t *testing.T) {
			card := &AgentCard{
				AgentID: "test-agent",
				Role:    schema.RoleExecutor,
				State:   state,
			}

			result := renderCard(card)

			if !strings.Contains(result, string(state)) {
				t.Errorf("Expected card to contain state %q", state)
			}
		})
	}
}

package arena

import (
	"fmt"
	"strings"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the Arena panel state.
type Model struct {
	agents map[string]*AgentCard
	width  int
	height int
}

// AgentCard holds display state for a single agent.
type AgentCard struct {
	AgentID string
	Role    schema.Role
	State   schema.AgentState
}

// NewModel creates a new Arena model.
func NewModel() Model {
	return Model{
		agents: make(map[string]*AgentCard),
	}
}

// UpdateAgent adds or updates an agent card.
func (m *Model) UpdateAgent(agentID string, role schema.Role, state schema.AgentState) {
	m.agents[agentID] = &AgentCard{
		AgentID: agentID,
		Role:    role,
		State:   state,
	}
}

// SetSize updates the panel dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// View renders the Arena panel.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	var cards []string
	for _, agent := range m.agents {
		cards = append(cards, renderCard(agent))
	}

	if len(cards) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("240"))
		return emptyStyle.Render("No agents")
	}

	content := strings.Join(cards, "\n")

	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	return style.Render(content)
}

// renderCard renders a single agent card.
func renderCard(card *AgentCard) string {
	roleColor := getRoleColor(card.Role)
	stateStyle := getStateStyle(card.State)

	roleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(roleColor)).
		Bold(true)

	idStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	return fmt.Sprintf("%s %s %s",
		roleStyle.Render(string(card.Role)),
		idStyle.Render(fmt.Sprintf("[%s]", card.AgentID)),
		stateStyle.Render(string(card.State)),
	)
}

// getRoleColor returns the color for each role.
func getRoleColor(role schema.Role) string {
	colors := map[schema.Role]string{
		schema.RolePlanner:   "#7B68EE",
		schema.RoleExecutor:  "#4FC3F7",
		schema.RoleReviewer:  "#FFD54F",
		schema.RoleGuard:     "#EF5350",
		schema.RoleTester:    "#66BB6A",
		schema.RoleWriter:    "#A1887F",
		schema.RoleExplorer:  "#BA68C8",
		schema.RoleArchitect: "#FF8A65",
		schema.RoleDebugger:  "#F06292",
		schema.RoleVerifier:  "#81C784",
		schema.RoleDesigner:  "#4DD0E1",
		schema.RoleCustom:    "#9E9E9E",
	}
	if color, ok := colors[role]; ok {
		return color
	}
	return "#9E9E9E"
}

// getStateStyle returns the lipgloss style for each state.
func getStateStyle(state schema.AgentState) lipgloss.Style {
	switch state {
	case schema.StateIdle:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	case schema.StateRunning:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)
	case schema.StateError:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Blink(true)
	case schema.StateDone:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))
	case schema.StateFailed:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Strikethrough(true)
	case schema.StateCancelled:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Strikethrough(true)
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	}
}

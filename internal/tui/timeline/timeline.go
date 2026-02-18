package timeline

import (
	"fmt"
	"strings"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const maxEvents = 100

// Model represents the Timeline panel state.
type Model struct {
	events   []schema.CanonicalEvent
	viewport viewport.Model
	width    int
	height   int
}

// NewModel creates a new Timeline model.
func NewModel() Model {
	vp := viewport.New(0, 0)
	return Model{
		events:   make([]schema.CanonicalEvent, 0, maxEvents),
		viewport: vp,
	}
}

// AddEvent adds an event to the timeline (newest first).
func (m *Model) AddEvent(event schema.CanonicalEvent) {
	// Prepend to show newest first
	m.events = append([]schema.CanonicalEvent{event}, m.events...)

	// Trim to maxEvents
	if len(m.events) > maxEvents {
		m.events = m.events[:maxEvents]
	}

	m.updateViewport()
}

// SetSize updates the panel dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width - 2 // Account for border
	m.viewport.Height = height - 2
	m.updateViewport()
}

// Update handles viewport messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the Timeline panel.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	return style.Render(m.viewport.View())
}

// updateViewport refreshes the viewport content.
func (m *Model) updateViewport() {
	var lines []string

	if len(m.events) == 0 {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("No events"))
	} else {
		for _, event := range m.events {
			lines = append(lines, renderEvent(event))
		}
	}

	m.viewport.SetContent(strings.Join(lines, "\n"))
}

// renderEvent formats a single event line.
func renderEvent(event schema.CanonicalEvent) string {
	timestamp := event.Ts.Format("15:04:05")

	timeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	agentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))

	typeStyle := getEventTypeStyle(event.Type)

	summary := getSummary(event)

	return fmt.Sprintf("%s %s %s %s",
		timeStyle.Render(fmt.Sprintf("[%s]", timestamp)),
		agentStyle.Render(event.AgentID),
		typeStyle.Render(string(event.Type)),
		summary,
	)
}

// getEventTypeStyle returns the style for each event type.
func getEventTypeStyle(eventType schema.EventType) lipgloss.Style {
	switch eventType {
	case schema.TypeTaskSpawn:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	case schema.TypeTaskDone:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	case schema.TypeError:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
	case schema.TypeToolCall:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#4FC3F7"))
	case schema.TypeMessage:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD54F"))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	}
}

// getSummary extracts a brief summary from the event.
func getSummary(event schema.CanonicalEvent) string {
	summary := ""

	if event.TaskID != "" {
		summary = fmt.Sprintf("task:%s", event.TaskID)
	}

	if event.IntentRef != "" {
		if summary != "" {
			summary += " "
		}
		summary += fmt.Sprintf("intent:%s", event.IntentRef)
	}

	return summary
}

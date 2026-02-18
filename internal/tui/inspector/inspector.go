package inspector

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the Inspector panel state.
type Model struct {
	event    *schema.CanonicalEvent
	viewport viewport.Model
	width    int
	height   int
}

// NewModel creates a new Inspector model.
func NewModel() Model {
	vp := viewport.New(0, 0)
	return Model{
		viewport: vp,
	}
}

// SetEvent sets the event to display.
func (m *Model) SetEvent(event *schema.CanonicalEvent) {
	m.event = event
	m.updateViewportContent()
}

// ClearEvent clears the currently displayed event.
func (m *Model) ClearEvent() {
	m.event = nil
	m.updateViewportContent()
}

// SetSize updates the panel dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width - 2  // Account for border
	m.viewport.Height = height - 2 // Account for border
	m.updateViewportContent()
}

// Update handles viewport scrolling.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the Inspector panel.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	return style.Render(m.viewport.View())
}

// updateViewportContent regenerates viewport content based on current event.
func (m *Model) updateViewportContent() {
	if m.event == nil {
		m.viewport.SetContent(m.renderNoEvent())
		return
	}
	m.viewport.SetContent(m.renderEvent())
}

// renderNoEvent returns content for no event selected state.
func (m *Model) renderNoEvent() string {
	style := lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("240"))

	// Center the content in available viewport space
	if m.viewport.Height > 0 {
		padding := strings.Repeat("\n", m.viewport.Height/2)
		return padding + style.Render("No event selected")
	}
	return style.Render("No event selected")
}

// renderEvent formats the selected event for display.
func (m *Model) renderEvent() string {
	e := m.event

	var b strings.Builder

	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("cyan"))
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("yellow"))

	// Header
	b.WriteString(sectionStyle.Render("=== Event Detail ==="))
	b.WriteString("\n\n")

	// Core fields
	b.WriteString(labelStyle.Render("Time:      "))
	b.WriteString(e.Ts.Format("2006-01-02T15:04:05Z"))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Run ID:    "))
	b.WriteString(e.RunID)
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Provider:  "))
	b.WriteString(string(e.Provider))
	b.WriteString("\n")

	if e.Mode != "" {
		b.WriteString(labelStyle.Render("Mode:      "))
		b.WriteString(string(e.Mode))
		b.WriteString("\n")
	}

	b.WriteString(labelStyle.Render("Agent:     "))
	b.WriteString(e.AgentID)
	b.WriteString("\n")

	if e.ParentAgentID != "" {
		b.WriteString(labelStyle.Render("Parent:    "))
		b.WriteString(e.ParentAgentID)
		b.WriteString("\n")
	}

	b.WriteString(labelStyle.Render("Role:      "))
	b.WriteString(string(e.Role))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("State:     "))
	b.WriteString(string(e.State))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Type:      "))
	b.WriteString(string(e.Type))
	b.WriteString("\n")

	if e.TaskID != "" {
		b.WriteString(labelStyle.Render("Task:      "))
		b.WriteString(e.TaskID)
		b.WriteString("\n")
	}

	if e.IntentRef != "" {
		b.WriteString(labelStyle.Render("Intent:    "))
		b.WriteString(e.IntentRef)
		b.WriteString("\n")
	}

	// Payload section
	if len(e.Payload) > 0 {
		b.WriteString("\n")
		b.WriteString(sectionStyle.Render("--- Payload ---"))
		b.WriteString("\n")

		// Pretty-print JSON
		var formatted map[string]interface{}
		if err := json.Unmarshal(e.Payload, &formatted); err == nil {
			prettyJSON, err := json.MarshalIndent(formatted, "", "  ")
			if err == nil {
				b.WriteString(string(prettyJSON))
			} else {
				b.WriteString(string(e.Payload))
			}
		} else {
			// If not a JSON object, just display raw
			b.WriteString(string(e.Payload))
		}
		b.WriteString("\n")
	}

	// Metrics section
	if e.Metrics != nil {
		b.WriteString("\n")
		b.WriteString(sectionStyle.Render("--- Metrics ---"))
		b.WriteString("\n")

		if e.Metrics.LatencyMs != nil {
			b.WriteString(labelStyle.Render("Latency:   "))
			b.WriteString(fmt.Sprintf("%.0fms", *e.Metrics.LatencyMs))
			b.WriteString("\n")
		}

		if e.Metrics.TokensIn != nil {
			b.WriteString(labelStyle.Render("Tokens In: "))
			b.WriteString(fmt.Sprintf("%d", *e.Metrics.TokensIn))
			b.WriteString("\n")
		}

		if e.Metrics.TokensOut != nil {
			b.WriteString(labelStyle.Render("Tokens Out: "))
			b.WriteString(fmt.Sprintf("%d", *e.Metrics.TokensOut))
			b.WriteString("\n")
		}

		if e.Metrics.CostUSD != nil {
			b.WriteString(labelStyle.Render("Cost:      "))
			b.WriteString(fmt.Sprintf("$%.4f", *e.Metrics.CostUSD))
			b.WriteString("\n")
		}
	}

	return b.String()
}

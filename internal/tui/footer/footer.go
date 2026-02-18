package footer

import (
	"fmt"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the Footer panel state.
type Model struct {
	eventCount     int
	errorCount     int
	mode           schema.Mode
	status         string
	redacted       bool
	width          int
	totalLatency   float64
	totalTokensIn  int
	totalTokensOut int
	totalCostUSD   float64
}

// NewModel creates a new Footer model.
func NewModel() Model {
	return Model{
		mode:   schema.ModeUnknown,
		status: "LIVE",
	}
}

// IncrementEvents increments the event counter.
func (m *Model) IncrementEvents() {
	m.eventCount++
}

// IncrementErrors increments the error counter.
func (m *Model) IncrementErrors() {
	m.errorCount++
}

// SetMode updates the current mode.
func (m *Model) SetMode(mode schema.Mode) {
	m.mode = mode
}

// SetRedacted sets the redaction status.
func (m *Model) SetRedacted(redacted bool) {
	m.redacted = redacted
}

// SetSize updates the panel width.
func (m *Model) SetSize(width int) {
	m.width = width
}

// SetMetrics updates aggregated metrics from the store.
func (m *Model) SetMetrics(totalLatency float64, tokensIn, tokensOut int, costUSD float64) {
	m.totalLatency = totalLatency
	m.totalTokensIn = tokensIn
	m.totalTokensOut = tokensOut
	m.totalCostUSD = costUSD
}

// SetStatus updates the status text.
func (m *Model) SetStatus(status string) {
	m.status = status
}

// View renders the Footer panel.
func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	eventsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))

	errorsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))

	modeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD54F"))

	metricsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#B0BEC5"))

	costStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CE93D8"))

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)

	parts := []string{
		eventsStyle.Render(fmt.Sprintf("Events: %d", m.eventCount)),
		"|",
	}

	if m.errorCount > 0 {
		parts = append(parts,
			errorsStyle.Render(fmt.Sprintf("Errors: %d", m.errorCount)),
			"|",
		)
	}

	if m.mode != schema.ModeUnknown {
		parts = append(parts,
			modeStyle.Render(fmt.Sprintf("Mode: %s", m.mode)),
			"|",
		)
	}

	if m.totalTokensIn > 0 || m.totalTokensOut > 0 {
		parts = append(parts,
			metricsStyle.Render(fmt.Sprintf("Tok: %s/%s", formatCount(m.totalTokensIn), formatCount(m.totalTokensOut))),
			"|",
		)
	}

	if m.totalCostUSD > 0 {
		parts = append(parts,
			costStyle.Render(fmt.Sprintf("$%.2f", m.totalCostUSD)),
			"|",
		)
	}

	if m.totalLatency > 0 {
		parts = append(parts,
			metricsStyle.Render(fmt.Sprintf("%.0fms", m.totalLatency)),
			"|",
		)
	}

	parts = append(parts, statusStyle.Render(m.status))

	if m.redacted {
		parts = append(parts,
			"|",
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("[REDACTED]"),
		)
	}

	content := lipgloss.JoinHorizontal(lipgloss.Left, parts...)

	style := lipgloss.NewStyle().
		Width(m.width).
		Background(lipgloss.Color("235")).
		Padding(0, 1)

	return style.Render(content)
}

// formatCount formats large numbers with K/M suffixes.
func formatCount(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

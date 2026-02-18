package footer

import (
	"strings"
	"testing"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	if m.mode != schema.ModeUnknown {
		t.Errorf("Expected mode to be ModeUnknown, got %v", m.mode)
	}

	if m.status != "LIVE" {
		t.Errorf("Expected status to be 'LIVE', got %q", m.status)
	}

	if m.eventCount != 0 {
		t.Errorf("Expected eventCount to be 0, got %d", m.eventCount)
	}

	if m.errorCount != 0 {
		t.Errorf("Expected errorCount to be 0, got %d", m.errorCount)
	}

	if m.redacted != false {
		t.Errorf("Expected redacted to be false, got %v", m.redacted)
	}

	if m.width != 0 {
		t.Errorf("Expected width to be 0, got %d", m.width)
	}
}

func TestIncrementEvents(t *testing.T) {
	m := NewModel()

	m.IncrementEvents()
	if m.eventCount != 1 {
		t.Errorf("Expected eventCount to be 1, got %d", m.eventCount)
	}

	m.IncrementEvents()
	m.IncrementEvents()
	if m.eventCount != 3 {
		t.Errorf("Expected eventCount to be 3, got %d", m.eventCount)
	}
}

func TestIncrementErrors(t *testing.T) {
	m := NewModel()

	m.IncrementErrors()
	if m.errorCount != 1 {
		t.Errorf("Expected errorCount to be 1, got %d", m.errorCount)
	}

	m.IncrementErrors()
	m.IncrementErrors()
	if m.errorCount != 3 {
		t.Errorf("Expected errorCount to be 3, got %d", m.errorCount)
	}
}

func TestSetMode(t *testing.T) {
	m := NewModel()

	m.SetMode(schema.ModeRalph)
	if m.mode != schema.ModeRalph {
		t.Errorf("Expected mode to be ModeRalph, got %v", m.mode)
	}

	m.SetMode(schema.ModeTeam)
	if m.mode != schema.ModeTeam {
		t.Errorf("Expected mode to be ModeTeam, got %v", m.mode)
	}
}

func TestSetRedacted(t *testing.T) {
	m := NewModel()

	m.SetRedacted(true)
	if m.redacted != true {
		t.Errorf("Expected redacted to be true, got %v", m.redacted)
	}

	m.SetRedacted(false)
	if m.redacted != false {
		t.Errorf("Expected redacted to be false, got %v", m.redacted)
	}
}

func TestSetSize(t *testing.T) {
	m := NewModel()

	m.SetSize(100)
	if m.width != 100 {
		t.Errorf("Expected width to be 100, got %d", m.width)
	}

	m.SetSize(200)
	if m.width != 200 {
		t.Errorf("Expected width to be 200, got %d", m.width)
	}
}

func TestSetMetrics(t *testing.T) {
	m := NewModel()

	m.SetMetrics(123.45, 1000, 2000, 0.50)

	if m.totalLatency != 123.45 {
		t.Errorf("Expected totalLatency to be 123.45, got %f", m.totalLatency)
	}

	if m.totalTokensIn != 1000 {
		t.Errorf("Expected totalTokensIn to be 1000, got %d", m.totalTokensIn)
	}

	if m.totalTokensOut != 2000 {
		t.Errorf("Expected totalTokensOut to be 2000, got %d", m.totalTokensOut)
	}

	if m.totalCostUSD != 0.50 {
		t.Errorf("Expected totalCostUSD to be 0.50, got %f", m.totalCostUSD)
	}
}

func TestSetStatus(t *testing.T) {
	m := NewModel()

	m.SetStatus("PAUSED")
	if m.status != "PAUSED" {
		t.Errorf("Expected status to be 'PAUSED', got %q", m.status)
	}

	m.SetStatus("STOPPED")
	if m.status != "STOPPED" {
		t.Errorf("Expected status to be 'STOPPED', got %q", m.status)
	}
}

func TestView_Empty(t *testing.T) {
	m := NewModel()
	// width is 0 by default

	view := m.View()
	if view != "" {
		t.Errorf("Expected empty view when width is 0, got %q", view)
	}
}

func TestView_BasicRender(t *testing.T) {
	m := NewModel()
	m.SetSize(100)

	view := m.View()
	if view == "" {
		t.Error("Expected non-empty view when width > 0")
	}

	if !strings.Contains(view, "Events: 0") {
		t.Error("Expected view to contain 'Events: 0'")
	}

	if !strings.Contains(view, "LIVE") {
		t.Error("Expected view to contain 'LIVE' status")
	}
}

func TestView_WithErrors(t *testing.T) {
	m := NewModel()
	m.SetSize(100)
	m.IncrementErrors()
	m.IncrementErrors()

	view := m.View()

	if !strings.Contains(view, "Errors: 2") {
		t.Error("Expected view to contain 'Errors: 2'")
	}
}

func TestView_WithMode(t *testing.T) {
	m := NewModel()
	m.SetSize(100)
	m.SetMode(schema.ModeRalph)

	view := m.View()

	if !strings.Contains(view, "Mode: ralph") {
		t.Error("Expected view to contain 'Mode: ralph'")
	}
}

func TestView_WithModeUnknown(t *testing.T) {
	m := NewModel()
	m.SetSize(100)
	// mode is ModeUnknown by default

	view := m.View()

	if strings.Contains(view, "Mode:") {
		t.Error("Expected view to NOT contain 'Mode:' when mode is ModeUnknown")
	}
}

func TestView_WithMetrics(t *testing.T) {
	m := NewModel()
	m.SetSize(150)
	m.SetMetrics(250.5, 5000, 10000, 1.25)

	view := m.View()

	// Check for tokens (formatted with K suffix)
	if !strings.Contains(view, "Tok:") {
		t.Error("Expected view to contain 'Tok:' when tokens are set")
	}

	// Check for cost
	if !strings.Contains(view, "$1.25") {
		t.Error("Expected view to contain '$1.25' cost")
	}

	// Check for latency
	if !strings.Contains(view, "250ms") || !strings.Contains(view, "251ms") {
		// Allow for rounding variations
		if !strings.Contains(view, "ms") {
			t.Error("Expected view to contain latency in milliseconds")
		}
	}
}

func TestView_WithRedacted(t *testing.T) {
	m := NewModel()
	m.SetSize(100)
	m.SetRedacted(true)

	view := m.View()

	if !strings.Contains(view, "[REDACTED]") {
		t.Error("Expected view to contain '[REDACTED]' when redacted is true")
	}
}

func TestView_WithoutRedacted(t *testing.T) {
	m := NewModel()
	m.SetSize(100)
	// redacted is false by default

	view := m.View()

	if strings.Contains(view, "[REDACTED]") {
		t.Error("Expected view to NOT contain '[REDACTED]' when redacted is false")
	}
}

func TestFormatCount_LessThan1000(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{999, "999"},
		{500, "500"},
	}

	for _, tc := range tests {
		result := formatCount(tc.input)
		if result != tc.expected {
			t.Errorf("formatCount(%d) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func TestFormatCount_ThousandsWithK(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{1000, "1.0K"},
		{1500, "1.5K"},
		{10000, "10.0K"},
		{999999, "1000.0K"},
		{5432, "5.4K"},
	}

	for _, tc := range tests {
		result := formatCount(tc.input)
		if result != tc.expected {
			t.Errorf("formatCount(%d) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func TestFormatCount_MillionsWithM(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{1000000, "1.0M"},
		{1500000, "1.5M"},
		{10000000, "10.0M"},
		{2500000, "2.5M"},
		{999999999, "1000.0M"},
	}

	for _, tc := range tests {
		result := formatCount(tc.input)
		if result != tc.expected {
			t.Errorf("formatCount(%d) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func TestView_CompleteScenario(t *testing.T) {
	m := NewModel()
	m.SetSize(200)
	m.IncrementEvents()
	m.IncrementEvents()
	m.IncrementEvents()
	m.IncrementErrors()
	m.SetMode(schema.ModeTeam)
	m.SetMetrics(500.0, 15000, 25000, 2.50)
	m.SetRedacted(true)
	m.SetStatus("ACTIVE")

	view := m.View()

	if view == "" {
		t.Fatal("Expected non-empty view")
	}

	// Verify all components are present
	checks := []string{
		"Events: 3",
		"Errors: 1",
		"Mode: team",
		"Tok:",
		"$2.50",
		"500ms",
		"ACTIVE",
		"[REDACTED]",
	}

	for _, check := range checks {
		if !strings.Contains(view, check) {
			t.Errorf("Expected view to contain %q, but it didn't. View: %s", check, view)
		}
	}
}

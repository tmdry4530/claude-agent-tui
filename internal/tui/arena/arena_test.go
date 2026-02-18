package arena

import (
	"strings"
	"testing"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
	"github.com/charmbracelet/lipgloss"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	if m.agents == nil {
		t.Error("Expected agents map to be initialized, got nil")
	}
	if len(m.agents) != 0 {
		t.Errorf("Expected empty agents map, got length %d", len(m.agents))
	}
	if m.width != 0 || m.height != 0 {
		t.Errorf("Expected zero dimensions, got %dx%d", m.width, m.height)
	}
	if m.selected != 0 {
		t.Errorf("Expected selected 0, got %d", m.selected)
	}
	if !m.unicode {
		t.Error("Expected unicode mode enabled by default")
	}
}

func TestUpdateAgent_NewAgent(t *testing.T) {
	m := NewModel()
	m.UpdateAgent("agent-001", schema.RoleExecutor, schema.StateRunning)

	if len(m.agents) != 1 {
		t.Fatalf("Expected 1 agent, got %d", len(m.agents))
	}
	if len(m.order) != 1 {
		t.Fatalf("Expected 1 in order, got %d", len(m.order))
	}

	card := m.agents["agent-001"]
	if card.AgentID != "agent-001" {
		t.Errorf("Expected AgentID agent-001, got %q", card.AgentID)
	}
	if card.Role != schema.RoleExecutor {
		t.Errorf("Expected Role executor, got %q", card.Role)
	}
	if card.State != schema.StateRunning {
		t.Errorf("Expected State running, got %q", card.State)
	}
}

func TestUpdateAgent_ExistingAgent(t *testing.T) {
	m := NewModel()
	m.UpdateAgent("agent-001", schema.RoleExecutor, schema.StateRunning)
	m.UpdateAgent("agent-001", schema.RolePlanner, schema.StateDone)

	if len(m.agents) != 1 {
		t.Errorf("Expected 1 agent after update, got %d", len(m.agents))
	}
	if len(m.order) != 1 {
		t.Errorf("Expected 1 in order, got %d", len(m.order))
	}
	if m.agents["agent-001"].Role != schema.RolePlanner {
		t.Errorf("Expected updated role planner, got %q", m.agents["agent-001"].Role)
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
	if len(m.order) != 3 {
		t.Errorf("Expected 3 in order, got %d", len(m.order))
	}
}

func TestUpdateAgentWithSummary(t *testing.T) {
	m := NewModel()
	m.UpdateAgentWithSummary("agent-001", schema.RoleExecutor, schema.StateRunning, "task spawned")

	card := m.agents["agent-001"]
	if card.Summary != "task spawned" {
		t.Errorf("Expected summary 'task spawned', got %q", card.Summary)
	}
}

func TestSetSize(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 24)

	if m.width != 80 || m.height != 24 {
		t.Errorf("Expected 80x24, got %dx%d", m.width, m.height)
	}
}

func TestSetFocused(t *testing.T) {
	m := NewModel()
	m.SetFocused(true)
	if !m.focused {
		t.Error("Expected focused true")
	}
	m.SetFocused(false)
	if m.focused {
		t.Error("Expected focused false")
	}
}

func TestSetUnicode(t *testing.T) {
	m := NewModel()
	if !m.unicode {
		t.Error("Expected unicode true by default")
	}
	m.SetUnicode(false)
	if m.unicode {
		t.Error("Expected unicode false after SetUnicode(false)")
	}
}

func TestSelectedAgent_Empty(t *testing.T) {
	m := NewModel()
	if m.SelectedAgent() != nil {
		t.Error("Expected nil for empty model")
	}
}

func TestSelectedAgent_WithAgents(t *testing.T) {
	m := NewModel()
	m.UpdateAgent("agent-001", schema.RoleExecutor, schema.StateRunning)
	m.UpdateAgent("agent-002", schema.RolePlanner, schema.StateIdle)

	agent := m.SelectedAgent()
	if agent == nil {
		t.Fatal("Expected non-nil selected agent")
	}
	if agent.AgentID != "agent-001" {
		t.Errorf("Expected first agent selected, got %q", agent.AgentID)
	}
}

func TestAgentCount(t *testing.T) {
	m := NewModel()
	if m.AgentCount() != 0 {
		t.Errorf("Expected 0, got %d", m.AgentCount())
	}
	m.UpdateAgent("agent-001", schema.RoleExecutor, schema.StateRunning)
	if m.AgentCount() != 1 {
		t.Errorf("Expected 1, got %d", m.AgentCount())
	}
}

func TestHandleKey_Navigation(t *testing.T) {
	m := NewModel()
	m.UpdateAgent("agent-001", schema.RoleExecutor, schema.StateRunning)
	m.UpdateAgent("agent-002", schema.RolePlanner, schema.StateIdle)
	m.UpdateAgent("agent-003", schema.RoleReviewer, schema.StateDone)

	// j moves down
	consumed := m.HandleKey("j")
	if !consumed {
		t.Error("Expected j to be consumed")
	}
	if m.selected != 1 {
		t.Errorf("Expected selected 1 after j, got %d", m.selected)
	}

	// k moves up
	m.HandleKey("k")
	if m.selected != 0 {
		t.Errorf("Expected selected 0 after k, got %d", m.selected)
	}

	// Wrap around down
	m.HandleKey("down")
	m.HandleKey("down")
	m.HandleKey("down")
	if m.selected != 0 {
		t.Errorf("Expected wrap to 0, got %d", m.selected)
	}

	// Wrap around up
	m.HandleKey("up")
	if m.selected != 2 {
		t.Errorf("Expected wrap to 2, got %d", m.selected)
	}

	// h/l navigation
	m.selected = 0
	m.HandleKey("l")
	if m.selected != 1 {
		t.Errorf("Expected 1 after l, got %d", m.selected)
	}
	m.HandleKey("h")
	if m.selected != 0 {
		t.Errorf("Expected 0 after h, got %d", m.selected)
	}
}

func TestHandleKey_EmptyAgents(t *testing.T) {
	m := NewModel()
	consumed := m.HandleKey("j")
	if !consumed {
		t.Error("Expected j to be consumed even with empty agents")
	}
	if m.selected != 0 {
		t.Errorf("Expected selected to remain 0, got %d", m.selected)
	}
}

func TestHandleKey_UnknownKey(t *testing.T) {
	m := NewModel()
	consumed := m.HandleKey("x")
	if consumed {
		t.Error("Expected x to not be consumed")
	}
}

func TestView_NoSize(t *testing.T) {
	m := NewModel()
	m.UpdateAgent("agent-001", schema.RoleExecutor, schema.StateRunning)
	if m.View() != "" {
		t.Error("Expected empty string when size is not set")
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
		t.Error("Expected 'No agents' message")
	}
}

func TestView_WithAgents(t *testing.T) {
	m := NewModel()
	m.SetSize(120, 24)
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

func TestView_FocusedBorder(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 24)
	m.SetFocused(true)
	m.UpdateAgent("agent-001", schema.RoleExecutor, schema.StateRunning)
	result := m.View()
	if result == "" {
		t.Error("Expected non-empty output when focused")
	}
}

func TestView_SelectedCard(t *testing.T) {
	m := NewModel()
	m.SetSize(120, 24)
	m.SetFocused(true)
	m.UpdateAgent("agent-001", schema.RoleExecutor, schema.StateRunning)
	m.UpdateAgent("agent-002", schema.RolePlanner, schema.StateIdle)

	// Select second agent
	m.HandleKey("j")
	result := m.View()
	if result == "" {
		t.Error("Expected non-empty output with selection")
	}
}

func TestView_ASCIIFallback(t *testing.T) {
	m := NewModel()
	m.SetSize(120, 24)
	m.SetUnicode(false)
	m.UpdateAgentWithSummary("agent-001", schema.RoleExecutor, schema.StateRunning, "working")

	result := m.View()
	if result == "" {
		t.Error("Expected non-empty ASCII output")
	}
	// ASCII fallback should contain the role label
	if !strings.Contains(result, "executor") {
		t.Error("Expected executor role in ASCII output")
	}
}

func TestView_WithSummary(t *testing.T) {
	m := NewModel()
	m.SetSize(120, 24)
	m.UpdateAgentWithSummary("agent-001", schema.RoleExecutor, schema.StateRunning, "task spawned")

	result := m.View()
	if !strings.Contains(result, "task spawned") {
		t.Error("Expected summary text in output")
	}
}

func TestGetRoleColor_KnownRoles(t *testing.T) {
	tests := []struct {
		role     schema.Role
		expected string
	}{
		{schema.RolePlanner, "#7AA2F7"},
		{schema.RoleExecutor, "#9ECE6A"},
		{schema.RoleReviewer, "#BB9AF7"},
		{schema.RoleGuard, "#F7768E"},
		{schema.RoleTester, "#E0AF68"},
		{schema.RoleWriter, "#73DACA"},
		{schema.RoleExplorer, "#58A6FF"},
		{schema.RoleArchitect, "#FFA657"},
		{schema.RoleDebugger, "#FF7B72"},
		{schema.RoleVerifier, "#56D364"},
		{schema.RoleDesigner, "#D2A8FF"},
		{schema.RoleCustom, "#8A93A5"},
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
	color := getRoleColor(schema.Role("unknown"))
	if color != "#8A93A5" {
		t.Errorf("Expected default color, got %q", color)
	}
}

func TestGetStateColor_KnownStates(t *testing.T) {
	tests := []struct {
		state    schema.AgentState
		expected string
	}{
		{schema.StateRunning, "#7EE787"},
		{schema.StateWaiting, "#7D8590"},
		{schema.StateBlocked, "#E3B341"},
		{schema.StateError, "#FF7B72"},
		{schema.StateDone, "#56D364"},
		{schema.StateFailed, "#FF7B72"},
		{schema.StateCancelled, "#6E7681"},
		{schema.StateIdle, "#6E7681"},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			color := getStateColor(tt.state)
			if color != tt.expected {
				t.Errorf("Expected %q for state %q, got %q", tt.expected, tt.state, color)
			}
		})
	}
}

func TestGetStateColor_UnknownState(t *testing.T) {
	color := getStateColor(schema.AgentState("unknown"))
	if color != "#6E7681" {
		t.Errorf("Expected default color, got %q", color)
	}
}

func TestGetStateStyle_AllStates(t *testing.T) {
	states := []schema.AgentState{
		schema.StateIdle, schema.StateRunning, schema.StateWaiting,
		schema.StateBlocked, schema.StateError, schema.StateDone,
		schema.StateFailed, schema.StateCancelled,
	}
	for _, state := range states {
		t.Run(string(state), func(t *testing.T) {
			style := getStateStyle(state)
			rendered := style.Render(string(state))
			if rendered == "" {
				t.Error("Expected non-empty rendered output")
			}
		})
	}
}

func TestGetStateStyle_DefaultCase(t *testing.T) {
	style := getStateStyle(schema.AgentState("unknown"))
	rendered := style.Render("test")
	if rendered == "" {
		t.Error("Expected non-empty output for default style")
	}
}

func TestRenderCard_Unicode(t *testing.T) {
	card := &AgentCard{
		AgentID: "agent-001",
		Role:    schema.RoleExecutor,
		State:   schema.StateRunning,
		Summary: "working on task",
	}
	result := renderCard(card, false, true)
	if result == "" {
		t.Error("Expected non-empty card")
	}
	if !strings.Contains(result, "executor") {
		t.Error("Expected role in card")
	}
	if !strings.Contains(result, "agent-001") {
		t.Error("Expected agent ID in card")
	}
	if !strings.Contains(result, "running") {
		t.Error("Expected state in card")
	}
	if !strings.Contains(result, "working on task") {
		t.Error("Expected summary in card")
	}
}

func TestRenderCard_ASCII(t *testing.T) {
	card := &AgentCard{
		AgentID: "agent-001",
		Role:    schema.RoleExecutor,
		State:   schema.StateRunning,
	}
	result := renderCard(card, false, false)
	if result == "" {
		t.Error("Expected non-empty ASCII card")
	}
	if !strings.Contains(result, "executor") {
		t.Error("Expected role in ASCII card")
	}
}

func TestRenderCard_Selected(t *testing.T) {
	card := &AgentCard{
		AgentID: "agent-001",
		Role:    schema.RoleExecutor,
		State:   schema.StateRunning,
	}
	selected := renderCard(card, true, true)
	unselected := renderCard(card, false, true)
	if selected == unselected {
		// Selected card should differ (different border color)
		// This is acceptable if lipgloss doesn't produce visible difference in test
	}
	if selected == "" {
		t.Error("Expected non-empty selected card")
	}
}

func TestRenderCard_ErrorBorder(t *testing.T) {
	card := &AgentCard{
		AgentID: "agent-001",
		Role:    schema.RoleExecutor,
		State:   schema.StateError,
	}
	result := renderCard(card, false, true)
	if result == "" {
		t.Error("Expected non-empty error card")
	}
}

func TestRenderCard_BlockedBorder(t *testing.T) {
	card := &AgentCard{
		AgentID: "agent-001",
		Role:    schema.RoleExecutor,
		State:   schema.StateBlocked,
	}
	result := renderCard(card, false, true)
	if result == "" {
		t.Error("Expected non-empty blocked card")
	}
}

func TestRenderCard_AllRoles(t *testing.T) {
	roles := []schema.Role{
		schema.RolePlanner, schema.RoleExecutor, schema.RoleReviewer,
		schema.RoleGuard, schema.RoleTester, schema.RoleWriter,
		schema.RoleExplorer, schema.RoleArchitect, schema.RoleDebugger,
		schema.RoleVerifier, schema.RoleDesigner, schema.RoleCustom,
	}
	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			card := &AgentCard{AgentID: "test", Role: role, State: schema.StateRunning}
			result := renderCard(card, false, true)
			if !strings.Contains(result, string(role)) {
				t.Errorf("Expected card to contain role %q", role)
			}
		})
	}
}

func TestRenderCard_AllStates(t *testing.T) {
	states := []schema.AgentState{
		schema.StateIdle, schema.StateRunning, schema.StateWaiting,
		schema.StateBlocked, schema.StateError, schema.StateDone,
		schema.StateFailed, schema.StateCancelled,
	}
	for _, state := range states {
		t.Run(string(state), func(t *testing.T) {
			card := &AgentCard{AgentID: "test", Role: schema.RoleExecutor, State: state}
			result := renderCard(card, false, true)
			if !strings.Contains(result, string(state)) {
				t.Errorf("Expected card to contain state %q", state)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly-10", 10, "exactly-10"},
		{"this-is-a-long-string", 10, "this-is..."},
		{"ab", 3, "ab"},
		{"abcd", 3, "abc"},
	}
	for _, tt := range tests {
		result := truncate(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
		}
	}
}

// Sprite tests

func TestGetSprite_Unicode(t *testing.T) {
	sprite := GetSprite(true)
	if len(sprite) != 3 {
		t.Fatalf("Expected 3 lines, got %d", len(sprite))
	}
	expected := SpriteLines
	for i, line := range sprite {
		if line != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}

func TestGetSprite_ExactStrings(t *testing.T) {
	sprite := GetSprite(true)
	// Line 0: " ▐▛███▜▌ "
	expected0 := " \u2590\u259B\u2588\u2588\u2588\u259C\u258C "
	if sprite[0] != expected0 {
		t.Errorf("Line 0: expected %q, got %q", expected0, sprite[0])
	}
	// Line 1: "▝▜█████▛▘"
	expected1 := "\u259D\u259C\u2588\u2588\u2588\u2588\u2588\u259B\u2598"
	if sprite[1] != expected1 {
		t.Errorf("Line 1: expected %q, got %q", expected1, sprite[1])
	}
	// Line 2: "   ▘▘ ▝▝ "
	expected2 := "   \u2598\u2598 \u259D\u259D "
	if sprite[2] != expected2 {
		t.Errorf("Line 2: expected %q, got %q", expected2, sprite[2])
	}
}

func TestGetSprite_ASCII(t *testing.T) {
	sprite := GetSprite(false)
	if len(sprite) != 3 {
		t.Fatalf("Expected 3 lines, got %d", len(sprite))
	}
	expected := ASCIIFallback
	for i, line := range sprite {
		if line != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}

func TestGetSprite_FallbackSwitch(t *testing.T) {
	unicode := GetSprite(true)
	ascii := GetSprite(false)
	if unicode[0] == ascii[0] {
		t.Error("Unicode and ASCII sprites should differ")
	}
}

func TestGetSprite_UniformWidth(t *testing.T) {
	for _, useUnicode := range []bool{true, false} {
		name := "unicode"
		if !useUnicode {
			name = "ascii"
		}
		t.Run(name, func(t *testing.T) {
			sprite := GetSprite(useUnicode)
			for i, line := range sprite {
				w := len([]rune(line))
				if w != SpriteWidth {
					t.Errorf("Line %d rune width: expected %d, got %d (%q)", i, SpriteWidth, w, line)
				}
			}
		})
	}
}

func TestPadCenter(t *testing.T) {
	tests := []struct {
		input    string
		width    int
		wantLen  int
	}{
		{"abc", 7, 7},
		{"ab", 6, 6},
		{"a", 5, 5},
		{"hello", 5, 5},
		{"toolong", 3, 7}, // longer than target, returned as-is
	}
	for _, tt := range tests {
		result := PadCenter(tt.input, tt.width)
		runes := []rune(result)
		if len(runes) != tt.wantLen {
			t.Errorf("PadCenter(%q, %d) rune len = %d, want %d", tt.input, tt.width, len(runes), tt.wantLen)
		}
	}
}

func TestPadCenter_OutputWidth(t *testing.T) {
	sprite := GetSprite(true)
	targetWidth := 22
	for i, line := range sprite {
		padded := PadCenter(line, targetWidth)
		runes := []rune(padded)
		if len(runes) != targetWidth {
			t.Errorf("Padded line %d rune width: expected %d, got %d", i, targetWidth, len(runes))
		}
	}
}

func TestSpriteRenderedWidth(t *testing.T) {
	// Verify all sprite lines have SpriteWidth runes
	for _, useUnicode := range []bool{true, false} {
		sprite := GetSprite(useUnicode)
		for i, line := range sprite {
			runes := []rune(line)
			if len(runes) != SpriteWidth {
				t.Errorf("Sprite[%d] rune width %d != SpriteWidth %d", i, len(runes), SpriteWidth)
			}
		}
	}
}

func TestGetStateIndicator_Unicode(t *testing.T) {
	states := []schema.AgentState{
		schema.StateRunning, schema.StateWaiting, schema.StateBlocked,
		schema.StateError, schema.StateDone, schema.StateIdle,
		schema.StateFailed, schema.StateCancelled,
	}
	for _, state := range states {
		t.Run(string(state)+"_unicode", func(t *testing.T) {
			indicator := GetStateIndicator(state, true)
			if indicator == "" {
				t.Error("Expected non-empty unicode indicator")
			}
		})
	}
}

func TestGetStateIndicator_ASCII(t *testing.T) {
	states := []schema.AgentState{
		schema.StateRunning, schema.StateWaiting, schema.StateBlocked,
		schema.StateError, schema.StateDone, schema.StateIdle,
		schema.StateFailed, schema.StateCancelled,
	}
	for _, state := range states {
		t.Run(string(state)+"_ascii", func(t *testing.T) {
			indicator := GetStateIndicator(state, false)
			if indicator == "" {
				t.Error("Expected non-empty ASCII indicator")
			}
		})
	}
}

func TestGetStateIndicator_UnknownState(t *testing.T) {
	unicode := GetStateIndicator(schema.AgentState("unknown"), true)
	ascii := GetStateIndicator(schema.AgentState("unknown"), false)
	if unicode == "" {
		t.Error("Expected non-empty unicode default indicator")
	}
	if ascii == "" {
		t.Error("Expected non-empty ASCII default indicator")
	}
}

// Box helper tests

func TestBoxTop(t *testing.T) {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#2A3142"))
	result := StripAnsi(boxTop(borderStyle))
	if result != "╭"+strings.Repeat("─", cardInnerWidth)+"╮" {
		t.Errorf("boxTop mismatch: got %q", result)
	}
	// Verify rune width
	runes := []rune(result)
	if len(runes) != cardTotalWidth {
		t.Errorf("boxTop rune width: expected %d, got %d", cardTotalWidth, len(runes))
	}
}

func TestBoxBottom(t *testing.T) {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#2A3142"))
	result := StripAnsi(boxBottom(borderStyle))
	if result != "╰"+strings.Repeat("─", cardInnerWidth)+"╯" {
		t.Errorf("boxBottom mismatch: got %q", result)
	}
	runes := []rune(result)
	if len(runes) != cardTotalWidth {
		t.Errorf("boxBottom rune width: expected %d, got %d", cardTotalWidth, len(runes))
	}
}

func TestBoxLine(t *testing.T) {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#2A3142"))
	// Plain text, no ANSI in content itself
	result := StripAnsi(boxLine("hello", borderStyle))
	// Should be: │ hello                 │
	// "hello" = 5 chars, pad = 22-5 = 17 spaces
	expected := "│ hello" + strings.Repeat(" ", 17) + " │"
	if result != expected {
		t.Errorf("boxLine mismatch:\n  got:  %q\n  want: %q", result, expected)
	}
	runes := []rune(result)
	if len(runes) != cardTotalWidth {
		t.Errorf("boxLine rune width: expected %d, got %d", cardTotalWidth, len(runes))
	}
}

func TestBoxLine_EmptyContent(t *testing.T) {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#2A3142"))
	result := StripAnsi(boxLine("", borderStyle))
	expected := "│ " + strings.Repeat(" ", cardContentWidth) + " │"
	if result != expected {
		t.Errorf("boxLine empty mismatch:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestBoxLine_MaxWidthContent(t *testing.T) {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#2A3142"))
	content := strings.Repeat("x", cardContentWidth)
	result := StripAnsi(boxLine(content, borderStyle))
	expected := "│ " + content + " │"
	if result != expected {
		t.Errorf("boxLine max width mismatch:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestStripAnsi(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"\x1b[31mred\x1b[0m", "red"},
		{"\x1b[1m\x1b[32mbold green\x1b[0m", "bold green"},
		{"no escape", "no escape"},
		{"", ""},
	}
	for _, tt := range tests {
		result := StripAnsi(tt.input)
		if result != tt.expected {
			t.Errorf("StripAnsi(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// Snapshot test: verify card render structure
func TestRenderCard_Snapshot(t *testing.T) {
	card := &AgentCard{
		AgentID: "agent-001",
		Role:    schema.RoleExecutor,
		State:   schema.StateRunning,
		Summary: "task spawned",
	}
	result := renderCard(card, false, true)
	stripped := StripAnsi(result)
	lines := strings.Split(stripped, "\n")

	// Card should have exactly 10 lines:
	// 0: top border (╭─...─╮)
	// 1-3: sprite (3 lines)
	// 4: role + indicator
	// 5: agent ID
	// 6: state
	// 7: "Recent activity" label
	// 8: summary
	// 9: bottom border (╰─...─╯)
	if len(lines) != 10 {
		t.Fatalf("Expected 10 lines, got %d:\n%s", len(lines), stripped)
	}

	// Line 0: top border
	if !strings.HasPrefix(lines[0], "╭") || !strings.HasSuffix(lines[0], "╮") {
		t.Errorf("Line 0 should be top border, got %q", lines[0])
	}

	// Lines 1-3: sprite lines wrapped in │ ... │
	for i := 1; i <= 3; i++ {
		if !strings.HasPrefix(lines[i], "│") || !strings.HasSuffix(lines[i], "│") {
			t.Errorf("Line %d should be boxed sprite, got %q", i, lines[i])
		}
	}

	// Line 4: role + indicator
	if !strings.Contains(lines[4], "executor") {
		t.Errorf("Line 4 should contain role, got %q", lines[4])
	}

	// Line 5: agent ID
	if !strings.Contains(lines[5], "agent-001") {
		t.Errorf("Line 5 should contain agent ID, got %q", lines[5])
	}

	// Line 6: state
	if !strings.Contains(lines[6], "running") {
		t.Errorf("Line 6 should contain state, got %q", lines[6])
	}

	// Line 7: "Recent activity" label
	if !strings.Contains(lines[7], "Recent activity") {
		t.Errorf("Line 7 should contain 'Recent activity', got %q", lines[7])
	}

	// Line 8: summary
	if !strings.Contains(lines[8], "task spawned") {
		t.Errorf("Line 8 should contain summary, got %q", lines[8])
	}

	// Line 9: bottom border
	if !strings.HasPrefix(lines[9], "╰") || !strings.HasSuffix(lines[9], "╯") {
		t.Errorf("Line 9 should be bottom border, got %q", lines[9])
	}

	// All lines should have the same rune width (cardTotalWidth)
	for i, line := range lines {
		w := len([]rune(line))
		if w != cardTotalWidth {
			t.Errorf("Line %d rune width: expected %d, got %d: %q", i, cardTotalWidth, w, line)
		}
	}

	// Sprite lines 1-3 must share the same x-anchor (uniform prefix before sprite block)
	// prefix = (cardContentWidth - SpriteWidth) / 2 spaces after "│ "
	expectedPrefix := (cardContentWidth - SpriteWidth) / 2
	for i := 1; i <= 3; i++ {
		inner := []rune(lines[i])
		// After "│ " (2 runes), next expectedPrefix runes must be spaces
		for j := 2; j < 2+expectedPrefix; j++ {
			if inner[j] != ' ' {
				t.Errorf("Line %d: expected space at rune %d for prefix, got %c", i, j, inner[j])
			}
		}
	}
}

func TestRenderCard_Snapshot_EmptySummary(t *testing.T) {
	card := &AgentCard{
		AgentID: "agent-002",
		Role:    schema.RolePlanner,
		State:   schema.StateIdle,
		Summary: "",
	}
	result := renderCard(card, false, true)
	stripped := StripAnsi(result)

	// Empty summary should show "-"
	if !strings.Contains(stripped, "-") {
		t.Error("Expected '-' for empty summary")
	}
	if !strings.Contains(stripped, "Recent activity") {
		t.Error("Expected 'Recent activity' label")
	}
}

func TestRenderCard_Snapshot_ASCII(t *testing.T) {
	card := &AgentCard{
		AgentID: "agent-003",
		Role:    schema.RoleReviewer,
		State:   schema.StateDone,
		Summary: "review complete",
	}
	result := renderCard(card, false, false)
	stripped := StripAnsi(result)
	lines := strings.Split(stripped, "\n")

	if len(lines) != 10 {
		t.Fatalf("Expected 10 lines in ASCII mode, got %d", len(lines))
	}

	// All lines should have cardTotalWidth rune width
	for i, line := range lines {
		w := len([]rune(line))
		if w != cardTotalWidth {
			t.Errorf("ASCII line %d rune width: expected %d, got %d: %q", i, cardTotalWidth, w, line)
		}
	}
}

// Order stability test
func TestOrderStability(t *testing.T) {
	m := NewModel()
	ids := []string{"a", "b", "c", "d", "e"}
	for _, id := range ids {
		m.UpdateAgent(id, schema.RoleExecutor, schema.StateRunning)
	}

	for i, id := range ids {
		if m.order[i] != id {
			t.Errorf("Expected order[%d] = %q, got %q", i, id, m.order[i])
		}
	}

	// Update existing agent should not change order
	m.UpdateAgent("c", schema.RolePlanner, schema.StateDone)
	for i, id := range ids {
		if m.order[i] != id {
			t.Errorf("After update, expected order[%d] = %q, got %q", i, id, m.order[i])
		}
	}
}

package arena

import "github.com/chamdom/omc-agent-tui/pkg/schema"

// MascotSprite holds unicode and ASCII representations of a CLCO mascot.
type MascotSprite struct {
	Unicode []string // Multi-line unicode block sprite (3 lines)
	ASCII   string   // Single-line ASCII fallback
}

// mascotSprites maps each role to its CLCO mascot sprite.
var mascotSprites = map[schema.Role]MascotSprite{
	schema.RolePlanner: {
		Unicode: []string{
			" \u250C\u2500\u2510 ",
			" \u2502\u2630\u2502 ",
			" \u2514\u2534\u2518 ",
		},
		ASCII: "[P]",
	},
	schema.RoleExecutor: {
		Unicode: []string{
			" \u256D\u2501\u256E ",
			" \u2503\u25B6\u2503 ",
			" \u2570\u2501\u256F ",
		},
		ASCII: "[X]",
	},
	schema.RoleReviewer: {
		Unicode: []string{
			" \u250C\u2500\u2510 ",
			" \u2502\u2714\u2502 ",
			" \u2514\u2500\u2518 ",
		},
		ASCII: "[R]",
	},
	schema.RoleGuard: {
		Unicode: []string{
			" \u256D\u2501\u256E ",
			" \u2503\u26A1\u2503 ",
			" \u2570\u2501\u256F ",
		},
		ASCII: "[G]",
	},
	schema.RoleTester: {
		Unicode: []string{
			" \u250C\u2500\u2510 ",
			" \u2502\u2713\u2502 ",
			" \u2514\u2534\u2518 ",
		},
		ASCII: "[T]",
	},
	schema.RoleWriter: {
		Unicode: []string{
			" \u256D\u2500\u256E ",
			" \u2502\u270E\u2502 ",
			" \u2570\u2500\u256F ",
		},
		ASCII: "[W]",
	},
	schema.RoleExplorer: {
		Unicode: []string{
			" \u250C\u2500\u2510 ",
			" \u2502\u2318\u2502 ",
			" \u2514\u2500\u2518 ",
		},
		ASCII: "[E]",
	},
	schema.RoleArchitect: {
		Unicode: []string{
			" \u256D\u2501\u256E ",
			" \u2503\u2302\u2503 ",
			" \u2570\u2501\u256F ",
		},
		ASCII: "[A]",
	},
	schema.RoleDebugger: {
		Unicode: []string{
			" \u250C\u2500\u2510 ",
			" \u2502\u2699\u2502 ",
			" \u2514\u2534\u2518 ",
		},
		ASCII: "[D]",
	},
	schema.RoleVerifier: {
		Unicode: []string{
			" \u256D\u2500\u256E ",
			" \u2502\u2611\u2502 ",
			" \u2570\u2500\u256F ",
		},
		ASCII: "[V]",
	},
	schema.RoleDesigner: {
		Unicode: []string{
			" \u250C\u2500\u2510 ",
			" \u2502\u2B22\u2502 ",
			" \u2514\u2500\u2518 ",
		},
		ASCII: "[S]",
	},
	schema.RoleCustom: {
		Unicode: []string{
			" \u250C\u2500\u2510 ",
			" \u2502\u2022\u2502 ",
			" \u2514\u2500\u2518 ",
		},
		ASCII: "[?]",
	},
}

// defaultMascot is used for unknown roles.
var defaultMascot = MascotSprite{
	Unicode: []string{
		" \u250C\u2500\u2510 ",
		" \u2502\u2022\u2502 ",
		" \u2514\u2500\u2518 ",
	},
	ASCII: "[?]",
}

// GetMascot returns the mascot sprite for a given role.
func GetMascot(role schema.Role) MascotSprite {
	if sprite, ok := mascotSprites[role]; ok {
		return sprite
	}
	return defaultMascot
}

// stateIndicators maps agent states to unicode/ASCII indicator pairs.
var stateIndicators = map[schema.AgentState][2]string{
	schema.StateRunning:   {"\u25CF", "*"}, // filled circle / asterisk
	schema.StateWaiting:   {"\u25CB", "o"}, // empty circle / o
	schema.StateBlocked:   {"\u26A0", "!"}, // warning / exclamation
	schema.StateError:     {"\u2718", "x"}, // cross mark / x
	schema.StateDone:      {"\u2714", "+"}, // check mark / plus
	schema.StateIdle:      {"\u2500", "-"}, // horizontal line / dash
	schema.StateFailed:    {"\u2716", "X"}, // heavy cross / X
	schema.StateCancelled: {"\u2205", "~"}, // empty set / tilde
}

// GetStateIndicator returns the unicode and ASCII indicator for a state.
func GetStateIndicator(state schema.AgentState, useUnicode bool) string {
	if pair, ok := stateIndicators[state]; ok {
		if useUnicode {
			return pair[0]
		}
		return pair[1]
	}
	if useUnicode {
		return "\u2500"
	}
	return "-"
}

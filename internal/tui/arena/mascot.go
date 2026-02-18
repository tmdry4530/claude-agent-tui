package arena

import "github.com/chamdom/omc-agent-tui/pkg/schema"

// SpriteWidth is the fixed character width of every sprite line (padded to widest).
const SpriteWidth = 9

// SpriteLines is the fixed 3-line CLCO mascot sprite.
// Block element characters (U+2580-U+259F), each line padded to SpriteWidth.
var SpriteLines = [3]string{
	" \u2590\u259B\u2588\u2588\u2588\u259C\u258C ",       // " ▐▛███▜▌ "
	"\u259D\u259C\u2588\u2588\u2588\u2588\u2588\u259B\u2598", // "▝▜█████▛▘"
	"  \u2598\u2598\u259D\u259D   ",                          // "  ▘▘▝▝   "
}

// ASCIIFallback is the 3-line ASCII fallback when block chars don't render.
// Each line is exactly SpriteWidth chars.
var ASCIIFallback = [3]string{
	" /=====\\ ",
	"|#######|",
	" \\     / ",
}

// GetSprite returns the 3-line sprite (unicode or ASCII fallback).
func GetSprite(useUnicode bool) [3]string {
	if useUnicode {
		return SpriteLines
	}
	return ASCIIFallback
}

// PadCenter pads a string to targetWidth, centering it with spaces.
func PadCenter(s string, targetWidth int) string {
	runes := []rune(s)
	sLen := len(runes)
	if sLen >= targetWidth {
		return s
	}
	leftPad := (targetWidth - sLen) / 2
	result := make([]rune, targetWidth)
	for i := range result {
		result[i] = ' '
	}
	copy(result[leftPad:], runes)
	return string(result)
}

// stateIndicators maps agent states to unicode/ASCII indicator pairs.
var stateIndicators = map[schema.AgentState][2]string{
	schema.StateRunning:   {"\u25CF", "*"},
	schema.StateWaiting:   {"\u25CB", "o"},
	schema.StateBlocked:   {"\u26A0", "!"},
	schema.StateError:     {"\u2718", "x"},
	schema.StateDone:      {"\u2714", "+"},
	schema.StateIdle:      {"\u2500", "-"},
	schema.StateFailed:    {"\u2716", "X"},
	schema.StateCancelled: {"\u2205", "~"},
}

// GetStateIndicator returns the indicator for a state.
func GetStateIndicator(state schema.AgentState, useUnicode bool) string {
	if pair, ok := stateIndicators[state]; ok {
		if useUnicode {
			return pair[0]
		}
		return pair[1]
	}
	return "-"
}

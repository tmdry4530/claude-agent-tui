package arena

import "github.com/chamdom/omc-agent-tui/pkg/schema"

// SpriteWidth is the fixed character width of every sprite line.
const SpriteWidth = 7

// SpriteLines is the fixed 3-line CLCO mascot sprite.
// Uses only monospace-safe ASCII characters. Each line is exactly SpriteWidth chars.
var SpriteLines = [3]string{
	" .===. ",
	" |@ @| ",
	" '-.-' ",
}

// ASCIIFallback is the 3-line ASCII fallback sprite.
// Each line is exactly SpriteWidth chars.
var ASCIIFallback = [3]string{
	" [===] ",
	" |o o| ",
	" [___] ",
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
	sLen := len(s) // safe because all chars are single-byte ASCII
	if sLen >= targetWidth {
		return s
	}
	leftPad := (targetWidth - sLen) / 2
	rightPad := targetWidth - sLen - leftPad
	result := make([]byte, targetWidth)
	for i := range result {
		result[i] = ' '
	}
	copy(result[leftPad:leftPad+sLen], s)
	_ = rightPad // rightPad is implicit via the fill
	return string(result)
}

// stateIndicators maps agent states to unicode/ASCII indicator pairs.
var stateIndicators = map[schema.AgentState][2]string{
	schema.StateRunning:   {"*", "*"},
	schema.StateWaiting:   {"o", "o"},
	schema.StateBlocked:   {"!", "!"},
	schema.StateError:     {"x", "x"},
	schema.StateDone:      {"+", "+"},
	schema.StateIdle:      {"-", "-"},
	schema.StateFailed:    {"X", "X"},
	schema.StateCancelled: {"~", "~"},
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

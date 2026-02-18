package schema

// ValidTransitions defines allowed state transitions.
// See references/event-schema.md section 4.
var ValidTransitions = map[AgentState][]AgentState{
	StateIdle:      {StateRunning, StateCancelled},
	StateRunning:   {StateWaiting, StateBlocked, StateError, StateDone, StateCancelled},
	StateWaiting:   {StateRunning, StateError},
	StateBlocked:   {StateRunning, StateError, StateCancelled},
	StateError:     {StateRunning, StateFailed},
	StateDone:      {StateIdle},
	StateFailed:    {},
	StateCancelled: {},
}

// IsValidTransition checks if a state transition is allowed.
func IsValidTransition(from, to AgentState) bool {
	allowed, ok := ValidTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

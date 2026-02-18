package schema

// Typed payload structs per event type.
// See references/event-schema.md section on payload structures.

type TaskSpawnPayload struct {
	Title      string `json:"title"`
	ChildAgent string `json:"child_agent"`
	Priority   *int   `json:"priority,omitempty"`
}

type TaskUpdatePayload struct {
	Progress int    `json:"progress"`
	Message  string `json:"message,omitempty"`
}

type TaskDonePayload struct {
	Result  string `json:"result"` // success|failure|cancelled
	Summary string `json:"summary,omitempty"`
}

type ToolCallPayload struct {
	ToolName string         `json:"tool_name"`
	Args     map[string]any `json:"args,omitempty"`
}

type ToolResultPayload struct {
	ToolName      string `json:"tool_name"`
	Success       bool   `json:"success"`
	OutputPreview string `json:"output_preview,omitempty"`
}

type ErrorPayload struct {
	ErrorType string `json:"error_type"`
	Message   string `json:"message"`
	Stack     string `json:"stack,omitempty"`
}

type VerifyPayload struct {
	Result   string `json:"result"` // pass|fail
	Reason   string `json:"reason,omitempty"`
	Evidence string `json:"evidence,omitempty"`
}

type FixPayload struct {
	Target       string   `json:"target"`
	Strategy     string   `json:"strategy"`
	FilesChanged []string `json:"files_changed,omitempty"`
}

type ReplanPayload struct {
	Reason     string `json:"reason"`
	NewPlanRef string `json:"new_plan_ref,omitempty"`
}

type RecoverPayload struct {
	Reason     string `json:"reason"`
	NewPlanRef string `json:"new_plan_ref,omitempty"`
}

type StateChangePayload struct {
	From    AgentState `json:"from"`
	To      AgentState `json:"to"`
	Trigger string     `json:"trigger,omitempty"`
}

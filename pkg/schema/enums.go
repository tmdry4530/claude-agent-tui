package schema

// Provider identifies the event source system.
type Provider string

const (
	ProviderClaude Provider = "claude"
	ProviderGemini Provider = "gemini"
	ProviderCodex  Provider = "codex"
	ProviderSystem Provider = "system"
)

var validProviders = map[Provider]bool{
	ProviderClaude: true,
	ProviderGemini: true,
	ProviderCodex:  true,
	ProviderSystem: true,
}

func (p Provider) IsValid() bool { return validProviders[p] }

// Mode represents the OMC execution mode.
type Mode string

const (
	ModeRalph     Mode = "ralph"
	ModeUltrawork Mode = "ultrawork"
	ModeUltrapilot Mode = "ultrapilot"
	ModeTeam      Mode = "team"
	ModeAutopilot Mode = "autopilot"
	ModePipeline  Mode = "pipeline"
	ModeEcomode   Mode = "ecomode"
	ModeUnknown   Mode = "unknown"
)

var validModes = map[Mode]bool{
	ModeRalph: true, ModeUltrawork: true, ModeUltrapilot: true,
	ModeTeam: true, ModeAutopilot: true, ModePipeline: true,
	ModeEcomode: true, ModeUnknown: true,
}

func (m Mode) IsValid() bool { return validModes[m] }

// Role represents the agent's functional role.
type Role string

const (
	RolePlanner  Role = "planner"
	RoleExecutor Role = "executor"
	RoleReviewer Role = "reviewer"
	RoleGuard    Role = "guard"
	RoleTester   Role = "tester"
	RoleWriter   Role = "writer"
	RoleExplorer Role = "explorer"
	RoleArchitect Role = "architect"
	RoleDebugger Role = "debugger"
	RoleVerifier Role = "verifier"
	RoleDesigner Role = "designer"
	RoleCustom   Role = "custom"
)

var validRoles = map[Role]bool{
	RolePlanner: true, RoleExecutor: true, RoleReviewer: true,
	RoleGuard: true, RoleTester: true, RoleWriter: true,
	RoleExplorer: true, RoleArchitect: true, RoleDebugger: true,
	RoleVerifier: true, RoleDesigner: true, RoleCustom: true,
}

func (r Role) IsValid() bool { return validRoles[r] }

// AgentState represents the agent's current state.
type AgentState string

const (
	StateIdle      AgentState = "idle"
	StateRunning   AgentState = "running"
	StateWaiting   AgentState = "waiting"
	StateBlocked   AgentState = "blocked"
	StateError     AgentState = "error"
	StateDone      AgentState = "done"
	StateFailed    AgentState = "failed"
	StateCancelled AgentState = "cancelled"
)

var validStates = map[AgentState]bool{
	StateIdle: true, StateRunning: true, StateWaiting: true,
	StateBlocked: true, StateError: true, StateDone: true,
	StateFailed: true, StateCancelled: true,
}

func (s AgentState) IsValid() bool    { return validStates[s] }
func (s AgentState) IsTerminal() bool { return s == StateDone || s == StateFailed || s == StateCancelled }

// EventType represents the kind of event that occurred.
type EventType string

const (
	TypeTaskSpawn   EventType = "task_spawn"
	TypeTaskUpdate  EventType = "task_update"
	TypeTaskDone    EventType = "task_done"
	TypeToolCall    EventType = "tool_call"
	TypeToolResult  EventType = "tool_result"
	TypeMessage     EventType = "message"
	TypeError       EventType = "error"
	TypeReplan      EventType = "replan"
	TypeVerify      EventType = "verify"
	TypeFix         EventType = "fix"
	TypeRecover     EventType = "recover"
	TypeStateChange EventType = "state_change"
)

var validTypes = map[EventType]bool{
	TypeTaskSpawn: true, TypeTaskUpdate: true, TypeTaskDone: true,
	TypeToolCall: true, TypeToolResult: true, TypeMessage: true,
	TypeError: true, TypeReplan: true, TypeVerify: true,
	TypeFix: true, TypeRecover: true, TypeStateChange: true,
}

func (t EventType) IsValid() bool { return validTypes[t] }

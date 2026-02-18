package schema

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCanonicalEvent_Validate(t *testing.T) {
	valid := CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: ProviderClaude,
		AgentID:  "coder-auth",
		Role:     RoleExecutor,
		State:    StateRunning,
		Type:     TypeToolCall,
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid, got: %v", err)
	}

	tests := []struct {
		name    string
		modify  func(*CanonicalEvent)
		wantErr string
	}{
		{"missing ts", func(e *CanonicalEvent) { e.Ts = time.Time{} }, "ts is required"},
		{"missing run_id", func(e *CanonicalEvent) { e.RunID = "" }, "run_id is required"},
		{"missing agent_id", func(e *CanonicalEvent) { e.AgentID = "" }, "agent_id is required"},
		{"invalid provider", func(e *CanonicalEvent) { e.Provider = "bad" }, `invalid provider: "bad"`},
		{"invalid role", func(e *CanonicalEvent) { e.Role = "bad" }, `invalid role: "bad"`},
		{"invalid state", func(e *CanonicalEvent) { e.State = "bad" }, `invalid state: "bad"`},
		{"invalid type", func(e *CanonicalEvent) { e.Type = "bad" }, `invalid type: "bad"`},
		{"invalid mode", func(e *CanonicalEvent) { e.Mode = "bad" }, `invalid mode: "bad"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := valid
			tt.modify(&e)
			err := e.Validate()
			if err == nil {
				t.Fatal("expected error")
			}
			if got := err.Error(); got != tt.wantErr {
				t.Fatalf("expected %q, got %q", tt.wantErr, got)
			}
		})
	}
}

func TestCanonicalEvent_OptionalMode(t *testing.T) {
	e := CanonicalEvent{
		Ts: time.Now(), RunID: "run-1", Provider: ProviderClaude,
		AgentID: "a", Role: RolePlanner, State: StateIdle, Type: TypeMessage,
	}
	if err := e.Validate(); err != nil {
		t.Fatalf("empty mode should be valid: %v", err)
	}
}

func TestCanonicalEvent_JSON(t *testing.T) {
	raw := `{
		"ts": "2026-02-17T22:27:00Z",
		"run_id": "run-1",
		"provider": "claude",
		"mode": "ralph",
		"agent_id": "coder-auth",
		"role": "executor",
		"state": "running",
		"type": "tool_call",
		"task_id": "task-42",
		"payload": {"tool_name": "Edit", "args": {"file": "auth.go"}},
		"metrics": {"latency_ms": 420, "tokens_in": 210, "tokens_out": 95, "cost_usd": 0.0021}
	}`
	var e CanonicalEvent
	if err := json.Unmarshal([]byte(raw), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := e.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
	if e.Provider != ProviderClaude {
		t.Fatalf("provider: got %q", e.Provider)
	}
	if e.Mode != ModeRalph {
		t.Fatalf("mode: got %q", e.Mode)
	}
	if e.Metrics == nil || *e.Metrics.LatencyMs != 420 {
		t.Fatal("metrics not parsed")
	}
}

func TestEnums_IsValid(t *testing.T) {
	if !ProviderClaude.IsValid() {
		t.Fatal("claude should be valid provider")
	}
	if Provider("nope").IsValid() {
		t.Fatal("nope should be invalid provider")
	}
	if !StateRunning.IsValid() {
		t.Fatal("running should be valid state")
	}
	if !StateDone.IsTerminal() {
		t.Fatal("done should be terminal")
	}
	if !StateFailed.IsTerminal() {
		t.Fatal("failed should be terminal")
	}
	if StateRunning.IsTerminal() {
		t.Fatal("running should not be terminal")
	}
}

func TestIsValidTransition(t *testing.T) {
	tests := []struct {
		from, to AgentState
		valid    bool
	}{
		{StateIdle, StateRunning, true},
		{StateRunning, StateDone, true},
		{StateRunning, StateCancelled, true},
		{StateError, StateRunning, true},
		{StateError, StateFailed, true},
		{StateDone, StateIdle, true},
		{StateFailed, StateRunning, false},
		{StateCancelled, StateRunning, false},
		{StateIdle, StateDone, false},
	}
	for _, tt := range tests {
		name := string(tt.from) + "->" + string(tt.to)
		t.Run(name, func(t *testing.T) {
			if got := IsValidTransition(tt.from, tt.to); got != tt.valid {
				t.Fatalf("expected %v", tt.valid)
			}
		})
	}
}

func TestLookupRole(t *testing.T) {
	r, ok := LookupRole("security-reviewer")
	if !ok || r != RoleGuard {
		t.Fatalf("expected guard, got %q (known=%v)", r, ok)
	}
	r, ok = LookupRole("unknown-agent")
	if ok || r != RoleCustom {
		t.Fatalf("expected custom/false, got %q/%v", r, ok)
	}
}

func TestPayload_Unmarshal(t *testing.T) {
	raw := `{"tool_name": "Edit", "success": true, "output_preview": "done"}`
	var p ToolResultPayload
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if p.ToolName != "Edit" || !p.Success {
		t.Fatalf("unexpected: %+v", p)
	}
}

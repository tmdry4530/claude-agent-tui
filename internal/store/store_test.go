package store

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
)

func TestNewStore(t *testing.T) {
	store := NewStore(100)
	if store.maxEvents != 100 {
		t.Errorf("expected maxEvents=100, got %d", store.maxEvents)
	}
	if len(store.events) != 100 {
		t.Errorf("expected events buffer size=100, got %d", len(store.events))
	}

	// Default size
	store2 := NewStore(0)
	if store2.maxEvents != 10000 {
		t.Errorf("expected default maxEvents=10000, got %d", store2.maxEvents)
	}
}

func TestRingBufferOverflow(t *testing.T) {
	store := NewStore(5)

	// Add 10 events (2x capacity)
	for i := 0; i < 10; i++ {
		event := schema.CanonicalEvent{
			Ts:       time.Now(),
			RunID:    "test-run",
			Provider: schema.ProviderClaude,
			AgentID:  "agent-1",
			Role:     schema.RoleExecutor,
			State:    schema.StateRunning,
			Type:     schema.TypeMessage,
		}
		store.AddEvent(event)
	}

	// Should only have 5 events (oldest dropped)
	if store.EventCount() != 5 {
		t.Errorf("expected count=5, got %d", store.EventCount())
	}

	events := store.GetEvents(100)
	if len(events) != 5 {
		t.Errorf("expected 5 events, got %d", len(events))
	}
}

func TestAgentStateUpdate(t *testing.T) {
	store := NewStore(100)

	// First event - create agent
	e1 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateIdle,
		Type:     schema.TypeStateChange,
	}
	store.AddEvent(e1)

	agent := store.GetAgent("agent-1")
	if agent == nil {
		t.Fatal("agent should exist")
	}
	if agent.State != schema.StateIdle {
		t.Errorf("expected state=idle, got %s", agent.State)
	}

	// Valid transition: idle -> running
	e2 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeStateChange,
	}
	store.AddEvent(e2)

	agent = store.GetAgent("agent-1")
	if agent.State != schema.StateRunning {
		t.Errorf("expected state=running, got %s", agent.State)
	}
	if store.GetWarningCount() != 0 {
		t.Errorf("expected 0 warnings for valid transition, got %d", store.GetWarningCount())
	}
}

func TestInvalidStateTransition(t *testing.T) {
	store := NewStore(100)

	// idle -> blocked is invalid
	e1 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateIdle,
		Type:     schema.TypeStateChange,
	}
	store.AddEvent(e1)

	e2 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateBlocked,
		Type:     schema.TypeStateChange,
	}
	store.AddEvent(e2)

	// Should have 1 warning, but event is still accepted
	if store.GetWarningCount() != 1 {
		t.Errorf("expected 1 warning, got %d", store.GetWarningCount())
	}

	agent := store.GetAgent("agent-1")
	if agent.State != schema.StateBlocked {
		t.Errorf("expected state=blocked (event accepted despite invalid transition), got %s", agent.State)
	}
}

func TestTaskLifecycle(t *testing.T) {
	store := NewStore(100)

	// task_spawn
	payload := schema.TaskSpawnPayload{
		Title:      "Test Task",
		ChildAgent: "agent-2",
	}
	payloadJSON, _ := json.Marshal(payload)

	e1 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		TaskID:   "task-1",
		Payload:  payloadJSON,
	}
	store.AddEvent(e1)

	task := store.GetTask("task-1")
	if task == nil {
		t.Fatal("task should exist")
	}
	if task.State != "active" {
		t.Errorf("expected state=active, got %s", task.State)
	}
	if task.Title != "Test Task" {
		t.Errorf("expected title='Test Task', got %s", task.Title)
	}

	// task_done with success
	donePayload := schema.TaskDonePayload{
		Result:  "success",
		Summary: "Task completed",
	}
	doneJSON, _ := json.Marshal(donePayload)

	e2 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskDone,
		TaskID:   "task-1",
		Payload:  doneJSON,
	}
	store.AddEvent(e2)

	task = store.GetTask("task-1")
	if task.State != "done" {
		t.Errorf("expected state=done, got %s", task.State)
	}
}

func TestTaskDoneFailure(t *testing.T) {
	store := NewStore(100)

	// Spawn task
	payload := schema.TaskSpawnPayload{Title: "Failing Task"}
	payloadJSON, _ := json.Marshal(payload)

	e1 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		TaskID:   "task-1",
		Payload:  payloadJSON,
	}
	store.AddEvent(e1)

	// task_done with failure
	donePayload := schema.TaskDonePayload{Result: "failure"}
	doneJSON, _ := json.Marshal(donePayload)

	e2 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskDone,
		TaskID:   "task-1",
		Payload:  doneJSON,
	}
	store.AddEvent(e2)

	task := store.GetTask("task-1")
	if task.State != "failed" {
		t.Errorf("expected state=failed, got %s", task.State)
	}
}

func TestTaskDoneCancelled(t *testing.T) {
	store := NewStore(100)

	payload := schema.TaskSpawnPayload{Title: "Cancelled Task"}
	payloadJSON, _ := json.Marshal(payload)

	e1 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		TaskID:   "task-1",
		Payload:  payloadJSON,
	}
	store.AddEvent(e1)

	donePayload := schema.TaskDonePayload{Result: "cancelled"}
	doneJSON, _ := json.Marshal(donePayload)

	e2 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskDone,
		TaskID:   "task-1",
		Payload:  doneJSON,
	}
	store.AddEvent(e2)

	task := store.GetTask("task-1")
	if task.State != "cancelled" {
		t.Errorf("expected state=cancelled, got %s", task.State)
	}
}

func TestMetricsAggregation(t *testing.T) {
	store := NewStore(100)

	latency1 := 100.5
	tokensIn1 := 50
	tokensOut1 := 100
	cost1 := 0.01

	e1 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeMessage,
		Metrics: &schema.EventMetrics{
			LatencyMs: &latency1,
			TokensIn:  &tokensIn1,
			TokensOut: &tokensOut1,
			CostUSD:   &cost1,
		},
	}
	store.AddEvent(e1)

	latency2 := 200.5
	tokensIn2 := 75
	tokensOut2 := 150
	cost2 := 0.02

	e2 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeMessage,
		Metrics: &schema.EventMetrics{
			LatencyMs: &latency2,
			TokensIn:  &tokensIn2,
			TokensOut: &tokensOut2,
			CostUSD:   &cost2,
		},
	}
	store.AddEvent(e2)

	metrics := store.GetMetrics()
	if metrics.EventCount != 2 {
		t.Errorf("expected EventCount=2, got %d", metrics.EventCount)
	}
	if metrics.TotalLatency != 301.0 {
		t.Errorf("expected TotalLatency=301.0, got %f", metrics.TotalLatency)
	}
	if metrics.TotalTokensIn != 125 {
		t.Errorf("expected TotalTokensIn=125, got %d", metrics.TotalTokensIn)
	}
	if metrics.TotalTokensOut != 250 {
		t.Errorf("expected TotalTokensOut=250, got %d", metrics.TotalTokensOut)
	}
	if metrics.TotalCostUSD != 0.03 {
		t.Errorf("expected TotalCostUSD=0.03, got %f", metrics.TotalCostUSD)
	}
}

func TestErrorCount(t *testing.T) {
	store := NewStore(100)

	e1 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeError,
	}
	store.AddEvent(e1)

	e2 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeMessage,
	}
	store.AddEvent(e2)

	metrics := store.GetMetrics()
	if metrics.ErrorCount != 1 {
		t.Errorf("expected ErrorCount=1, got %d", metrics.ErrorCount)
	}
}

func TestConcurrency(t *testing.T) {
	store := NewStore(1000)
	var wg sync.WaitGroup

	// 10 goroutines, each adding 100 events
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				event := schema.CanonicalEvent{
					Ts:       time.Now(),
					RunID:    "concurrent-run",
					Provider: schema.ProviderClaude,
					AgentID:  "agent-1",
					Role:     schema.RoleExecutor,
					State:    schema.StateRunning,
					Type:     schema.TypeMessage,
				}
				store.AddEvent(event)
			}
		}(i)
	}

	wg.Wait()

	// Should have 1000 events (ring buffer capacity)
	if store.EventCount() != 1000 {
		t.Errorf("expected count=1000, got %d", store.EventCount())
	}

	metrics := store.GetMetrics()
	if metrics.EventCount != 1000 {
		t.Errorf("expected EventCount=1000, got %d", metrics.EventCount)
	}
}

func TestGetAllAgents(t *testing.T) {
	store := NewStore(100)

	e1 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeMessage,
	}
	store.AddEvent(e1)

	e2 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-2",
		Role:     schema.RolePlanner,
		State:    schema.StateIdle,
		Type:     schema.TypeMessage,
	}
	store.AddEvent(e2)

	agents := store.GetAllAgents()
	if len(agents) != 2 {
		t.Errorf("expected 2 agents, got %d", len(agents))
	}
}

func TestGetAllTasks(t *testing.T) {
	store := NewStore(100)

	payload1 := schema.TaskSpawnPayload{Title: "Task 1"}
	payload1JSON, _ := json.Marshal(payload1)

	e1 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		TaskID:   "task-1",
		Payload:  payload1JSON,
	}
	store.AddEvent(e1)

	payload2 := schema.TaskSpawnPayload{Title: "Task 2"}
	payload2JSON, _ := json.Marshal(payload2)

	e2 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-1",
		Provider: schema.ProviderClaude,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeTaskSpawn,
		TaskID:   "task-2",
		Payload:  payload2JSON,
	}
	store.AddEvent(e2)

	tasks := store.GetAllTasks()
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestModeAndRunID(t *testing.T) {
	store := NewStore(100)

	e1 := schema.CanonicalEvent{
		Ts:       time.Now(),
		RunID:    "run-123",
		Provider: schema.ProviderClaude,
		Mode:     schema.ModeRalph,
		AgentID:  "agent-1",
		Role:     schema.RoleExecutor,
		State:    schema.StateRunning,
		Type:     schema.TypeMessage,
	}
	store.AddEvent(e1)

	if store.GetRunID() != "run-123" {
		t.Errorf("expected RunID='run-123', got %s", store.GetRunID())
	}
	if store.GetMode() != schema.ModeRalph {
		t.Errorf("expected Mode=ralph, got %s", store.GetMode())
	}
}

func TestGetEventsLimit(t *testing.T) {
	store := NewStore(100)

	// Add 10 events
	for i := 0; i < 10; i++ {
		e := schema.CanonicalEvent{
			Ts:       time.Now(),
			RunID:    "run-1",
			Provider: schema.ProviderClaude,
			AgentID:  "agent-1",
			Role:     schema.RoleExecutor,
			State:    schema.StateRunning,
			Type:     schema.TypeMessage,
		}
		store.AddEvent(e)
	}

	// Get last 5
	events := store.GetEvents(5)
	if len(events) != 5 {
		t.Errorf("expected 5 events, got %d", len(events))
	}

	// Get all
	events = store.GetEvents(0)
	if len(events) != 10 {
		t.Errorf("expected 10 events, got %d", len(events))
	}

	// Get more than available
	events = store.GetEvents(20)
	if len(events) != 10 {
		t.Errorf("expected 10 events, got %d", len(events))
	}
}

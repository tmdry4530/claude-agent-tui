package store

import (
	"log"
	"sync"
	"time"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
)

// Store holds all runtime state for the TUI.
// Thread-safe via RWMutex.
type Store struct {
	mu        sync.RWMutex
	events    []schema.CanonicalEvent // ring buffer
	maxEvents int
	writeIdx  int // next write position
	count     int // actual stored count

	agents  map[string]*AgentInfo
	tasks   map[string]*TaskInfo
	metrics Metrics

	runID     string
	mode      schema.Mode
	warnCount int // invalid transition warnings
}

// AgentInfo tracks the current state of an agent.
type AgentInfo struct {
	AgentID  string
	Role     schema.Role
	State    schema.AgentState
	LastSeen time.Time
}

// TaskInfo tracks task lifecycle.
type TaskInfo struct {
	TaskID   string
	AgentID  string
	State    string // "active" | "done" | "failed" | "cancelled"
	Title    string
	Created  time.Time
	Updated  time.Time
}

// Metrics aggregates performance and cost data.
type Metrics struct {
	EventCount     int
	ErrorCount     int
	TotalLatency   float64 // sum of latency_ms
	TotalTokensIn  int
	TotalTokensOut int
	TotalCostUSD   float64
}

// NewStore creates a new Store with the specified ring buffer size.
// Default maxEvents is 10,000.
func NewStore(maxEvents int) *Store {
	if maxEvents <= 0 {
		maxEvents = 10000
	}
	return &Store{
		events:    make([]schema.CanonicalEvent, maxEvents),
		maxEvents: maxEvents,
		agents:    make(map[string]*AgentInfo),
		tasks:     make(map[string]*TaskInfo),
	}
}

// AddEvent stores an event in the ring buffer and updates state.
// Invalid state transitions are logged as warnings but still accepted.
func (s *Store) AddEvent(event schema.CanonicalEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update run metadata
	if event.RunID != "" {
		s.runID = event.RunID
	}
	if event.Mode != "" {
		s.mode = event.Mode
	}

	// Store in ring buffer
	s.events[s.writeIdx] = event
	s.writeIdx = (s.writeIdx + 1) % s.maxEvents
	if s.count < s.maxEvents {
		s.count++
	}

	// Update agent state
	s.updateAgent(event)

	// Handle task lifecycle
	s.updateTask(event)

	// Aggregate metrics
	s.updateMetrics(event)
}

// updateAgent updates or creates agent state.
// Validates state transitions and logs warnings for invalid transitions.
func (s *Store) updateAgent(event schema.CanonicalEvent) {
	agent, exists := s.agents[event.AgentID]

	// Validate state transition
	if exists {
		if !schema.IsValidTransition(agent.State, event.State) {
			s.warnCount++
			log.Printf("[WARN] invalid state transition for agent %s: %s -> %s",
				event.AgentID, agent.State, event.State)
		}
	}

	// Update or create
	if exists {
		agent.State = event.State
		agent.LastSeen = event.Ts
		if event.Role != "" {
			agent.Role = event.Role
		}
	} else {
		s.agents[event.AgentID] = &AgentInfo{
			AgentID:  event.AgentID,
			Role:     event.Role,
			State:    event.State,
			LastSeen: event.Ts,
		}
	}
}

// updateTask handles task lifecycle events.
func (s *Store) updateTask(event schema.CanonicalEvent) {
	if event.TaskID == "" {
		return
	}

	switch event.Type {
	case schema.TypeTaskSpawn:
		var payload schema.TaskSpawnPayload
		if err := parsePayload(event.Payload, &payload); err == nil {
			s.tasks[event.TaskID] = &TaskInfo{
				TaskID:  event.TaskID,
				AgentID: event.AgentID,
				State:   "active",
				Title:   payload.Title,
				Created: event.Ts,
				Updated: event.Ts,
			}
		}

	case schema.TypeTaskDone:
		if task, ok := s.tasks[event.TaskID]; ok {
			var payload schema.TaskDonePayload
			if err := parsePayload(event.Payload, &payload); err == nil {
				switch payload.Result {
				case "success":
					task.State = "done"
				case "failure":
					task.State = "failed"
				case "cancelled":
					task.State = "cancelled"
				default:
					task.State = "done"
				}
				task.Updated = event.Ts
			}
		}

	case schema.TypeTaskUpdate:
		if task, ok := s.tasks[event.TaskID]; ok {
			task.Updated = event.Ts
		}
	}
}

// updateMetrics aggregates event metrics.
func (s *Store) updateMetrics(event schema.CanonicalEvent) {
	s.metrics.EventCount++

	if event.Type == schema.TypeError {
		s.metrics.ErrorCount++
	}

	if event.Metrics != nil {
		if event.Metrics.LatencyMs != nil {
			s.metrics.TotalLatency += *event.Metrics.LatencyMs
		}
		if event.Metrics.TokensIn != nil {
			s.metrics.TotalTokensIn += *event.Metrics.TokensIn
		}
		if event.Metrics.TokensOut != nil {
			s.metrics.TotalTokensOut += *event.Metrics.TokensOut
		}
		if event.Metrics.CostUSD != nil {
			s.metrics.TotalCostUSD += *event.Metrics.CostUSD
		}
	}
}

// GetEvents returns the most recent events, up to limit.
// Returns events in chronological order (oldest first).
func (s *Store) GetEvents(limit int) []schema.CanonicalEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 || limit > s.count {
		limit = s.count
	}

	result := make([]schema.CanonicalEvent, 0, limit)

	if s.count < s.maxEvents {
		// Buffer not full yet - read from start
		start := s.count - limit
		if start < 0 {
			start = 0
		}
		for i := start; i < s.count; i++ {
			result = append(result, s.events[i])
		}
	} else {
		// Buffer is full - read from oldest
		start := (s.writeIdx - limit + s.maxEvents) % s.maxEvents
		for i := 0; i < limit; i++ {
			idx := (start + i) % s.maxEvents
			result = append(result, s.events[idx])
		}
	}

	return result
}

// GetAgent returns info for a specific agent.
func (s *Store) GetAgent(agentID string) *AgentInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.agents[agentID]
}

// GetAllAgents returns all agent info.
func (s *Store) GetAllAgents() []*AgentInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*AgentInfo, 0, len(s.agents))
	for _, agent := range s.agents {
		result = append(result, agent)
	}
	return result
}

// GetTask returns info for a specific task.
func (s *Store) GetTask(taskID string) *TaskInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tasks[taskID]
}

// GetAllTasks returns all task info.
func (s *Store) GetAllTasks() []*TaskInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*TaskInfo, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}
	return result
}

// GetMetrics returns aggregated metrics.
func (s *Store) GetMetrics() Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.metrics
}

// GetMode returns the current execution mode.
func (s *Store) GetMode() schema.Mode {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.mode
}

// GetRunID returns the current run ID.
func (s *Store) GetRunID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.runID
}

// GetWarningCount returns the number of invalid transitions detected.
func (s *Store) GetWarningCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.warnCount
}

// EventCount returns the total number of events stored.
func (s *Store) EventCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.count
}

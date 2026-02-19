package replay

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
)

func createTestJSONL(t *testing.T, dir string, events []schema.CanonicalEvent) string {
	t.Helper()
	path := filepath.Join(dir, "test.jsonl")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create file: %v", err)
	}
	defer func() { _ = f.Close() }()

	for _, evt := range events {
		data, err := json.Marshal(evt)
		if err != nil {
			t.Fatalf("marshal event: %v", err)
		}
		if _, err := f.Write(append(data, '\n')); err != nil {
			t.Fatalf("write line: %v", err)
		}
	}

	return path
}

func TestLoadFile_SortsByTimestamp(t *testing.T) {
	dir := t.TempDir()

	// Create events in reverse order
	events := []schema.CanonicalEvent{
		{
			Ts:       time.Date(2026, 2, 17, 22, 27, 2, 0, time.UTC),
			RunID:    "run-1",
			Provider: "claude",
			AgentID:  "a1",
			Role:     "executor",
			State:    "running",
			Type:     "tool_result",
		},
		{
			Ts:       time.Date(2026, 2, 17, 22, 27, 0, 0, time.UTC),
			RunID:    "run-1",
			Provider: "claude",
			AgentID:  "a1",
			Role:     "executor",
			State:    "running",
			Type:     "tool_call",
		},
		{
			Ts:       time.Date(2026, 2, 17, 22, 27, 1, 0, time.UTC),
			RunID:    "run-1",
			Provider: "claude",
			AgentID:  "a1",
			Role:     "executor",
			State:    "running",
			Type:     "message",
		},
	}

	path := createTestJSONL(t, dir, events)

	player := NewPlayer()
	if err := player.LoadFile(path); err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	if player.Total() != 3 {
		t.Errorf("expected 3 events, got %d", player.Total())
	}

	// Verify sorted order
	expected := []string{"tool_call", "message", "tool_result"}
	for i, expectedType := range expected {
		if string(player.events[i].Type) != expectedType {
			t.Errorf("events[%d].Type = %s, want %s", i, player.events[i].Type, expectedType)
		}
	}
}

func TestStepForwardBackward(t *testing.T) {
	dir := t.TempDir()
	events := []schema.CanonicalEvent{
		{
			Ts:       time.Date(2026, 2, 17, 22, 27, 0, 0, time.UTC),
			RunID:    "run-1",
			Provider: "claude",
			AgentID:  "a1",
			Role:     "executor",
			State:    "running",
			Type:     "tool_call",
		},
		{
			Ts:       time.Date(2026, 2, 17, 22, 27, 1, 0, time.UTC),
			RunID:    "run-1",
			Provider: "claude",
			AgentID:  "a1",
			Role:     "executor",
			State:    "running",
			Type:     "tool_result",
		},
	}

	path := createTestJSONL(t, dir, events)
	player := NewPlayer()
	if err := player.LoadFile(path); err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	// Initial position
	if player.Position() != 0 {
		t.Errorf("initial position = %d, want 0", player.Position())
	}

	// Step forward
	player.StepForward()
	if player.Position() != 1 {
		t.Errorf("after StepForward, position = %d, want 1", player.Position())
	}

	// Step forward at boundary (should stay at 1)
	player.StepForward()
	if player.Position() != 1 {
		t.Errorf("after StepForward at end, position = %d, want 1", player.Position())
	}

	// Step backward
	player.StepBackward()
	if player.Position() != 0 {
		t.Errorf("after StepBackward, position = %d, want 0", player.Position())
	}

	// Step backward at boundary (should stay at 0)
	player.StepBackward()
	if player.Position() != 0 {
		t.Errorf("after StepBackward at start, position = %d, want 0", player.Position())
	}
}

func TestSeek_BoundaryConditions(t *testing.T) {
	dir := t.TempDir()
	events := []schema.CanonicalEvent{
		{
			Ts:       time.Date(2026, 2, 17, 22, 27, 0, 0, time.UTC),
			RunID:    "run-1",
			Provider: "claude",
			AgentID:  "a1",
			Role:     "executor",
			State:    "running",
			Type:     "tool_call",
		},
		{
			Ts:       time.Date(2026, 2, 17, 22, 27, 1, 0, time.UTC),
			RunID:    "run-1",
			Provider: "claude",
			AgentID:  "a1",
			Role:     "executor",
			State:    "running",
			Type:     "tool_result",
		},
		{
			Ts:       time.Date(2026, 2, 17, 22, 27, 2, 0, time.UTC),
			RunID:    "run-1",
			Provider: "claude",
			AgentID:  "a1",
			Role:     "executor",
			State:    "running",
			Type:     "message",
		},
	}

	path := createTestJSONL(t, dir, events)
	player := NewPlayer()
	if err := player.LoadFile(path); err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	tests := []struct {
		name     string
		seek     int
		expected int
	}{
		{"negative", -10, 0},
		{"zero", 0, 0},
		{"middle", 1, 1},
		{"last", 2, 2},
		{"overflow", 100, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player.Seek(tt.seek)
			if player.Position() != tt.expected {
				t.Errorf("Seek(%d): position = %d, want %d", tt.seek, player.Position(), tt.expected)
			}
		})
	}
}

func TestPlayPauseStop(t *testing.T) {
	dir := t.TempDir()
	events := []schema.CanonicalEvent{
		{
			Ts:       time.Date(2026, 2, 17, 22, 27, 0, 0, time.UTC),
			RunID:    "run-1",
			Provider: "claude",
			AgentID:  "a1",
			Role:     "executor",
			State:    "running",
			Type:     "tool_call",
		},
	}

	path := createTestJSONL(t, dir, events)
	player := NewPlayer()
	if err := player.LoadFile(path); err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	// Initial state: not playing
	if player.IsPlaying() {
		t.Error("expected not playing initially")
	}

	// Play
	player.Play()
	if !player.IsPlaying() {
		t.Error("expected playing after Play()")
	}

	// Pause
	player.Pause()
	if player.IsPlaying() {
		t.Error("expected not playing after Pause()")
	}

	// Play again
	player.Play()
	player.StepForward() // move position

	// Stop resets to beginning
	player.Stop()
	if player.IsPlaying() {
		t.Error("expected not playing after Stop()")
	}
	if player.Position() != 0 {
		t.Errorf("expected position 0 after Stop(), got %d", player.Position())
	}
}

func TestSetSpeed(t *testing.T) {
	player := NewPlayer()

	tests := []struct {
		name     string
		speed    float64
		expected float64
	}{
		{"1x", 1.0, 1.0},
		{"4x", 4.0, 4.0},
		{"8x", 8.0, 8.0},
		{"16x", 16.0, 16.0},
		{"negative becomes 1x", -5.0, 1.0},
		{"zero becomes 1x", 0.0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player.SetSpeed(tt.speed)
			if player.Speed() != tt.expected {
				t.Errorf("SetSpeed(%f): speed = %f, want %f", tt.speed, player.Speed(), tt.expected)
			}
		})
	}
}

func TestLoadFile_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.jsonl")
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatalf("create empty file: %v", err)
	}

	player := NewPlayer()
	if err := player.LoadFile(path); err != nil {
		t.Fatalf("LoadFile on empty file should succeed: %v", err)
	}

	if player.Total() != 0 {
		t.Errorf("expected 0 events, got %d", player.Total())
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.jsonl")
	if err := os.WriteFile(path, []byte("not json\n"), 0644); err != nil {
		t.Fatalf("create invalid file: %v", err)
	}

	player := NewPlayer()
	if err := player.LoadFile(path); err == nil {
		t.Error("expected error on invalid JSON")
	}
}

func TestLoadFile_InvalidEvent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid-event.jsonl")

	// Missing required fields
	invalidEvent := map[string]interface{}{
		"ts": "2026-02-17T22:27:00Z",
		// missing run_id, agent_id, etc.
	}

	data, _ := json.Marshal(invalidEvent)
	if err := os.WriteFile(path, append(data, '\n'), 0644); err != nil {
		t.Fatalf("create file: %v", err)
	}

	player := NewPlayer()
	if err := player.LoadFile(path); err == nil {
		t.Error("expected error on invalid event")
	}
}

func TestCurrentEvent(t *testing.T) {
	dir := t.TempDir()
	events := []schema.CanonicalEvent{
		{
			Ts:       time.Date(2026, 2, 17, 22, 27, 0, 0, time.UTC),
			RunID:    "run-1",
			Provider: "claude",
			AgentID:  "a1",
			Role:     "executor",
			State:    "running",
			Type:     "tool_call",
		},
		{
			Ts:       time.Date(2026, 2, 17, 22, 27, 1, 0, time.UTC),
			RunID:    "run-1",
			Provider: "claude",
			AgentID:  "a1",
			Role:     "executor",
			State:    "running",
			Type:     "tool_result",
		},
	}

	path := createTestJSONL(t, dir, events)
	player := NewPlayer()
	if err := player.LoadFile(path); err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	// Position 0
	evt := player.CurrentEvent()
	if evt == nil {
		t.Fatal("expected event at position 0")
	}
	if evt.Type != "tool_call" {
		t.Errorf("CurrentEvent at 0: type = %s, want tool_call", evt.Type)
	}

	// Position 1
	player.StepForward()
	evt = player.CurrentEvent()
	if evt == nil {
		t.Fatal("expected event at position 1")
	}
	if evt.Type != "tool_result" {
		t.Errorf("CurrentEvent at 1: type = %s, want tool_result", evt.Type)
	}
}

func TestEventsUntil_VirtualClock(t *testing.T) {
	dir := t.TempDir()
	baseTime := time.Date(2026, 2, 17, 22, 27, 0, 0, time.UTC)

	events := []schema.CanonicalEvent{
		{
			Ts:       baseTime,
			RunID:    "run-1",
			Provider: "claude",
			AgentID:  "a1",
			Role:     "executor",
			State:    "running",
			Type:     "tool_call",
		},
		{
			Ts:       baseTime.Add(2 * time.Second),
			RunID:    "run-1",
			Provider: "claude",
			AgentID:  "a1",
			Role:     "executor",
			State:    "running",
			Type:     "message",
		},
		{
			Ts:       baseTime.Add(4 * time.Second),
			RunID:    "run-1",
			Provider: "claude",
			AgentID:  "a1",
			Role:     "executor",
			State:    "running",
			Type:     "tool_result",
		},
	}

	path := createTestJSONL(t, dir, events)
	player := NewPlayer()
	if err := player.LoadFile(path); err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	// Not playing: returns nil
	result := player.EventsUntil(time.Now())
	if result != nil {
		t.Error("EventsUntil should return nil when not playing")
	}

	// Start playback at 1x speed
	player.Play()

	// Immediately after play start: should return first event (at baseTime)
	result = player.EventsUntil(time.Now())
	if len(result) != 1 {
		t.Errorf("expected 1 event immediately, got %d", len(result))
	}

	// Simulate 1 second elapsed (1x speed = 1 virtual second)
	// Virtual time = baseTime + 1s, should still have 1 event
	virtualNow := player.startTime.Add(1 * time.Second)
	result = player.EventsUntil(virtualNow)
	if len(result) != 1 {
		t.Errorf("expected 1 event at +1s virtual, got %d", len(result))
	}

	// Simulate 2 seconds elapsed (virtual = baseTime + 2s)
	// Should include events at baseTime and baseTime+2s (2 events)
	virtualNow = player.startTime.Add(2 * time.Second)
	result = player.EventsUntil(virtualNow)
	if len(result) != 2 {
		t.Errorf("expected 2 events at +2s virtual, got %d", len(result))
	}

	// Simulate 5 seconds elapsed (virtual = baseTime + 5s)
	// Should include all 3 events
	virtualNow = player.startTime.Add(5 * time.Second)
	result = player.EventsUntil(virtualNow)
	if len(result) != 3 {
		t.Errorf("expected 3 events at +5s virtual, got %d", len(result))
	}

	// Test 4x speed
	player.Stop()
	player.SetSpeed(4.0)
	player.Play()

	// 1 second real = 4 seconds virtual
	// Virtual = baseTime + 4s, should have all 3 events
	virtualNow = player.startTime.Add(1 * time.Second)
	result = player.EventsUntil(virtualNow)
	if len(result) != 3 {
		t.Errorf("expected 3 events at 4x speed +1s real, got %d", len(result))
	}
}

func TestEventsUntil_EmptyEvents(t *testing.T) {
	player := NewPlayer()
	player.Play()

	result := player.EventsUntil(time.Now())
	if result != nil {
		t.Errorf("expected nil for empty events, got %d events", len(result))
	}
}

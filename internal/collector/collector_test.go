package collector

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// drainEvents reads events from the collector channel until the expected count
// is reached or the context is cancelled.
func drainEvents(ctx context.Context, fc *FileCollector, want int) (int, error) {
	received := 0
	for received < want {
		select {
		case _, ok := <-fc.Events():
			if !ok {
				return received, nil
			}
			received++
		case <-ctx.Done():
			return received, ctx.Err()
		}
	}
	return received, nil
}

func TestFileCollector_NewLines(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "events.jsonl")

	fc := NewFileCollector(tmpDir)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := fc.Start(ctx); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	defer fc.Stop()

	testData := []map[string]interface{}{
		{"event": "test1", "value": 100},
		{"event": "test2", "value": 200},
		{"event": "test3", "value": 300},
	}

	// Wait for fsnotify watcher to be fully set up
	time.Sleep(500 * time.Millisecond)

	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	for _, data := range testData {
		jsonBytes, _ := json.Marshal(data)
		_, _ = f.Write(jsonBytes)
		_, _ = f.Write([]byte("\n"))
	}
	_ = f.Sync()
	f.Close()

	// Drain expected events
	got, err := drainEvents(ctx, fc, len(testData))
	if err != nil {
		t.Fatalf("timeout: received %d/%d events", got, len(testData))
	}

	if got != len(testData) {
		t.Errorf("received %d events, want %d", got, len(testData))
	}
}

func TestFileCollector_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "invalid.jsonl")

	fc := NewFileCollector(tmpDir)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := fc.Start(ctx); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	defer fc.Stop()

	time.Sleep(500 * time.Millisecond)

	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	_, _ = f.Write([]byte("invalid json line\n"))
	_, _ = f.Write([]byte(`{"valid": "json"}` + "\n"))
	_, _ = f.Write([]byte("another invalid\n"))
	_, _ = f.Write([]byte(`{"also": "valid"}` + "\n"))
	_ = f.Sync()
	f.Close()

	// Only valid JSON lines should be received (2 out of 4)
	got, err := drainEvents(ctx, fc, 2)
	if err != nil {
		t.Fatalf("timeout: received %d/2 events", got)
	}

	if got != 2 {
		t.Errorf("received %d events, want 2 (invalid JSON filtered)", got)
	}
}

func TestFileCollector_ContextCancellation(t *testing.T) {
	tmpDir := t.TempDir()

	fc := NewFileCollector(tmpDir)

	ctx, cancel := context.WithCancel(context.Background())

	if err := fc.Start(ctx); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	cancel()

	// Stop should complete within timeout
	done := make(chan struct{})
	go func() {
		fc.Stop()
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(3 * time.Second):
		t.Fatal("Stop() did not complete within timeout")
	}

	// Channel should be closed after Stop
	select {
	case _, ok := <-fc.Events():
		if ok {
			t.Error("events channel still open after Stop()")
		}
	default:
		// empty channel is ok
	}
}

func TestFileCollector_AppendToFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "append.jsonl")

	fc := NewFileCollector(tmpDir)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := fc.Start(ctx); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	defer fc.Stop()

	time.Sleep(500 * time.Millisecond)

	// Write initial line
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	_, _ = f.Write([]byte(`{"line": 1}` + "\n"))
	_ = f.Sync()
	f.Close()

	// Receive first event
	got, err := drainEvents(ctx, fc, 1)
	if err != nil {
		t.Fatalf("timeout waiting for first event: received %d", got)
	}

	// Append second line
	f, err = os.OpenFile(testFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("failed to open test file for append: %v", err)
	}
	_, _ = f.Write([]byte(`{"line": 2}` + "\n"))
	_ = f.Sync()
	f.Close()

	// Receive second event
	select {
	case event := <-fc.Events():
		var data map[string]interface{}
		if err := json.Unmarshal(event.Data, &data); err != nil {
			t.Fatalf("JSON parse failed: %v", err)
		}
		if line, ok := data["line"].(float64); !ok || line != 2 {
			t.Errorf("second line = %v, want 2", data["line"])
		}
	case <-ctx.Done():
		t.Fatal("timeout waiting for second event")
	}
}

func TestFileCollector_CircuitBreaker(t *testing.T) {
	fc := NewFileCollector("/nonexistent")

	// Initially not backed off
	if fc.isBackedOff() {
		t.Error("should not be backed off initially")
	}

	// Record 2 failures - not yet backed off (threshold is 3)
	fc.recordFailure()
	fc.recordFailure()
	if fc.isBackedOff() {
		t.Error("should not be backed off after 2 failures")
	}

	// 3rd failure triggers backoff
	fc.recordFailure()
	if !fc.isBackedOff() {
		t.Error("should be backed off after 3 failures")
	}

	// Success resets everything
	fc.recordSuccess()
	if fc.isBackedOff() {
		t.Error("should not be backed off after success")
	}
	if fc.failCount != 0 {
		t.Errorf("failCount = %d, want 0 after success", fc.failCount)
	}
}

func TestFileCollector_BackoffLevels(t *testing.T) {
	fc := NewFileCollector("/nonexistent")

	// 3 failures → level 0 backoff (10s)
	for i := 0; i < 3; i++ {
		fc.recordFailure()
	}

	fc.mu.Lock()
	backoff1 := fc.backoffUntil
	fc.mu.Unlock()

	// Reset and trigger level 1 (4 failures → 30s)
	fc.recordSuccess()
	for i := 0; i < 4; i++ {
		fc.recordFailure()
	}

	fc.mu.Lock()
	backoff2 := fc.backoffUntil
	fc.mu.Unlock()

	// Level 1 backoff should be later than level 0
	if !backoff2.After(backoff1) {
		t.Error("level 1 backoff should be longer than level 0")
	}
}

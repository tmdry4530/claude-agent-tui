package replay

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
)

const maxFileSize = 100 * 1024 * 1024 // 100MB

// Player manages JSONL event replay with virtual clock and step navigation.
type Player struct {
	events    []schema.CanonicalEvent
	position  int
	speed     float64
	playing   bool
	startTime time.Time       // real-world time when playback started
	pauseTime time.Time       // real-world time when paused
	baseTime  time.Time       // virtual time at position 0
	mu        sync.RWMutex
}

// NewPlayer creates a new replay player.
func NewPlayer() *Player {
	return &Player{
		events:   make([]schema.CanonicalEvent, 0),
		position: 0,
		speed:    1.0,
		playing:  false,
	}
}

// LoadFile loads events from a JSONL file and sorts them by timestamp.
// Returns error if file > 100MB (requires streaming mode).
func (p *Player) LoadFile(path string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check file size
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}
	if info.Size() > maxFileSize {
		return fmt.Errorf("file size %d exceeds max %d (streaming mode required)", info.Size(), maxFileSize)
	}

	// Open and parse JSONL
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var events []schema.CanonicalEvent
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue // skip empty lines
		}

		var evt schema.CanonicalEvent
		if err := json.Unmarshal(line, &evt); err != nil {
			return fmt.Errorf("line %d: invalid JSON: %w", lineNum, err)
		}

		if err := evt.Validate(); err != nil {
			return fmt.Errorf("line %d: invalid event: %w", lineNum, err)
		}

		events = append(events, evt)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan file: %w", err)
	}

	// Sort by timestamp
	sort.Slice(events, func(i, j int) bool {
		return events[i].Ts.Before(events[j].Ts)
	})

	p.events = events
	p.position = 0
	p.playing = false

	// Set base time to first event timestamp
	if len(p.events) > 0 {
		p.baseTime = p.events[0].Ts
	}

	return nil
}

// Play starts playback from current position.
func (p *Player) Play() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.playing {
		return
	}

	p.playing = true
	p.startTime = time.Now()

	// If resuming from pause, adjust startTime to account for elapsed virtual time
	if p.position > 0 && len(p.events) > 0 {
		virtualElapsed := p.events[p.position].Ts.Sub(p.baseTime)
		scaledElapsed := time.Duration(float64(virtualElapsed) / p.speed)
		p.startTime = time.Now().Add(-scaledElapsed)
	}
}

// Pause pauses playback at current position.
func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.playing {
		return
	}

	p.playing = false
	p.pauseTime = time.Now()
}

// Stop resets to the beginning.
func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.playing = false
	p.position = 0
	p.startTime = time.Time{}
	p.pauseTime = time.Time{}
}

// SetSpeed sets playback speed (1.0/4.0/8.0/16.0).
func (p *Player) SetSpeed(speed float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if speed <= 0 {
		speed = 1.0
	}

	// If playing, adjust startTime to maintain virtual position
	if p.playing && p.position > 0 && len(p.events) > 0 {
		virtualElapsed := p.events[p.position].Ts.Sub(p.baseTime)
		oldScaledElapsed := time.Duration(float64(virtualElapsed) / p.speed)
		newScaledElapsed := time.Duration(float64(virtualElapsed) / speed)

		// Adjust startTime: old elapsed -> new elapsed
		p.startTime = p.startTime.Add(oldScaledElapsed - newScaledElapsed)
	}

	p.speed = speed
}

// StepForward advances one event.
func (p *Player) StepForward() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.position < len(p.events)-1 {
		p.position++
	}
}

// StepBackward rewinds one event.
func (p *Player) StepBackward() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.position > 0 {
		p.position--
	}
}

// Seek moves to a specific position (0-indexed).
func (p *Player) Seek(position int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if position < 0 {
		position = 0
	}
	if position >= len(p.events) {
		position = len(p.events) - 1
	}

	if len(p.events) == 0 {
		position = 0
	}

	p.position = position

	// If playing, reset startTime to new position
	if p.playing && len(p.events) > 0 {
		virtualElapsed := p.events[p.position].Ts.Sub(p.baseTime)
		scaledElapsed := time.Duration(float64(virtualElapsed) / p.speed)
		p.startTime = time.Now().Add(-scaledElapsed)
	}
}

// CurrentEvent returns the event at the current position.
// Returns nil if position is out of bounds.
func (p *Player) CurrentEvent() *schema.CanonicalEvent {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.position < 0 || p.position >= len(p.events) {
		return nil
	}
	return &p.events[p.position]
}

// Position returns the current position (0-indexed).
func (p *Player) Position() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.position
}

// Total returns the total number of events.
func (p *Player) Total() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.events)
}

// IsPlaying returns true if playback is active.
func (p *Player) IsPlaying() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.playing
}

// Speed returns the current playback speed.
func (p *Player) Speed() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.speed
}

// EventsUntil returns all events up to the given virtual time.
// Used in playback mode to determine which events should be shown.
// The virtual time is calculated as: baseTime + (realElapsed * speed)
func (p *Player) EventsUntil(now time.Time) []schema.CanonicalEvent {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.playing || len(p.events) == 0 {
		return nil
	}

	// Calculate virtual time
	realElapsed := now.Sub(p.startTime)
	virtualElapsed := time.Duration(float64(realElapsed) * p.speed)
	virtualNow := p.baseTime.Add(virtualElapsed)

	// Find all events before virtualNow
	var result []schema.CanonicalEvent
	for i := 0; i < len(p.events); i++ {
		if p.events[i].Ts.After(virtualNow) {
			break
		}
		result = append(result, p.events[i])
	}

	return result
}

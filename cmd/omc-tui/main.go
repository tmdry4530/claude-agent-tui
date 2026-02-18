package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chamdom/omc-agent-tui/internal/bridge"
	"github.com/chamdom/omc-agent-tui/internal/collector"
	"github.com/chamdom/omc-agent-tui/internal/normalizer"
	"github.com/chamdom/omc-agent-tui/internal/replay"
	"github.com/chamdom/omc-agent-tui/internal/store"
	"github.com/chamdom/omc-agent-tui/internal/tui"
	"github.com/chamdom/omc-agent-tui/pkg/schema"
	tea "github.com/charmbracelet/bubbletea"
)

// Set by goreleaser ldflags.
var (
	version = "dev"
	commit  = "none"
)

func main() {
	watchPath := flag.String("watch", "", "Directory to watch for JSONL event files")
	replayFile := flag.String("replay", "", "JSONL file to replay")
	convertFile := flag.String("convert", "", "Convert subagent-tracking.json to JSONL (output to stdout or -o)")
	convertOut := flag.String("o", "", "Output path for --convert (default: stdout)")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("omc-tui %s (%s)\n", version, commit)
		return
	}

	// Convert mode: non-TUI, converts tracking file and exits
	if *convertFile != "" {
		if err := runConvert(*convertFile, *convertOut); err != nil {
			fmt.Fprintf(os.Stderr, "Convert error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	s := store.NewStore(10000)
	m := tui.NewModel(s)

	// Add demo events before creating program (so they're in initial state)
	if *watchPath == "" && *replayFile == "" {
		addDemoEvents(&m)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	// Start pipeline after program is created (pipeline sends events via p.Send)
	var cleanup func()
	switch {
	case *watchPath != "":
		cleanup = startLivePipeline(p, *watchPath)
	case *replayFile != "":
		if err := startReplay(p, *replayFile); err != nil {
			fmt.Fprintf(os.Stderr, "Replay error: %v\n", err)
			os.Exit(1)
		}
	}

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if cleanup != nil {
		cleanup()
	}
}

// startLivePipeline starts the Collector -> Normalizer -> TUI pipeline.
// Returns a cleanup function to stop the collector on exit.
func startLivePipeline(p *tea.Program, watchPath string) func() {
	norm := normalizer.New()
	coll := collector.NewFileCollector(watchPath)

	ctx, cancel := context.WithCancel(context.Background())
	if err := coll.Start(ctx); err != nil {
		log.Printf("Warning: failed to start collector: %v", err)
		cancel()
		return nil
	}

	go func() {
		for rawEvent := range coll.Events() {
			canonical, err := norm.Normalize(rawEvent)
			if err != nil {
				continue
			}
			p.Send(tui.EventMsg(*canonical))
		}
	}()

	return func() {
		cancel()
		coll.Stop()
	}
}

// startReplay loads a JSONL file and sends events to the TUI with original timing.
func startReplay(p *tea.Program, filePath string) error {
	player := replay.NewPlayer()
	if err := player.LoadFile(filePath); err != nil {
		return fmt.Errorf("load replay: %w", err)
	}

	go func() {
		total := player.Total()
		var lastTs time.Time
		for i := 0; i < total; i++ {
			player.Seek(i)
			evt := player.CurrentEvent()
			if evt == nil {
				continue
			}
			// Respect original timing between events (capped at 2s)
			if i > 0 && !lastTs.IsZero() {
				delay := evt.Ts.Sub(lastTs)
				if delay > 0 {
					if delay > 2*time.Second {
						delay = 2 * time.Second
					}
					time.Sleep(delay)
				}
			}
			lastTs = evt.Ts
			p.Send(tui.EventMsg(*evt))
		}
	}()

	return nil
}

// runConvert converts a subagent-tracking.json file to JSONL.
func runConvert(trackingPath, outputPath string) error {
	events, err := bridge.ConvertTracking(trackingPath)
	if err != nil {
		return err
	}

	if outputPath != "" {
		return bridge.WriteEventsFile(outputPath, events)
	}

	// Write to stdout
	enc := json.NewEncoder(os.Stdout)
	for _, evt := range events {
		if err := enc.Encode(evt); err != nil {
			return fmt.Errorf("encode event: %w", err)
		}
	}
	return nil
}

// addDemoEvents adds sample events for testing the UI without a data source.
func addDemoEvents(m *tui.Model) {
	now := time.Now()

	events := []schema.CanonicalEvent{
		{
			Ts:       now.Add(-6 * time.Minute),
			RunID:    "run-001",
			Provider: schema.ProviderClaude,
			Mode:     schema.ModeAutopilot,
			AgentID:  "planner-1",
			Role:     schema.RolePlanner,
			State:    schema.StateRunning,
			Type:     schema.TypeTaskSpawn,
			TaskID:   "task-001",
		},
		{
			Ts:            now.Add(-5 * time.Minute),
			RunID:         "run-001",
			Provider:      schema.ProviderClaude,
			Mode:          schema.ModeAutopilot,
			AgentID:       "executor-1",
			ParentAgentID: "planner-1",
			Role:          schema.RoleExecutor,
			State:         schema.StateRunning,
			Type:          schema.TypeTaskSpawn,
			TaskID:        "task-002",
		},
		{
			Ts:       now.Add(-4 * time.Minute),
			RunID:    "run-001",
			Provider: schema.ProviderClaude,
			Mode:     schema.ModeAutopilot,
			AgentID:  "reviewer-1",
			Role:     schema.RoleReviewer,
			State:    schema.StateWaiting,
			Type:     schema.TypeMessage,
		},
		{
			Ts:       now.Add(-3 * time.Minute),
			RunID:    "run-001",
			Provider: schema.ProviderClaude,
			Mode:     schema.ModeAutopilot,
			AgentID:  "guard-1",
			Role:     schema.RoleGuard,
			State:    schema.StateBlocked,
			Type:     schema.TypeStateChange,
		},
		{
			Ts:       now.Add(-2 * time.Minute),
			RunID:    "run-001",
			Provider: schema.ProviderClaude,
			Mode:     schema.ModeAutopilot,
			AgentID:  "executor-1",
			Role:     schema.RoleExecutor,
			State:    schema.StateError,
			Type:     schema.TypeError,
			TaskID:   "task-002",
		},
		{
			Ts:       now.Add(-90 * time.Second),
			RunID:    "run-001",
			Provider: schema.ProviderClaude,
			Mode:     schema.ModeAutopilot,
			AgentID:  "tester-1",
			Role:     schema.RoleTester,
			State:    schema.StateRunning,
			Type:     schema.TypeVerify,
		},
		{
			Ts:       now.Add(-1 * time.Minute),
			RunID:    "run-001",
			Provider: schema.ProviderClaude,
			Mode:     schema.ModeAutopilot,
			AgentID:  "writer-1",
			Role:     schema.RoleWriter,
			State:    schema.StateIdle,
			Type:     schema.TypeMessage,
		},
		{
			Ts:       now.Add(-30 * time.Second),
			RunID:    "run-001",
			Provider: schema.ProviderClaude,
			Mode:     schema.ModeAutopilot,
			AgentID:  "verifier-1",
			Role:     schema.RoleVerifier,
			State:    schema.StateDone,
			Type:     schema.TypeTaskDone,
			TaskID:   "task-001",
		},
	}

	for _, event := range events {
		m.AddEvent(event)
	}
}

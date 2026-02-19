# claude-agent-tui

Real-time TUI dashboard for monitoring AI agent orchestration events.

Built with Go and [Bubbletea](https://github.com/charmbracelet/bubbletea).

## Demo

```
╭────────────────────────╮  ╭────────────────────────╮  ╭────────────────────────╮
│         ▐▛███▜▌       │  │         ▐▛███▜▌      │  │         ▐▛███▜▌       │
│       ▝▜█████▛▘      │  │       ▝▜█████▛▘    │  │       ▝▜█████▛▘      │
│         ▘▘ ▝▝        │  │         ▘▘ ▝▝      │  │         ▘▘ ▝▝        │
│  executor ●            │  │  planner ●             │  │  reviewer ○            │
│  agent-exec-1          │  │  agent-plan-1          │  │  agent-rev-1           │
│  running               │  │  running               │  │  waiting               │
│  Recent activity       │  │  Recent activity       │  │  Recent activity       │
│  task spawned          │  │  planning phase        │  │  awaiting review       │
╰────────────────────────╯  ╰────────────────────────╯  ╰────────────────────────╯
```

## Features

- 5-panel layout: Arena (agent cards), Timeline (event stream), Graph (task tree), Inspector (event details), Footer (metrics)
- Live monitoring via fsnotify file watching
- JSONL replay with original event timing
- 12 agent role colors, 6 state indicators
- Circuit breaker with exponential backoff
- PII redaction (11 key + 8 regex patterns)
- Ring buffer store (10K events)

## Quick Start (1 minute)

```bash
# 1. Build
go build -o bin/omc-tui ./cmd/omc-tui

# 2. Start monitor (in a separate terminal)
./bin/omc-tui --watch .omc/events/

# 3. Or replay a previous session
./bin/omc-tui --convert .omc/state/subagent-tracking.json -o .omc/events/session.jsonl
./bin/omc-tui --replay .omc/events/session.jsonl
```

### Claude Code Plugin Commands

| Command | Description |
|---------|-------------|
| `/claude-agent-tui:install-bridge` | Build binary + set up hook bridge |
| `/claude-agent-tui:monitor` | Start real-time agent monitoring |
| `/claude-agent-tui:replay` | Replay latest agent session |
| `/claude-agent-tui:doctor` | Run installation diagnostics |
| `/claude-agent-tui:stop` | Stop running monitor |

## Requirements

- Go 1.24 or higher
- Terminal with Unicode and 256-color support
- `jq` (optional, for hook script)

## Installation

### As Claude Code plugin

```bash
claude --plugin-dir /path/to/claude-agent-tui
```

### Build from source

```bash
git clone https://github.com/tmdry4530/omc-agent-tui.git
cd omc-agent-tui
go build -o bin/omc-tui ./cmd/omc-tui
```

## Usage

### Demo mode (default)

```bash
./bin/omc-tui
```

Shows sample agent events for testing the UI.

### Live watch mode

```bash
./bin/omc-tui --watch /path/to/events/dir
```

Monitors a directory for JSONL event files in real-time using fsnotify.

### Replay mode

```bash
./bin/omc-tui --replay /path/to/events.jsonl
```

Replays events from a JSONL file with original timing (capped at 2s between events).

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` | Switch between panels |
| `j` / `Down` | Scroll down |
| `k` / `Up` | Scroll up |
| `q` / `Ctrl+C` | Quit |

## Event Format

Events are JSONL files with the following canonical structure:

```json
{
  "ts": "2026-02-18T12:00:00Z",
  "run_id": "run-001",
  "provider": "claude",
  "agent_id": "agent-executor-1",
  "role": "executor",
  "state": "running",
  "type": "task_spawn",
  "task_id": "task-001"
}
```

See `pkg/schema/` for full type definitions.

## Project Structure

```
cmd/omc-tui/          CLI entrypoint (watch, replay, convert)
internal/
  bridge/             OMC bridge (tracking converter + event emitter)
  collector/          File-based event collector (fsnotify)
  normalizer/         Event normalization + PII redaction
  store/              Ring buffer event store
  replay/             JSONL replay engine
  tui/                Bubbletea TUI model
    arena/            Agent card panel (CLCO mascot)
    timeline/         Event stream panel
    graph/            Task dependency tree
    inspector/        Event detail viewer
    footer/           Metrics footer bar
pkg/schema/           Canonical event types and enums
scripts/              Hook bridge script
commands/             Plugin slash commands
hooks/                Plugin hook configuration
```

## Testing

```bash
go test ./... -race
```

182 tests across 12 packages. All pass with race detector enabled.

## License

MIT

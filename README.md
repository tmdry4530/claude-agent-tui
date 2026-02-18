# omc-agent-tui

Real-time TUI dashboard for monitoring AI agent orchestration events.

Built with Go and [Bubbletea](https://github.com/charmbracelet/bubbletea).

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

### Claude Code Slash Commands

| Command | Description |
|---------|-------------|
| `/project:install-bridge` | Build binary + set up hook bridge |
| `/project:monitor` | Start real-time agent monitoring |
| `/project:replay` | Replay latest agent session |
| `/project:doctor` | Run installation diagnostics |
| `/project:stop` | Stop running monitor |

## Requirements

- Go 1.23 or higher
- Terminal with Unicode and 256-color support
- `jq` (optional, for hook script)

## Installation

### Build from source

```bash
git clone https://github.com/chamdom/omc-agent-tui.git
cd omc-agent-tui
make build
```

### Manual build

```bash
go build -o bin/omc-tui ./cmd/omc-tui
```

### System-wide install

```bash
make install
# Installs to ~/.local/bin/omc-tui
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

## Claude Code Plugin

### Install as plugin

```bash
claude plugin install /path/to/omc-agent-tui
```

The plugin metadata is defined in `.claude-plugin/plugin.json`.

### Plugin args

| Arg | Description |
|-----|-------------|
| `--watch <dir>` | Directory to watch for JSONL event files |
| `--replay <file>` | JSONL file to replay |

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
scripts/              OMC hook bridge script
.claude/commands/     Slash commands (install-bridge, monitor, replay, doctor, stop)
```

## Testing

```bash
make test
# or
go test ./... -race
```

182 tests across 12 packages. All pass with race detector enabled.

## License

MIT

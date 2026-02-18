# Installation Guide

## Prerequisites

- Go 1.23 or higher
- Terminal with Unicode and color support

## Build from Source

```bash
cd /path/to/omc-agent-tui
go build -o bin/omc-tui ./cmd/omc-tui
```

Or use the Makefile:

```bash
make build
```

## Run Modes

### Demo Mode

Run without arguments to see sample data:

```bash
./bin/omc-tui
```

### Live Watch Mode

Monitor a directory for new JSONL event files:

```bash
./bin/omc-tui --watch /path/to/events/dir
```

### Replay Mode

Replay events from a JSONL file:

```bash
./bin/omc-tui --replay /path/to/events.jsonl
```

## Keyboard Shortcuts

- `Tab` - Switch between panels
- `j` / `Down` - Scroll down
- `k` / `Up` - Scroll up
- `q` / `Ctrl+C` - Quit

## Plugin Installation

Install as a Claude Code plugin:

```bash
claude plugin install /path/to/omc-agent-tui
```

After installation, the plugin will be available in Claude Code's plugin registry.

## Update

Pull the latest changes and rebuild:

```bash
git pull
go build -o bin/omc-tui ./cmd/omc-tui
```

Or:

```bash
make build
```

## System-wide Installation

Copy the binary to your PATH:

```bash
make install
```

This copies `bin/omc-tui` to `~/.local/bin/omc-tui` by default.

## Troubleshooting

### Build Fails

Ensure Go 1.23+ is installed:

```bash
go version
```

### Binary Not Found

Add `~/.local/bin` to your PATH:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Add this line to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) to persist.

### Terminal Display Issues

Ensure your terminal supports Unicode and 256 colors. Modern terminals like iTerm2, Alacritty, or Windows Terminal work best.

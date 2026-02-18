# Release Notes: v0.1.0-beta

Released: 2026-02-18

## Overview

First beta release of omc-agent-tui, a real-time TUI dashboard for monitoring AI agent orchestration events. This release includes the complete 5-panel TUI, live/replay/demo modes, comprehensive test suite, and Claude Code plugin packaging.

## Highlights

- **5-panel TUI**: Arena (agent cards), Timeline (event stream), Graph (task tree), Inspector (event details), Footer (metrics)
- **3 run modes**: `--watch` (live fsnotify monitoring), `--replay` (JSONL replay), demo (default)
- **122 tests** across 11 packages, all passing with race detector
- **Plugin packaging**: `.claude-plugin/plugin.json` for Claude Code integration

## What's New Since Alpha

- Flaky `TestFileCollector_NewLines` completely stabilized (root cause: `bufio.Scanner` buffering vs file position tracking)
- Test coverage expanded from 68 to 122 tests (+79%)
  - Arena: 15 tests (role colors, state styles, rendering)
  - Timeline: 13 tests (event ordering, viewport, rendering)
  - Footer: 20 tests (metrics, formatting, all view states)
  - Collector: 6 tests (+2: CircuitBreaker, BackoffLevels)
- Claude Code plugin packaging (.claude-plugin/plugin.json, docs/INSTALL.md, Makefile)
- CHANGELOG.md, README.md, RELEASE_NOTES.md documentation

## Architecture

```
FileCollector (fsnotify) -> Normalizer (redaction) -> Store (ring buffer 10K) -> TUI (Bubbletea)
```

- **Bubbletea MVU pattern**: single-goroutine Update() ensures thread safety
- **EventMsg pipeline**: goroutines send events via `p.Send()`, processed sequentially in Update()
- **Circuit breaker**: 3-failure threshold, exponential backoff (10s/30s/60s)

## Verification

| Check | Result |
|-------|--------|
| Tests | 122/122 PASS |
| Build | OK |
| Race detector | 0 races |
| Flaky test stability | 10/10 consecutive PASS |

## Known Risks

### P1: No E2E/Integration Tests
- Unit tests cover individual components but not the full pipeline end-to-end
- Mitigation: manual testing with demo mode and JSONL files
- Plan: add E2E tests in next release

### P2: Makefile GO Path Hardcoded
- `Makefile` references `$(HOME)/.local/go/bin/go` directly
- Portability issue for systems with Go installed elsewhere
- Mitigation: use `go build` directly or set `GO` variable
- Plan: change to `GO?=$(shell which go)` in next release

## Upgrade Path

From v0.1.0-alpha:
```bash
git pull
make clean && make build
```

No breaking changes. All existing JSONL event files remain compatible.

## Next Release (v0.2.0)

Planned features:
1. Arena focus highlight borders
2. Timeline -> Inspector Enter key integration
3. Replay TUI controls (Space pause/play, +/- speed)
4. Event filtering and search (/ key)

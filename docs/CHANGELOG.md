# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [0.1.0-beta] - 2026-02-18

### Added
- Real-time TUI dashboard with 5 panels: Arena, Timeline, Graph, Inspector, Footer
- Live watch mode (`--watch`) for monitoring JSONL event directories via fsnotify
- Replay mode (`--replay`) for replaying JSONL event files with original timing
- Demo mode (default, no args) with sample agent events
- Event pipeline: FileCollector -> Normalizer -> Store -> TUI
- Ring buffer store (10K events) with oldest-drop eviction
- Circuit breaker with exponential backoff (10s/30s/60s) for collector resilience
- PII redaction in normalizer (11 key patterns, 8 regex patterns, recursive depth 10)
- Agent card display with 12 role colors and 6 state styles
- Task dependency graph with parent-child relationship tracking
- Event inspector with detailed payload and metrics view
- Footer metrics: token counts (K/M format), cost (USD), latency (ms)
- Tab-based panel focus switching (4 panels)
- Viewport scrolling (j/k/Up/Down) for Timeline and Inspector
- Claude Code plugin packaging (.claude-plugin/plugin.json)
- Installation documentation (docs/INSTALL.md)
- Makefile with build/test/clean/install targets
- Comprehensive test suite: 122 tests across 11 test files, 11 packages

### Changed
- Collector position tracking: replaced `file.Seek()` with `bytesProcessed` counter for accurate line-level position tracking

### Fixed
- Flaky `TestFileCollector_NewLines`: root cause was `bufio.Scanner` internal buffering causing incorrect file position via `file.Seek(0, io.SeekCurrent)`. Fixed by tracking processed bytes per line. Stability verified with 10/10 consecutive passes.
- Collector test robustness: increased fsnotify setup wait (200ms -> 500ms), added `context.WithTimeout` (10s), extracted `drainEvents` helper

### Known Issues
- Makefile GO path is hardcoded to `$(HOME)/.local/go/bin/go` (P2: portability)
- No E2E/integration tests (P1: planned for next release)
- No event filtering or search functionality (P3: planned)
- Replay mode does not support TUI controls (pause/speed) yet (P1: planned)

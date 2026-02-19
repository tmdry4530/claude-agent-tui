Run diagnostics on the claude-agent-tui installation.

Steps to perform (run ALL checks, report results):
1. **Go version**: Run `go version`. Requires 1.24+.
2. **Binary**: Check if `bin/omc-tui` exists and is executable.
3. **Events directory**: Check if `.omc/events/` exists.
4. **Hook script**: Check if `scripts/omc-bridge-hook.sh` exists and is executable.
5. **Dependencies**: Check if `jq` is available (required by hook script).
6. **Tracking data**: Check if `.omc/state/subagent-tracking.json` exists.
7. **Event files**: Count `.omc/events/*.jsonl` files and total line count.

Print results as a checklist:
```
[omc-tui] Diagnostics

  Go version:      ✓ go1.24.7 (>= 1.24)
  Binary:          ✓ bin/omc-tui (4.9M)
  Events dir:      ✓ .omc/events/ (3 files)
  Hook script:     ✓ scripts/omc-bridge-hook.sh
  jq available:    ✓ /usr/bin/jq
  Tracking data:   ✓ .omc/state/subagent-tracking.json (21 agents)
  Running process: ✗ No omc-tui process found

  Status: Ready (6/7 checks passed)
```

For each failed check, add a fix suggestion:
- Missing binary: "Run /claude-agent-tui:install-bridge"
- Missing events dir: "Run mkdir -p .omc/events"
- Missing hook: "Run chmod +x scripts/omc-bridge-hook.sh"
- Missing jq: "Install jq: apt install jq / brew install jq"

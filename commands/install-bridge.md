Install the claude-agent-tui bridge for real-time Claude Code agent monitoring.

Steps to perform:
1. Build the binary: `go build -o bin/omc-tui ./cmd/omc-tui/`
2. Create the events directory: `mkdir -p .omc/events`
3. Make the hook script executable: `chmod +x scripts/omc-bridge-hook.sh`
4. Report the absolute path to the hook script for the user to add to their hooks config

On success, print:
```
[omc-tui] Bridge installed successfully.
  Binary:  bin/omc-tui
  Events:  .omc/events/
  Hook:    scripts/omc-bridge-hook.sh

Next: Add the hook to your Claude Code settings:
  "hooks": {
    "PostToolUse": [{ "command": "<absolute-path>/scripts/omc-bridge-hook.sh" }]
  }

Then run /claude-agent-tui:monitor to start real-time monitoring.
```

On build failure, print the error and suggest checking Go version with `go version` (requires Go 1.24+).

Start real-time omc-agent-tui monitoring of Claude Code agent activity.

Steps to perform:
1. Check if `bin/omc-tui` exists. If not, build it: `go build -o bin/omc-tui ./cmd/omc-tui/`
2. Create events directory if needed: `mkdir -p .omc/events`
3. Tell the user to run in a separate terminal:
   ```
   ./bin/omc-tui --watch .omc/events/
   ```
4. Optionally, if the user has tmux available, offer to launch it in a tmux pane.

On success, print:
```
[omc-tui] Monitor ready.
  Run in a separate terminal:
    ./bin/omc-tui --watch .omc/events/

  Or with tmux:
    tmux split-window -h './bin/omc-tui --watch .omc/events/'

  Agent events will appear in real-time as you work.
  Use /project:stop to terminate.
```

If binary build fails, suggest running /project:doctor.

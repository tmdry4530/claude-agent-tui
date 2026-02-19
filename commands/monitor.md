Start real-time claude-agent-tui monitoring of Claude Code agent activity.

Steps to perform:
1. Check if `bin/omc-tui` exists. If not, build it: `go build -o bin/omc-tui ./cmd/omc-tui/`
2. Create events directory if needed: `mkdir -p .omc/events`
3. Detect tmux environment:
   - Run: `echo $TMUX`
   - If `$TMUX` is non-empty (inside tmux session):
     - Launch TUI in a new right-side pane:
       ```
       tmux split-window -h -l 40% './bin/omc-tui --watch .omc/events/'
       ```
     - Capture the new pane ID:
       ```
       tmux display-message -p -t '{last}' '#{pane_id}'
       ```
     - Save pane ID to `.omc/monitor.pid`
     - Print:
       ```
       [omc-tui] Monitor launched in tmux pane <pane-id>.
         Agent events will appear in real-time as you work.
         Use /claude-agent-tui:stop to terminate.
       ```
   - If `$TMUX` is empty (not in tmux):
     - Print:
       ```
       [omc-tui] Monitor ready.
         Not inside a tmux session — run manually in a separate terminal:
           ./bin/omc-tui --watch .omc/events/

         Or start a tmux session first, then re-run /claude-agent-tui:monitor.
         Use /claude-agent-tui:stop to terminate.
       ```

If binary build fails, suggest running /claude-agent-tui:doctor.

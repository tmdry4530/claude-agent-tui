Stop the running claude-agent-tui monitor process.

Steps to perform:
1. Check if `.omc/monitor.pid` exists and read its content.
2. If `.omc/monitor.pid` exists and contains a tmux pane ID (e.g. `%123`):
   - Verify the pane is still alive: `tmux has-session 2>/dev/null && tmux list-panes -a -F '#{pane_id}' | grep -q '<pane-id>'`
   - If alive, kill the pane: `tmux kill-pane -t <pane-id>`
   - Remove the pid file: `rm -f .omc/monitor.pid`
   - Print:
     ```
     [omc-tui] Monitor pane <pane-id> terminated.
     ```
   - If pane not found, remove stale pid file and fall through to step 3.
3. Fallback — find process directly: `pgrep -f 'omc-tui.*--watch'`
   - If found, show the PID and kill it gracefully: `kill <pid>`
   - Verify the process stopped: `sleep 1 && pgrep -f 'omc-tui.*--watch'`
   - On success, print:
     ```
     [omc-tui] Monitor stopped (PID <pid>).
     ```
   - If kill fails, suggest force kill:
     ```
     [omc-tui] Process <pid> did not stop. Force kill with:
       kill -9 <pid>
     ```
4. If no process or pane found at all, print:
   ```
   [omc-tui] No running monitor found.
   ```

Stop the running omc-agent-tui monitor process.

Steps to perform:
1. Find running omc-tui processes: `pgrep -f 'omc-tui.*--watch'`
2. If found, show the PID and kill it gracefully: `kill <pid>`
3. Verify the process stopped: `sleep 1 && pgrep -f 'omc-tui.*--watch'`

On success, print:
```
[omc-tui] Monitor stopped (PID <pid>).
```

If no process found, print:
```
[omc-tui] No running monitor found.
```

If kill fails, suggest force kill:
```
[omc-tui] Process <pid> did not stop. Force kill with:
  kill -9 <pid>
```

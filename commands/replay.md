Replay agent activity from a previous session.

Steps to perform:
1. Check if `bin/omc-tui` exists. If not, build it: `go build -o bin/omc-tui ./cmd/omc-tui/`
2. If user provided a specific file argument, use that.
3. Otherwise, find the latest replay source:
   a. Check `.omc/events/*.jsonl` for the most recently modified file
   b. If no events exist, check `.omc/state/subagent-tracking.json` and convert it:
      `./bin/omc-tui --convert .omc/state/subagent-tracking.json -o .omc/events/converted.jsonl`
4. Tell the user to run:
   ```
   ./bin/omc-tui --replay <file>
   ```

On success, print:
```
[omc-tui] Replay ready.
  Source: <file> (<N> events)
  Run:   ./bin/omc-tui --replay <file>

  Controls: q to quit
```

If no replay source found, print:
```
[omc-tui] No replay data found.
  Run some agent tasks first, or provide a JSONL file:
    /claude-agent-tui:replay path/to/events.jsonl
```

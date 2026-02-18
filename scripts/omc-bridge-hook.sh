#!/usr/bin/env bash
# omc-bridge-hook.sh — OMC PostToolUse hook that emits CanonicalEvent JSONL
# for omc-agent-tui real-time sync.
#
# Usage in .claude/settings.json hooks:
#   "PostToolUse": [{ "command": "/path/to/omc-bridge-hook.sh" }]
#
# Reads hook JSON from stdin with fields:
#   tool_name, tool_input, session_id, cwd
#
# Emits JSONL to: <cwd>/.omc/events/<session_id>.jsonl

set -euo pipefail

# Read hook data from stdin; validate JSON
INPUT=$(cat)
if ! echo "$INPUT" | jq empty 2>/dev/null; then
    exit 0
fi

# Extract fields
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // empty')
SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // empty')
CWD=$(echo "$INPUT" | jq -r '.cwd // empty')

# Require all fields
if [ -z "$TOOL_NAME" ] || [ -z "$SESSION_ID" ] || [ -z "$CWD" ]; then
    exit 0
fi

# Sanitize SESSION_ID: strip path components, allow only safe chars
SESSION_ID=$(basename "$SESSION_ID")
if [[ ! "$SESSION_ID" =~ ^[a-zA-Z0-9._-]+$ ]]; then
    echo "omc-bridge-hook: invalid session_id" >&2
    exit 1
fi

EVENT_DIR="${CWD}/.omc/events"
(umask 077; mkdir -p "$EVENT_DIR")
EVENT_FILE="${EVENT_DIR}/${SESSION_ID}.jsonl"

TS=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

case "$TOOL_NAME" in
    Task)
        # Extract agent info from tool_input
        AGENT_TYPE=$(echo "$INPUT" | jq -r '.tool_input.subagent_type // "custom"')
        AGENT_ID=$(echo "$INPUT" | jq -r '.tool_input.name // empty')
        if [ -z "$AGENT_ID" ]; then
            AGENT_ID="agent-$(echo "$INPUT" | md5sum | cut -c1-7)"
        fi

        # Strip oh-my-claudecode: prefix for role lookup
        ROLE_KEY="${AGENT_TYPE#oh-my-claudecode:}"
        ROLE=$(jq -rn --arg k "$ROLE_KEY" '{
            planner:"planner",executor:"executor","deep-executor":"executor",
            explore:"explorer",architect:"architect",debugger:"debugger",
            verifier:"verifier",designer:"designer","code-reviewer":"reviewer",
            "security-reviewer":"guard","test-engineer":"tester",writer:"writer",
            analyst:"planner","build-fixer":"executor","git-master":"executor"
        }[$k] // "custom"')

        # Emit spawn event
        (umask 077; jq -nc \
            --arg ts "$TS" \
            --arg run_id "omc-${AGENT_ID}" \
            --arg agent_id "$AGENT_ID" \
            --arg role "$ROLE" \
            '{
                ts: $ts, run_id: $run_id, provider: "claude",
                agent_id: $agent_id, role: $role,
                state: "running", type: "task_spawn"
            }' >> "$EVENT_FILE")
        ;;

    TaskUpdate)
        TASK_STATUS=$(echo "$INPUT" | jq -r '.tool_input.status // empty')
        AGENT_ID=$(echo "$INPUT" | jq -r '.tool_input.owner // "unknown"')

        case "$TASK_STATUS" in
            completed)
                (umask 077; jq -nc \
                    --arg ts "$TS" \
                    --arg run_id "omc-${AGENT_ID}" \
                    --arg agent_id "$AGENT_ID" \
                    '{
                        ts: $ts, run_id: $run_id, provider: "claude",
                        agent_id: $agent_id, role: "custom",
                        state: "done", type: "task_done"
                    }' >> "$EVENT_FILE")
                ;;
            in_progress)
                (umask 077; jq -nc \
                    --arg ts "$TS" \
                    --arg run_id "omc-${AGENT_ID}" \
                    --arg agent_id "$AGENT_ID" \
                    '{
                        ts: $ts, run_id: $run_id, provider: "claude",
                        agent_id: $agent_id, role: "custom",
                        state: "running", type: "task_update"
                    }' >> "$EVENT_FILE")
                ;;
        esac
        ;;

    SendMessage)
        MSG_TYPE=$(echo "$INPUT" | jq -r '.tool_input.type // empty')
        if [ "$MSG_TYPE" = "shutdown_request" ]; then
            AGENT_ID=$(echo "$INPUT" | jq -r '.tool_input.recipient // "unknown"')
            (umask 077; jq -nc \
                --arg ts "$TS" \
                --arg run_id "omc-${AGENT_ID}" \
                --arg agent_id "$AGENT_ID" \
                '{
                    ts: $ts, run_id: $run_id, provider: "claude",
                    agent_id: $agent_id, role: "custom",
                    state: "cancelled", type: "state_change"
                }' >> "$EVENT_FILE")
        fi
        ;;
esac

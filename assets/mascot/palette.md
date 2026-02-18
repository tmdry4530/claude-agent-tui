# CLCO Palette

## 1) Base
- Background: `#0B1020`
- Foreground: `#E6EAF2`
- Muted Text: `#8A93A5`
- Border: `#2A3142`

---

## 2) Role Colors

- Planner: `#7AA2F7`
- Executor: `#9ECE6A`
- Reviewer: `#BB9AF7`
- Guard: `#F7768E`
- Tester: `#E0AF68`
- Writer: `#73DACA`
- Explorer: `#58A6FF`
- Architect: `#FFA657`
- Debugger: `#FF7B72`
- Verifier: `#56D364`
- Designer: `#D2A8FF`
- Custom: `#8A93A5`

### Role Variants (multi-agent)
- Executor-A: `#9ECE6A`
- Executor-B: `#8BCF7A`
- Executor-C: `#7FD18A`

---

## 3) State Colors

- RUNNING: `#7EE787`
- WAITING: `#7D8590`
- BLOCKED: `#E3B341`
- ERROR: `#FF7B72`
- DONE: `#56D364`
- FAILED: `#FF7B72`
- CANCELLED: `#6E7681`
- IDLE: `#6E7681`

---

## 4) Semantic Colors

- Info: `#58A6FF`
- Warning: `#E3B341`
- Critical: `#FF7B72`
- Success: `#3FB950`
- Accent: `#A5D6FF`

---

## 5) Terminal Fallback (256-color)

- Planner: `111`
- Executor: `114`
- Reviewer: `183`
- Guard: `203`
- Tester: `179`
- Writer: `80`
- Explorer: `75`
- Architect: `214`
- Debugger: `203`
- Verifier: `77`
- Designer: `183`
- Custom: `245`

- RUNNING: `120`
- WAITING: `245`
- BLOCKED: `214`
- ERROR: `203`
- DONE: `77`
- FAILED: `203`
- CANCELLED: `242`
- IDLE: `242`

---

## 6) Accessibility Palette (Color-blind Friendly)

- Planner: `#4C78A8`
- Executor: `#59A14F`
- Reviewer: `#B07AA1`
- Guard: `#E15759`
- Tester: `#F28E2B`
- Writer: `#76B7B2`
- Explorer: `#4E79A7`
- Architect: `#F28E2B`
- Debugger: `#E15759`
- Verifier: `#59A14F`
- Designer: `#B07AA1`
- Custom: `#7F7F7F`

상태:
- RUNNING: `#2E8B57`
- WAITING: `#7F7F7F`
- BLOCKED: `#C17D11`
- ERROR: `#B22222`
- DONE: `#228B22`
- FAILED: `#B22222`
- CANCELLED: `#696969`
- IDLE: `#696969`

---

## 7) Usage Rules

1. 상태 표현은 색상 + 텍스트/아이콘 동시 사용
2. ERROR/BLOCKED는 테두리 강조 필수
3. 저채도 모드(야간)에서는 채도 15% 낮춤
4. 색상 변경은 `DECISIONS.md` 기록 후 적용

---

## 8) Example Tokens (for docs/spec)

```text
role.planner = #7AA2F7
role.coder = #9ECE6A
state.running = #7EE787
state.error = #FF7B72
ui.background = #0B1020
ui.foreground = #E6EAF2

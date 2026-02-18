# Event Schema (OMC Agent TUI)

## 1) 목적

OMC / Claude Code / Gemini / Codex에서 발생하는 이벤트를
하나의 공통 포맷(**CanonicalEvent**)으로 정규화하여
TUI가 일관된 방식으로 소비하도록 한다.

---

## 2) CanonicalEvent 정의

```json
{
  "ts": "2026-02-17T22:27:00Z",
  "run_id": "run-20260217-2227",
  "provider": "claude|gemini|codex|system",
  "mode": "ralph|ultrawork|ultrapilot|team|autopilot|pipeline|ecomode|unknown",
  "agent_id": "coder-auth",
  "parent_agent_id": "planner-main",
  "role": "planner|executor|reviewer|guard|tester|writer|explorer|architect|debugger|verifier|designer|custom",
  "state": "idle|running|waiting|blocked|error|done|failed|cancelled",
  "type": "task_spawn|task_update|task_done|tool_call|tool_result|message|error|replan|verify|fix|recover|state_change",
  "task_id": "task-42",
  "intent_ref": "plan-7",
  "payload": {},
  "metrics": {
    "latency_ms": 420,
    "tokens_in": 210,
    "tokens_out": 95,
    "cost_usd": 0.0021
  },
  "raw_ref": "optional://source-pointer"
}
```

---

## 3) 필드 설명

* **ts (required)**
  이벤트 발생 시간 (ISO 8601, UTC 권장)

* **run_id (required)**
  실행 세션 식별자

* **provider (required)**
  이벤트 유래 모델/시스템

* **mode (optional)**
  OMC 실행 모드

* **agent_id (required)**
  이벤트를 발생시킨 에이전트

* **parent_agent_id (optional)**
  상위 에이전트 ID

* **role (required)**
  에이전트 역할

* **state (required)**
  에이전트 상태

* **type (required)**
  이벤트 타입

* **task_id (optional)**
  관련 태스크 ID

* **intent_ref (optional)**
  계획 단계 참조 ID

* **payload (optional)**
  원본 이벤트 세부 데이터

* **metrics (optional)**
  성능/비용 메타데이터

* **raw_ref (optional)**
  원본 로그 포인터 (디버깅용)

---

## 4) 상태 전이 규칙

### 정규 전이
* `idle → running` : 작업 시작
* `running → waiting` : 외부 결과 대기
* `waiting → running` : 결과 수신 후 재개
* `running → blocked` : 의존성 미충족
* `blocked → running` : 의존성 해소
* `running → error` : 실패
* `error → running` : 재시도
* `running → done` : 정상 완료

### 종료 전이
* `error → failed` : max retry 초과, 최종 실패 확정
* `idle → cancelled` : 시작 전 취소
* `running → cancelled` : 실행 중 취소
* `blocked → cancelled` : blocked 상태에서 취소
* `blocked → error` : blocked timeout 발생
* `waiting → error` : 대기 중 외부 호출 timeout/실패
* `done → idle` : 에이전트 재활용 (새 태스크 할당)

### 원칙
* `done`/`failed`/`cancelled`은 terminal state (done→idle 재활용 제외)
* `done`은 성공, `failed`은 실패, `cancelled`은 외부 취소를 의미
* 유효하지 않은 상태 전이 발생 시: warning 로그 + 이벤트는 수용 (drop하지 않음)

---

## 5) 이벤트 타입 규약

### state와 type의 직교성 원칙
* **state**: 에이전트의 현재 상태 (who is doing what)
* **type**: 이벤트의 종류 (what happened)
* 예: `state=running, type=verify` = "에이전트가 실행 중이며 검증을 수행함"
* 예: `state=error, type=verify` = "에이전트가 검증 수행 중 에러 발생" (verify 결과 fail은 `state=running, payload.result=fail`로 표현)

### 타입 목록
* **task_spawn** : 자식 태스크 생성
* **task_update** : 진행/상태 업데이트
* **task_done** : 태스크 종료
* **tool_call** : 도구 호출 시작
* **tool_result** : 도구 호출 결과
* **message** : 에이전트 메시지
* **error** : 예외/오류
* **replan** : 계획 재수립
* **verify** : 검증 단계 진입/결과
* **fix** : 수정 단계 진입/결과
* **recover** : 장애 복구 시도 (Ultrapilot 모드)
* **state_change** : 에이전트 상태 전이 기록

---

## 6) OMC 모드별 기대 패턴

### Ralph

* `verify → fix → verify` 반복 이벤트 빈도 높음
* `done` 전까지 loop 지속

### Ultrawork

* `task_spawn` fan-out 빈도 높음
* 병렬 agent 다수 `running` 상태

### Ultrapilot

* 단일 목표 중심의 긴 연속 실행
* 중간 `replan` + `recover` 이벤트 중요

---

## 7) Provider 매핑 가이드 (초안)

* Claude Code 내부 이벤트 → `provider = claude`
* Gemini MCP 호출 결과 → `provider = gemini`
* Codex MCP 호출 결과 → `provider = codex`
* 시스템 / 파서 / 스토어 이벤트 → `provider = system`

---

## 8) Redaction 규칙 (기본 ON)

### 마스킹 대상 키/패턴

**키명:**

* `api_key`, `token`, `secret`, `password`, `authorization`
* `credential`, `private_key`, `access_key`, `secret_key`, `conn_string`, `passwd`

**값 패턴:**

* `sk-...` (OpenAI/Anthropic API 키)
* `AKIA...` (AWS Access Key)
* `AIza...` (GCP API Key)
* `ghp_...`, `gho_...`, `ghu_...` (GitHub 토큰)
* `Bearer ...` 토큰
* `-----BEGIN (RSA|EC|OPENSSH) PRIVATE KEY-----`
* 긴 base64/hex 토큰 (40자 이상 연속 hex/base64)

### 처리 규칙

* 원문 보존 금지 (기본)
* 출력은 `***REDACTED***`
* payload 내 중첩 객체에도 재귀 적용 (depth 10 제한)
* 마스킹 해제: CLI `--no-redact` 플래그 (명시적 동의)
* 오탐(false positive) 발생 시: redaction 로그 기록 + 화이트리스트 지원 예정(v0.2)

---

## 9) 필드 포맷 규칙

| 필드 | 타입 | 포맷 | 예시 |
|------|------|------|------|
| `ts` | string | ISO 8601 UTC, 밀리초 허용 | `2026-02-17T22:27:00.123Z` |
| `run_id` | string | `^run-[a-zA-Z0-9_-]+$` | `run-20260217-2227` |
| `provider` | string | enum | `claude` |
| `mode` | string | enum, optional | `ralph` |
| `agent_id` | string | `^[a-zA-Z0-9_-]{1,64}$` | `coder-auth` |
| `parent_agent_id` | string | 동일 포맷, optional | `planner-main` |
| `role` | string | enum | `executor` |
| `state` | string | enum | `running` |
| `type` | string | enum | `task_spawn` |
| `task_id` | string | `^task-[a-zA-Z0-9_-]+$`, optional | `task-42` |
| `intent_ref` | string | `^plan-[a-zA-Z0-9_-]+$`, optional | `plan-7` |
| `payload` | object | 타입별 구조 (아래 참조), optional | `{}` |
| `metrics` | object | 고정 구조, optional | 아래 참조 |
| `raw_ref` | string | URI, optional | `file:///path/to/log` |

### metrics 필드 구조

| 필드 | 타입 | 설명 |
|------|------|------|
| `latency_ms` | number (>= 0) | 왕복 지연 (밀리초) |
| `tokens_in` | integer (>= 0) | 입력 토큰 수 |
| `tokens_out` | integer (>= 0) | 출력 토큰 수 |
| `cost_usd` | number (>= 0) | 비용 (USD) |

provider별 일부 필드 unavailable 시 `null` 허용.

### payload 타입별 구조

#### task_spawn
```json
{ "title": "string", "child_agent": "string", "priority": "number (optional)" }
```

#### task_update
```json
{ "progress": "number (0-100)", "message": "string (optional)" }
```

#### task_done
```json
{ "result": "success|failure|cancelled", "summary": "string (optional)" }
```

#### tool_call
```json
{ "tool_name": "string", "args": "object (optional)" }
```

#### tool_result
```json
{ "tool_name": "string", "success": "boolean", "output_preview": "string (max 500 chars)" }
```

#### error
```json
{ "error_type": "string", "message": "string", "stack": "string (optional)" }
```

#### verify
```json
{ "result": "pass|fail", "reason": "string (optional)", "evidence": "string (optional)" }
```

#### fix
```json
{ "target": "string", "strategy": "string", "files_changed": "string[] (optional)" }
```

#### replan / recover
```json
{ "reason": "string", "new_plan_ref": "string (optional)" }
```

#### state_change
```json
{ "from": "state enum", "to": "state enum", "trigger": "string (optional)" }
```

---

## 10) 유효성 검증 규칙

### Required 필드

* `ts`
* `run_id`
* `provider`
* `agent_id`
* `role`
* `state`
* `type`

### Enum 검증

* `provider / mode / role / state / type`는 정의된 enum만 허용
* 미등록 값은 `unknown`으로 강등 + warning 기록

### 오류 처리

* 파싱 실패 이벤트는 drop하되 카운터 증가
* store / TUI는 파싱 실패로 종료되면 안 됨

---

## 11) Role 매핑 테이블

OMC agent catalog → CanonicalEvent role 매핑:

| OMC Agent | role | 비고 |
|-----------|------|------|
| planner | planner | |
| executor, deep-executor | executor | |
| explore | explorer | |
| architect | architect | |
| debugger | debugger | |
| verifier | verifier | |
| designer | designer | |
| code-reviewer, style-reviewer, quality-reviewer, api-reviewer, performance-reviewer | reviewer | |
| security-reviewer | guard | 보안 역할 |
| test-engineer | tester | |
| writer | writer | |
| analyst, product-manager, product-analyst, ux-researcher, information-architect | planner | 분석/기획 역할은 planner로 통합 |
| build-fixer | executor | 빌드 수정은 실행 역할 |
| scientist | explorer | 데이터 분석은 탐색 역할 |
| dependency-expert | explorer | 외부 조사는 탐색 역할 |
| git-master | executor | git 작업은 실행 역할 |
| qa-tester | tester | |
| critic | reviewer | |
| 미등록 역할 | custom | warning 로그 기록 |

---

## 12) 샘플 이벤트

### Sample: `task_spawn`

```json
{
  "ts": "2026-02-17T22:28:10Z",
  "run_id": "run-1",
  "provider": "claude",
  "mode": "ultrawork",
  "agent_id": "planner-main",
  "role": "planner",
  "state": "running",
  "type": "task_spawn",
  "task_id": "task-100",
  "payload": {
    "title": "Fix auth flow",
    "child_agent": "coder-auth"
  }
}
```

---

### Sample: `verify_fail` (Ralph 모드)

```json
{
  "ts": "2026-02-17T22:29:00Z",
  "run_id": "run-1",
  "provider": "claude",
  "mode": "ralph",
  "agent_id": "reviewer-1",
  "role": "reviewer",
  "state": "running",
  "type": "verify",
  "task_id": "task-100",
  "payload": {
    "result": "fail",
    "reason": "test regression"
  }
}
```

> state=running (에이전트가 실행 중), type=verify (검증 수행), payload.result=fail (검증 결과 실패)

### Sample: `tool_call` + `tool_result` 쌍

```json
{
  "ts": "2026-02-17T22:30:00Z",
  "run_id": "run-1",
  "provider": "claude",
  "agent_id": "coder-auth",
  "role": "executor",
  "state": "running",
  "type": "tool_call",
  "task_id": "task-100",
  "payload": { "tool_name": "Edit", "args": { "file": "auth.go" } }
}
```

```json
{
  "ts": "2026-02-17T22:30:02Z",
  "run_id": "run-1",
  "provider": "claude",
  "agent_id": "coder-auth",
  "role": "executor",
  "state": "running",
  "type": "tool_result",
  "task_id": "task-100",
  "payload": { "tool_name": "Edit", "success": true, "output_preview": "File updated" },
  "metrics": { "latency_ms": 2000, "tokens_in": 150, "tokens_out": 80, "cost_usd": 0.0015 }
}
```

### Sample: `fix` (Ralph verify-fix 루프)

```json
{
  "ts": "2026-02-17T22:31:00Z",
  "run_id": "run-1",
  "provider": "claude",
  "mode": "ralph",
  "agent_id": "coder-auth",
  "role": "executor",
  "state": "running",
  "type": "fix",
  "task_id": "task-100",
  "payload": { "target": "auth.go:42", "strategy": "fix test regression", "files_changed": ["auth.go", "auth_test.go"] }
}
```

### Sample: `error` (최종 실패)

```json
{
  "ts": "2026-02-17T22:35:00Z",
  "run_id": "run-1",
  "provider": "claude",
  "agent_id": "coder-auth",
  "role": "executor",
  "state": "failed",
  "type": "error",
  "task_id": "task-100",
  "payload": { "error_type": "max_retry_exceeded", "message": "3회 재시도 초과" }
}
```

---



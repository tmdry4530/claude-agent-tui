# OMC Agent TUI 컴포넌트 명세 v0.1

## 1) Panel A — Agent Arena
목적: 에이전트 상태를 즉시 파악

표시:
- 마스코트(색상+역할 배지)
- agent id/name
- state badge
- progress(0~100)
- blocked_by

정렬 우선순위:
`error > failed > blocked > running > waiting > idle > done > cancelled`

상태 스타일:
- running: pulse (밝은 색 + 약한 애니메이션)
- waiting: dim (저채도/밝기 감소)
- blocked: yellow border (`#E3B341`)
- error: red border (`#FF7B72`)
- done: check badge (`#56D364`)
- failed: red border + `X` badge (`#FF7B72`, done과 구분)
- cancelled: grey border + `-` badge (`#6E7681`)
- idle: 기본색, 모션 없음

---

## 2) Panel B — Live Timeline
목적: 이벤트 흐름 추적

행 포맷:
`HH:MM:SS | agent | type | summary`

필터:
- `/agent:`
- `/type:`
- `/provider:`
- `/state:`
- `/mode:`

동작:
- Enter: 상세 보기
- y: payload 복사(선택)

---

## 3) Panel C — Task Graph
목적: 태스크 관계/병목 시각화

표시:
- parent-child DAG
- critical path
- blocker chain

규칙:
- 기본 depth 3 (키 `[`/`]`로 조절)
- 50노드 초과 시 그래프 축약 모드 자동 전환
- v0.1-alpha: 미포함, v0.1-beta: 트리 목록으로 구현, v0.2: DAG 전환

---

## 4) Panel D — Inspector
목적: 선택 에이전트/태스크 상세 분석

섹션:
1. Summary
2. Intent(plan)
3. Action(tool calls/results)
4. Diff(intent vs action)
5. Verify/Fix history

핵심 지표:
- avg latency
- retries
- error count
- success ratio

---

## 5) Panel E — Footer Metrics
표시:
- tokens in/out
- avg latency
- error rate
- cost estimate
- live/replay 상태

---

## 6) Replay Controls (Panel F)
표시 (replay 모드 활성 시):
- 현재 배속 (1x/4x/8x/16x)
- 진행바 (현재 시점 / 전체 시간)
- step 표시 (현재 이벤트 번호 / 전체)
- play/pause 상태

동작:
- space: 재생/일시정지
- 좌/우 화살표: 한 스텝 뒤/앞
- 1/4/8/16: 배속 변경 (Inspector에 포커스 없을 때만)
- t: 특정 시각 점프
- n/N: 다음/이전 에러 이벤트 점프

---

## 7) Footer Metrics (Panel E)
표시:
- tokens in/out
- avg latency
- error rate
- cost estimate
- live/replay 상태
- redaction: ON/OFF indicator

---

## 8) 키바인딩 핵심
- q 종료
- tab 포커스 이동 (Footer 제외 — 표시 전용)
- / 필터
- ? 단축키 도움말 (오버레이)
- r replay 토글
- m mascot 토글
- p pause
- space (replay play/pause)

> 키 충돌 규칙: Replay 배속키(1/4/8/16)는 Inspector에 포커스가 없을 때만 동작.
> Inspector 포커스 시 숫자키는 탭 전환(1~5)으로 동작.

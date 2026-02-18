# MASTER PROMPT — OMC Agent TUI (Claude Code + OMC)

## EXECUTION CONTRACT (강제)
- 중간 상태 보고 후 대기 금지.
- 아래 완료 조건을 모두 만족할 때까지 연속 실행한다.
- 진행 중 질문/확인 요청은 금지(치명적 차단 이슈 제외).
- 각 단계 완료 시 즉시 다음 단계로 자동 진입한다.

### 완료 조건 (ALL REQUIRED)
1) TASKS.md의 남은 핵심 마일스톤 100% 완료
2) 테스트 전체 PASS
3) 빌드 성공
4) replay 포함 통합 검증 PASS
5) 문서(status/PROGRESS, NEXT_ACTIONS, DECISIONS) 최종 업데이트
6) 최종 결과 보고서 작성(변경 파일, 검증 증거, 미해결 리스크 0개 또는 명시)

### 중단 허용 조건 (ONLY)
- Preflight FAIL (Gemini/Codex 호출 실패)
- 동일 치명 오류 3회 연속 재현 + 대체경로 모두 실패
- 사용자의 명시적 STOP 지시

### 보고 정책
- 중간 보고 금지
- 최종 완료 보고 1회만 출력
- 단, 중단 허용 조건 발생 시 즉시 실패 보고

---

너는 Claude Code 환경에서 OMC(Oh My ClaudeCode)를 사용해 멀티에이전트 작업을 수행한다.

## 0) 절대 규칙 (가장 먼저 적용)
1. **구현 시작 전 Preflight를 반드시 수행**한다.
2. Preflight에서 **Gemini + Codex MCP 호출이 모두 성공**해야만 다음 단계로 진행한다.
3. 둘 중 하나라도 실패하면:
   - 즉시 작업 중단
   - 원인/증거를 `status/FAILURES.md`에 기록
   - 사용자에게 “호출 실패로 작업 취소”를 보고
   - 임의 우회/대체 구현 금지

---

## 1) Preflight 체크 (필수 게이트)
아래 순서를 정확히 수행:

### 1-1. 환경 확인
- OMC 활성 상태 확인
- MCP 목록에서 Gemini/Codex 사용 가능 상태 확인
- 필요한 ENV가 설정되었다고 가정하되, 누락 징후가 있으면 즉시 보고

### 1-2. 실제 호출 테스트
- Gemini에 최소 테스트 요청 1회
- Codex에 최소 테스트 요청 1회
- 각 호출에 대해:
  - 성공/실패
  - 왕복 시간
  - 에러 메시지(실패 시)
  - request/task 식별자(가능 시)
  를 기록

### 1-3. 게이트 판정
- 둘 다 성공: `PRECHECK: PASS`
- 하나라도 실패: `PRECHECK: FAIL` + **작업 취소**

### 1-4. 기록
- `status/PROGRESS.md`에 preflight 결과 요약
- 실패 시 `status/FAILURES.md`에 상세 기록

---

## 2) 프로젝트 목표
OMC 멀티에이전트 실행 상태를 실시간으로 보여주는 TUI 관제 도구를 설계/구현한다.

핵심 패널:
1. Agent Arena (마스코트 기반)
2. Live Timeline
3. Task Graph + Blocker Chain
4. Inspector (Intent vs Action diff)
5. Replay (Time Machine)

---

## 3) 강제 실행 원칙
- OMC **Ralph / Ultrawork / Ultrapilot** 중심으로 A→Z 중단 없이 진행
- 실패 시 verify-fix 루프 기반 복구 시도
- 완료 기준: 구현 + 검증 + 결과 보고
- 부분 완료 상태로 종료 금지
- 환경변수 값 생성/노출 금지 (사용자 세팅 전제)

---

## 4) 입력 문서 (반드시 참조)
- PRD.md
- COMPONENT_SPEC.md
- ARCHITECTURE.md
- TASKS.md
- references/event-schema.md
- references/keybindings.md
- references/mascot-guidelines.md
- status/PROGRESS.md
- status/FAILURES.md
- status/NEXT_ACTIONS.md
- status/DECISIONS.md

---

## 5) 작업 절차
1. Preflight 수행 및 게이트 판정
2. TASKS/STATUS 기준으로 이번 턴 목표 1~3개 선택
3. 작업 수행
4. 검증 증거 수집(테스트/로그/체크)
5. status 파일 업데이트
6. 다음 액션 제시(최대 3개)

---

## 6) 실패 처리 규칙
실패 시 `status/FAILURES.md`에 다음 항목 필수 기록:
- 증상
- 재현 방법
- 추정 원인
- 시도한 해결
- 대안 2개
- 다음 액션
- 상태(Open/Mitigating/Closed)

동일 실패 2회 반복 시 접근 방식 변경(도구/경로/단계 분할).

---

## 7) 출력 형식 (매 턴 고정)
1) 이번에 한 일
2) 검증 결과(증거)
3) 업데이트 파일
4) 남은 리스크
5) 다음 액션(최대 3개)

---

## 8) 금지 사항
- Preflight FAIL 상태에서 작업 진행 금지
- 근거 없는 완료 보고 금지
- 민감정보 출력 금지


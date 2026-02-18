# PHASE 0 — DISCOVERY (OMC Agent TUI)

## 0) Preflight 게이트 (필수)
디스커버리 시작 전 Gemini/Codex MCP 호출 가능 여부를 확인한다.

절차:
1. OMC 활성 확인
2. Gemini 테스트 호출 1회
3. Codex 테스트 호출 1회
4. 판정:
   - 둘 다 성공: `PRECHECK: PASS`
   - 하나라도 실패: `PRECHECK: FAIL` + 작업 취소

실패 시:
- `status/FAILURES.md` 기록
- 사용자에게 취소 사유 보고
- 우회 진행 금지

---

## 1) 목표
문제/요구/제약/성공지표를 확정 가능한 수준으로 정리한다.

## 2) 입력 문서
- PRD.md
- TASKS.md
- references/event-schema.md
- status/PROGRESS.md
- status/FAILURES.md
- status/DECISIONS.md

## 3) 수행 지침
1. 불명확 요구사항 식별
2. 사용자/시나리오/문제 3개/성공지표 정리
3. OMC 모드별(ralph/ultrawork/ultrapilot) 관찰 포인트 분리
4. event-schema 필수/선택 필드 제안
5. 불확실 사항은 Assumption으로 명시

## 4) 출력 형식
1) Discovery Summary
2) Confirmed Requirements
3) Open Questions
4) Assumptions
5) Next Actions(최대 5개)


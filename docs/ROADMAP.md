# ROADMAP (v0.1 ~ v0.3)

## v0.1-alpha — 핵심 관찰 가능성 확보
목표:
- 에이전트 상태와 이벤트 흐름의 기본 관찰 가능

범위:
- Agent Arena (12종 role, 8종 state)
- Live Timeline (필터, 자동 스크롤)
- Security Redaction (기본 ON)
- Footer Metrics

완료 기준:
- UX S1(실시간 관제) 시나리오 pass

---

## v0.1-beta — 분석 기능 추가
목표:
- 병목 분석과 장애 재현 기능 확보

범위:
- Task Graph (트리 목록, critical path, blocker chain)
- Inspector (Summary, Intent, Action, Diff, Verify/Fix history)
- Replay (JSONL 기반, 배속 조절, step 탐색)

완료 기준:
- UX S1~S3, F1~F3, R1 시나리오 pass

---

## v0.2 — 분석 품질 강화
목표:
- 원인 분석 속도 개선, 비용 가시화

범위:
- Intent vs Action diff 고도화
- FR-7 Cost/Token Analytics (히트맵)
- blocker chain 시각 개선
- 필터 UX 강화
- 알림 (선택)
- 성능/안정성 튜닝

완료 기준:
- 병목 식별 시간 1분 이내 달성

---

## v0.3 — 운영 확장
목표:
- 실사용 운영 레벨 준비

범위:
- 플러그인형 provider 확장 포인트 정리
- 리포트/export 강화 (JSON/CSV/OTel)
- 접근성/테마/설정 고도화
- 운영 리포트 자동 요약

완료 기준:
- 운영 문서 + 릴리즈 가이드 + 데모 패키지 완성

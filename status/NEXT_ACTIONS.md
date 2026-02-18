# NEXT ACTIONS

## 완료 (세션 1-4)
- [x] 설계: PRD, event-schema, ARCHITECTURE, COMPONENT_SPEC
- [x] 구현: schema, collector, normalizer, tui, store, replay, graph, inspector
- [x] 통합: 파이프라인 Collector→Normalizer→Store→TUI
- [x] 안정화: flaky 테스트 수정, Arena/Timeline/Footer 테스트 추가
- [x] 패키징: .claude-plugin/plugin.json, docs/INSTALL.md, Makefile

## 완료 (세션 5 — 릴리즈 문서화)
- [x] CHANGELOG.md (v0.1.0-beta, Added/Changed/Fixed/Known Issues)
- [x] README.md (설치/실행/플러그인/이벤트 포맷/프로젝트 구조)
- [x] RELEASE_NOTES.md (하이라이트/검증/리스크/업그레이드)
- [x] status 문서 최종 갱신

## 완료 (세션 6 — CLCO 마스코트 Agent Arena)
- [x] Arena 패널 전면 리팩토링 (mascot, palette, selection, focus)
- [x] CLCO 마스코트 유니코드 + ASCII fallback
- [x] palette.md 기준 역할/상태 색상 적용
- [x] 키바인딩: hjkl/화살표 + Enter Inspector 연동
- [x] 테스트 32개 추가 (154 total)
- [x] 데모 이벤트 확장 (8 에이전트, 다양한 상태)

## 완료 (세션 7 — OMC Bridge 실시간 동기화)
- [x] internal/bridge 패키지 (tracking 변환 + event emitter)
- [x] --convert 플래그 (subagent-tracking.json → JSONL)
- [x] OMC hook 스크립트 (scripts/omc-bridge-hook.sh)
- [x] spawn/update/done/error 4종 이벤트 매핑
- [x] CanonicalEvent 스키마 100% 준수
- [x] bridge 테스트 14개 추가

## 완료 (세션 8 — Plugin UX)
- [x] 슬래시 명령 5종 (.claude/commands/)
- [x] plugin.json v0.2.0
- [x] README Quick Start

## v0.3.0 (다음 릴리즈)
1. Timeline → Inspector Enter 키 연동
2. Replay TUI 컨트롤 (Space 일시정지, +/- 속도)
3. Arena 에이전트 필터링 (역할/상태별)

## P2: 안정화
1. E2E/통합 테스트 추가
2. 대용량 부하 테스트 (10K+ 이벤트)
3. Makefile GO 경로 일반화 (`GO?=$(shell which go)`)

## P3: 기능 확장
1. 에이전트/이벤트 필터링
2. 검색 기능 (/ 키)
3. 키보드 도움말 (? 키)
4. Arena 에이전트 정렬 (상태/역할/최근 활동 기준)

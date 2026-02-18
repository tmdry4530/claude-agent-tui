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

## v0.2.0 (다음 릴리즈)
1. Arena 포커스 하이라이트 보더 + 시각적 구분
2. Timeline → Inspector Enter 키 연동
3. Replay TUI 컨트롤 (Space 일시정지, +/- 속도)

## P2: 안정화
1. E2E/통합 테스트 추가
2. 대용량 부하 테스트 (10K+ 이벤트)
3. Makefile GO 경로 일반화 (`GO?=$(shell which go)`)

## P3: 기능 확장
1. 에이전트/이벤트 필터링
2. 검색 기능 (/ 키)
3. 키보드 도움말 (? 키)

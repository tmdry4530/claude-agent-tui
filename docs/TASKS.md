# OMC Agent TUI TASKS

## M0. 부트스트랩
- [x] 프로젝트 문서 구조 확정
- [x] status/* 초기화
- [x] references/* 핵심 스키마 채우기

## M1. 이벤트 스키마/정규화
- [x] CanonicalEvent 필드 확정 (타입 정의, state/type 직교성 원칙 추가)
- [x] OMC mode 매핑 규칙 확정 (pipeline/ecomode 추가)
- [x] redaction 규칙 확정 (클라우드 자격증명/재귀 마스킹 추가)
- [x] Role 매핑 테이블 작성 (25+ agent → 12 role)
- [x] 상태 전이 모델 확장 (done/failed/cancelled terminal state)

## M2. 컴포넌트 설계
- [x] Agent Arena 상세 명세
- [x] Timeline 상세 명세
- [x] Task Graph 상세 명세
- [x] Inspector 상세 명세

## M3. 아키텍처 고정
- [x] 모듈 경계 확정(collector/normalizer/store/tui/replay)
- [x] 장애/복구 전략 확정
- [x] 성능 목표 수치 확정 (50 events/s, 10 FPS, RSS 500MB)
- [x] 기술 스택 결정: Go + Bubbletea (ADR-006)

## M4. UX 시나리오
- [x] 정상 시나리오 3개 (S1~S3)
- [x] 실패 시나리오 3개 (F1~F3)
- [x] replay 시나리오 2개 (R1~R2)

## M5. 실행 준비
- [x] MASTER_PROMPT 최종 고정
- [x] PHASE 0/1/2 프롬프트 고정
- [x] 우선순위 Top 5 재정렬
- [x] PRD P1 보완 완료 (각 FR에 acceptance criteria 추가)
- [x] event-schema P1 보완 완료 (payload 구조, 필드 포맷, 샘플 확장)
- [x] ARCHITECTURE.md v0.2 (Go+Bubbletea, 프로젝트 구조, goroutine 파이프라인)
- [x] COMPONENT_SPEC failed/cancelled state 반영
- [x] Architect 검증 1차: NEEDS_REVISION (MAJOR 3건)
- [x] MAJOR Fix: mascot-guidelines 12종 role 시각 정의 추가
- [x] MAJOR Fix: mascot-guidelines + palette FAILED/CANCELLED 상태 추가
- [x] MAJOR Fix: ROADMAP alpha/beta 분할 + FR-7 반영
- [x] MINOR Fix: event-schema 섹션 번호 정리 (10→10/11/12)
- [x] MINOR Fix: COMPONENT_SPEC Replay Controls + Footer + 키 충돌 규칙
- [x] MINOR Fix: RISK_REGISTER R-08~R-10 추가
- [x] Architect 재검증: APPROVED (MAJOR 3건 전체 RESOLVED)

## M6. v0.1-alpha 구현
- [x] Go 프로젝트 초기화 (go mod init, 디렉토리 구조)
- [x] Go 1.23.6 설치 + 의존성 (bubbletea/lipgloss/bubbles/fsnotify/testify)
- [x] pkg/schema: CanonicalEvent struct + 5 enums + transitions + role_map + 11 payloads
- [x] internal/collector: FileCollector + fsnotify + circuit-breaker (3→10s/30s/60s)
- [x] internal/normalizer: Normalize + Redactor (키 11 + 패턴 8 + 재귀 depth 10)
- [x] internal/tui: Bubbletea Model (Arena + Timeline + Footer)
- [x] cmd/omc-tui: main.go CLI 엔트리포인트
- [x] 전체 테스트: 23/23 PASS
- [x] 바이너리 빌드: omc-tui BUILD OK
- [x] Architect 검증: APPROVED (MAJOR 0, MINOR 5)

## M7. v0.1-alpha 고도화 (다음 세션)
- [ ] internal/store: ring buffer + RunState/AgentState/TaskState
- [ ] 파이프라인 통합: Collector → Normalizer → Store → TUI (live 모드)
- [ ] Footer metrics 집계 (latency/tokens/cost)
- [ ] State transition validation 적용

## M8. v0.1-beta (Graph + Inspector + Replay)
- [ ] internal/replay: JSONL 로딩 + 가상 시계
- [ ] internal/tui/graph: Task Graph 패널
- [ ] internal/tui/inspector: Inspector 패널
- [ ] 키바인딩 통합 (keybindings.md 전체 적용)

---

## 강제 운영 원칙
- [ ] OMC Ralph/Ultrawork/Ultrapilot로 A→Z 완주
- [ ] 실패 시 verify-fix 루프
- [ ] 부분완료 금지 (구현+검증+보고)
- [ ] 환경변수는 사용자 세팅 전제

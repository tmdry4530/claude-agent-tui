# DECISIONS (ADR)

## ADR-001: 운영 환경
- 결정: Claude Code + OMC(Oh My ClaudeCode) 기반으로 진행
- 근거: 현재 사용자 실사용 환경과 일치
- 상태: Accepted

## ADR-002: 실행 전략
- 결정: Ralph/Ultrawork/Ultrapilot 모드 중심 A→Z 무중단 진행
- 근거: 사용자의 명시적 요구
- 상태: Accepted

## ADR-003: 문서 우선 원칙
- 결정: 현재 단계에서는 코드 작성 없이 문서/설계 우선
- 근거: 사용자 지시(코드 작성 금지)
- 상태: Accepted

## ADR-004: 보안 기본값
- 결정: 이벤트/로그 redaction 기본 ON
- 근거: 민감정보 노출 리스크 최소화
- 상태: Accepted

## ADR-005: 시각 아이덴티티
- 결정: CLCO 마스코트 기반 Agent Arena 구성
- 근거: 사용자 요청 + 차별화 요소
- 상태: Accepted

## ADR-006: 기술 스택
- 결정: Go + Bubbletea (charmbracelet/bubbletea)
- 근거:
  - 개발 생산성 최우수 (Charmbracelet 생태계: Lip Gloss, Bubbles)
  - MVU 패턴이 Collector→Normalizer→Store→Renderer 파이프라인과 자연 매핑
  - goroutine/channel이 다중 소스 수집에 적합
  - 단일 정적 바이너리 배포 (런타임 의존 없음)
  - Lip Gloss 내장 TrueColor/256/ANSI 자동 감지
  - 가중합산 점수: Go 4.35 > Rust 4.30 > TS 3.25
- 대안: Rust+Ratatui (성능 극한 필요 시 v0.3+ 재검토), TS+Ink (비추천)
- 상태: Accepted

## ADR-007: MVP 범위 분할
- 결정: v0.1을 alpha/beta 2단계로 분할
  - v0.1-alpha: Agent Arena + Timeline + Redaction
  - v0.1-beta: Task Graph(트리) + Inspector(기본) + Replay(기본)
- 근거: 6개 기능 동시 MVP는 범위 과다, 단계적 출시로 리스크 감소
- 상태: Accepted

## ADR-008: 상태 전이 모델 확장
- 결정: terminal state를 done/failed/cancelled 3종으로 분리
- 근거: Ralph max retry 초과(→failed), 사용자 취소(→cancelled) 표현 불가 해소
- 상태: Accepted

## ADR-009: Role 매핑 전략
- 결정: 추상 role 12종 + 매핑 테이블 방식 (A안 확장)
- 근거: OMC 25+ agent를 7개로 압축하면 Arena 필터링 가치 감소, 12종이 적정 균형
- 상태: Accepted

## ADR-010: CLCO 마스코트 Arena 구현
- 결정: 유니코드 블록 기반 미니 스프라이트 + ASCII fallback 이중 렌더링
- 근거:
  - 유니코드 지원 터미널에서 시각적 차별화 (역할별 고유 아이콘)
  - ASCII fallback으로 256-color 터미널 호환성 보장
  - palette.md 색상을 arena에 직접 적용하여 시각 일관성 확보
- 대안: 이미지 기반 렌더링 (터미널 호환성 부족으로 기각)
- 상태: Accepted

## ADR-011: Arena 키바인딩 설계
- 결정: hjkl/화살표로 에이전트 선택, Enter로 Inspector 연동
- 근거:
  - Vim 스타일 내비게이션이 TUI 사용자에게 익숙
  - Enter 키로 선택한 에이전트의 최신 이벤트를 Inspector에 표시
  - Tab은 패널 간 전환, hjkl은 패널 내 내비게이션으로 계층 분리
- 상태: Accepted

## ADR-012: 상태별 시각 규칙
- 결정: palette.md 상태 색상 + CSS-like 시각 규칙 조합
  - running: bold + #7EE787 (강조)
  - waiting: faint + #7D8590 (저채도)
  - blocked: bold + #E3B341 (경고, 테두리 강조)
  - error: bold+blink + #FF7B72 (오류, 테두리 강조)
  - done: normal + #56D364 (완료)
  - failed: strikethrough + #FF7B72 (취소선)
- 근거: palette.md 색상 시스템과 일관성 유지, 상태 구분 시인성 극대화
- 상태: Accepted

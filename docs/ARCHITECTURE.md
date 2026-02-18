# OMC Agent TUI 아키텍처 v0.2

## 1) 목적
OMC(Oh My ClaudeCode) 멀티에이전트 실행 이벤트를 수집/정규화해서 TUI로 실시간 시각화한다.

## 2) 기술 스택 (ADR-006)
- **언어**: Go 1.22+
- **TUI 프레임워크**: Bubbletea (charmbracelet/bubbletea)
- **스타일링**: Lip Gloss (charmbracelet/lipgloss) — TrueColor/256/ANSI 자동 감지
- **컴포넌트**: Bubbles (charmbracelet/bubbles) — spinner, progress, table, viewport
- **파일 감시**: fsnotify/fsnotify
- **JSON 파싱**: encoding/json (표준 라이브러리)
- **테스트**: go test + testify
- **린트**: golangci-lint

## 3) 아키텍처 개요

### 파이프라인
```
OMC hooks/logs → Collector (goroutine) → channel → Normalizer/Redactor → Store → TUI Renderer (Bubbletea)
```

### MVU 패턴 매핑
```
Model:  Store (RunState, AgentState, TaskState, EventBuffer)
View:   Renderer (Arena, Timeline, Graph, Inspector, Footer)
Update: Bubbletea Msg/Cmd (이벤트 수신, 키 입력, 타이머)
```

## 4) 모듈 경계

### 4.1 Collector (`internal/collector/`)
- 입력: OMC 훅, stdout/stderr, 파일 로그(JSONL/txt)
- 출력: RawEvent → Go channel
- 구현: goroutine per source, fsnotify로 `.omc/state/` 감시
- 책임: 입력 수집, 타임스탬프 보정, 소스 장애 격리
- 장애 격리: 소스별 circuit-breaker (3회 연속 실패 시 backoff)

### 4.2 Normalizer (`internal/normalizer/`)
- RawEvent → CanonicalEvent 변환 (Go struct with JSON tags)
- provider/mode/role/state/type 정규화 (enum validation)
- role 매핑 테이블 적용 (event-schema.md 섹션 11)
- 민감정보 마스킹(redaction) — 재귀 depth 10 제한

### 4.3 Store (`internal/store/`)
- RunState/AgentState/TaskState 유지
- 이벤트 ring buffer (기본 10,000건, oldest drop)
- 메트릭 집계(latency/tokens/error/cost)
- 상태 전이 검증 (invalid transition → warning 로그, 이벤트 수용)

### 4.4 TUI (`internal/tui/`)
- Bubbletea Program + Lip Gloss 스타일
- 5개 패널: Arena / Timeline / Graph / Inspector / Footer
- 필터/포커스/키바인딩 처리 (keybindings.md 참조)
- live/replay 모드 전환
- 80x24 fallback 레이아웃

### 4.5 Replay (`internal/replay/`)
- JSONL 로딩 (최대 100MB, 초과 시 스트리밍)
- 가상 시계 기반 재생(1x/4x/8x/16x)
- step 탐색 (이벤트 1건 단위)

## 5) 데이터 흐름
```
1. Collector goroutine: 소스 감시 → RawEvent → eventCh (buffered channel)
2. Normalizer goroutine: eventCh → CanonicalEvent → storeCh
3. Store: storeCh → 상태 갱신 → Bubbletea Msg 발행
4. TUI Update: Msg 수신 → Model 갱신 → View 재렌더
5. Footer: Store 메트릭 → 주기적 갱신 (1초 간격)
```

## 6) 장애/복구 전략
- 입력 소스 실패: 소스 단위 circuit-breaker (3회 → 10s backoff → 30s → 60s)
- 파싱 실패: 이벤트 drop + warning count 증가 (연속 100건 시 경고 배지)
- 렌더 실패: 패널 단위 recover (프로세스 종료 금지)
- burst 트래픽: 50 events/s 초과 시 샘플링 모드 전환 + 사용자 안내 배지
- 채널 백프레셔: buffered channel (1,000건), 초과 시 oldest drop

## 7) 보안
- redaction 기본 ON (event-schema.md 섹션 8 참조)
- 확장된 패턴: token/key/secret/password + 클라우드 자격증명 + SSH 키
- raw 로그 저장은 opt-in (`--no-redact` CLI 플래그)
- payload 재귀 마스킹 depth 10

## 8) OMC 모드 반영
- Ralph: verify-fix 루프 강조 (Inspector verify/fix history 탭)
- Ultrawork: fan-out 병렬 작업 강조 (Arena 다중 running 카드)
- Ultrapilot: A→Z 연속 실행 추적 강조 (Timeline 연속 스트림)
- Pipeline: 단계별 체인 시각화 (Graph에서 sequential flow)
- Ecomode: 모델 라우팅 표시 (haiku/sonnet 배지)

## 9) 프로젝트 구조 (예정)
```
omc-agent-tui/
├── cmd/
│   └── omc-tui/          # main.go (CLI 엔트리포인트)
├── internal/
│   ├── collector/         # 이벤트 수집
│   ├── normalizer/        # 정규화 + redaction
│   ├── store/             # 상태 관리
│   ├── tui/               # Bubbletea 모델/뷰
│   │   ├── arena/         # Agent Arena 패널
│   │   ├── timeline/      # Live Timeline 패널
│   │   ├── graph/         # Task Graph 패널
│   │   ├── inspector/     # Inspector 패널
│   │   └── footer/        # Footer 메트릭
│   └── replay/            # JSONL 재생
├── pkg/
│   └── schema/            # CanonicalEvent 타입 정의
└── go.mod
```

## 10) 확장 포인트
- provider plugin: Go interface (`EventSource`) 구현으로 신규 provider 추가
- exporter: JSON/CSV/OTel — Store 데이터를 외부 포맷으로 내보내기
- alert sink: webhook/desktop notification

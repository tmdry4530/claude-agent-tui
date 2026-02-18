# PRD — OMC Agent TUI

## 1. 문제 정의
OMC(Oh My ClaudeCode) 기반 멀티에이전트 실행은 강력하지만, 실행 중 “누가/무엇을/왜” 하는지 즉시 파악하기 어렵다.  
현재는 로그 중심 관찰이라 병목, 실패 원인, verify-fix 루프의 상태를 빠르게 추적하기 힘들다.

### 현재 Pain Point
- 에이전트별 현재 상태(실행/대기/블록/오류) 가시성 부족
- 계획(Intent)과 실제 실행(Action)의 괴리 추적 어려움
- 장애/재시도/병목이 텍스트 로그에 묻힘
- 세션 종료 후 재현(replay) 난이도 높음

---

## 2. 목표
### 제품 목표
1. 사용자(개발자)가 **2초 이내**에 현재 멀티에이전트 상태를 이해한다.
   - 측정: 5개 이상 에이전트 동시 실행 중, 특정 에이전트의 현재 state를 정확히 답하는 시간 중앙값
2. 사용자(개발자)가 **1분 이내**에 병목 원인을 특정한다.
   - 측정: blocker chain 2단계 이상 시나리오(UX S2)에서 병목 에이전트를 특정하는 시간
3. 실행 세션을 replay하여 실패 원인을 **재현 가능**하게 만든다.
   - 측정: JSONL 로드 후 원본 실행의 error 이벤트와 동일 순서/내용이 Inspector에서 확인 가능하면 PASS

### 운영 목표
- OMC Ralph/Ultrawork/Ultrapilot 모드에서 A→Z 실행 관찰 지원
- 실패 발생 시 verify-fix 루프 상태를 명확히 시각화

### 비목표(초기)
- TUI에서 코드 직접 편집/실행 제어
- 브라우저 UI 우선 개발
- 벤더별 API 심층 통합(초기에는 이벤트 레벨 통합)

---

## 3. 핵심 사용자
- Claude Code + OMC로 실제 개발하는 개인/팀 개발자
- 멀티에이전트 자동화 결과를 검증/관제해야 하는 Tech Lead

---

## 4. 핵심 기능 (Functional Requirements)
### FR-1 Agent Arena
- 에이전트별 마스코트 카드 표시 (최소: role 배지 + agent_id + state badge + progress)
- 상태: idle/running/waiting/blocked/error/done/failed/cancelled
- 역할/진행률/blocked-by 표시
- 최대 동시 표시 12카드, 초과 시 스크롤
- AC: 에이전트 상태 변경 후 500ms 이내에 TUI 카드 반영

### FR-2 Live Timeline
- 실시간 이벤트 스트림 표시 (ring buffer 10,000건)
- 필터: agent/type/provider/state/mode (200ms 이내 결과 표시)
- 상세 이벤트 inspect 지원
- 자동 스크롤 ON (수동 스크롤 시 일시 해제, 맨 아래 이동 시 재활성)
- AC: 필터 적용/해제 시 200ms 이내 반영, 50 events/s에서 프레임 드랍 없음

### FR-3 Task Graph
- parent-child 태스크 DAG 시각화 (v0.1-beta: 트리 목록, v0.2: DAG 전환)
- critical path 표시
- blocker chain 표시
- 기본 depth 3, 50노드 초과 시 축약 모드 자동 전환
- AC: 10개 이하 태스크에서 겹침 없이 표시, 50개 이상에서 축약 자동 전환

### FR-4 Inspector (Intent vs Action)
- 선택 에이전트의 계획(Intent)과 실제 행동(Action) 비교
- Intent 데이터 소스: CanonicalEvent의 `intent_ref` 필드 → OMC plan/task description에서 추출
- Intent 데이터 부재 시 "Intent 데이터 없음" 메시지 표시 (Action history만 표시)
- 이탈(diff), 재시도, verify-fix 히스토리 표시
- AC: 선택 에이전트 변경 후 300ms 이내에 Inspector 갱신

### FR-5 Replay
- JSONL 기반 세션 재생 (사용자가 파일 경로를 CLI 인자로 지정)
- 속도 조절(1x/4x/8x/16x), step 탐색 (이벤트 1건 단위)
- 최대 파일 크기: 100MB (초과 시 스트리밍 파싱)
- AC: 배속 변경 시 UI 프리징 없이 전환, step 탐색 시 모든 패널이 해당 시점으로 동기화

### FR-6 Security Redaction
- token/key/secret/password/credential 자동 마스킹 기본 ON
- 클라우드 자격증명 패턴: `AKIA*`, `AIza*`, `ghp_*`, `sk-*`
- 마스킹 해제(opt-out): CLI 플래그 `--no-redact` (명시적 동의 필요)
- AC: redaction 적용 후 원문이 TUI 표시/메모리 버퍼에 노출되지 않음

---

## 5. 비기능 요구사항 (NFR)
- 성능: 50 events/s 처리 시 렌더 10 FPS 이상 유지, 1시간 세션 누적 RSS 500MB 이하
- 안정성: 파서 오류 시 해당 이벤트 drop + 연속 100건 실패 시 경고 표시, 프로세스 유지
- 가독성: 80x24 터미널에서 Agent Arena + Timeline이 기능 저하 없이 표시
- 보안: redaction 기본 ON (event-schema.md 참조), 민감정보 노출 방지
- 이식성: Linux/macOS 터미널 (xterm, iTerm2, gnome-terminal, Windows Terminal), SSH 원격 지원
- 시작 시간: TUI 첫 화면 렌더까지 2초 이내 (cold start)

---

## 6. 성공 지표 (Success Metrics)
- 병목 식별 평균 시간: 1분 이내 (UX S2 시나리오 기준, 3명 이상 테스터 측정)
- 장애 원인 회귀 분석 성공률: 80% 이상 (replay에서 error root cause agent를 정확히 특정하면 성공)
- 디버깅 시간 30% 단축: v0.1 출시 전 baseline 측정 후 비교 (v0.2 시점에 정식 평가)

---

## 7. 제약사항
- 환경변수는 사용자 직접 세팅(본 프로젝트는 ENV 존재를 가정)
- OMC 및 외부 MCP 이벤트 형식 변동 가능성 존재
- 초기 단계는 “관찰/분석” 중심, 실행 제어는 제외

---

## 8. 운영 원칙
- OMC Ralph/Ultrawork/Ultrapilot 기반 A→Z 완주 추적
- 실패 시 중단하지 않고 verify-fix 루프 중심으로 복구 상태 표시
- 부분완료 보고 금지: 구현+검증+결과 보고까지 완료 기준

---

## 9. 릴리즈 범위

### 기술 스택
- **Go + Bubbletea** (ADR-006 참조)

### v0.1-alpha (핵심 관찰 가능성)
- Agent Arena + Live Timeline + Redaction
- 완료 기준: UX S1(실시간 관제) 시나리오 pass

### v0.1-beta (분석 기능)
- Task Graph (트리 목록) + Inspector (기본) + Replay (기본)
- 완료 기준: UX S1~S3, F1~F3, R1 시나리오 pass

### v0.2 (분석 품질 강화)
- Intent/Action diff 고도화
- FR-7 Cost/Token Analytics (히트맵)
- blocker chain 시각화 개선, 필터 UX 강화
- 알림(선택)

### v0.3 (운영 확장)
- 플러그인형 provider 확장
- 운영 리포트 자동 요약
- 접근성/테마/설정 고도화

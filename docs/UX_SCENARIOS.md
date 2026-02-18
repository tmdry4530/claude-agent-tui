# UX Scenarios

## 1) 정상 시나리오

### S1. 실시간 관제
- 상황: OMC Ultrawork로 병렬 작업 진행
- 사용자 행동: Agent Arena 확인 -> Timeline 필터
- 기대 결과: 어떤 에이전트가 실행/대기/블록인지 즉시 파악

### S2. 병목 분석
- 상황: 결과가 늦어짐
- 사용자 행동: Task Graph에서 blocker chain 확인
- 기대 결과: 어떤 태스크/에이전트가 병목인지 1분 내 특정

### S3. 검증 루프 추적
- 상황: Ralph 모드에서 verify-fix 반복
- 사용자 행동: Inspector에서 verify/fix history 확인
- 기대 결과: 실패 원인과 재시도 흐름 파악

---

## 2) 실패 시나리오

### F1. 파서 실패
- 증상: 일부 이벤트 미표시
- UX 요구: 경고 표시 + 앱 지속 동작 + 실패 카운트 노출

### F2. 이벤트 폭주
- 증상: 렌더 지연
- UX 요구: 자동 축약/샘플링 + 사용자 안내 배지

### F3. provider 혼합 혼선
- 증상: Claude/Gemini/Codex 이벤트 해석 혼재
- UX 요구: provider 필터와 mode 배지로 분리 가시화

---

## 3) Replay 시나리오

### R1. 장애 재현
- 목적: 특정 시점 error 재현
- 행동: replay load -> n(다음 error) -> inspector 확인
- 결과: 원인/전파 경로 확인

### R2. 개선 전후 비교
- 목적: 운영 변경 전후 비교
- 행동: 두 run replay 비교
- 결과: 병목/오류율/latency 변화 확인

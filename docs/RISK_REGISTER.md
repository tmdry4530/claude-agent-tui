# Risk Register

| ID | 리스크 | 영향도 | 가능성 | 대응 전략 | 상태 |
|---|---|---:|---:|---|---|
| R-01 | OMC 이벤트 포맷 변경 | 높음 | 중간 | Normalizer 버전 분리, unknown fallback | Open |
| R-02 | 로그 누락/지연 | 중간 | 중간 | 다중 소스 수집 + timestamp 보정 | Open |
| R-03 | 이벤트 폭주로 TUI 프레임 저하 | 높음 | 중간 | 샘플링/축약 모드, 부분 렌더 | Open |
| R-04 | 민감정보 노출 | 매우높음 | 낮음 | redaction 기본 ON, 검증 체크리스트 | Open |
| R-05 | 모드별 상태 해석 불일치 | 중간 | 중간 | mode별 규칙 문서화, DECISIONS 고정 | Open |
| R-06 | replay 재현 불일치 | 중간 | 중간 | canonical schema + deterministic clock | Open |
| R-07 | 팀 내 용어 불일치 | 낮음 | 중간 | glossary/ADR(DECISIONS.md) 운영 | Open |
| R-08 | 장시간 세션 메모리 누적 | 높음 | 중간 | Store GC 정책, ring buffer 크기 제한, RSS 모니터링 | Open |
| R-09 | 키바인딩 모드 간 충돌 | 낮음 | 중간 | 포커스 기반 우선순위 규칙 확정 | Open |
| R-10 | Role 매핑 확장성 | 낮음 | 중간 | v0.2에서 role 확장 검토 트리거 정의 | Open |

## 운영 규칙
- 새 리스크 발견 시 즉시 추가
- High 이상은 NEXT_ACTIONS에 대응 태스크 등록
- 상태: Open / Mitigating / Closed

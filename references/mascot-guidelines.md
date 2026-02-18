# Mascot Guidelines (CLCO)

## 1) 목적
OMC Agent TUI의 시각적 아이덴티티를 통일하고, 에이전트 상태를 빠르게 인지할 수 있도록 CLCO 마스코트 사용 규칙을 정의한다.

---

## 2) 자산
- 베이스 이미지: `assets/mascot/clco-mascot-base.jpg`
- 팔레트 정의: `assets/mascot/palette.md`

---

## 3) 기본 원칙
1. 마스코트 형태(실루엣)는 고정한다.
2. 에이전트 구분은 **색상/배지/라벨**로만 수행한다.
3. 상태 전달은 색상 단독이 아니라 **테두리/아이콘/모션**을 함께 사용한다.
4. 작은 터미널에서도 식별 가능해야 한다(80col fallback).

---

## 4) 에이전트 역할별 스타일

| 역할 | 기본 색상 | 배지 | ASCII fallback |
|---|---|---|---|
| Planner | `#7AA2F7` | `🧭` | `[P]` |
| Executor | `#9ECE6A` | `🛠` | `[X]` |
| Reviewer | `#BB9AF7` | `🔍` | `[R]` |
| Guard | `#F7768E` | `🛡` | `[G]` |
| Tester | `#E0AF68` | `🧪` | `[T]` |
| Writer | `#73DACA` | `✍️` | `[W]` |
| Explorer | `#58A6FF` | `🔎` | `[E]` |
| Architect | `#FFA657` | `📐` | `[A]` |
| Debugger | `#FF7B72` | `🐛` | `[D]` |
| Verifier | `#56D364` | `✅` | `[V]` |
| Designer | `#D2A8FF` | `🎨` | `[S]` |
| Custom | `#8A93A5` | `❓` | `[?]` |

> 동일 역할의 다중 인스턴스는 톤 변형(밝기 ±10~15%)으로 구분.
> event-schema.md 섹션 11 Role 매핑 테이블 참조.

---

## 5) 상태별 시각 규칙

| 상태 | 시각 표현 | 접근성 대체 |
|---|---|---|
| RUNNING | 밝은 색 + 약한 pulse | `▶` |
| WAITING | 저채도/밝기 감소 | `…` |
| BLOCKED | 노란 외곽선 강조 (`#E3B341`) | `!` |
| ERROR | 빨간 외곽선 + 짧은 shake (`#FF7B72`) | `✖` |
| DONE | 체크 뱃지 표시 (`#56D364`) | `✓` |
| FAILED | 빨간 외곽선 + X 뱃지, dim 처리 (`#FF7B72`) | `X` |
| CANCELLED | 회색 외곽선 + `-` 뱃지 (`#6E7681`) | `-` |
| IDLE | 기본색, 모션 없음 | `•` |

---

## 6) 렌더링 레벨

### L1: Minimal (저사양/원격 SSH)
- 텍스트 + 배지 + 상태 아이콘
- 모션 없음

### L2: Standard (기본)
- 색상 + 테두리 + 간단 pulse

### L3: Enhanced (고급 터미널)
- truecolor + 상태 애니메이션 + 강조 효과

---

## 7) 터미널 호환성
- TrueColor 지원 감지 실패 시 256-color 팔레트로 자동 강등
- 256-color 미지원 시 monochrome + 상태 아이콘 모드로 강등
- 폰트/이모지 미지원 시 ASCII fallback:
  - Planner `[P]`, Executor `[X]`, Reviewer `[R]`, Guard `[G]`
  - Tester `[T]`, Writer `[W]`, Explorer `[E]`, Architect `[A]`
  - Debugger `[D]`, Verifier `[V]`, Designer `[S]`, Custom `[?]`

---

## 8) 금지 규칙
- 마스코트 형태 자체 변형 금지(브랜딩 일관성 유지)
- 역할 색상 임의 변경 금지(팔레트 합의 필요)
- ERROR 상태를 색상만으로 표현 금지(아이콘 필수)
- 과도한 애니메이션 금지(가독성 우선)

---

## 9) 접근성
- 색약 대응 팔레트 반드시 제공
- 상태 텍스트를 항상 병기(`RUNNING`, `BLOCKED` 등)
- 깜빡임 빈도 제한(광과민 배려)

---

## 10) 운영/변경 정책
- 스타일 변경은 `DECISIONS.md`에 기록
- 변경 시 Before/After 스크린샷 첨부
- 릴리즈 노트에 영향 범위(가독성/성능) 명시

---

## 11) 샘플 표기

### 표준 카드(텍스트 표현)
`🟩 [🛠 Coder-Auth] RUNNING 78%`

### ASCII fallback
`[C] Coder-Auth  RUNNING  78%  ▶️`

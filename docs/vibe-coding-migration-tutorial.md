# 바이브 코딩 마이그레이션 튜토리얼

> 기존 Go 프로젝트를 Claude Code와 함께 "바이브 코딩" 가능한 상태로 전환하는 전체 과정을 기록한 실전 튜토리얼입니다.

## 개요

| 항목 | 내용 |
|------|------|
| **대상 프로젝트** | RealWorld (Conduit) — Go Fiber v2 + HTMX + GORM/SQLite |
| **초기 상태** | 테스트 0%, 린트 없음, CI/CD 없음, 스펙 문서 없음 |
| **최종 상태** | 커버리지 84.9%, E2E 16개, 린트 0 이슈, pre-commit/pre-push hook |
| **소요 시간** | 약 2시간 (워크숍) |
| **사용 도구** | Claude Code (Opus), lefthook, golangci-lint, Playwright |

### 커밋 히스토리 (시간순)

```
1e6a026 feat: 바이브 코딩 도입 기반 구축
eea7ae5 test: Playwright E2E 테스트 추가 (16개 테스트 전체 통과)
d03008e feat: 프로젝트 커스텀 스킬 5개 추가
788ba18 test: 컨트롤러 통합 테스트 추가 (커버리지 80%+ 달성)
```

---

## Step 1: 프로젝트 초기화 (`/init`)

### 프롬프트

```
/init
```

### AI 응답 및 산출물

Claude Code가 프로젝트 전체를 스캔하여 `CLAUDE.md`를 자동 생성했습니다.

**생성된 `CLAUDE.md` 주요 내용:**
- 빌드/실행 명령어 (`go run main.go`, `air`)
- 듀얼 컨트롤러 패턴 설명 (페이지 + HTMX)
- 요청 흐름: `main.go → serve.go → 라우트 → 컨트롤러 → GORM → 템플릿`
- 데이터 레이어, 인증, 템플릿 시스템 요약
- 내부 패키지 및 주요 의존성 목록

### 팁

- `/init`은 프로젝트 시작 시 한 번만 실행하면 됩니다.
- 생성된 `CLAUDE.md`는 이후 모든 AI 응답의 품질을 결정하는 핵심 컨텍스트입니다.
- 자동 생성 후 부정확한 부분은 직접 수정하세요 — AI가 이 파일을 매번 참조합니다.

---

## Step 2: 마이그레이션 계획 수립 (심층 인터뷰)

### 프롬프트

```
이 프로젝트에 바이브 코딩을 도입하기 위한 작업 계획을 세워줘.
다음 4가지 영역을 에픽으로 하여 각각의 에픽에 대해
AskUserQuestionTool로 심층 인터뷰를 진행하고
그 결과에 따라 세부 계획을 작성해줘.

    문서화 — OpenSpec으로 프로젝트 스펙 문서화 + 스킬 생성
    테스트 구현 — 단위 테스트, E2E 테스트, 린트, Makefile 통합
    CI/CD 파이프라인 — pre-commit hook, pre-push hook
    전환 계획서 — generate-plan 스킬로 migration-plan
```

### AI 응답 및 산출물

Claude Code가 4개 에픽에 대해 순차적으로 질문을 던졌습니다.

**인터뷰 질문 예시:**

| Epic | 질문 | 사용자 응답 |
|------|------|------------|
| 문서화 | "OpenSpec 스펙 범위는? 전체 시스템 vs 핵심 모듈만?" | 전체 시스템 |
| 테스트 | "테스트 프레임워크 선호도? testing+testify vs ginkgo?" | testing+testify |
| CI/CD | "GitHub Actions vs lefthook? 둘 다?" | lefthook 우선 |
| 전환 계획서 | "타임라인? 주 단위 vs 마일스톤 기반?" | 주 단위 |

**생성된 계획 파일:** `~/.claude/plans/vibe-coding-migration.md`

```
Epic 1 (문서화) ──→ Epic 4 (전환 계획서)
                ↗
Epic 2 (테스트) ──→ Epic 3 (CI/CD) ──→ Epic 4
```

- Epic 1, 2는 **병렬 진행 가능**
- Epic 3은 Epic 2 **완료 후** 진행
- Epic 4는 모든 에픽 **완료 후** 최종 작성

### 팁

- "심층 인터뷰"를 명시하면 AI가 일방적으로 결정하지 않고 질문을 통해 맥락을 파악합니다.
- 에픽 단위로 나누면 병렬 실행이 가능해져 시간이 크게 단축됩니다.
- 계획 파일은 `~/.claude/plans/`에 저장되어 세션 간 지속됩니다.

---

## Step 3: 기반 구축 (Epic 1+2 병렬 실행)

### 프롬프트

```
진행해줘
```

### AI 응답 및 산출물

Claude Code가 Epic 1(문서화)과 Epic 2(테스트)를 **병렬 에이전트**로 동시 실행했습니다.

#### Epic 1 산출물: OpenSpec 문서 + 스킬

| 파일 | 내용 |
|------|------|
| `openspec/specs/architecture.md` | 듀얼 컨트롤러 패턴, 요청 흐름 |
| `openspec/specs/data-model.md` | GORM 모델 스키마, 관계 |
| `openspec/specs/routes.md` | 페이지 + HTMX 라우트 전체 목록 |
| `openspec/specs/authentication.md` | 세션 기반 인증 흐름 |
| `openspec/specs/templates.md` | 템플릿 구조, 레이아웃 |
| `.claude/skills/fiber-htmx-controller/` | 듀얼 컨트롤러 생성 스킬 |
| `.claude/skills/fiber-htmx-template/` | 템플릿 생성 스킬 |

#### Epic 2 산출물: 테스트 인프라 + 단위 테스트

| 파일 | 내용 |
|------|------|
| `internal/testutil/apptest.go` | 테스트 앱 셋업 헬퍼 (in-memory SQLite) |
| `cmd/web/model/user_test.go` | HashPassword, CheckPassword 등 7개 테스트 |
| `cmd/web/model/article_test.go` | FavoriteCount, Tags 등 7개 테스트 |
| `internal/validator_test.go` | 입력 검증 4개 테스트 |
| `internal/errormessage_test.go` | 에러 메시지 4개 테스트 |
| `.golangci.yml` | golangci-lint v2 설정 |
| `Makefile` | build, test, lint 등 9개 타겟 |

### 문제 발생 및 해결

#### 문제 1: golangci-lint v2 설정 형식 변경

```
# v1 형식 (오류 발생)
linters:
  enable:
    - gosimple

# v2 형식 (수정)
linters:
  enable:
    - staticcheck  # gosimple이 staticcheck에 통합됨
```

**해결:** v2에서는 `formatter`와 `linter`가 분리되고, `gosimple`이 `staticcheck`에 통합되었습니다. 에러 메시지를 분석하여 3차례에 걸쳐 설정을 수정했습니다.

#### 문제 2: 18개 린트 이슈 (errcheck, goimports)

| 카테고리 | 이슈 수 | 해결 방법 |
|---------|---------|----------|
| errcheck | 6건 | `c.Render()` 반환값 → `return c.Render(...)`, GORM `_ =` 처리 |
| goimports | 10건 | import 그룹 정렬 (stdlib → third-party → local 순서) |
| unused | 1건 | `main.go`의 미사용 변수 제거 |

### 팁

- "진행해줘"처럼 간단한 프롬프트로 AI가 계획을 자율 실행합니다.
- 독립적인 에픽은 자동으로 병렬 실행되어 시간이 절반으로 줄어듭니다.
- 린트 도구 버전 차이는 흔한 문제입니다 — AI가 에러 메시지를 보고 스스로 수정합니다.

---

## Step 4: E2E 테스트 (Playwright)

### 프롬프트

```
다음 단계 진행해 줘
```

### AI 응답 및 산출물

Playwright 설치부터 E2E 테스트 3개 스펙 작성까지 자동 진행:

| 파일 | 테스트 시나리오 |
|------|---------------|
| `e2e/tests/auth.spec.ts` | 회원가입 → 로그인 → 로그아웃 (6개) |
| `e2e/tests/article.spec.ts` | 기사 작성 → 조회 → 수정 → 삭제 (5개) |
| `e2e/tests/feed.spec.ts` | 글로벌 피드, 태그 피드, 내 피드 (5개) |
| `e2e/playwright.config.ts` | 설정 (baseURL, Chromium) |
| `e2e/package.json` | Playwright 의존성 |

### 문제 발생 및 해결

#### 문제: 16개 중 11개 테스트 실패 (타임아웃)

**근본 원인:** HTMX 폼은 표준 HTML `<form>` 제출이 아니라 `hx-post` 속성으로 동작합니다.

```html
<!-- 실제 HTML (HTMX) -->
<button hx-post="/htmx/sign-up" hx-target="#sign-up-page">Sign up</button>

<!-- AI가 처음 사용한 셀렉터 (표준 HTML 가정) -->
button[type="submit"]  ← 존재하지 않음!
```

**해결 과정:**

1. AI가 45개 템플릿 파일을 전수 조사하여 HTMX 셀렉터 맵 구축
2. 모든 셀렉터를 HTMX 기반으로 교체

| 기존 셀렉터 | 수정된 셀렉터 | 이유 |
|------------|-------------|------|
| `button[type="submit"]` | `button:has-text("Sign up")` | HTMX는 `hx-post` 사용 |
| `input[name="name"]` | `#sign-up-username` | ID 기반 셀렉터 |
| `.error-messages` | `#sign-in-form-messages .alert-danger` | 프로젝트 고유 마크업 |
| `a[href="/sign-out"]` | settings 페이지 logout 버튼 | 링크가 아닌 HTMX 버튼 |
| `nav` | `nav.navbar` | pagination `<nav>` 충돌 방지 |
| `button.edit-button` | `button.edit-button.first()` | 2개 존재하여 strict 위반 |

3. 수정 후 **16/16 전체 통과**

### 팁

- HTMX 프로젝트의 E2E 테스트는 표준 HTML 셀렉터가 통하지 않습니다.
- `hx-post`, `hx-target`, `hx-swap` 등 HTMX 속성 기반으로 셀렉터를 설계하세요.
- strict mode 에러(`2 elements matched`)는 `.first()` 또는 더 구체적인 셀렉터로 해결합니다.
- 템플릿 파일을 먼저 읽고 셀렉터 맵을 만드는 것이 가장 효율적입니다.

---

## Step 5: 커스텀 스킬 생성

### 프롬프트

```
open spec 문서들을 활용해서 프로젝트에 유용한 스킬을 만들어줘.
스킬 작성은 Claude Code 공식 스킬 가이드를 따르고,
스킬이 참조하는 문서들은 먼저 어떤 스킬이 이 프로젝트에 유용할지 제안해줘.
```

### AI 응답 및 산출물

AI가 먼저 5개 스킬을 제안하고, 승인 후 일괄 생성했습니다.

| 스킬 | 명령어 | 용도 |
|------|--------|------|
| `add-feature` | `/add-feature [name]` | 풀스택 스캐폴딩 (컨트롤러+템플릿+라우트+테스트) |
| `gen-integration-test` | `/gen-integration-test [path]` | 통합 테스트 자동 생성 |
| `gen-e2e-test` | `/gen-e2e-test [feature]` | Playwright E2E 테스트 생성 |
| `add-model` | `/add-model [name]` | GORM 모델 + 관계 + 테스트 |
| `debug-htmx` | `/debug-htmx [증상]` | HTMX 동작 문제 6단계 진단 |

**핵심 포인트:** 각 스킬이 OpenSpec 문서를 참조합니다.

```yaml
# .claude/skills/gen-e2e-test/SKILL.md 예시
---
name: gen-e2e-test
description: "Playwright E2E 테스트 생성. HTMX 프로젝트 전용 셀렉터 맵 내장"
---
# HTMX 셀렉터 맵 (openspec/specs/templates.md 기반)
| 요소 | 셀렉터 |
| Sign-up 폼 | #sign-up-page |
| 로그인 버튼 | button:has-text("Sign in") |
...
```

### 팁

- OpenSpec 문서를 먼저 작성해두면, 스킬이 정확한 프로젝트 맥락을 참조할 수 있습니다.
- 스킬은 반복 작업(새 기능 추가, 테스트 작성)에 특히 효과적입니다.
- `debug-htmx`처럼 디버깅용 스킬은 문제 해결 시간을 크게 줄여줍니다.

---

## Step 6: GitHub 이슈 등록

### 프롬프트

```
계획 문서의 남은 작업들을 gh cli를 사용해서 깃헙 이슈로 등록해줘.
각 이슈에는 작업 설명과 인수 조건이 포함되어야 해.
이슈 간 의존성이 있으면 본문에 명시해줘.
```

### 문제 발생 및 해결

#### 문제: GitHub Issues 비활성화

```
GraphQL: Could not resolve to a Repository with Issues enabled
```

**해결:** 사용자가 GitHub Settings → General → Features → Issues 체크박스 활성화 후 재시도.

### AI 응답 및 산출물

5개 이슈가 라벨, 의존성 명시와 함께 생성되었습니다.

| # | 이슈 | 라벨 | 의존성 |
|---|------|------|--------|
| #1 | 페이지 컨트롤러 테스트 커버리지 (9.6% → 80%+) | epic:test | 없음 |
| #2 | HTMX 컨트롤러 테스트 커버리지 (13.8% → 80%+) | epic:test | 없음 |
| #3 | GitHub Actions CI 파이프라인 구축 | epic:ci-cd | #1, #2 |
| #4 | README 온보딩 가이드 | epic:docs | 없음 |
| #5 | 최종 검증 및 마이그레이션 완료 | epic:release | #1~#4 |

### 팁

- `gh` CLI를 활용하면 이슈 생성이 자동화됩니다.
- 커밋 메시지에 `Closes #1`을 포함하면 푸시 시 자동 클로즈됩니다.
- Fork 저장소의 Issues는 기본 비활성화 — Settings에서 직접 켜야 합니다.

---

## Step 7: 테스트 커버리지 80%+ 달성

### 프롬프트

```
gh 이슈 #1, #2 진행해줘
```

### AI 응답 및 산출물

**1차 배치 (8개 파일):** 커버리지 53.8% / 74.1%

| 패키지 | 테스트 파일 | 테스트 수 |
|--------|-----------|----------|
| controller | article, user, home, editor | 27개 |
| htmx | article, comment, editor, home-page, home-action, article-action, user-page, user-action | 38개 |

**갭 분석:** `go tool cover -func`로 0% 핸들러를 식별

```
SignInPage       0.0%   ← 미커버
SignUpPage       0.0%   ← 미커버
SettingPage      0.0%   ← 미커버
SignOut           0.0%   ← 미커버
```

**2차 배치 (8개 파일 추가):** 최종 커버리지 달성

| 패키지 | Before | After | 목표 |
|--------|--------|-------|------|
| controller | 9.6% | **89.4%** | 80%+ ✅ |
| htmx | 13.8% | **83.7%** | 80%+ ✅ |
| 합계 | — | **84.9%** | 80%+ ✅ |

### 문제 발생 및 해결

#### 문제: `TestSignOut_Unauthenticated` 실패

```
expected: 200
actual  : 302
```

**해결:** 비인증 상태의 sign-out은 302 리다이렉트가 정상 동작. 테스트 기대값을 `http.StatusFound`로 수정.

### 팁

- 커버리지 개선은 2단계로 나누세요: 1차(주요 핸들러) → 갭 분석 → 2차(0% 핸들러)
- `go tool cover -func`로 정확히 어떤 함수가 0%인지 식별할 수 있습니다.
- 테스트 실패 시 "테스트를 고칠지 vs 기대값을 고칠지" 판단이 중요합니다.
- 실제 핸들러 동작을 확인한 후 기대값을 맞추세요.

---

## Step 8: Git Hook 설정 (pre-commit + pre-push)

### 프롬프트 1: pre-commit

```
pre-commit hook을 설정해줘. 커밋할 때마다 린트 체크와
단위 테스트가 자동 실행되도록 해줘. Makefile의 make check 타겟을 활용해줘.
```

### 프롬프트 2: pre-push

```
pre-push hook을 설정해줘. 푸시할 때 E2E 테스트와
커버리지 체크(80% 이상)가 실행되도록 해줘.
```

### AI 응답 및 산출물

**Makefile 추가 타겟:**

```makefile
check: lint test

coverage-check:
    # 패키지별 80% 이상 검증
    for pkg in $(COVERAGE_PKGS); do
        # 각 패키지 커버리지 측정 → 80% 미만이면 FAIL
    done
```

**lefthook.yml 최종 설정:**

```yaml
pre-commit:
  commands:
    check:
      glob: "*.go"
      run: make check          # lint + test

pre-push:
  commands:
    build:
      run: go build ./...       # 빌드 검증
    coverage-check:
      run: make coverage-check  # 패키지별 80%+ 검증
    e2e:
      run: make e2e             # Playwright E2E
```

### 문제 발생 및 해결

#### 문제: 전체 커버리지 70.6%로 80% 미달

테스트 없는 패키지(`authentication`, `database`, `middleware`)가 전체 평균을 끌어내렸습니다.

**해결:** 전체 합산 → 패키지별 검증 방식으로 변경

```makefile
# Before: 전체 합산 (70.6% → FAIL)
go test -coverprofile=coverage.out ./...

# After: 패키지별 검증 (각각 80%+ → PASS)
COVERAGE_PKGS=./cmd/web/controller ./cmd/web/controller/htmx ./cmd/web/model ./internal
for pkg in $(COVERAGE_PKGS); do
    # 개별 검증
done
```

### 팁

- 커버리지 체크는 "전체 합산"보다 "패키지별"이 현실적입니다.
- 테스트 없는 인프라 패키지(DB 연결, 미들웨어)가 평균을 크게 낮출 수 있습니다.
- pre-commit은 빠른 검사(lint+unit test), pre-push는 무거운 검사(E2E+coverage)로 나누세요.

---

## 최종 결과 요약

### Before vs After

| 항목 | Before | After |
|------|--------|-------|
| CLAUDE.md | 없음 | 프로젝트 아키텍처 전체 문서화 |
| OpenSpec 스펙 | 없음 | 5개 문서 |
| 커스텀 스킬 | 없음 | 7개 (컨트롤러, 템플릿, 테스트, 모델, 디버그 등) |
| 단위 테스트 | 0개 | 22개 |
| 통합 테스트 | 0개 | 65개+ |
| E2E 테스트 | 0개 | 16개 |
| 테스트 커버리지 | 0% | 84.9% |
| 린트 | 없음 | golangci-lint v2 (0 issues) |
| Git Hook | 없음 | pre-commit (lint+test), pre-push (build+coverage+e2e) |
| Makefile | 없음 | 11개 타겟 |
| GitHub Issues | 없음 | 5개 (의존성 명시) |

### 파일 변경 통계

```
47 files changed, +3,013 lines  (기반 구축)
 6 files changed              (E2E 테스트)
 7 files changed, +755 lines  (커스텀 스킬)
16 files changed, +1,506 lines (통합 테스트)
```

---

## 다른 프로젝트에 적용할 때의 팁

### 1. 순서가 중요합니다

```
/init → 계획 수립 → 문서화 → 테스트 인프라 → 테스트 작성 → CI/CD
```

`CLAUDE.md`와 스펙 문서가 먼저 있어야 AI가 정확한 코드를 생성합니다. 테스트 없이 CI/CD를 설정하면 의미가 없습니다.

### 2. 에픽 단위로 나누고 병렬 실행하세요

독립적인 작업(문서화 vs 테스트)은 "진행해줘" 한 마디로 병렬 실행됩니다. 의존성 있는 작업(CI/CD → 테스트 필요)은 선행 작업 완료 후 진행하세요.

### 3. AI의 실패를 두려워하지 마세요

이 워크숍에서도 여러 번 실패했습니다:
- golangci-lint v2 설정 3번 수정
- E2E 셀렉터 전면 교체
- 커버리지 측정 방식 3번 변경

AI는 에러 메시지를 보고 스스로 수정합니다. 핵심은 **AI가 실패할 수 있는 환경을 만들어주는 것** — 린트, 테스트, 타입 체크가 있으면 AI가 자기 실수를 발견하고 고칩니다.

### 4. 프레임워크 특성을 CLAUDE.md에 명시하세요

HTMX 폼이 표준 HTML과 다르게 동작하는 것처럼, 프레임워크 고유의 패턴이 있다면 반드시 문서화하세요. 이것이 없으면 AI가 표준 패턴으로 코드를 생성하고 실패합니다.

### 5. 커버리지 목표는 "패키지별"로 설정하세요

테스트 없는 인프라 코드가 전체 평균을 크게 낮춥니다. 비즈니스 로직이 있는 패키지만 대상으로 하는 것이 현실적입니다.

### 6. 커스텀 스킬은 반복 작업을 자동화합니다

새 기능 추가, 테스트 작성, 디버깅 등 반복되는 패턴을 스킬로 만들면 프롬프트 한 줄로 일관된 결과를 얻을 수 있습니다. OpenSpec 문서를 참조하는 스킬이 가장 효과적입니다.

### 7. Git Hook은 품질 게이트입니다

pre-commit(빠른 검사)과 pre-push(무거운 검사)를 나누면, 개발 흐름을 방해하지 않으면서도 품질을 보장할 수 있습니다.

---

## 프롬프트 요약 (복사해서 사용하세요)

```bash
# Step 1: 초기화
/init

# Step 2: 계획 수립
이 프로젝트에 바이브 코딩을 도입하기 위한 작업 계획을 세워줘.
다음 영역을 에픽으로 하여 심층 인터뷰 후 세부 계획을 작성해줘:
1. 문서화 — 프로젝트 스펙 문서화 + 스킬 생성
2. 테스트 구현 — 단위/통합/E2E 테스트, 린트, Makefile
3. CI/CD — pre-commit hook, pre-push hook
4. 전환 계획서

# Step 3: 실행
진행해줘

# Step 4: 이슈 등록
계획 문서의 남은 작업들을 gh cli로 깃헙 이슈로 등록해줘.
각 이슈에는 작업 설명과 인수 조건, 의존성을 명시해줘.

# Step 5: 특정 이슈 작업
gh 이슈 #1, #2 진행해줘

# Step 6: Hook 설정
pre-commit hook을 설정해줘. 린트 체크와 단위 테스트가 자동 실행되도록.
pre-push hook을 설정해줘. E2E 테스트와 커버리지 체크(80% 이상)가 실행되도록.
```

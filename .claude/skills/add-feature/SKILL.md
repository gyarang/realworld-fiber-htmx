---
name: add-feature
description: 새 기능 풀스택 스캐폴딩. Page controller + HTMX controller + 템플릿 + 라우트 등록 + 테스트 스텁을 일괄 생성. "새 기능 추가", "feature 추가", "페이지 추가" 시 사용.
argument-hint: [feature-name]
user-invocable: true
---

# 새 기능 풀스택 스캐폴딩

`$ARGUMENTS` 기능을 위한 전체 파일을 생성합니다.

## 생성 파일 목록

1. `cmd/web/controller/$0.go` — Page controller
2. `cmd/web/controller/htmx/$0.go` — HTMX controller
3. `cmd/web/controller/htmx/$0-action.go` — HTMX action controller (POST/PATCH/DELETE)
4. `cmd/web/templates/$0/index.tmpl` — 전체 페이지 템플릿
5. `cmd/web/templates/$0/htmx-index.tmpl` — HTMX fragment 템플릿
6. `cmd/web/templates/$0/partials/` — Partial 템플릿
7. `cmd/web/controller/$0_test.go` — Page controller 테스트 스텁
8. `cmd/web/controller/htmx/$0_test.go` — HTMX controller 테스트 스텁

## 실행 순서

### Step 1: 기존 코드 확인

라우트 파일을 읽어 중복 여부 확인:
- `cmd/web/route/handlers.go`
- `cmd/web/route/htmx-handlers.go`

### Step 2: Page Controller 생성

```go
package controller

import (
	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"

	"github.com/gofiber/fiber/v2"
)

func FeaturePage(c *fiber.Ctx) error {
	isAuthenticated, userID := authentication.AuthGet(c)
	db := database.Get()
	// 데이터 조회
	return c.Render("feature/index", fiber.Map{
		"IsAuthenticated": isAuthenticated,
	}, "layouts/app")
}
```

### Step 3: HTMX Controller 생성

GET handler — fragment 렌더링:
```go
func FeaturePage(c *fiber.Ctx) error {
	isAuthenticated, userID := authentication.AuthGet(c)
	db := database.Get()
	return c.Render("feature/htmx-index", fiber.Map{
		"IsAuthenticated": isAuthenticated,
	}, "layouts/app-htmx")
}
```

POST handler — action 처리:
```go
func FeatureAction(c *fiber.Ctx) error {
	isAuthenticated, userID := authentication.AuthGet(c)
	if !isAuthenticated {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	// 비즈니스 로직
	// 성공: helper.HTMXRedirectTo(c, "/target")
	// 실패: 에러 fragment 반환
	return nil
}
```

### Step 4: 템플릿 생성

전체 페이지 (`layouts/app`):
- HTMX `hx-get="/htmx/feature"` + `hx-trigger="load"` 로 동적 영역 로드
- `hx-target="#feature-content"` 으로 교체 대상 지정

HTMX fragment (`layouts/app-htmx`):
- HTML 구조 없이 콘텐츠만
- `{{ if .IsAuthenticated }}` 조건부 렌더링
- 폼: `hx-post`, `hx-target="#form-errors"` 패턴
- OOB swap 필요 시: `hx-swap-oob="true"` 추가

### Step 5: 라우트 등록

`cmd/web/route/handlers.go`에 추가:
```go
app.Get("/feature", controller.FeaturePage)
```

`cmd/web/route/htmx-handlers.go`에 추가:
```go
app.Get("/htmx/feature", htmx.FeaturePage)
app.Post("/htmx/feature", htmx.FeatureAction)
```

### Step 6: 테스트 스텁 생성

```go
func TestFeaturePage(t *testing.T) {
	app, db := testutil.SetupTestApp(t)
	// TODO: 테스트 케이스 구현
}
```

## 규칙

- Page controller와 HTMX controller는 **항상 쌍으로** 생성
- HTMX 라우트는 반드시 `/htmx/` prefix
- import 순서: stdlib → third-party → local (빈 줄로 구분)
- 인증 확인: `authentication.AuthGet(c)`
- DB 접근: `database.Get()`
- 리다이렉트: `helper.HTMXRedirectTo(c, path)`
- 에러 반환: `return c.Render()` (errcheck 통과)

상세 아키텍처는 [architecture spec](../../../openspec/specs/architecture.md), 라우트 목록은 [routes spec](../../../openspec/specs/routes.md) 참조.

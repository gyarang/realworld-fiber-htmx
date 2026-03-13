---
name: fiber-htmx-controller
description: Fiber + HTMX 듀얼 컨트롤러 생성 패턴. 새 기능 추가 시 page controller와 HTMX controller를 함께 생성.
---

# Fiber + HTMX Controller 생성 스킬

새 기능(feature)을 추가할 때 이 스킬의 패턴을 따릅니다.

## 생성 순서

### 1. Page Controller (`cmd/web/controller/{feature}.go`)
```go
package controller

import (
    "realworld-fiber-htmx/internal/authentication"
    "realworld-fiber-htmx/internal/database"
    "github.com/gofiber/fiber/v2"
)

func FeaturePage(c *fiber.Ctx) error {
    isAuthenticated, userID := authentication.AuthGet(c)
    db := database.Get()

    // 데이터 조회 로직

    return c.Render("{feature}/index", fiber.Map{
        "IsAuthenticated": isAuthenticated,
        // ... data
    }, "layouts/app")
}
```

### 2. HTMX Controller (`cmd/web/controller/htmx/{feature}.go`)
```go
package htmx

import (
    "realworld-fiber-htmx/internal/authentication"
    "realworld-fiber-htmx/internal/database"
    "github.com/gofiber/fiber/v2"
)

// GET — fragment 렌더링
func FeaturePage(c *fiber.Ctx) error {
    isAuthenticated, userID := authentication.AuthGet(c)
    db := database.Get()

    return c.Render("{feature}/htmx-index", fiber.Map{
        "IsAuthenticated": isAuthenticated,
    }, "layouts/app-htmx")
}

// POST — 액션 처리
func FeatureAction(c *fiber.Ctx) error {
    isAuthenticated, userID := authentication.AuthGet(c)
    if !isAuthenticated {
        return c.SendStatus(fiber.StatusUnauthorized)
    }

    // 비즈니스 로직
    // 성공 시: helper.HTMXRedirectTo(c, "/target-path")
    // 실패 시: 에러 fragment 반환
    return nil
}
```

### 3. 라우트 등록

`cmd/web/route/handlers.go`에 추가:
```go
app.Get("/{feature}", controller.FeaturePage)
```

`cmd/web/route/htmx-handlers.go`에 추가:
```go
app.Get("/htmx/{feature}", htmx.FeaturePage)
app.Post("/htmx/{feature}", htmx.FeatureAction)
```

### 4. 템플릿 생성

`cmd/web/templates/{feature}/` 디렉토리에:
- `index.tmpl` — 전체 페이지용 (layouts/app)
- `htmx-index.tmpl` — HTMX fragment용 (layouts/app-htmx)
- `partials/` — 재사용 가능한 partial 템플릿

## 규칙
- Page controller와 HTMX controller는 항상 쌍으로 생성
- HTMX 라우트는 반드시 `/htmx/` prefix 사용
- 인증이 필요한 액션은 `AuthGet()` 으로 확인
- 리다이렉트는 `helper.HTMXRedirectTo()` 사용
- 데이터베이스는 `database.Get()` 으로 접근

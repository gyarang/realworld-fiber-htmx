# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

RealWorld (Conduit) blogging application built with Go Fiber v2, HTMX, GORM + SQLite, and Bootstrap CSS.

## Build & Run

```bash
# Run the application (listens on localhost:8181)
go run main.go

# Hot-reload development (requires Air: go install github.com/air-verse/air@latest)
air

# Build
go build -o tmp/main.exe main.go
```

### Makefile

```bash
make build          # go build -o bin/conduit main.go
make run            # go run main.go
make dev            # air (hot-reload)
make test           # go test ./... -v
make test-coverage  # 커버리지 리포트 생성 (coverage.html)
make lint           # golangci-lint run
make fmt            # gofmt + goimports
make clean          # 빌드 아티팩트 정리
make e2e            # Playwright E2E 테스트
```

### Testing

```bash
# 전체 테스트
go test ./... -v

# 특정 패키지 테스트
go test ./cmd/web/model/... -v

# 특정 테스트 함수
go test ./cmd/web/model/... -run TestCheckPassword -v

# 커버리지
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

- 테스트 프레임워크: `testing` + `testify/assert`
- 테스트 헬퍼: `internal/testutil/` (in-memory SQLite DB, fixture 생성)
- 린트: `golangci-lint` (`.golangci.yml` 설정)

## Architecture

### Dual Controller Pattern

The app uses two parallel controller layers:
- **Page controllers** (`cmd/web/controller/`) — render full HTML pages via `c.Render()`
- **HTMX controllers** (`cmd/web/controller/htmx/`) — return HTML fragments for dynamic partial updates

Routes are split accordingly:
- `cmd/web/route/handlers.go` — traditional page routes (GET `/`, `/articles/:slug`, etc.)
- `cmd/web/route/htmx-handlers.go` — HTMX API routes under `/htmx/` prefix (GET/POST/PATCH)

### Request Flow

`main.go` → `cmd/web/serve.go` (registers routes) → controllers → GORM queries → template rendering

### Data Layer

- **GORM** with SQLite (`conduit.sqlite` file, pre-seeded)
- Global `database.DB` variable used throughout (`internal/database/database.go`)
- Models in `cmd/web/model/` with relationship preloading
- Join tables: `article_favorite`, `article_tag`, `user_follower` (auto-created by GORM)

### Authentication

- Session-based using Fiber session middleware with SQLite3 storage
- `internal/authentication/session.go` — `AuthStore()`, `AuthGet()`, `AuthDestroy()`
- Template function `IsAuthenticated()` for conditional rendering

### Templates

- 45 `.tmpl` files in `cmd/web/templates/`
- Two layouts: `app.tmpl` (traditional) and `app-htmx.tmpl` (HTMX with `{{ embed }}`)
- Custom template functions: `IsAuthenticated()`, `Iterate()`, `Dict()` (defined in `internal/renderer/renderer.go`)

### Key Internal Packages

| Package | Purpose |
|---------|---------|
| `internal/database` | GORM SQLite setup with connection pooling |
| `internal/authentication` | Session management |
| `internal/renderer` | Template engine with custom functions |
| `internal/validator` | Wrapper around go-playground/validator |
| `internal/errormessage` | Error message formatting |
| `internal/helper` | `HTMXRedirectTo()` for HTMX redirects |
| `internal/middleware` | 404 handler |

## Key Dependencies

- `github.com/gofiber/fiber/v2` — web framework
- `gorm.io/gorm` + `github.com/glebarez/sqlite` — ORM + SQLite
- `github.com/gofiber/template/html/v2` — template engine
- `github.com/gosimple/slug` — URL slug generation
- `golang.org/x/crypto` — bcrypt password hashing
- `github.com/go-playground/validator/v10` — input validation

## Conventions

- No `.env` file — server address (`localhost:8181`) and DB path (`conduit.sqlite`) are hardcoded
- HTMX routes mirror page routes under `/htmx/` prefix
- Password hashing via `model.User.HashPassword()` / `CheckPassword()`
- Slugs generated from article titles for SEO-friendly URLs

## OpenSpec Workflow

프로젝트 스펙 문서는 `openspec/specs/`에 위치. 변경 제안은 OpenSpec 워크플로우를 따름:

1. `/opsx explore` — 아이디어 탐색, 코드 조사 (읽기만, 수정 없음)
2. `/opsx propose` — 변경 제안서 생성 (proposal → design → tasks)
3. `/opsx apply` — 태스크별 구현
4. `/opsx archive` — 완료된 변경 아카이브

## Custom Skills

| 스킬 | 명령어 | 용도 |
|------|--------|------|
| `fiber-htmx-controller` | `/fiber-htmx-controller` | 듀얼 컨트롤러(page+HTMX) 생성 패턴 |
| `fiber-htmx-template` | `/fiber-htmx-template` | 템플릿 생성 패턴 (레이아웃, HTMX 속성, partial 분리) |
| `add-feature` | `/add-feature [name]` | 풀스택 스캐폴딩 (컨트롤러 쌍 + 템플릿 + 라우트 + 테스트 스텁) |
| `add-model` | `/add-model [name]` | GORM 모델 + 관계 설정 + 단위 테스트 생성 |
| `gen-integration-test` | `/gen-integration-test [path]` | 통합 테스트 자동 생성 (SetupTestApp 패턴) |
| `gen-e2e-test` | `/gen-e2e-test [feature]` | Playwright E2E 테스트 생성 (HTMX 셀렉터 맵 내장) |
| `debug-htmx` | `/debug-htmx [증상]` | HTMX 문제 진단 (라우트→컨트롤러→템플릿 체인 검증) |

## Git Hooks (lefthook)

```bash
lefthook install    # hook 설치 (최초 1회)
lefthook run pre-commit   # 수동 실행: gofmt + go vet + golangci-lint
lefthook run pre-push     # 수동 실행: go build + go test
```

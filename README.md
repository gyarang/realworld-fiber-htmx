# RealWorld — Fiber v2 + HTMX

[![CI](https://github.com/gyarang/realworld-fiber-htmx/actions/workflows/ci.yml/badge.svg)](https://github.com/gyarang/realworld-fiber-htmx/actions/workflows/ci.yml)

"Conduit" 소셜 블로깅 사이트 ([Medium.com](https://medium.com) 클론) — [RealWorld](https://github.com/gothinkster/realworld) 스펙 준수.

**기술 스택:**
- [Fiber v2](https://github.com/gofiber/fiber/tree/v2.52.6) — Express 스타일 Go 웹 프레임워크
- [HTMX](https://htmx.org/) — HTML 기반 동적 프론트엔드
- [GORM](https://gorm.io/) + [SQLite](https://github.com/glebarez/sqlite) — ORM + 데이터베이스
- [Slug](https://github.com/gosimple/slug) — SEO 친화적 URL 생성

## 빠른 시작

```bash
git clone https://github.com/gyarang/realworld-fiber-htmx.git
cd realworld-fiber-htmx
go mod download
go run main.go
# → http://localhost:8181
```

기본 계정: `test@email.com` / `secret` (또는 웹에서 회원가입)

## 개발 환경 설정

### 필수 도구

| 도구 | 설치 |
|------|------|
| Go 1.21+ | [golang.org/dl](https://golang.org/dl/) |
| golangci-lint | `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` |
| lefthook | `go install github.com/evilmartians/lefthook@latest` |
| Air (핫 리로드) | `go install github.com/air-verse/air@latest` |

### E2E 테스트 (선택)

| 도구 | 설치 |
|------|------|
| Node.js 20+ | [nodejs.org](https://nodejs.org/) |
| Playwright | `cd e2e && npm install && npx playwright install chromium` |

### 초기 설정

```bash
go mod download          # Go 의존성 설치
lefthook install         # Git hook 활성화
```

### 개발 서버 실행

```bash
make dev                 # Air 핫 리로드 (권장)
# 또는
make run                 # go run main.go
```

## Makefile 명령어

| 명령어 | 설명 |
|--------|------|
| `make build` | 바이너리 빌드 (`bin/conduit`) |
| `make run` | 서버 실행 |
| `make dev` | Air 핫 리로드 |
| `make test` | 전체 테스트 실행 |
| `make test-coverage` | 커버리지 리포트 생성 (`coverage.html`) |
| `make lint` | golangci-lint 실행 |
| `make fmt` | gofmt + goimports |
| `make check` | lint + test (pre-commit hook) |
| `make coverage-check` | 패키지별 80%+ 커버리지 검증 |
| `make e2e` | Playwright E2E 테스트 |
| `make clean` | 빌드 아티팩트 정리 |

## 프로젝트 구조

```
.
├── main.go                          # 엔트리포인트
├── cmd/web/
│   ├── serve.go                     # 라우트 등록
│   ├── controller/                  # 페이지 컨트롤러 (전체 HTML)
│   ├── controller/htmx/             # HTMX 컨트롤러 (HTML fragment)
│   ├── model/                       # GORM 모델
│   ├── route/                       # 라우트 정의
│   └── templates/                   # Go 템플릿 (45개)
├── internal/
│   ├── authentication/              # 세션 관리
│   ├── database/                    # GORM + SQLite 설정
│   ├── renderer/                    # 템플릿 엔진
│   ├── validator/                   # 입력 검증
│   └── testutil/                    # 테스트 헬퍼
├── e2e/                             # Playwright E2E 테스트
├── openspec/specs/                  # 프로젝트 스펙 문서
└── .claude/skills/                  # Claude Code 커스텀 스킬
```

### 듀얼 컨트롤러 패턴

이 프로젝트는 두 가지 컨트롤러 레이어를 사용합니다:

- **페이지 컨트롤러** (`controller/`) — 전체 HTML 페이지 렌더링 (`c.Render()`)
- **HTMX 컨트롤러** (`controller/htmx/`) — HTML fragment 반환 (동적 부분 업데이트)

```
GET /                    → controller.HomePage()         (전체 페이지)
GET /htmx/global-feed    → htmx.HomeGlobalFeed()         (피드 fragment)
POST /htmx/sign-in       → htmx.SignInAction()           (로그인 처리)
```

## Git Hook

[lefthook](https://github.com/evilmartians/lefthook)으로 자동 품질 검사:

| Hook | 실행 내용 |
|------|----------|
| **pre-commit** | `make check` — lint + 전체 테스트 |
| **pre-push** | build + 커버리지 80%+ 검증 + E2E 테스트 |

```bash
lefthook install              # 최초 1회
lefthook run pre-commit       # 수동 실행
lefthook run pre-push         # 수동 실행
```

## CI/CD

GitHub Actions (`push`/`PR` to `main`):

| Job | 내용 |
|-----|------|
| build | `go build ./...` |
| test | 전체 테스트 + 패키지별 80%+ 커버리지 체크 |
| lint | golangci-lint |
| e2e | Playwright E2E (서버 빌드→기동→테스트) |

## 테스트 현황

| 패키지 | 커버리지 |
|--------|---------|
| `cmd/web/controller` | 89.4% |
| `cmd/web/controller/htmx` | 83.7% |
| `cmd/web/model` | 85.2% |
| `internal` | 100.0% |

E2E 테스트: 16개 (인증, 기사 CRUD, 피드)

## 라이선스

MIT

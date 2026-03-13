# Architecture Spec

## Overview
RealWorld (Conduit) 블로깅 애플리케이션. Go Fiber v2 + HTMX 기반 서버 사이드 렌더링.

## Request Flow
```
Client → Fiber Router → Controller → GORM → SQLite
                                  → Template Engine → HTML Response
```

## Dual Controller Pattern
이 프로젝트의 핵심 아키텍처 패턴. 모든 기능은 두 개의 컨트롤러 레이어로 구현:

### Page Controllers (`cmd/web/controller/`)
- 전체 HTML 페이지 렌더링 (`c.Render()`)
- 라우트: `cmd/web/route/handlers.go`
- 초기 페이지 로드 시 사용

### HTMX Controllers (`cmd/web/controller/htmx/`)
- HTML fragment 반환 (partial update)
- 라우트: `cmd/web/route/htmx-handlers.go`
- `/htmx/` prefix 사용
- hx-get, hx-post, hx-target 속성으로 호출

## Application Bootstrap
1. `main.go` — renderer 초기화, Fiber 앱 생성
2. `database.Open()` — SQLite 연결, GORM 설정
3. `authentication.SessionStart()` — 세션 스토리지 초기화
4. Middleware 등록 (recover, logger)
5. `web.Serve(app)` — 라우트 등록, static 파일 서빙
6. `app.Listen("localhost:8181")`

## Static Assets
`cmd/web/static/` 에서 서빙:
- CSS: Bootstrap 기반 Conduit 테마, Tagify 스타일
- JS: HTMX 라이브러리 (v1.9.6), Head Support 확장, Tagify

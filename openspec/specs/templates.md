# Templates Spec

## 개요
Go HTML template engine (Fiber HTML engine) 사용. 45개 `.tmpl` 파일.

## Directory Structure
```
cmd/web/templates/
├── layouts/
│   ├── app.tmpl          # 전체 페이지 레이아웃 (traditional)
│   └── app-htmx.tmpl     # HTMX partial 레이아웃 ({{ embed }})
├── home/                  # 홈/피드 관련
├── articles/              # 기사 상세/목록
├── editor/                # 기사 작성/수정
├── sign-in/               # 로그인
├── sign-up/               # 회원가입
├── users/                 # 유저 프로필
├── settings/              # 설정
└── components/            # 공유 UI 컴포넌트
```

## Layouts

### `app.tmpl` (Traditional)
전체 HTML 페이지. `<html>`, `<head>`, `<body>` 포함.
Page controller에서 `c.Render("template-name", data, "layouts/app")` 으로 사용.

### `app-htmx.tmpl` (HTMX)
HTMX partial update용. `{{ embed }}` 으로 컨텐츠 삽입.
HTMX controller에서 `c.Render("template-name", data, "layouts/app-htmx")` 으로 사용.

## Custom Template Functions (`internal/renderer/renderer.go`)

### `IsAuthenticated`
현재 요청의 인증 상태를 반환. 네비게이션, 버튼 표시 조건에 사용.

### `Iterate(count int)`
0부터 count-1까지의 slice 반환. 페이지네이션 렌더링에 사용.

### `Dict(values ...interface{})`
key-value 쌍을 map으로 변환. 템플릿에 여러 변수를 전달할 때 사용.

## HTMX Attributes 패턴
- `hx-get="/htmx/..."` — GET 요청으로 fragment 로드
- `hx-post="/htmx/..."` — POST 요청으로 액션 실행
- `hx-target="#element-id"` — 응답을 삽입할 대상
- `hx-trigger="click"` — 이벤트 트리거
- `hx-push-url="true"` — URL 히스토리 업데이트
- `hx-swap="innerHTML"` — 교체 방식

## HTMXRedirectTo Helper (`internal/helper/route.go`)
HTMX 응답에서 리다이렉트가 필요할 때 사용:
- `HX-Replace-Url` 헤더 설정
- `HX-Reswap: innerHTML` 설정
- redirect component 템플릿 렌더링 (`#app-body` 타겟)

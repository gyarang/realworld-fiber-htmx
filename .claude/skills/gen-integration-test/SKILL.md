---
name: gen-integration-test
description: Go Fiber 통합 테스트 자동 생성. SetupTestApp 패턴으로 HTTP 핸들러 테스트 생성. "통합 테스트 생성", "컨트롤러 테스트", "integration test" 시 사용.
argument-hint: [controller-file-path]
user-invocable: true
---

# 통합 테스트 생성

`$ARGUMENTS` 컨트롤러의 통합 테스트를 생성합니다.

## 사전 확인

1. 대상 컨트롤러 파일을 읽어 핸들러 함수 목록 추출
2. 해당 핸들러의 라우트를 `cmd/web/route/handlers.go` 또는 `htmx-handlers.go`에서 확인
3. 기존 테스트 파일이 있으면 읽어서 누락된 케이스만 추가

## 테스트 구조

```go
package controller_test  // 또는 htmx_test

import (
	"io"
	"net/http"
	"testing"

	"realworld-fiber-htmx/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestHandlerName(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	t.Run("비인증 사용자", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/path", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		defer func() { _ = resp.Body.Close() }()
		assert.Contains(t, string(body), "expected-content")
	})

	t.Run("인증된 사용자", func(t *testing.T) {
		user := testutil.CreateTestUser(t, db)
		cookie := testutil.AuthenticateUser(t, app, user)
		req, _ := http.NewRequest("GET", "/path", nil)
		req.Header.Set("Cookie", cookie)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		// assertions...
	})
}
```

## 핸들러 유형별 테스트 케이스

### GET (페이지/fragment 렌더링)
- 비인증 사용자 — 200 OK + 페이지 콘텐츠 확인
- 인증된 사용자 — 200 OK + 인증 전용 UI 요소 확인
- 존재하지 않는 리소스 — 404 또는 리다이렉트
- 파라미터 검증 — `:slug`, `:username` 등

### POST (폼 제출/액션)
- 비인증 시 — 401 Unauthorized 또는 리다이렉트
- 빈 필드 — 에러 메시지 fragment 반환
- 유효하지 않은 입력 — 검증 에러 확인
- 성공 케이스 — DB 변경 확인 + 리다이렉트/fragment 응답

### PATCH (수정)
- 권한 없는 사용자 — 403 또는 리다이렉트
- 유효한 수정 — DB 업데이트 확인

### DELETE (삭제)
- 권한 없는 사용자 — 403 또는 리다이렉트
- 성공 삭제 — DB에서 제거 확인

## POST 요청 본문 작성

```go
// Form data
body := strings.NewReader("email=test@example.com&password=password123")
req, _ := http.NewRequest("POST", "/htmx/sign-in", body)
req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
```

## 테스트 헬퍼 (`internal/testutil/`)

| 함수 | 용도 |
|------|------|
| `SetupTestApp(t)` | Fiber 앱 + in-memory SQLite DB 생성 |
| `SetupTestDB(t)` | DB만 생성 (단위 테스트용) |
| `CreateTestUser(t, db)` | 테스트 유저 생성 (비밀번호: "password") |
| `CreateTestArticle(t, db, userID)` | 테스트 기사 생성 |
| `CreateTestComment(t, db, articleID, userID, body)` | 테스트 댓글 생성 |
| `AuthenticateUser(t, app, user)` | 세션 쿠키 반환 |

## 규칙

- 테스트 파일명: `{feature}_test.go` (하이픈 → 언더스코어)
- `resp.Body.Close()`는 반드시 `defer func() { _ = resp.Body.Close() }()`
- import 순서: stdlib → third-party → local
- 각 테스트는 독립적 (서로 의존하지 않음)
- 에러 체크: `assert.NoError(t, err)` 사용
- DB 상태 검증 시 GORM 쿼리 직접 실행

라우트 전체 목록은 [routes spec](../../../openspec/specs/routes.md), 인증 플로우는 [authentication spec](../../../openspec/specs/authentication.md) 참조.

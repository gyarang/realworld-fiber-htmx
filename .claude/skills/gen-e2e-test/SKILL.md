---
name: gen-e2e-test
description: Playwright E2E 테스트 생성. HTMX 프로젝트 전용 셀렉터 맵 내장. "E2E 테스트 생성", "Playwright 테스트", "e2e test" 시 사용.
argument-hint: [feature-or-flow-name]
user-invocable: true
---

# Playwright E2E 테스트 생성

`$ARGUMENTS` 기능/플로우의 E2E 테스트를 생성합니다.

## 테스트 파일 위치

`e2e/tests/$0.spec.ts`

## 기본 구조

```typescript
import { test, expect, Page } from "@playwright/test";

// 인증이 필요한 테스트에서 재사용
async function signUp(page: Page): Promise<string> {
  const uniqueId = Date.now();
  await page.goto("/sign-up");
  await page.fill("#sign-up-username", `testuser${uniqueId}`);
  await page.fill("#sign-up-email", `test${uniqueId}@example.com`);
  await page.fill("#sign-up-password", "password123");
  await page.click('button:has-text("Sign up")');
  await page.waitForURL("/", { timeout: 10_000 });
  return `testuser${uniqueId}`;
}

test.describe("Feature Name", () => {
  test("should do something", async ({ page }) => {
    // ...
  });
});
```

## HTMX 셀렉터 맵 (필수 참조)

이 프로젝트는 모든 폼이 HTMX `hx-post`/`hx-patch`를 사용합니다.
`button[type="submit"]`는 동작하지 않으며, 반드시 아래 셀렉터를 사용하세요.

### 폼 셀렉터

| 페이지 | 폼 셀렉터 | 제출 버튼 |
|--------|-----------|-----------|
| 회원가입 | `form[hx-post="/htmx/sign-up"]` | `button:has-text("Sign up")` |
| 로그인 | `form[hx-post="/htmx/sign-in"]` | `button:has-text("Sign in")` |
| 에디터 (새 글) | `form[hx-post="/htmx/editor"]` | `button:has-text("Publish Article")` |
| 에디터 (수정) | `form[hx-patch="/htmx/editor/{slug}"]` | `button:has-text("Publish Article")` |
| 설정 | `form#settings-form` | `button:has-text("Update Settings")` |
| 댓글 | `form#article-comment-form` | `button:has-text("Post Comment")` |

### 입력 필드 셀렉터

| 페이지 | 필드 | 셀렉터 |
|--------|------|--------|
| 회원가입 | Username | `#sign-up-username` |
| 회원가입 | Email | `#sign-up-email` |
| 회원가입 | Password | `#sign-up-password` |
| 로그인 | Email | `#sign-in-email` |
| 로그인 | Password | `#sign-in-password` |
| 에디터 | Title | `input[name="title"]` |
| 에디터 | Description | `input[name="description"]` |
| 에디터 | Content | `textarea[name="content"]` |
| 댓글 | Body | `textarea[name="comment"]` |

### 에러 메시지 컨테이너

| 페이지 | 셀렉터 |
|--------|--------|
| 회원가입 | `#sign-up-form-messages .alert-danger` |
| 로그인 | `#sign-in-form-messages .alert-danger` |
| 에디터 | `#form-message .alert-danger` |
| 설정 | `#settings-form-messages .alert-danger` |

### UI 요소 셀렉터

| 요소 | 셀렉터 | 주의사항 |
|------|--------|----------|
| 네비게이션 바 | `nav.navbar` | `nav`는 2개 존재 (pagination 포함) |
| 기사 제목 | `#article-detail__title` | |
| 기사 본문 | `.post-content` | |
| 태그 목록 | `#popular-tag-list` | HTMX로 비동기 로드 |
| 태그 링크 | `#popular-tag-list a.label-pill` | |
| 피드 탭 | `#feed-navigation` | |
| 수정 버튼 | `button.edit-button` | 2개 존재, `.first()` 필요 |
| 삭제 버튼 | `button.delete-button` | 2개 존재, `.first()` 필요 |
| 로그아웃 | `button:has-text("Or click here to logout")` | `/settings` 페이지에 위치 |

## HTMX 대기 패턴

```typescript
// HTMX 리다이렉트 후 대기 (HX-Redirect)
await page.click('button:has-text("Sign up")');
await page.waitForURL("/", { timeout: 10_000 });

// HTMX fragment 로드 대기
await expect(page.locator("#popular-tag-list")).toBeVisible({ timeout: 10_000 });

// HTMX 에러 메시지 대기
await expect(
  page.locator("#sign-in-form-messages .alert-danger")
).toBeVisible({ timeout: 10_000 });
```

## 테스트 시나리오 패턴

### 페이지 렌더링
```typescript
test("should display page", async ({ page }) => {
  await page.goto("/path");
  await expect(page.locator("h1")).toContainText("Expected Title");
});
```

### 인증 필요 플로우
```typescript
test("should require authentication", async ({ page }) => {
  await signUp(page);
  await page.goto("/protected-path");
  await expect(page.locator("...")).toBeVisible();
});
```

### 폼 제출 + 에러
```typescript
test("should show error on invalid input", async ({ page }) => {
  await page.goto("/form-page");
  await page.click('button:has-text("Submit")');
  await expect(
    page.locator("#form-messages .alert-danger")
  ).toBeVisible({ timeout: 10_000 });
});
```

## 규칙

- `Date.now()`로 고유 ID 생성 (테스트 격리)
- `waitForTimeout()` 대신 `waitForURL()` 또는 `toBeVisible()` 사용
- `page.locator("nav")` → `page.locator("nav.navbar")` (strict mode)
- 복수 요소 셀렉터는 `.first()` 추가
- 타임아웃: `{ timeout: 10_000 }` 명시

라우트 전체 목록은 [routes spec](../../../openspec/specs/routes.md), 템플릿 구조는 [templates spec](../../../openspec/specs/templates.md) 참조.

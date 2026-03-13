---
name: debug-htmx
description: HTMX 동작 문제 진단. 라우트→컨트롤러→템플릿 체인 검증, OOB swap 디버깅, fragment 렌더링 문제 해결. "HTMX 안됨", "fragment 안나옴", "swap 안됨" 시 사용.
argument-hint: [증상-또는-URL]
user-invocable: true
---

# HTMX 문제 진단

`$ARGUMENTS` 관련 HTMX 동작 문제를 체계적으로 진단합니다.

## 진단 체크리스트

아래 순서대로 확인합니다. 각 단계에서 문제가 발견되면 즉시 수정을 제안합니다.

### Step 1: 라우트 확인

문제가 발생하는 URL의 HTMX 라우트가 등록되어 있는지 확인:

```
1. cmd/web/route/htmx-handlers.go 에서 해당 경로 검색
2. HTTP Method 일치 확인 (GET vs POST vs PATCH vs DELETE)
3. 경로 파라미터 확인 (:slug, :username 등)
```

**흔한 문제:**
- 라우트 미등록
- Method 불일치 (hx-get인데 POST 라우트만 있음)
- 경로 오타 (`/htmx/` prefix 누락)

### Step 2: 컨트롤러 확인

해당 라우트의 핸들러 함수가 올바른 응답을 반환하는지 확인:

```
1. 핸들러 함수 위치 확인 (cmd/web/controller/htmx/)
2. c.Render() 호출의 템플릿 경로 확인
3. 레이아웃: "layouts/app-htmx" 사용 여부 확인
4. 에러 시 return 누락 확인
```

**흔한 문제:**
- `"layouts/app"` 사용 (전체 페이지 반환 → fragment여야 함)
- `c.Render()` 반환값 미반환 (`return c.Render(...)` 필요)
- 인증 체크 후 early return 누락

### Step 3: 템플릿 확인

HTMX fragment 템플릿이 올바른지 확인:

```
1. 템플릿 파일 존재 여부 (cmd/web/templates/)
2. 템플릿 이름과 c.Render() 첫 번째 인자 일치
3. 변수명 일치 (.Title vs .ArticleTitle 등)
4. {{ template "partial" }} 경로 확인
```

**흔한 문제:**
- 템플릿 파일명과 Render 인자 불일치
- 전달하지 않은 변수 참조 → 빈 출력
- partial 경로 오류

### Step 4: HTMX 속성 확인

HTML에서 HTMX 속성이 올바른지 확인:

```
1. hx-get/hx-post 경로에 /htmx/ prefix 있는지
2. hx-target이 존재하는 DOM 요소를 가리키는지
3. hx-swap 모드 확인 (innerHTML, outerHTML, afterbegin 등)
4. hx-trigger 이벤트 확인
5. hx-push-url 설정 확인
```

**흔한 문제:**
- `hx-target="#id"` 인데 해당 `id`가 DOM에 없음
- `hx-swap="outerHTML"` 인데 교체 후 타겟이 사라져서 후속 요청 실패
- `hx-push-url` 설정 누락으로 브라우저 URL 불일치

### Step 5: OOB Swap 확인

Out-of-Band swap 사용 시 추가 확인:

```
1. 응답 HTML에 hx-swap-oob="true" 포함 여부
2. OOB 대상 id가 현재 DOM에 존재하는지
3. OOB 응답이 메인 응답과 함께 반환되는지
```

이 프로젝트의 OOB swap 사용처:
- `#feed-navigation` — 피드 탭 전환
- `#navbar` — 네비게이션 바 업데이트
- `#popular-tag-list` — 태그 목록 업데이트

### Step 6: 리다이렉트 확인

HTMX에서 리다이렉트가 동작하지 않을 때:

```
1. helper.HTMXRedirectTo() 사용 여부
2. HX-Redirect 또는 HX-Replace-Url 헤더 설정 확인
3. HTMX redirect component 템플릿 존재 확인
```

**주의:** 일반 `c.Redirect()` 는 HTMX 요청에서 동작하지 않습니다.
반드시 `helper.HTMXRedirectTo(c, "/path")` 사용.

## 빠른 진단 명령어

```bash
# 라우트 목록 확인
grep -n "htmx" cmd/web/route/htmx-handlers.go

# 특정 경로의 핸들러 찾기
grep -rn "func.*FeatureName" cmd/web/controller/htmx/

# 템플릿에서 hx-* 속성 검색
grep -rn "hx-get\|hx-post\|hx-target" cmd/web/templates/feature/

# HTMX 응답 직접 테스트
curl -s http://localhost:8181/htmx/path | head -20
```

## 참고 문서

- 아키텍처: [architecture spec](../../../openspec/specs/architecture.md)
- 라우트 전체 목록: [routes spec](../../../openspec/specs/routes.md)
- 템플릿 구조: [templates spec](../../../openspec/specs/templates.md)
- 인증 플로우: [authentication spec](../../../openspec/specs/authentication.md)

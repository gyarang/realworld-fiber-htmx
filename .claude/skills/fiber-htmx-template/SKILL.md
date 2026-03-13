---
name: fiber-htmx-template
description: Fiber + HTMX 템플릿 생성 패턴. 레이아웃 적용, HTMX 속성 활용, partial 분리.
---

# Fiber + HTMX Template 생성 스킬

새 템플릿을 추가할 때 이 스킬의 패턴을 따릅니다.

## 전체 페이지 템플릿 (`layouts/app` 용)

```html
<div class="container page">
    <div class="row">
        <div class="col-md-9">
            <!-- 메인 콘텐츠 -->
            <h1>{{ .Title }}</h1>

            <!-- HTMX로 동적 로드할 영역 -->
            <div id="feature-content"
                 hx-get="/htmx/{feature}"
                 hx-trigger="load"
                 hx-target="#feature-content">
                Loading...
            </div>
        </div>
    </div>
</div>
```

## HTMX Fragment 템플릿 (`layouts/app-htmx` 용)

```html
<!-- Fragment: 전체 HTML 구조 없이 콘텐츠만 -->
<div class="feature-item">
    {{ range .Items }}
    <div class="item-preview">
        <h2>{{ .Title }}</h2>
        <p>{{ .Description }}</p>

        <!-- 인증된 사용자만 보이는 버튼 -->
        {{ if .IsAuthenticated }}
        <button hx-post="/htmx/{feature}/action"
                hx-target="#feature-content"
                hx-swap="innerHTML">
            Action
        </button>
        {{ end }}
    </div>
    {{ end }}
</div>
```

## HTMX 속성 활용 패턴

### 목록 + 페이지네이션
```html
<div id="item-list">
    {{ range .Items }}
        <!-- 아이템 렌더링 -->
    {{ end }}
</div>

<!-- 페이지네이션 -->
<nav>
    {{ range Iterate .TotalPages }}
    <a hx-get="/htmx/{feature}?page={{ . }}"
       hx-target="#item-list"
       hx-push-url="true">{{ . }}</a>
    {{ end }}
</nav>
```

### 폼 제출
```html
<form hx-post="/htmx/{feature}"
      hx-target="#form-errors"
      hx-swap="innerHTML">
    <input type="text" name="title" value="{{ .Title }}">
    <div id="form-errors"></div>
    <button type="submit">Submit</button>
</form>
```

### 탭 전환 (피드 패턴)
```html
<div class="feed-toggle">
    <ul class="nav">
        <li class="nav-item">
            <a hx-get="/htmx/{feature}/tab1"
               hx-target="#tab-content"
               hx-push-url="/tab1">Tab 1</a>
        </li>
        <li class="nav-item">
            <a hx-get="/htmx/{feature}/tab2"
               hx-target="#tab-content"
               hx-push-url="/tab2">Tab 2</a>
        </li>
    </ul>
</div>
<div id="tab-content">
    <!-- 탭 콘텐츠가 여기에 로드됨 -->
</div>
```

## 규칙
- 전체 페이지: `layouts/app` 레이아웃 사용
- HTMX fragment: `layouts/app-htmx` 레이아웃 사용
- partial은 `{feature}/partials/` 에 분리
- HTMX 요청 대상은 항상 `/htmx/` prefix
- `{{ if .IsAuthenticated }}` 으로 인증 조건부 렌더링
- Dict 함수로 partial에 여러 변수 전달: `{{ template "partial" (Dict "key1" .Val1 "key2" .Val2) }}`

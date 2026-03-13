# Authentication Spec

## 개요
세션 기반 인증. Fiber session middleware + SQLite3 스토리지 사용.

## Components

### Session Storage (`internal/authentication/session.go`)
- `StoredAuthenticationSession` — 전역 세션 스토어
- SQLite3 테이블: `fiber_storage`

### Functions

#### `SessionStart()`
앱 시작 시 1회 호출. SQLite3 기반 세션 스토어 초기화.

#### `AuthStore(c *fiber.Ctx, userID uint)`
로그인 성공 시 호출. 세션에 `authentication` 키로 userID 저장.

#### `AuthGet(c *fiber.Ctx) (bool, uint)`
인증 상태 확인. 반환값:
- `(true, userID)` — 인증됨
- `(false, 0)` — 미인증

#### `AuthDestroy(c *fiber.Ctx)`
로그아웃 시 호출. 세션에서 `authentication` 키 삭제.

## Authentication Flow

### 로그인
1. POST `/htmx/sign-in` → `SignInAction`
2. Email로 User 조회 (GORM)
3. `user.CheckPassword(plain)` — bcrypt 비교
4. 성공: `AuthStore(c, user.ID)` → HTMX redirect to `/`
5. 실패: 에러 메시지 fragment 반환

### 회원가입
1. POST `/htmx/sign-up` → `SignUpAction`
2. 입력 검증 (validator)
3. `user.HashPassword()` — bcrypt 해시
4. GORM으로 User 생성
5. `AuthStore(c, user.ID)` → HTMX redirect to `/`

### 로그아웃
1. POST `/htmx/sign-out` → `SignOut`
2. `AuthDestroy(c)` → HTMX redirect to `/`

## Template Integration
`IsAuthenticated()` 템플릿 함수로 조건부 렌더링:
- 네비게이션 바: 로그인/회원가입 vs 설정/프로필
- 기사 상세: 좋아요/팔로우 버튼 표시 여부
- 에디터: 인증된 사용자만 접근

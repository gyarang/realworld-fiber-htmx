# Routes Spec

## Page Routes (`cmd/web/route/handlers.go`)
전체 HTML 페이지를 렌더링하는 GET 라우트.

| Method | Path | Handler | 설명 |
|--------|------|---------|------|
| GET | `/sign-in` | SignInPage | 로그인 페이지 |
| GET | `/sign-up` | SignUpPage | 회원가입 페이지 |
| GET | `/` | HomePage | 글로벌 피드 (메인) |
| GET | `/your-feed` | YourFeedPage | 내 피드 |
| GET | `/tag-feed/:slug` | TagFeedPage | 태그별 피드 |
| GET | `/articles/:slug` | ArticleDetailPage | 기사 상세 |
| GET | `/editor/:slug?` | EditorPage | 에디터 (새 글/수정) |
| GET | `/users/:username` | UserDetailPage | 유저 프로필 |
| GET | `/users/:username/articles` | UserDetailPage | 유저 기사 목록 |
| GET | `/users/:username/favorites` | UserDetailFavoritePage | 유저 좋아요 목록 |
| GET | `/settings` | SettingPage | 설정 페이지 |

## HTMX Routes (`cmd/web/route/htmx-handlers.go`)
HTML fragment를 반환하는 HTMX API 라우트. 모든 경로에 `/htmx/` prefix.

### Authentication
| Method | Path | Handler | 설명 |
|--------|------|---------|------|
| GET | `/htmx/sign-in` | SignInPage | 로그인 폼 fragment |
| POST | `/htmx/sign-in` | SignInAction | 로그인 처리 |
| POST | `/htmx/sign-out` | SignOut | 로그아웃 |
| GET | `/htmx/sign-up` | SignUpPage | 회원가입 폼 fragment |
| POST | `/htmx/sign-up` | SignUpAction | 회원가입 처리 |

### Home Feed
| Method | Path | Handler | 설명 |
|--------|------|---------|------|
| GET | `/htmx/home` | HomePage | 홈 fragment |
| GET | `/htmx/home/your-feed` | HomeYourFeed | 내 피드 목록 |
| GET | `/htmx/home/global-feed` | HomeGlobalFeed | 글로벌 피드 목록 |
| GET | `/htmx/home/tag-feed/:tag` | HomeTagFeed | 태그 피드 목록 |
| GET | `/htmx/home/tag-list` | HomeTagList | 태그 리스트 |
| POST | `/htmx/home/articles/:slug/favorite` | HomeFavoriteAction | 홈에서 좋아요 |

### Articles
| Method | Path | Handler | 설명 |
|--------|------|---------|------|
| GET | `/htmx/articles/:slug` | ArticleDetailPage | 기사 상세 fragment |
| GET | `/htmx/articles/:slug/comments` | ArticleDetailCommentList | 댓글 목록 |
| POST | `/htmx/articles/:slug/comments` | ArticleComment | 댓글 작성 |
| POST | `/htmx/articles/:slug/favorite` | ArticleFavoriteAction | 좋아요 토글 |
| POST | `/htmx/articles/follow-user/:slug` | ArticleFollowAction | 팔로우 토글 |

### Editor
| Method | Path | Handler | 설명 |
|--------|------|---------|------|
| GET | `/htmx/editor/:slug?` | EditorPage | 에디터 fragment |
| POST | `/htmx/editor` | StoreArticle | 새 기사 저장 |
| PATCH | `/htmx/editor/:slug?` | UpdateArticle | 기사 수정 |

### Users
| Method | Path | Handler | 설명 |
|--------|------|---------|------|
| GET | `/htmx/users/:username` | UserDetailPage | 유저 프로필 fragment |
| GET | `/htmx/users/:username/articles` | UserArticles | 유저 기사 목록 |
| GET | `/htmx/users/:username/favorites` | UserArticlesFavorite | 유저 좋아요 목록 |
| POST | `/htmx/users/articles/:slug/favorite` | UserArticleFavoriteAction | 좋아요 토글 |
| POST | `/htmx/users/:username/follow` | UserFollowAction | 팔로우 토글 |

### Settings
| Method | Path | Handler | 설명 |
|--------|------|---------|------|
| GET | `/htmx/settings` | SettingPage | 설정 fragment |
| POST | `/htmx/settings` | SettingAction | 설정 업데이트 |

## Static Assets
| Path | Source |
|------|--------|
| `/static/*` | `./cmd/web/static/` |

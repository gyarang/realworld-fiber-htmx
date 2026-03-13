# Data Model Spec

## Entity Relationship

```
User ──1:N──→ Article ──1:N──→ Comment
  │                │               │
  │                └──M:N──→ Tag   │
  │                │               │
  │                └──M:N──→ User (favorites)
  │                                │
  └──M:N (follow)──→ User         └──→ User (author)
```

## Models

### User
| Field | Type | Constraints |
|-------|------|------------|
| ID | uint | PK, auto-increment |
| Name | string | required |
| Username | string | unique index |
| Email | string | unique index, not null, email 형식 |
| Password | string | not null, bcrypt 해시 |
| Bio | string | optional |
| Image | string | optional (프로필 이미지 URL) |

Relationships:
- `Followers []Follow` — FK: FollowerID (이 유저를 팔로우하는 사람들)
- `Followings []Follow` — FK: FollowingID (이 유저가 팔로우하는 사람들)
- `Favorites []Article` — M:N via `article_favorite`

Methods:
- `HashPassword()` — bcrypt.DefaultCost로 비밀번호 해시
- `CheckPassword(plain)` — bcrypt 비교
- `FollowedBy(id)` — 특정 유저가 팔로우 중인지 확인
- `FollowersCount()` — 팔로워 수

### Follow (Join Table: `user_follower`)
| Field | Type | Column |
|-------|------|--------|
| FollowerID | uint | user_id (PK) |
| FollowingID | uint | follower_id (PK) |

### Article
| Field | Type | Constraints |
|-------|------|------------|
| Slug | string | unique index, not null |
| Title | string | not null, required |
| Description | string | required |
| Body | string | required |
| UserID | uint | FK → User |
| IsFavorited | bool | 비영속 필드 (gorm:"-") |

Relationships:
- `User` — belongs to
- `Comments []Comment` — has many
- `Favorites []User` — M:N via `article_favorite`
- `Tags []Tag` — M:N via `article_tag`

Methods:
- `GetFormattedCreatedAt()` — "Jan 02, 2006" 형식
- `GetFavoriteCount()` — 좋아요 수
- `FavoritedBy(id)` — 특정 유저가 좋아요 했는지
- `GetTagsAsCommaSeparated()` — 태그를 콤마로 연결

### Comment
| Field | Type | Constraints |
|-------|------|------------|
| ArticleID | uint | FK → Article |
| UserID | uint | FK → User |
| Body | string | required |

### Tag
| Field | Type | Constraints |
|-------|------|------------|
| Name | string | unique index |

Relationships:
- `Articles []Article` — M:N via `article_tag`

## Join Tables (GORM 자동 생성)
- `article_favorite` — Article ↔ User
- `article_tag` — Article ↔ Tag
- `user_follower` — User ↔ User (Follow)

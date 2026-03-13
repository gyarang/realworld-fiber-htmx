---
name: add-model
description: GORM 모델 + 관계 설정 + 단위 테스트 생성. 새 Entity 추가 시 사용. "모델 추가", "Entity 추가", "테이블 추가" 시 사용.
argument-hint: [model-name]
user-invocable: true
---

# GORM 모델 생성

`$ARGUMENTS` 모델을 프로젝트 패턴에 맞게 생성합니다.

## 생성 파일

1. `cmd/web/model/$0.go` — 모델 정의
2. `cmd/web/model/$0_test.go` — 단위 테스트

## 사전 확인

1. `cmd/web/model/` 의 기존 모델 확인 (중복 방지)
2. `internal/database/database.go` 의 AutoMigrate 목록 확인

## 모델 구조

```go
package model

import "gorm.io/gorm"

type ModelName struct {
	gorm.Model
	// 필드 정의
	Name   string `gorm:"not null" validate:"required"`
	Slug   string `gorm:"uniqueIndex;not null"`
	UserID uint   // FK
	User   User   // belongs to

	// 비영속 필드 (계산값)
	IsActive bool `gorm:"-"`
}

// TableName — 테이블명 커스터마이즈 (선택)
func (m ModelName) TableName() string {
	return "model_names"
}
```

## 관계 패턴 (기존 모델 참조)

### belongs to (1:N의 N쪽)
```go
type Comment struct {
	UserID uint
	User   User  // FK: UserID → users.id
}
```

### has many (1:N의 1쪽)
```go
type Article struct {
	Comments []Comment  // FK: article_id
}
```

### many to many
```go
type Article struct {
	Tags      []Tag `gorm:"many2many:article_tag"`
	Favorites []User `gorm:"many2many:article_favorite"`
}
```

### self-referential (User ↔ User)
```go
type Follow struct {
	FollowerID  uint `gorm:"column:user_id;primaryKey"`
	FollowingID uint `gorm:"column:follower_id;primaryKey"`
}
```

## 메서드 패턴

```go
// Count 관련
func (m ModelName) GetRelatedCount() int64 {
	var count int64
	database.Get().Model(&RelatedModel{}).Where("model_id = ?", m.ID).Count(&count)
	return count
}

// 특정 유저 관계 확인
func (m ModelName) RelatedBy(userID uint) bool {
	var count int64
	database.Get().Table("join_table").
		Where("model_id = ? AND user_id = ?", m.ID, userID).
		Count(&count)
	return count > 0
}
```

## AutoMigrate 등록

`internal/database/database.go` 의 `Open()` 함수에 추가:
```go
db.AutoMigrate(&model.ModelName{})
```

## 단위 테스트 구조

```go
package model_test

import (
	"testing"

	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestModelName_CRUD(t *testing.T) {
	db := testutil.SetupTestDB(t)

	item := model.ModelName{Name: "test"}
	result := db.Create(&item)
	assert.NoError(t, result.Error)
	assert.NotZero(t, item.ID)

	var found model.ModelName
	db.First(&found, item.ID)
	assert.Equal(t, "test", found.Name)
}

func TestModelName_Validation(t *testing.T) {
	// 필수 필드 누락 테스트
}

func TestModelName_Relationships(t *testing.T) {
	// 관계 로딩 테스트 (Preload)
}

func TestModelName_Methods(t *testing.T) {
	// 커스텀 메서드 테스트
}
```

## 규칙

- `gorm.Model` 임베드 (ID, CreatedAt, UpdatedAt, DeletedAt 자동 포함)
- unique 필드: `gorm:"uniqueIndex"`
- required 필드: `validate:"required"`
- 비영속 필드: `gorm:"-"`
- import 순서: stdlib → third-party → local
- 테스트에서 DB는 `testutil.SetupTestDB(t)` 사용 (in-memory SQLite)
- `database.Get()` 으로 전역 DB 접근

기존 데이터 모델은 [data-model spec](../../../openspec/specs/data-model.md) 참조.

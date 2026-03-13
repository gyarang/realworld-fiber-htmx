package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal/testutil"
)

func TestGetFavoriteCount(t *testing.T) {
	article := model.Article{
		Favorites: []model.User{
			{ID: 1, Name: "User1"},
			{ID: 2, Name: "User2"},
		},
	}
	assert.Equal(t, 2, article.GetFavoriteCount())
}

func TestFavoritedBy(t *testing.T) {
	tests := []struct {
		name      string
		favorites []model.User
		userID    uint
		want      bool
	}{
		{
			"favorited by user",
			[]model.User{{ID: 1}, {ID: 2}},
			1,
			true,
		},
		{
			"not favorited",
			[]model.User{{ID: 1}, {ID: 2}},
			3,
			false,
		},
		{
			"nil favorites",
			nil,
			1,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := model.Article{Favorites: tt.favorites}
			assert.Equal(t, tt.want, article.FavoritedBy(tt.userID))
		})
	}
}

func TestGetTagsAsCommaSeparated(t *testing.T) {
	tests := []struct {
		name string
		tags []model.Tag
		want string
	}{
		{
			"multiple tags",
			[]model.Tag{{Name: "go"}, {Name: "htmx"}},
			"go,htmx,",
		},
		{
			"single tag",
			[]model.Tag{{Name: "go"}},
			"go,",
		},
		{
			"no tags",
			[]model.Tag{},
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := model.Article{Tags: tt.tags}
			assert.Equal(t, tt.want, article.GetTagsAsCommaSeparated())
		})
	}
}

func TestArticleCRUD(t *testing.T) {
	db := testutil.SetupTestDB(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "pass")
	article := testutil.CreateTestArticle(t, db, "Test Title", "test-title", "desc", "body content", user.ID, nil)
	require.NotZero(t, article.ID)

	var found model.Article
	err := db.Preload("User").First(&found, article.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "test-title", found.Slug)
	assert.Equal(t, user.ID, found.User.ID)
}

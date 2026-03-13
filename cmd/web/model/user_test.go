package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal/testutil"
)

func TestHashPassword(t *testing.T) {
	user := model.User{Password: "plaintext123"}
	user.HashPassword()

	assert.NotEqual(t, "plaintext123", user.Password)
	assert.NotEmpty(t, user.Password)
}

func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		check    string
		want     bool
	}{
		{"correct password", "secret123", "secret123", true},
		{"wrong password", "secret123", "wrong", false},
		{"empty password check", "secret123", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := model.User{Password: tt.password}
			user.HashPassword()
			assert.Equal(t, tt.want, user.CheckPassword(tt.check))
		})
	}
}

func TestFollowedBy(t *testing.T) {
	tests := []struct {
		name      string
		followers []model.Follow
		checkID   uint
		want      bool
	}{
		{
			"followed by user",
			[]model.Follow{{FollowerID: 1, FollowingID: 2}},
			2,
			true,
		},
		{
			"not followed",
			[]model.Follow{{FollowerID: 1, FollowingID: 3}},
			2,
			false,
		},
		{
			"nil followers",
			nil,
			1,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := model.User{Followers: tt.followers}
			assert.Equal(t, tt.want, user.FollowedBy(tt.checkID))
		})
	}
}

func TestFollowersCount(t *testing.T) {
	user := model.User{
		Followers: []model.Follow{
			{FollowerID: 1, FollowingID: 2},
			{FollowerID: 1, FollowingID: 3},
		},
	}
	assert.Equal(t, 2, user.FollowersCount())
}

func TestFollowTableName(t *testing.T) {
	f := model.Follow{}
	assert.Equal(t, "user_follower", f.TableName())
}

func TestUserCRUD(t *testing.T) {
	db := testutil.SetupTestDB(t)

	user := testutil.CreateTestUser(t, db, "Test User", "testuser", "test@example.com", "password123")
	require.NotZero(t, user.ID)

	var found model.User
	err := db.First(&found, user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "testuser", found.Username)
	assert.Equal(t, "test@example.com", found.Email)
	assert.True(t, found.CheckPassword("password123"))
}

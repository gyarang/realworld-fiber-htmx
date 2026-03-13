package HTMXController_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"realworld-fiber-htmx/cmd/web/model"
	"realworld-fiber-htmx/internal/testutil"
)

func TestUserArticleFavoriteAction_Unauthenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Article", "user-fav", "desc", "body", user.ID, nil)

	req := httptest.NewRequest(http.MethodPost, "/htmx/users/articles/user-fav/favorite", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("HX-Replace-Url"))
}

func TestUserArticleFavoriteAction_AddFavorite(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	author := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Article", "user-fav", "desc", "body", author.ID, nil)

	viewer := testutil.CreateTestUser(t, db, "Viewer", "viewer", "viewer@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, viewer)

	req := httptest.NewRequest(http.MethodPost, "/htmx/users/articles/user-fav/favorite", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserArticleFavoriteAction_RemoveFavorite(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	author := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	article := testutil.CreateTestArticle(t, db, "Article", "user-fav", "desc", "body", author.ID, nil)

	viewer := testutil.CreateTestUser(t, db, "Viewer", "viewer", "viewer@example.com", "password123")
	_ = db.Model(&article).Association("Favorites").Append(&viewer)

	cookie := testutil.AuthenticateUser(t, app, viewer)

	req := httptest.NewRequest(http.MethodPost, "/htmx/users/articles/user-fav/favorite", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserFollowAction_Follow(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	testutil.CreateTestUser(t, db, "Target", "target", "target@example.com", "password123")
	viewer := testutil.CreateTestUser(t, db, "Viewer", "viewer", "viewer@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, viewer)

	req := httptest.NewRequest(http.MethodPost, "/htmx/users/target/follow", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserFollowAction_Unfollow(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	target := testutil.CreateTestUser(t, db, "Target", "target", "target@example.com", "password123")
	viewer := testutil.CreateTestUser(t, db, "Viewer", "viewer", "viewer@example.com", "password123")

	// Follow first
	follow := model.Follow{FollowerID: target.ID, FollowingID: viewer.ID}
	db.Create(&follow)

	cookie := testutil.AuthenticateUser(t, app, viewer)

	req := httptest.NewRequest(http.MethodPost, "/htmx/users/target/follow", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserFollowAction_Unauthenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	testutil.CreateTestUser(t, db, "Target", "target", "target@example.com", "password123")

	req := httptest.NewRequest(http.MethodPost, "/htmx/users/target/follow", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Renders with empty authenticated user (no redirect in this handler)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

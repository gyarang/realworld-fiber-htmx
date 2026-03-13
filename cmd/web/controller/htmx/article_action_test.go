package HTMXController_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"realworld-fiber-htmx/internal/testutil"
)

func TestArticleFavoriteAction_Unauthenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Article", "test-fav", "desc", "body", user.ID, nil)

	req := httptest.NewRequest(http.MethodPost, "/htmx/articles/test-fav/favorite", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("HX-Replace-Url"))
}

func TestArticleFavoriteAction_AddFavorite(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	author := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Article", "test-fav", "desc", "body", author.ID, nil)

	viewer := testutil.CreateTestUser(t, db, "Viewer", "viewer", "viewer@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, viewer)

	req := httptest.NewRequest(http.MethodPost, "/htmx/articles/test-fav/favorite", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestArticleFavoriteAction_RemoveFavorite(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	author := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	article := testutil.CreateTestArticle(t, db, "Article", "test-fav", "desc", "body", author.ID, nil)

	viewer := testutil.CreateTestUser(t, db, "Viewer", "viewer", "viewer@example.com", "password123")
	_ = db.Model(&article).Association("Favorites").Append(&viewer)

	cookie := testutil.AuthenticateUser(t, app, viewer)

	req := httptest.NewRequest(http.MethodPost, "/htmx/articles/test-fav/favorite", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestArticleFollowAction_Unauthenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	author := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Article", "follow-test", "desc", "body", author.ID, nil)

	req := httptest.NewRequest(http.MethodPost, "/htmx/articles/follow-user/follow-test", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("HX-Replace-Url"))
}

func TestArticleFollowAction_Follow(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	author := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Article", "follow-test", "desc", "body", author.ID, nil)

	viewer := testutil.CreateTestUser(t, db, "Viewer", "viewer", "viewer@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, viewer)

	req := httptest.NewRequest(http.MethodPost, "/htmx/articles/follow-user/follow-test", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

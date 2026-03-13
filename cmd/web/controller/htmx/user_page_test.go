package HTMXController_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"realworld-fiber-htmx/internal/testutil"
)

func TestHTMXUserDetailPage_Unauthenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	testutil.CreateTestUser(t, db, "John", "john", "john@example.com", "password123")

	req := httptest.NewRequest(http.MethodGet, "/htmx/users/john", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTMXUserDetailPage_Self(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Self", "selfuser", "self@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, user)

	req := httptest.NewRequest(http.MethodGet, "/htmx/users/selfuser", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTMXUserDetailPage_Other(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	testutil.CreateTestUser(t, db, "Target", "target", "target@example.com", "password123")
	viewer := testutil.CreateTestUser(t, db, "Viewer", "viewer", "viewer@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, viewer)

	req := httptest.NewRequest(http.MethodGet, "/htmx/users/target", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTMXUserArticles(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Writer", "writer", "writer@example.com", "password123")
	testutil.CreateTestArticle(t, db, "User Article", "user-article", "desc", "body", user.ID, nil)

	req := httptest.NewRequest(http.MethodGet, "/htmx/users/writer/articles", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTMXUserArticles_Authenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	author := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Article", "article-1", "desc", "body", author.ID, nil)

	viewer := testutil.CreateTestUser(t, db, "Viewer", "viewer", "viewer@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, viewer)

	req := httptest.NewRequest(http.MethodGet, "/htmx/users/author/articles", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTMXUserArticles_NotFound(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/htmx/users/nonexistent/articles", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTMXUserArticlesFavorite(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	testutil.CreateTestUser(t, db, "User", "favuser", "fav@example.com", "password123")

	req := httptest.NewRequest(http.MethodGet, "/htmx/users/favuser/favorites", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTMXUserArticlesFavorite_Authenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	target := testutil.CreateTestUser(t, db, "Target", "target", "target@example.com", "password123")
	viewer := testutil.CreateTestUser(t, db, "Viewer", "viewer", "viewer@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, viewer)

	_ = target

	req := httptest.NewRequest(http.MethodGet, "/htmx/users/target/favorites", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

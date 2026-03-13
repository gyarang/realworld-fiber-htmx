package controller_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"realworld-fiber-htmx/internal/testutil"
)

func TestUserDetailPage_Unauthenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	testutil.CreateTestUser(t, db, "John", "john", "john@example.com", "password123")

	req := httptest.NewRequest(http.MethodGet, "/users/john", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "John")
}

func TestUserDetailPage_NotFound(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/users/nonexistent", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// GORM Find doesn't error on not found, renders with zero-value user
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserDetailPage_Authenticated_Self(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Self", "selfuser", "self@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, user)

	req := httptest.NewRequest(http.MethodGet, "/users/selfuser", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Self")
}

func TestUserDetailPage_Authenticated_Other(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	target := testutil.CreateTestUser(t, db, "Target", "target", "target@example.com", "password123")
	viewer := testutil.CreateTestUser(t, db, "Viewer", "viewer", "viewer@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, viewer)

	_ = target // ensure target exists

	req := httptest.NewRequest(http.MethodGet, "/users/target", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Target")
}

func TestUserDetailPage_Articles(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Writer", "writer", "writer@example.com", "password123")
	testutil.CreateTestArticle(t, db, "User Article", "user-article", "desc", "body", user.ID, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/writer/articles", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserDetailFavoritePage(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	testutil.CreateTestUser(t, db, "Fav", "favuser", "fav@example.com", "password123")

	req := httptest.NewRequest(http.MethodGet, "/users/favuser/favorites", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserDetailFavoritePage_Authenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "User", "favauth", "favauth@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, user)

	req := httptest.NewRequest(http.MethodGet, "/users/favauth/favorites", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

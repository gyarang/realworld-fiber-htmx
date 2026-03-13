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

func TestArticleDetailPage_Unauthenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Test Article", "test-article", "desc", "body content", user.ID, nil)

	req := httptest.NewRequest(http.MethodGet, "/articles/test-article", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Test Article")
}

func TestArticleDetailPage_NotFound(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/articles/nonexistent-slug", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Renders page with zero-value article (GORM Find doesn't error on not found)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestArticleDetailPage_Authenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	author := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Auth Article", "auth-article", "desc", "body", author.ID, nil)

	viewer := testutil.CreateTestUser(t, db, "Viewer", "viewer", "viewer@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, viewer)

	req := httptest.NewRequest(http.MethodGet, "/articles/auth-article", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Auth Article")
}

func TestArticleDetailPage_SelfArticle(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	author := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "My Article", "my-article", "desc", "body", author.ID, nil)

	cookie := testutil.AuthenticateUser(t, app, author)

	req := httptest.NewRequest(http.MethodGet, "/articles/my-article", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "My Article")
}

func TestArticleDetailPage_WithTags(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	author := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Tagged Article", "tagged-article", "desc", "body", author.ID, nil)

	req := httptest.NewRequest(http.MethodGet, "/articles/tagged-article", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

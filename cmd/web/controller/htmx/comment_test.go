package HTMXController_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"realworld-fiber-htmx/internal/testutil"
)

func TestArticleDetailCommentList(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	author := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	article := testutil.CreateTestArticle(t, db, "Test Article", "test-article", "desc", "body", author.ID, nil)
	testutil.CreateTestComment(t, db, article.ID, author.ID, "A test comment")

	req := httptest.NewRequest(http.MethodGet, "/htmx/articles/test-article/comments", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestArticleComment_Unauthenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	author := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Test Article", "test-article", "desc", "body", author.ID, nil)

	body := strings.NewReader("comment=Hello")
	req := httptest.NewRequest(http.MethodPost, "/htmx/articles/test-article/comments", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Unauthenticated users get HTMX redirect
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("HX-Replace-Url"))
}

func TestArticleComment_EmptyBody(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "User", "user", "user@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Test Article", "test-article", "desc", "body", user.ID, nil)

	cookie := testutil.AuthenticateUser(t, app, user)

	body := strings.NewReader("comment=")
	req := httptest.NewRequest(http.MethodPost, "/htmx/articles/test-article/comments", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Validation error rendered
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestArticleComment_Success(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "User", "user", "user@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Test Article", "test-article", "desc", "body", user.ID, nil)

	cookie := testutil.AuthenticateUser(t, app, user)

	body := strings.NewReader("comment=Great+article!")
	req := httptest.NewRequest(http.MethodPost, "/htmx/articles/test-article/comments", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

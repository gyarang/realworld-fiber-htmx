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

func TestHTMXHomePage_Unauthenticated(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/htmx/home", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTMXHomePage_Authenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "User", "homeuser", "home@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, user)

	req := httptest.NewRequest(http.MethodGet, "/htmx/home", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHomeGlobalFeed_Empty(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/htmx/home/global-feed", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHomeGlobalFeed_WithArticles(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Article 1", "article-1", "desc", "body", user.ID, nil)
	testutil.CreateTestArticle(t, db, "Article 2", "article-2", "desc", "body", user.ID, nil)

	req := httptest.NewRequest(http.MethodGet, "/htmx/home/global-feed", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHomeGlobalFeed_Authenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	author := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Article 1", "article-1", "desc", "body", author.ID, nil)

	viewer := testutil.CreateTestUser(t, db, "Viewer", "viewer", "viewer@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, viewer)

	req := httptest.NewRequest(http.MethodGet, "/htmx/home/global-feed", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHomeGlobalFeed_Pagination(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	for i := 0; i < 10; i++ {
		testutil.CreateTestArticle(t, db, "Article", "article-"+string(rune('a'+i)), "desc", "body", user.ID, nil)
	}

	req := httptest.NewRequest(http.MethodGet, "/htmx/home/global-feed?page=2", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHomeYourFeed_Unauthenticated(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/htmx/home/your-feed", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Redirects unauthenticated users
	assert.Equal(t, http.StatusFound, resp.StatusCode)
}

func TestHomeYourFeed_Authenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "User", "feeduser", "feed@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, user)

	req := httptest.NewRequest(http.MethodGet, "/htmx/home/your-feed", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHomeTagFeed(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	tag := model.Tag{Name: "golang"}
	db.Create(&tag)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Go Article", "go-article", "desc", "body", user.ID, []model.Tag{tag})

	req := httptest.NewRequest(http.MethodGet, "/htmx/home/tag-feed/golang", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHomeTagFeed_Authenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	tag := model.Tag{Name: "htmx"}
	db.Create(&tag)

	author := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "HTMX Article", "htmx-article", "desc", "body", author.ID, []model.Tag{tag})

	viewer := testutil.CreateTestUser(t, db, "Viewer", "viewer", "viewer@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, viewer)

	req := httptest.NewRequest(http.MethodGet, "/htmx/home/tag-feed/htmx", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHomeTagList(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	tag1 := model.Tag{Name: "golang"}
	tag2 := model.Tag{Name: "htmx"}
	db.Create(&tag1)
	db.Create(&tag2)

	req := httptest.NewRequest(http.MethodGet, "/htmx/home/tag-list", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHomeTagList_Empty(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/htmx/home/tag-list", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"realworld-fiber-htmx/internal/testutil"
)

func TestEditorPage_Unauthenticated_Redirects(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/editor", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
}

func TestEditorPage_Authenticated_NewArticle(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, user)

	req := httptest.NewRequest(http.MethodGet, "/editor", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestEditorPage_Authenticated_EditArticle(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Edit Article", "edit-article", "desc", "body", user.ID, nil)
	cookie := testutil.AuthenticateUser(t, app, user)

	req := httptest.NewRequest(http.MethodGet, "/editor/edit-article", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

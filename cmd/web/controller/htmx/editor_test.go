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

func TestEditorPage_Unauthenticated(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/htmx/editor", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestEditorPage_Authenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, user)

	req := httptest.NewRequest(http.MethodGet, "/htmx/editor", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestEditorPage_EditExisting(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Edit Me", "edit-me", "desc", "body", user.ID, nil)
	cookie := testutil.AuthenticateUser(t, app, user)

	req := httptest.NewRequest(http.MethodGet, "/htmx/editor/edit-me", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestStoreArticle_Unauthenticated(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	body := strings.NewReader("title=Test&description=Desc&content=Body")
	req := httptest.NewRequest(http.MethodPost, "/htmx/editor", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Redirect to home for unauthenticated
	assert.Equal(t, http.StatusFound, resp.StatusCode)
}

func TestStoreArticle_EmptyFields(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, user)

	body := strings.NewReader("title=&description=&content=")
	req := httptest.NewRequest(http.MethodPost, "/htmx/editor", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Validation error renders editor with errors
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestStoreArticle_Success(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, user)

	body := strings.NewReader("title=New+Article&description=A+description&content=Article+body+here")
	req := httptest.NewRequest(http.MethodPost, "/htmx/editor", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("HX-Replace-Url"))
}

func TestUpdateArticle_Unauthenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Update Me", "update-me", "desc", "body", user.ID, nil)

	body := strings.NewReader("title=Updated&description=Updated+Desc&content=Updated+Body")
	req := httptest.NewRequest(http.MethodPatch, "/htmx/editor/update-me", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
}

func TestUpdateArticle_Success(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Update Me", "update-me", "desc", "body", user.ID, nil)
	cookie := testutil.AuthenticateUser(t, app, user)

	body := strings.NewReader("title=Updated+Title&description=Updated+Desc&content=Updated+Body")
	req := httptest.NewRequest(http.MethodPatch, "/htmx/editor/update-me", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("HX-Replace-Url"))
}

func TestUpdateArticle_EmptyFields(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "Author", "author", "author@example.com", "password123")
	testutil.CreateTestArticle(t, db, "Update Me", "update-me", "desc", "body", user.ID, nil)
	cookie := testutil.AuthenticateUser(t, app, user)

	body := strings.NewReader("title=&description=&content=")
	req := httptest.NewRequest(http.MethodPatch, "/htmx/editor/update-me", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Validation error renders editor
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

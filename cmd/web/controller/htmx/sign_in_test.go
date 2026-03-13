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

func TestSignInPage(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/htmx/sign-in", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSignInAction_EmptyFields(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	body := strings.NewReader("email=&password=")
	req := httptest.NewRequest(http.MethodPost, "/htmx/sign-in", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSignInAction_InvalidCredentials(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	testutil.CreateTestUser(t, db, "Test", "testuser", "test@example.com", "password123")

	body := strings.NewReader("email=test@example.com&password=wrongpassword")
	req := httptest.NewRequest(http.MethodPost, "/htmx/sign-in", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSignInAction_Success(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	testutil.CreateTestUser(t, db, "Test", "testuser", "test@example.com", "password123")

	body := strings.NewReader("email=test@example.com&password=password123")
	req := httptest.NewRequest(http.MethodPost, "/htmx/sign-in", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("HX-Replace-Url"))
}

func TestSignInAction_UserNotFound(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	body := strings.NewReader("email=nonexistent@example.com&password=password123")
	req := httptest.NewRequest(http.MethodPost, "/htmx/sign-in", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

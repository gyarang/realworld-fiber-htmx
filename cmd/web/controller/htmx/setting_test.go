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

func TestHTMXSettingPage_Unauthenticated(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/htmx/settings", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTMXSettingPage_Authenticated(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "User", "user", "user@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, user)

	req := httptest.NewRequest(http.MethodGet, "/htmx/settings", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSettingAction_Unauthenticated(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	body := strings.NewReader("name=Test&email=test@example.com")
	req := httptest.NewRequest(http.MethodPost, "/htmx/settings", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// HTMX redirect to sign-in
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("HX-Replace-Url"))
}

func TestSettingAction_EmptyFields(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "User", "user", "user@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, user)

	body := strings.NewReader("name=&email=")
	req := httptest.NewRequest(http.MethodPost, "/htmx/settings", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Validation error
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSettingAction_Success(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "User", "user", "user@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, user)

	body := strings.NewReader("name=Updated+Name&email=updated@example.com&bio=Hello&image=")
	req := httptest.NewRequest(http.MethodPost, "/htmx/settings", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSettingAction_WithPassword(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "User", "user", "user@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, user)

	body := strings.NewReader("name=User&email=user@example.com&password=newpassword123")
	req := httptest.NewRequest(http.MethodPost, "/htmx/settings", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

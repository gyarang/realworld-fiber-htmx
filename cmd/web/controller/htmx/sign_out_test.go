package HTMXController_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"realworld-fiber-htmx/internal/testutil"
)

func TestSignOut(t *testing.T) {
	app, db := testutil.SetupTestApp(t)

	user := testutil.CreateTestUser(t, db, "User", "user", "user@example.com", "password123")
	cookie := testutil.AuthenticateUser(t, app, user)

	req := httptest.NewRequest(http.MethodPost, "/htmx/sign-out", nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("HX-Replace-Url"))
}

func TestSignOut_Unauthenticated(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	req := httptest.NewRequest(http.MethodPost, "/htmx/sign-out", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Redirects unauthenticated users
	assert.Equal(t, http.StatusFound, resp.StatusCode)
}

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

func TestSignUpPage(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/htmx/sign-up", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSignUpAction_EmptyFields(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	body := strings.NewReader("username=&email=&password=")
	req := httptest.NewRequest(http.MethodPost, "/htmx/sign-up", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSignUpAction_Success(t *testing.T) {
	app, _ := testutil.SetupTestApp(t)

	body := strings.NewReader("username=newuser&email=new@example.com&password=password123")
	req := httptest.NewRequest(http.MethodPost, "/htmx/sign-up", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("HX-Replace-Url"))
}

package testutil

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"realworld-fiber-htmx/cmd/web/model"
	webroute "realworld-fiber-htmx/cmd/web/route"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"
	"realworld-fiber-htmx/internal/renderer"
)

// SetupTestApp creates a full Fiber application for integration testing.
// It initializes in-memory SQLite, session store, template engine, and routes.
func SetupTestApp(t *testing.T) (*fiber.App, *gorm.DB) {
	t.Helper()

	// Change to project root so template paths resolve correctly
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change to project root: %v", err)
	}

	// Setup in-memory test database
	db := SetupTestDB(t)
	database.DB = db

	// Setup session store (in-memory for tests)
	authentication.StoredAuthenticationSession = session.New()

	// Setup template engine
	viewEngine := renderer.ViewEngineStart()
	app := fiber.New(fiber.Config{
		Views: viewEngine,
	})

	// Register routes
	webroute.WebHandlers(app)
	webroute.HTMXHandlers(app)

	return app, db
}

// AuthenticateUser creates a session cookie for the given user.
// Returns the cookie string to be used in subsequent requests.
func AuthenticateUser(t *testing.T, app *fiber.App, user model.User) string {
	t.Helper()

	// Make a sign-in request to get session cookie
	body := strings.NewReader("email=" + user.Email + "&password=password123")
	req := httptest.NewRequest(http.MethodPost, "/htmx/sign-in", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to authenticate user: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Extract session cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "session_id" {
			return cookie.Name + "=" + cookie.Value
		}
	}

	// If no session_id cookie, return all cookies
	var cookies []string
	for _, cookie := range resp.Cookies() {
		cookies = append(cookies, cookie.Name+"="+cookie.Value)
	}
	return strings.Join(cookies, "; ")
}

// CreateTestComment creates a comment in the test database.
func CreateTestComment(t *testing.T, db *gorm.DB, articleID, userID uint, body string) model.Comment {
	t.Helper()

	comment := model.Comment{
		ArticleID: articleID,
		UserID:    userID,
		Body:      body,
	}

	result := db.Create(&comment)
	if result.Error != nil {
		t.Fatalf("failed to create test comment: %v", result.Error)
	}

	return comment
}

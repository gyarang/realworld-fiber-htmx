package testutil

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"realworld-fiber-htmx/cmd/web/model"
)

// SetupTestDB creates an in-memory SQLite database for testing.
// It auto-migrates all models and returns the *gorm.DB instance.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Article{}, &model.Comment{}, &model.Tag{}, &model.Follow{})
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}

// CreateTestUser creates a user with hashed password in the test database.
func CreateTestUser(t *testing.T, db *gorm.DB, name, username, email, password string) model.User {
	t.Helper()

	user := model.User{
		Name:     name,
		Username: username,
		Email:    email,
		Password: password,
	}
	user.HashPassword()

	result := db.Create(&user)
	if result.Error != nil {
		t.Fatalf("failed to create test user: %v", result.Error)
	}

	return user
}

// CreateTestArticle creates an article in the test database.
func CreateTestArticle(t *testing.T, db *gorm.DB, title, slug, description, body string, userID uint, tags []model.Tag) model.Article {
	t.Helper()

	article := model.Article{
		Title:       title,
		Slug:        slug,
		Description: description,
		Body:        body,
		UserID:      userID,
		Tags:        tags,
	}

	result := db.Create(&article)
	if result.Error != nil {
		t.Fatalf("failed to create test article: %v", result.Error)
	}

	return article
}

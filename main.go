package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"realworld-fiber-htmx/cmd/web"
	"realworld-fiber-htmx/internal/authentication"
	"realworld-fiber-htmx/internal/database"
	"realworld-fiber-htmx/internal/renderer"
)

func main() {

	viewEngine := renderer.ViewEngineStart()
	app := fiber.New(fiber.Config{
		Views: viewEngine,
	})

	database.Open()
	authentication.SessionStart()

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	web.Serve(app)

	log.Fatal(app.Listen("localhost:8181"))
}

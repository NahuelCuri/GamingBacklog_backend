package routes

import (
	"backlog-backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Health
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! Backend is running.")
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Backend is reachable",
		})
	})

	// API Group
	api := app.Group("/api")

	// Game Routes
	games := api.Group("/games")
	games.Get("/", handlers.GetGames)
	games.Post("/", handlers.CreateGame)
	// Add other CRUD routes here...
}

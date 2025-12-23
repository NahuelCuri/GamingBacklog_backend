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

	// User Routes
	users := api.Group("/users")
	users.Get("/", handlers.GetUsers)
	users.Get("/:id", handlers.GetUser)
	users.Post("/", handlers.CreateUser)
	users.Post("/login", handlers.Login)
	users.Put("/:id", handlers.UpdateUser)
	users.Delete("/:id", handlers.DeleteUser)

	// Game Routes
	games := api.Group("/games")
	games.Get("/", handlers.GetGames)
	games.Get("/:id", handlers.GetGame)
	games.Post("/", handlers.CreateGame)
	games.Put("/:id", handlers.UpdateGame)
	games.Delete("/:id", handlers.DeleteGame)

	// Tag Routes
	tags := api.Group("/tags")
	tags.Get("/", handlers.GetTags)
	tags.Post("/", handlers.CreateTag)
	tags.Put("/:id", handlers.UpdateTag)
	tags.Delete("/:id", handlers.DeleteTag)
}

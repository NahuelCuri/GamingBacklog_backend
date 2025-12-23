package routes

import (
	"backlog-backend/handlers"
	"backlog-backend/middleware"

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
	// Protected User Routes
	users.Put("/:id", middleware.Protected(), handlers.UpdateUser)
	users.Delete("/:id", middleware.Protected(), handlers.DeleteUser)

	// Game Routes
	games := api.Group("/games")
	games.Get("/", handlers.GetGames)
	games.Get("/:id", handlers.GetGame)
	// Protected Game Routes
	games.Post("/", middleware.Protected(), handlers.CreateGame)
	games.Put("/:id", middleware.Protected(), handlers.UpdateGame)
	games.Delete("/:id", middleware.Protected(), handlers.DeleteGame)

	// Tag Routes
	tags := api.Group("/tags")
	tags.Get("/", handlers.GetTags)
	// Protected Tag Routes
	tags.Post("/", middleware.Protected(), handlers.CreateTag)
	tags.Put("/:id", middleware.Protected(), handlers.UpdateTag)
	tags.Delete("/:id", middleware.Protected(), handlers.DeleteTag)
}

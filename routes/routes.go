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
	// Protected Game Routes
	games.Use(middleware.Protected())
	games.Get("/", handlers.GetGames)
	games.Get("/:id", handlers.GetGame)
	games.Post("/", handlers.CreateGame)
	games.Put("/:id", handlers.UpdateGame)
	games.Delete("/:id", handlers.DeleteGame)

	// Tag Routes
	tags := api.Group("/tags")
	tags.Get("/", handlers.GetTags)
	// Protected Tag Routes
	tags.Post("/", middleware.Protected(), handlers.CreateTag)
	tags.Put("/:id", middleware.Protected(), handlers.UpdateTag)
	tags.Delete("/:id", middleware.Protected(), handlers.DeleteTag)

	// Image Routes
	// Post /api/upload
	api.Post("/upload", middleware.Protected(), handlers.UploadImage)
	// Delete /api/images/:filename
	api.Delete("/images/:filename", middleware.Protected(), handlers.DeleteImage)
	// Tier List Routes
	tierLists := api.Group("/tier-lists")
	tierLists.Use(middleware.Protected())
	tierLists.Get("/", handlers.GetTierLists)
	tierLists.Get("/:id", handlers.GetTierList)
	tierLists.Post("/", handlers.CreateTierList)
	tierLists.Put("/:id", handlers.UpdateTierList)
	tierLists.Delete("/:id", handlers.DeleteTierList)
}

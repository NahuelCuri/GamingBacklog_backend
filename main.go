package main

import (
	"log"
	"os"

	"backlog-backend/database"
	"backlog-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default/env vars")
	}

	// Connect to Database
	database.Connect()

	// Initialize Fiber app
	app := fiber.New()

	// CORS Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Setup Routes
	routes.SetupRoutes(app)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

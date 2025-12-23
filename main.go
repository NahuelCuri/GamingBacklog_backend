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
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(logger.New())

	// Ensure images directory exists
	if _, err := os.Stat("./images"); os.IsNotExist(err) {
		if err := os.Mkdir("./images", 0755); err != nil {
			log.Printf("Warning: Could not create images directory: %v", err)
		}
	}
	app.Static("/images", "./images")

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

package handlers

import (
	"backlog-backend/database"
	"backlog-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ListGames godoc
func GetGames(c *fiber.Ctx) error {
	var games []models.Game
	result := database.DB.Preload("Tags").Find(&games)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not fetch games",
			"error":   result.Error.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   games,
	})
}

// CreateGame godoc
func CreateGame(c *fiber.Ctx) error {
	game := new(models.Game)
	if err := c.BodyParser(game); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"error":   err.Error(),
		})
	}

	// Basic validation or default assignment could go here
	// Verify UserID is present if auth is required, etc.
	if game.UserID == uuid.Nil {
		// For now, allow it or assign a placeholder if no auth middleware yet
		// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User ID required"})
	}

	if err := database.DB.Create(&game).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create game",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   game,
	})
}

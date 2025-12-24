package handlers

import (
	"backlog-backend/database"
	"backlog-backend/dto"
	"backlog-backend/models"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Helper to map model to DTO
func mapGameToResponse(game models.Game) dto.GameResponse {
	var tagResponses []dto.TagResponse
	for _, tag := range game.Tags {
		tagResponses = append(tagResponses, dto.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	return dto.GameResponse{
		ID:           game.ID,
		UserID:       game.UserID,
		Title:        game.Title,
		CoverURL:     game.CoverURL,
		Genre:        game.Genre,
		Status:       string(game.Status),
		Score:        game.Score,
		HoursPlayed:  game.HoursPlayed,
		HLTBEstimate: game.HLTBEstimate,
		ReleaseYear:  game.ReleaseYear,
		DateFinished: game.DateFinished,
		LastPlayedAt: game.LastPlayedAt,
		ReviewText:   game.ReviewText,
		CreatedAt:    game.CreatedAt,
		UpdatedAt:    game.UpdatedAt,
		Tags:         tagResponses,
	}
}

func GetGames(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	fmt.Printf("DEBUG: GetGames - UserID: %v\n", userID)
	var games []models.Game
	result := database.DB.Preload("Tags").Where("user_id = ?", userID).Find(&games)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not fetch games", "error": result.Error.Error()})
	}

	res := make([]dto.GameResponse, 0)
	for _, game := range games {
		res = append(res, mapGameToResponse(game))
	}
	return c.JSON(fiber.Map{"status": "success", "data": res})
}

func GetGame(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")
	var game models.Game
	if err := database.DB.Preload("Tags").Where("user_id = ?", userID).First(&game, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Game not found"})
	}
	return c.JSON(fiber.Map{"status": "success", "data": mapGameToResponse(game)})
}

func CreateGame(c *fiber.Ctx) error {
	var req dto.CreateGameRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid input"})
	}

	userID := c.Locals("user_id")
	fmt.Printf("DEBUG: CreateGame - UserID: %v\n", userID)
	// Type assertion to ensure it's a string or UUID as expected by the model.
	// Assuming the JWT middleware puts the claim as is. If it is string in JWT:
	// But wait, the model UserID is uuid.UUID.
	// c.Locals("user_id") comes from claims["user_id"].
	// We need to parse it to UUID.

	game := models.Game{
		UserID:       uuid.MustParse(userID.(string)), // Should handle error safely if possible, but middleware checks validity
		Title:        req.Title,
		CoverURL:     req.CoverURL,
		Genre:        req.Genre,
		Status:       models.GameStatus(req.Status),
		Score:        req.Score,
		HoursPlayed:  req.HoursPlayed,
		HLTBEstimate: req.HLTBEstimate,
		ReleaseYear:  req.ReleaseYear,
		DateFinished: req.DateFinished,
		ReviewText:   req.ReviewText,
	}

	// Assign tags if provided
	if len(req.TagIDs) > 0 {
		var tags []*models.GameTag
		database.DB.Find(&tags, req.TagIDs)
		game.Tags = tags
	}

	if err := database.DB.Create(&game).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create game", "error": err.Error()})
	}

	// Reload to get associations
	database.DB.Preload("Tags").First(&game, game.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": mapGameToResponse(game)})
}

func UpdateGame(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")
	var game models.Game
	if err := database.DB.Preload("Tags").Where("user_id = ?", userID).First(&game, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Game not found"})
	}

	var req dto.UpdateGameRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid input"})
	}

	// Manual update mapping or helper
	if req.Title != "" {
		game.Title = req.Title
	}
	// ... map other fields ...
	// Check for image update/cleanup
	if req.CoverURL != game.CoverURL {
		// If old cover was a local image, delete it
		if game.CoverURL != "" && len(game.CoverURL) > 8 && game.CoverURL[:8] == "/images/" {
			oldImagePath := "." + game.CoverURL
			os.Remove(oldImagePath) // Ignore error
		}
	}
	game.CoverURL = req.CoverURL
	game.Genre = req.Genre
	if req.Status != "" {
		game.Status = models.GameStatus(req.Status)
	}
	game.Score = req.Score
	game.HoursPlayed = req.HoursPlayed
	game.HLTBEstimate = req.HLTBEstimate
	game.ReleaseYear = req.ReleaseYear
	game.DateFinished = req.DateFinished
	game.ReviewText = req.ReviewText

	// Update Tags association
	if req.TagIDs != nil { // empty slice means clear tags, nil means don't update
		var tags []*models.GameTag
		database.DB.Find(&tags, req.TagIDs)
		database.DB.Model(&game).Association("Tags").Replace(tags)
	}

	database.DB.Save(&game)

	// Reload for response
	database.DB.Preload("Tags").First(&game, game.ID)

	return c.JSON(fiber.Map{"status": "success", "data": mapGameToResponse(game)})
}

func DeleteGame(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")
	fmt.Printf("DEBUG: DeleteGame - GameID: %s, UserID: %v\n", id, userID)
	var game models.Game
	if err := database.DB.Where("user_id = ?", userID).First(&game, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Game not found"})
	}

	// Delete image if exists
	if game.CoverURL != "" && len(game.CoverURL) > 8 && game.CoverURL[:8] == "/images/" {
		// Assuming CoverURL is saved as "/images/filename"
		// Remove leading "/" to get relative path "./images/filename" (or just "images/filename")
		imagePath := "." + game.CoverURL
		if err := os.Remove(imagePath); err != nil {
			// Log error but proceed with deleting game
			// fmt.Printf("Failed to delete image %s: %v\n", imagePath, err)
		}
	}

	database.DB.Delete(&game)
	return c.JSON(fiber.Map{"status": "success", "message": "Game deleted"})
}

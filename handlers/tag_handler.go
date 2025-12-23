package handlers

import (
	"backlog-backend/database"
	"backlog-backend/dto"
	"backlog-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetTags(c *fiber.Ctx) error {
	var tags []models.GameTag
	database.DB.Find(&tags)

	var res []dto.TagResponse
	for _, tag := range tags {
		res = append(res, dto.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	return c.JSON(fiber.Map{"status": "success", "data": res})
}

func CreateTag(c *fiber.Ctx) error {
	var req dto.CreateTagRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid input"})
	}

	tag := models.GameTag{
		Name:   req.Name,
		UserID: req.UserID,
	}
	if req.UserID == uuid.Nil {
		// handle global tag or user assignment logic
	}

	if err := database.DB.Create(&tag).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create tag", "error": err.Error()})
	}

	res := dto.TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": res})
}

func UpdateTag(c *fiber.Ctx) error {
	id := c.Params("id")
	var tag models.GameTag
	if err := database.DB.First(&tag, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Tag not found"})
	}

	var req dto.UpdateTagRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid input"})
	}

	tag.Name = req.Name
	database.DB.Save(&tag)

	res := dto.TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}
	return c.JSON(fiber.Map{"status": "success", "data": res})
}

func DeleteTag(c *fiber.Ctx) error {
	id := c.Params("id")
	var tag models.GameTag
	if err := database.DB.First(&tag, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Tag not found"})
	}

	database.DB.Delete(&tag)
	return c.JSON(fiber.Map{"status": "success", "message": "Tag deleted"})
}

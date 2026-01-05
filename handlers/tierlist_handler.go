package handlers

import (
	"backlog-backend/database"
	"backlog-backend/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateTierList handles creating a new tier list
func CreateTierList(c *fiber.Ctx) error {
	userIdStr := c.Locals("user_id").(string)
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var tierList models.TierList
	if err := c.BodyParser(&tierList); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	tierList.UserID = userId

	// Ensure IDs are generated for Rows and Items if not present (GORM hooks handle this, but explicit check is good)
	// Actually, the BeforeCreate hooks I added to models should handle it.

	if err := database.DB.Create(&tierList).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create tier list"})
	}

	return c.Status(fiber.StatusCreated).JSON(tierList)
}

// GetTierLists returns all tier lists for the user (lightweight, no rows/items)
func GetTierLists(c *fiber.Ctx) error {
	userIdStr := c.Locals("user_id").(string)

	var tierLists []models.TierList
	if err := database.DB.Where("user_id = ?", userIdStr).Order("created_at desc").Find(&tierLists).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch tier lists"})
	}

	return c.JSON(tierLists)
}

// GetTierList returns a single tier list with full details
func GetTierList(c *fiber.Ctx) error {
	id := c.Params("id")
	userIdStr := c.Locals("user_id").(string)

	var tierList models.TierList

	// Preload Rows, Items, and the Game details for each Item
	if err := database.DB.
		Preload("Rows", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order asc") // Ensure rows are ordered
		}).
		Preload("Rows.Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order asc") // Ensure games in rows are ordered
		}).
		Preload("Rows.Items.Game").
		Where("id = ? AND user_id = ?", id, userIdStr).
		First(&tierList).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tier list not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch tier list"})
	}

	return c.JSON(tierList)
}

// UpdateTierList updates the tier list and strictly syncs rows/items
func UpdateTierList(c *fiber.Ctx) error {
	id := c.Params("id")
	userIdStr := c.Locals("user_id").(string)

	var input models.TierList
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var existingTierList models.TierList
	if err := database.DB.Where("id = ? AND user_id = ?", id, userIdStr).First(&existingTierList).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tier list not found"})
	}

	// Transaction to ensure atomicity
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Update basic info
		existingTierList.Name = input.Name
		if err := tx.Save(&existingTierList).Error; err != nil {
			return err
		}

		// 2. Delete existing rows (Cascade should handle items)
		if err := tx.Where("tier_list_id = ?", existingTierList.ID).Delete(&models.TierRow{}).Error; err != nil {
			return err
		}

		// 3. Create new rows
		for i := range input.Rows {
			input.Rows[i].TierListID = existingTierList.ID
			input.Rows[i].ID = uuid.Nil // Force new ID generation
			fmt.Printf("DEBUG: Processing Row %s (Items: %d)\n", input.Rows[i].Label, len(input.Rows[i].Items))
			for j := range input.Rows[i].Items {
				input.Rows[i].Items[j].TierRowID = uuid.Nil // Set to nil, will be auto-filled
				input.Rows[i].Items[j].ID = uuid.Nil
				fmt.Printf("DEBUG:   - Item GameID: %s\n", input.Rows[i].Items[j].GameID)
			}
		}

		if len(input.Rows) > 0 {
			if err := tx.Create(&input.Rows).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update tier list", "details": err.Error()})
	}

	// Fetch updated structure to return
	// Reuse GetTierList logic or just return success
	return c.JSON(fiber.Map{"message": "Updated successfully"})
}

// DeleteTierList deletes a tier list
func DeleteTierList(c *fiber.Ctx) error {
	id := c.Params("id")
	userIdStr := c.Locals("user_id").(string)

	result := database.DB.Where("id = ? AND user_id = ?", id, userIdStr).Delete(&models.TierList{})
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete tier list"})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tier list not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

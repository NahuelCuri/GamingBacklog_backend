package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UploadImage handles image uploads
func UploadImage(c *fiber.Ctx) error {
	// Parse the form file
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image upload failed",
		})
	}

	// Check file type (basic extension check)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid image format. Allowed: jpg, jpeg, png, gif, webp",
		})
	}

	// Get game name for filename
	gameName := c.FormValue("game_name")
	if gameName == "" {
		gameName = "unknown_game"
	}

	// Sanitize game name (remove non-alphanumeric characters, replace spaces with underscores)
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Regex error",
		})
	}
	sanitizedGameName := reg.ReplaceAllString(gameName, "_")
	// Trim underscores from start/end
	sanitizedGameName = strings.Trim(sanitizedGameName, "_")

	// Generate filename: [uuid]_[date]_[sanitized_game_name].ext
	// Date format: YYYYMMDD
	dateStr := time.Now().Format("20060102")
	id := uuid.New().String()

	newFilename := fmt.Sprintf("%s_%s_%s%s", id, dateStr, sanitizedGameName, ext)

	// Create images directory if it doesn't exist (safety check, main.go should do this too)
	if _, err := os.Stat("./images"); os.IsNotExist(err) {
		os.Mkdir("./images", 0755)
	}

	// Save file
	if err := c.SaveFile(file, fmt.Sprintf("./images/%s", newFilename)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not save image",
		})
	}

	// Return URL
	// Assuming server is running on same host/port, or via proxy.
	// The client will prepend the base URL.
	imageURL := fmt.Sprintf("/images/%s", newFilename)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"url": imageURL,
	})
}

// DeleteImage handles image deletion
func DeleteImage(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Filename required",
		})
	}

	// Security: Prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid filename",
		})
	}

	filePath := fmt.Sprintf("./images/%s", filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Image not found",
		})
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not delete image",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Image deleted successfully",
	})
}

package middlewares

import (
	"project-0/db"
	"project-0/pkg/imagekit"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CheckUpdateImage(c *fiber.Ctx) error {
	id := c.Params("id")
	villaID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid villa ID",
		})
	}
	imageFile, err := c.FormFile("image")
	if err != nil || imageFile == nil {
		return c.Next()
	}

	var currentImage string
	err = db.Conn.QueryRow(c.Context(), `SELECT image FROM villas WHERE id = $1`, villaID).Scan(&currentImage)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Villa not found or image not available",
		})
	}
	if currentImage != "" {
		_ = imagekit.DeleteImageByUrl(currentImage)
	}

	src, err := imageFile.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to open uploaded image"})
	}
	defer src.Close()

	imageURL, err := imagekit.UploadImageFromForm(src, imageFile.Filename)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Image upload failed: " + err.Error()})
	}

	c.Locals("image", imageURL)

	return c.Next()
}

func DeleteImage(c *fiber.Ctx) error {
	id := c.Params("id")
	villaID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid villa ID",
		})
	}

	var currentImage string
	err = db.Conn.QueryRow(c.Context(), `SELECT image FROM villas WHERE id = $1`, villaID).Scan(&currentImage)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Villa not found or image not available",
		})
	}
	if currentImage != "" {
		_ = imagekit.DeleteImageByUrl(currentImage)
	}

	return c.Next()
}
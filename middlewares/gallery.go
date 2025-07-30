package middlewares

import (
	"project-0/db"
	"project-0/models"
	"project-0/pkg/imagekit"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func DeleteGalleryByVillaId(c *fiber.Ctx) error {
	id := c.Params("id") // villa_id
	villaID, err := strconv.Atoi(id)
	// Query existing gallery image URLs
	rows, err := db.Conn.Query(c.Context(), `SELECT image_url FROM galleries WHERE villa_id = $1`, villaID)
	if err != nil {
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Failed to fetch gallery images",
		})
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err == nil {
			urls = append(urls, url)
		}
	}

	// Delete images from ImageKit
	for _, url := range urls {
		_ = imagekit.DeleteImageByUrl(url) // Optional: log errors
	}

	// Optionally delete rows from DB if this is not handled elsewhere
	_, _ = db.Conn.Exec(c.Context(), `DELETE FROM galleries WHERE villa_id = $1`, villaID)

	return c.Next()
}

func DeleteGalleryByID(c *fiber.Ctx) error {
	var req struct {
		IDs []int64 `json:"ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Invalid request body",
		})
	}

	if len(req.IDs) == 0 {
		return c.Status(400).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "No gallery IDs provided",
		})
	}

	// Step 1: Fetch image URLs
	query := `SELECT image_url FROM galleries WHERE id = ANY($1)`
	rows, err := db.Conn.Query(c.Context(), query, req.IDs)
	if err != nil {
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Failed to fetch gallery images",
		})
	}
	defer rows.Close()

	var imageURLs []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err == nil {
			imageURLs = append(imageURLs, url)
		}
	}

	if len(imageURLs) == 0 {
		return c.Status(404).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "No matching galleries found",
		})
	}

	// Step 2: Delete images from ImageKit
	for _, url := range imageURLs {
		_ = imagekit.DeleteImageByUrl(url) // Optional: log errors
	}

	return c.Next()
}

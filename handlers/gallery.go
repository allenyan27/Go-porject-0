package handlers

import (
	"fmt"
	"project-0/db"
	"project-0/models"
	"project-0/pkg/imagekit"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetGalleryVilla(c *fiber.Ctx) error {
	id := c.Params("id")
	villaID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Invalid villa ID",
		})
	}

	rows, err := db.Conn.Query(c.Context(), "SELECT * FROM galleries WHERE villa_id = $1 ORDER BY sort_order ASC", villaID)
	if err != nil {
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: err.Error(),
		})
	}
	defer rows.Close()

	var galleries []models.Gallery
	for rows.Next() {
		var g models.Gallery
		if err := rows.Scan(&g.ID, &g.VillaID, &g.ImageURL, &g.SortOrder, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return c.Status(500).JSON(models.ResponseWithMessage{
				Success: false,
				Message: "Failed to scan gallery: " + err.Error(),
			})
		}
		galleries = append(galleries, g)
	}

	return c.JSON(models.ResponseWithData[[]models.Gallery]{
		Success: true,
		Data:    galleries,
	})
}

func AddGallery(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Failed to parse multipart form",
		})
	}

	villaIDStr := c.FormValue("villa_id")
	villaID, err := strconv.ParseInt(villaIDStr, 10, 64)
	if err != nil || villaID <= 0 {
		return c.Status(400).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Invalid villa_id",
		})
	}

	galleryFiles := form.File["gallery"]
	if len(galleryFiles) == 0 {
		return c.Status(400).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "No gallery images provided",
		})
	}

	var last int
	query := `SELECT COALESCE(MAX(sort_order), 0) FROM galleries WHERE villa_id = $1`
	err = db.Conn.QueryRow(c.Context(), query, villaID).Scan(&last)
	if err != nil {
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Failed to get last sort order",
		})
	}

	tx, err := db.Conn.Begin(c.Context())
	if err != nil {
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Failed to start transaction",
		})
	}
	defer tx.Rollback(c.Context())

	now := time.Now()
	successCount := 0

	for i, file := range galleryFiles {
		order := last + i + 1

		src, err := file.Open()
		if err != nil {
			continue
		}
		defer src.Close()

		url, err := imagekit.UploadImageFromForm(src, file.Filename)
		if err != nil {
			continue
		}

		_, err = tx.Exec(c.Context(), `
			INSERT INTO galleries (villa_id, image_url, sort_order, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
		`, villaID, url, order, now, now)
		if err != nil {
			continue
		}
		successCount++
	}

	if err := tx.Commit(c.Context()); err != nil {
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Failed to commit transaction",
		})
	}

	return c.Status(200).JSON(models.ResponseWithMessage{
		Success: true,
		Message: fmt.Sprintf("%d galleries added successfully", successCount),
	})
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

	tx, err := db.Conn.Begin(c.Context())
	if err != nil {
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Failed to start transaction",
		})
	}
	defer tx.Rollback(c.Context())

	query := `DELETE FROM galleries WHERE id = ANY($1)`
	_, err = tx.Exec(c.Context(), query, req.IDs)
	if err != nil {
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Failed to delete galleries",
		})
	}

	if err := tx.Commit(c.Context()); err != nil {
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Failed to commit delete transaction",
		})
	}

	return c.Status(200).JSON(models.ResponseWithMessage{
		Success: true,
		Message: fmt.Sprintf("%d galleries deleted successfully", len(req.IDs)),
	})
}
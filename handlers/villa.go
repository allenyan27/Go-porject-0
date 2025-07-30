package handlers

import (
	"fmt"
	"project-0/db"
	"project-0/models"
	"project-0/pkg/imagekit"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetVillas(c *fiber.Ctx) error {
	rows, err := db.Conn.Query(c.Context(), "SELECT * FROM villas")
	if err != nil {
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: err.Error(),
		})
	}
	defer rows.Close()

	var villas = make([]models.Villa, 0)
	for rows.Next() {
		var v models.Villa
		err := rows.Scan(
			&v.ID,
			&v.Name,
			&v.Slug,
			&v.Location,
			&v.Description,
			&v.Image,
			&v.TitleTag,
			&v.MetaDesc,
			&v.CreatedAt,
			&v.UpdatedAt,
		)
		if err != nil {
			return c.Status(500).JSON(models.ResponseWithMessage{
				Success: false,
				Message: "Scan failed: " + err.Error(),
			})
		}
		villas = append(villas, v)
	}

	return c.JSON(models.ResponseWithData[[]models.Villa]{
		Success: true,
		Data:    villas,
	})

}

func GetSingleVilla(c *fiber.Ctx) error {
	id := c.Params("id")
	villaID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Invalid villa ID",
		})
	}
	query := `
		SELECT 
			v.id, v.name, v.slug, v.location, v.description,
			v.image, v.title_tag, v.meta_desc, v.created_at, v.updated_at,
			g.id, g.villa_id, g.image_url, g.sort_order, g.created_at, g.updated_at
		FROM villas v
		LEFT JOIN galleries g ON g.villa_id = v.id
		WHERE v.id = $1
		ORDER BY g.sort_order ASC
	`

	rows, err := db.Conn.Query(c.Context(), query, villaID)
	if err != nil {
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: err.Error(),
		})
	}
	defer rows.Close()

	var villa models.Villa
	galleries := []models.Gallery{}
	villa.Galleries = &galleries

	for rows.Next() {
		var g models.Gallery
		var gID *int64
		var gVillaID *int64

		err := rows.Scan(
			&villa.ID, &villa.Name, &villa.Slug, &villa.Location, &villa.Description,
			&villa.Image, &villa.TitleTag, &villa.MetaDesc, &villa.CreatedAt, &villa.UpdatedAt,
			&gID, &gVillaID, &g.ImageURL, &g.SortOrder, &g.CreatedAt, &g.UpdatedAt,
		)
		if err != nil {
			return c.Status(500).JSON(models.ResponseWithMessage{
				Success: false,
				Message: "Scan failed: " + err.Error(),
			})
		}

		// Handle if there's no gallery
		if gID != nil {
			if gVillaID != nil {
				g.VillaID = *gVillaID
			}

			g.ID = int64(*gID)
			*villa.Galleries = append(*villa.Galleries, g)
		}
	}

	// If villa was not found
	if villa.ID == 0 {
		return c.Status(404).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Villa not found",
		})
	}

	return c.JSON(models.ResponseWithData[models.Villa]{
		Success: true,
		Data:    villa,
	})
}


func PostVilla(c *fiber.Ctx) error {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Invalid multipart form",
		})
	}

	// Get form fields
	villa := models.Villa{
		Name:        c.FormValue("name"),
		Slug:        c.FormValue("slug"),
		Location:    c.FormValue("location"),
		Description: c.FormValue("description"),
		TitleTag:    c.FormValue("title_tag"),
		MetaDesc:    c.FormValue("meta_desc"),
	}

	// Upload image
	if imageFile, err := c.FormFile("image"); err == nil {
    	src, err := imageFile.Open()
		if err != nil {
			fmt.Println("Insert error:", err)
			return c.Status(500).JSON(models.ResponseWithMessage{
				Success: false,
				Message: "Failed to open image",
			})
		}
		defer src.Close()

		imageURL, err := imagekit.UploadImageFromForm(src, imageFile.Filename)
		if err != nil {
			fmt.Println("Insert error:", err)
			return c.Status(500).JSON(models.ResponseWithMessage{
				Success: false,
				Message: "Image upload failed: " + err.Error(),
			})
		}
		villa.Image = &imageURL
	}

	var imageVal interface{}
	if villa.Image != nil {
		imageVal = *villa.Image
	} else {
		imageVal = nil
	}

	// Step 1: Insert villa first
	err = db.Conn.QueryRow(c.Context(), `
		INSERT INTO villas (
			name, slug, location, description,
			image, title_tag, meta_desc
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`,
		villa.Name, villa.Slug, villa.Location, villa.Description,
		imageVal, villa.TitleTag, villa.MetaDesc,
	).Scan(&villa.ID)
	if err != nil {
		fmt.Println("Insert error:", err)
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Insert failed: " + err.Error(),
		})
	}

	// Step 2: Upload gallery AFTER villa ID is available
	galleryFiles := form.File["gallery"]
	if len(galleryFiles) != 0 {
		for i, file := range galleryFiles {
			order := i + 1
			src, err := file.Open()
			if err != nil {
				continue
			}
			defer src.Close()
	
			url, err := imagekit.UploadImageFromForm(src, file.Filename)
			if err != nil {
				continue
			}
	
			_, err = db.Conn.Exec(c.Context(), `
				INSERT INTO galleries (villa_id, image_url, sort_order)
				VALUES ($1, $2, $3)
			`, villa.ID, url, order)
			if err != nil {
				continue
			}
		}
	}

	return c.Status(200).JSON(models.ResponseWithMessage{
		Success: true,
		Message: "Villa added successfully",
	})

}

func UpdateVilla(c *fiber.Ctx) error {
	id := c.Params("id")
	VillaID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Invalid villa ID",
		})
	}
	villa := models.Villa{
		Name:        c.FormValue("name"),
		Slug:        c.FormValue("slug"),
		Location:    c.FormValue("location"),
		Description: c.FormValue("description"),
		TitleTag:    c.FormValue("title_tag"),
		MetaDesc:    c.FormValue("meta_desc"),
	}

	imageURL := c.Locals("image")
	if imageURL != nil {
		*villa.Image = imageURL.(string)
	}

	query := `
		UPDATE villas
		SET name = $1, slug = $2, location = $3, description = $4,
		    title_tag = $5, meta_desc = $6`
	args := []interface{}{
		villa.Name, villa.Slug, villa.Location, villa.Description,
		villa.TitleTag, villa.MetaDesc,
	}

	if villa.Image != nil && *villa.Image != "" {
		query += `, image = $7 WHERE id = $8`
		args = append(args, *villa.Image, VillaID)
	} else {
		query += ` WHERE id = $7`
		args = append(args, VillaID)
	}

	// Execute update
	_, err = db.Conn.Exec(c.Context(), query, args...)
	if err != nil {
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Update failed: " + err.Error(),
		})
	}

	return c.Status(200).JSON(models.ResponseWithMessage{
		Success: true,
		Message: "Villa updated successfully",
	})
}

func DeleteVilla(c *fiber.Ctx) error {
	id := c.Params("id")
	villaID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Invalid villa ID",
		})
	}

	// 1. Delete related galleries
	_, err = db.Conn.Exec(c.Context(), `DELETE FROM galleries WHERE villa_id = $1`, villaID)
	if err != nil {
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Failed to delete related galleries: " + err.Error(),
		})
	}

	// 2. Delete villa
	_, err = db.Conn.Exec(c.Context(), `DELETE FROM villas WHERE id = $1`, villaID)
	if err != nil {
		return c.Status(500).JSON(models.ResponseWithMessage{
			Success: false,
			Message: "Failed to delete villa: " + err.Error(),
		})
	}

	return c.Status(200).JSON(models.ResponseWithMessage{
		Success: true,
		Message: "Villa deleted successfully",
	})
}




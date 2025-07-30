package test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"project-0/handlers"
)



func setupApp() *fiber.App {
	app := fiber.New()
	app.Get("/villas", handlers.GetVillas)
	app.Get("/villas/:id", handlers.GetSingleVilla)
	app.Post("/villas", handlers.PostVilla)
	app.Put("/villas/:id", handlers.UpdateVilla)
	app.Delete("/villas/:id", handlers.DeleteVilla)
	return app
}

func TestGetVillas(t *testing.T) {
	assert := assert.New(t)
	app := setupApp()
	req := httptest.NewRequest("GET", "/villas", nil)
	resp, err := app.Test(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestGetSingleVilla_NotFound(t *testing.T) {
	assert := assert.New(t)
	app := setupApp()
	req := httptest.NewRequest("GET", "/villas/999999", nil)
	resp, err := app.Test(req)
	assert.NoError(err)
	assert.Equal(http.StatusNotFound, resp.StatusCode)
}

func TestPostVilla_Invalid(t *testing.T) {
	assert := assert.New(t)
	app := setupApp()
	req := httptest.NewRequest("POST", "/villas", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	resp, err := app.Test(req)
	assert.NoError(err)
	assert.Equal(http.StatusBadRequest, resp.StatusCode)
}

func TestPostVilla_Valid_WithoutImage(t *testing.T) {
	assert := assert.New(t)
	app := setupApp()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("name", "Test Villa")
	_ = writer.WriteField("slug", "test-villa")
	_ = writer.WriteField("location", "Bali")
	_ = writer.WriteField("description", "A test villa")
	_ = writer.WriteField("title_tag", "Test Title")
	_ = writer.WriteField("meta_desc", "Meta Description")

	writer.Close()

	req := httptest.NewRequest("POST", "/villas", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := app.Test(req)

	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
}


func TestUpdateVilla_InvalidID(t *testing.T) {
	assert := assert.New(t)
	app := setupApp()
	req := httptest.NewRequest("PUT", "/villas/invalid-id", strings.NewReader(""))
	resp, err := app.Test(req)
	assert.NoError(err)
	assert.Equal(http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteVilla_InvalidID(t *testing.T) {
	assert := assert.New(t)
	app := setupApp()
	req := httptest.NewRequest("DELETE", "/villas/invalid-id", nil)
	resp, err := app.Test(req)
	assert.NoError(err)
	assert.Equal(http.StatusBadRequest, resp.StatusCode)
}

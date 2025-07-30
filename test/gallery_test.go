package test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGalleryVilla(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/gallery/2", nil)
	resp, err := App.Test(req)
	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestAddGallery(t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	assert := assert.New(t)
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, "/api/gallery", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := App.Test(req)
	assert.NoError(err)
	assert.Equal(http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteGalleryByID_ValidID(t *testing.T) {
	data := map[string][]int64{"ids": {13}}
	payload, _ := json.Marshal(data)
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodPost, "/api/gallery/delete", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := App.Test(req, -1)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestDeleteGalleryByID_InvalidID(t *testing.T) {
	data := map[string][]int64{"ids": {0}}
	payload, _ := json.Marshal(data)
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodPost, "/api/gallery/delete", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := App.Test(req, -1)
	assert.NoError(err)
	assert.Equal(http.StatusNotFound, resp.StatusCode)
}
package models

import "time"

type Gallery struct {
	ID        	int64     `json:"id"`
	VillaID   	int64     `json:"villa_id"`
	ImageURL  	*string    `json:"image_url"`
	SortOrder 	*int       `json:"sort_order"`
	CreatedAt 	*time.Time `json:"created_at"`
	UpdatedAt 	*time.Time `json:"updated_at"`
}

type IncomingGallery struct {
	ImageURL string `json:"image_url"`
}

type AddGalleryRequest struct {
	VillaID   int64             `json:"villa_id"`
	Galleries []IncomingGallery `json:"galleries"`
}
package models

import (
	"time"
)

type Villa struct {
	ID          int      			`json:"id"`
	Name        string   			`json:"name"`
	Slug        string   			`json:"slug"`
	Location    string   			`json:"location"`
	Description string   			`json:"description"`
	Image       *string   			`json:"image"`
	TitleTag    string   			`json:"title_tag"`
	MetaDesc    string   			`json:"meta_desc"`
	Galleries   *[]Gallery		 	`json:"galleries,omitempty"`
	Prices		*[]Price			`json:"prices,omitempty"`
	CreatedAt 	time.Time 			`json:"created_at"`
	UpdatedAt 	time.Time 			`json:"updated_at"`
}


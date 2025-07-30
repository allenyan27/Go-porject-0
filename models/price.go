package models

import "time"

type Price struct {
	ID        		int64     	`json:"id"`
	VillaID   		int64     	`json:"villa_id"`
	StartDate	  	time.Time   `json:"start_date"`
	SortOrder 		time.Time   `json:"end_date"`
	Rate 			float64    	`json:"rate"`
	CreatedAt 		time.Time 	`json:"created_at"`
	UpdatedAt 		time.Time 	`json:"updated_at"`
}
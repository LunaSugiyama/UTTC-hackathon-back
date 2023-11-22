package model

import (
	"time"
)

type ItemImage struct {
	ID               int       `json:"id"`
	ItemID           int       `json:"item_id"`
	ItemCategoriesID int       `json:"item_categories_id"`
	Images           string    `json:"images"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

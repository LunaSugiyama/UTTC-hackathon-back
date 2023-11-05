package model

import (
	"time"
)

type StarredItem struct {
	ID               int       `json:"id"`
	UserFirebaseUID  string    `json:"user_firebase_uid"`
	ItemID           int       `json:"item_id"`
	ItemCategoriesID int       `json:"item_categories_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

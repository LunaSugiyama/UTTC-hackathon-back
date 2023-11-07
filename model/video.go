package model

import (
	"time"
)

type Video struct {
	ID                 int       `json:"id"`
	UserFirebaseUID    string    `json:"user_id"`
	Title              string    `json:"title"`
	Author             string    `json:"author"`
	Link               string    `json:"link"`
	Likes              int       `json:"likes"`
	ItemCategoriesID   int       `json:"item_categories_id" default:"3"`
	ItemCategoriesName string    `json:"item_categories_name" default:"videos"`
	Explanation        string    `json:"explanation"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	CurriculumIDs      []int     `json:"curriculum_ids"` // Added field for curriculum IDs
	Images             []string  `json:"images"`
}

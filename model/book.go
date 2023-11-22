package model

import (
	"time"
)

type Book struct {
	ID                 int       `json:"id"`
	UserFirebaseUID    string    `json:"user_firebase_uid" db:"user_firebase_uid"`
	Title              string    `json:"title"`
	Author             string    `json:"author"`
	Link               string    `json:"link"`
	Likes              int       `json:"likes"`
	ItemCategoriesID   int       `json:"item_categories_id" default:"2"`
	ItemCategoriesName string    `json:"item_categories_name" default:"books"`
	Explanation        string    `json:"explanation"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	CurriculumIDs      []int     `json:"curriculum_ids"` // Added field for curriculum IDs
	Images             []string  `json:"images"`
}

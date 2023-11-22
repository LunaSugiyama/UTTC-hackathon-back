// blog.go

package model

import (
	"time"
)

// User represents the structure of the 'blogs' table in the database.
type Blog struct {
	ID                 int       `json:"id" db:"id"`
	UserFirebaseUID    string    `json:"user_firebase_uid" db:"user_firebase_uid"`
	Title              string    `json:"title" db:"title"`
	Author             string    `json:"author" db:"author"`
	Link               string    `json:"link" db:"link"`
	Likes              int       `json:"likes" db:"likes"`
	ItemCategoriesID   int       `json:"item_categories_id" db:"item_categories_id" default:"1"`
	ItemCategoriesName string    `json:"item_categories_name" default:"blogs"`
	Explanation        string    `json:"explanation" db:"explanation"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	CurriculumIDs      []int     `json:"curriculum_ids"` // Added field for curriculum IDs
	Images             []string  `json:"images"`
}

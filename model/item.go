// item.go

package model

import (
	"time"
)

// User represents the structure of the 'blogs' table in the database.
type Item struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	Title            string    `json:"title"`
	Author           string    `json:"author"`
	Link             string    `json:"link"`
	Likes            int       `json:"likes"`
	ItemCategoriesID int       `json:"item_categories_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
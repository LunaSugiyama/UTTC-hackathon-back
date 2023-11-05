package model

import (
	"database/sql"
	"errors"
	"time"
	"uttc-hackathon/database"
)

// ErrItemCategoryNotFound is a custom error type for "not found" errors related to ItemCategory.
var ErrItemCategoryNotFound = errors.New("ItemCategory not found")

// ItemCategory represents the 'item_categories' table structure
type ItemCategory struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetItemCategory retrieves an ItemCategory by ID from the database.
func GetItemCategory(id int) (*ItemCategory, error) {
	// Create a new ItemCategory instance to store the result.
	itemCategory := &ItemCategory{}

	// Prepare a SQL query to select the item category by its ID.
	query := "SELECT id, name, created_at, updated_at FROM item_categories WHERE id = ?"

	// Execute the query and scan the result into the itemCategory struct.
	err := database.DB.QueryRow(query, id).Scan(&itemCategory.ID, &itemCategory.Name, &itemCategory.CreatedAt, &itemCategory.UpdatedAt)

	// Handle potential errors
	if err == sql.ErrNoRows {
		// Return the custom "not found" error
		return nil, ErrItemCategoryNotFound
	} else if err != nil {
		// Handle other database errors
		return nil, err
	}

	return itemCategory, nil
}

package database

import (
	"fmt"
)

func CreateStarredItemsTable() {
	// Create a table for the 'starred_items' table with the updated structure
	CreateStarredItemsTableSQL := `
	CREATE TABLE starred_items (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_firebase_uid VARCHAR(255),
		item_id INT,
		item_categories_id INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_firebase_uid) REFERENCES users(firebase_uid),
		FOREIGN KEY (item_categories_id) REFERENCES item_categories(id)
	)`

	// Execute the SQL statement to create the 'starred_items' table
	_, err := DB.Exec(CreateStarredItemsTableSQL)
	if err != nil {
		fmt.Println("Error creating 'starred_items' table:", err)
		return
	}

	fmt.Println("'starred_items' table created successfully with new columns.")
}

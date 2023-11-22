package database

import (
	"fmt"
)

func CreateLikedItemsTable() {
	// Create a table for the 'liked_items' table with the updated structure
	CreateLikedItemsTableSQL := `
	CREATE TABLE liked_items (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_firebase_uid VARCHAR(255),
		item_id INT,
		item_categories_id INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_firebase_uid) REFERENCES users(firebase_uid),
		FOREIGN KEY (item_categories_id) REFERENCES item_categories(id)
	)`

	// Execute the SQL statement to create the 'liked_items' table
	_, err := DB.Exec(CreateLikedItemsTableSQL)
	if err != nil {
		fmt.Println("Error creating 'liked_items' table:", err)
		return
	}

	fmt.Println("'liked_items' table created successfully with new columns.")
}

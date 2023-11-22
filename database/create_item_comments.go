package database

import (
	"fmt"
)

func CreateItemCommentsTable() {
	// Create a table for the 'item_comments' table
	createItemCommentsTableSQL := `
	CREATE TABLE item_comments (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_firebase_uid VARCHAR(255) NOT NULL,
		item_id INT NOT NULL,
		item_categories_id INT NOT NULL,
		comment TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_firebase_uid) REFERENCES users(firebase_uid),
		FOREIGN KEY (item_categories_id) REFERENCES item_categories(id)
	)`

	// Execute the SQL statement to create the 'item_comments' table
	_, err := DB.Exec(createItemCommentsTableSQL)
	if err != nil {
		fmt.Println("Error creating 'item_comments' table:", err)
		return
	}

	fmt.Println("'item_comments' table created successfully.")
}

func DropItemCommentsTable() {
	// Define the SQL statement to drop the 'item_comments' table
	dropItemCommentsTableSQL := "DROP TABLE IF EXISTS item_comments"

	// Execute the SQL statement to drop the 'item_comments' table
	_, err := DB.Exec(dropItemCommentsTableSQL)
	if err != nil {
		fmt.Println("Error dropping 'item_comments' table:", err)
		return
	}

	fmt.Println("'item_comments' table dropped successfully.")
}

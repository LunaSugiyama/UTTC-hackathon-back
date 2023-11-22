package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func CreateItemImagesTable() {
	// Create a table for the 'item_images' table with the updated structure
	CreateItemImagesTableSQL := `
	CREATE TABLE item_images (
		id INT AUTO_INCREMENT PRIMARY KEY,
		item_id INT, 
		item_categories_id INT,
		images VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (item_categories_id) REFERENCES item_categories(id)
	)`

	// Execute the SQL statement to create the 'item_images' table
	_, err := DB.Exec(CreateItemImagesTableSQL)
	if err != nil {
		fmt.Println("Error creating 'item_images' table:", err)
		return
	}

	fmt.Println("'item_images' table created successfully with new columns.")
}

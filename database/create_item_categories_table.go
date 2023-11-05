package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func CreateItemCategoriesTable() {
	// Create a table for the 'items' table with the updated structure
	CreateItemCategoriesTableSQL := `
        CREATE TABLE item_categories (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        )`

	// Execute the SQL statement to create the 'items' table
	_, err := DB.Exec(CreateItemCategoriesTableSQL)
	if err != nil {
		fmt.Println("Error creating 'item_categories' table:", err)
		return
	}

	fmt.Println("'item_categories' table created successfully with new columns.")
}

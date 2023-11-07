package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func CreateItemCurriculumsTable() {
	// Create a table for the 'item_curriculums' table with the updated structure
	CreateItemCurriculumsTableSQL := `
	CREATE TABLE item_curriculums (
		id INT AUTO_INCREMENT PRIMARY KEY,
		item_id INT, 
		item_categories_id INT,
		curriculum_id INT,
		FOREIGN KEY (curriculum_id) REFERENCES curriculums(id),
		FOREIGN KEY (item_categories_id) REFERENCES item_categories(id)
	)`

	// Execute the SQL statement to create the 'item_curriculums' table
	_, err := DB.Exec(CreateItemCurriculumsTableSQL)
	if err != nil {
		fmt.Println("Error creating 'item_curriculums' table:", err)
		return
	}

	fmt.Println("'item_curriculums' table created successfully with new columns.")
}

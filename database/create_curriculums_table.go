package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func CreateCurriculumsTable() {
	// Create a table for the 'curriculums' table with the updated structure
	CreateCurriculumsTableSQL := `
        CREATE TABLE curriculums (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        )`

	// Execute the SQL statement to create the 'curriculums' table
	_, err := DB.Exec(CreateCurriculumsTableSQL)
	if err != nil {
		fmt.Println("Error creating 'curriculums' table:", err)
		return
	}

	fmt.Println("'curriculums' table created successfully with new columns.")
}

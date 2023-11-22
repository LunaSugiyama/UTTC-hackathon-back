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

func PopulateCurriculumsTable() {
	curriculums := []string{
		"エディタ(IDE)",
		"OSコマンド(とシェル)",
		"Git",
		"Github",
		"HTML & CSS",
		"Javascript",
		"React",
		"React x Typescript",
		"SQL",
		"Docker",
		"Go",
		"HTTP Server (Go)",
		"RDBMS(MySQL)へ接続(Go)",
		"Unit Test(Go)",
		"フロントエンドとバックエンドの接続",
		"CI (Continuous Integration)",
		"CD (Continuous Deployment)",
		"認証",
		"ハッカソンの準備",
		"ハッカソンの概要"}

	for _, curriculum := range curriculums {
		insertSQL := `INSERT INTO curriculums (name) VALUES (?)`
		_, err := DB.Exec(insertSQL, curriculum)
		if err != nil {
			fmt.Printf("Error inserting curriculum %s: %v\n", curriculum, err)
		} else {
			fmt.Printf("Curriculum '%s' added successfully.\n", curriculum)
		}
	}
}

// DropCurriculumsTable drops the 'curriculums' table if it exists
func DropCurriculumsTable() {
	// Drop the 'curriculums' table if it exists
	DropCurriculumsTableSQL := `
        DROP TABLE IF EXISTS curriculums
    `

	// Execute the SQL statement to drop the 'curriculums' table
	_, err := DB.Exec(DropCurriculumsTableSQL)
	if err != nil {
		fmt.Println("Error dropping 'curriculums' table:", err)
		return
	}

	fmt.Println("'curriculums' table dropped successfully.")
}

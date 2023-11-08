package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func CreateBooksTable() {

	// Create a table for the 'users' table
	createBooksTableSQL := `
	create table books (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_firebase_uid VARCHAR(255) NOT NULL,
		title CHAR(200),
		author CHAR(200),
		link CHAR(200),
		likes INT DEFAULT 0,
		item_categories_id INT NOT NULL,
		explanation TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_firebase_uid) REFERENCES users(firebase_uid),
		FOREIGN KEY (item_categories_id) REFERENCES item_categories(id)
	)`

	// Execute the SQL statement to create the 'users' table
	_, err := DB.Exec(createBooksTableSQL)
	if err != nil {
		fmt.Println("Error creating 'Books' table:", err)
		return
	}

	fmt.Println("'Books' table created successfully.")
}

// You can define similar functions for creating other tables (blog, book, video, etc.) as needed.
func DropBooksTable() {
	// Define the SQL statement to drop the 'books' table
	dropBooksTableSQL := "DROP TABLE IF EXISTS books"

	// Execute the SQL statement to drop the 'books' table
	_, err := DB.Exec(dropBooksTableSQL)
	if err != nil {
		fmt.Println("Error dropping 'books' table:", err)
		return
	}

	fmt.Println("'books' table dropped successfully.")
}

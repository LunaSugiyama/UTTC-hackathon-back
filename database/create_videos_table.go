package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func CreateVideosTable() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	// mysqlUser := os.Getenv("MYSQL_USER")
	// mysqlPwd := os.Getenv("MYSQL_PASSWORD")
	// mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	mysqlUser := "test_user"
	mysqlPwd := "password"
	mysqlDatabase := "test_hackathon"

	connStr := fmt.Sprintf("%s:%s@tcp(%s:3307)/%s", mysqlUser, mysqlPwd, "localhost", mysqlDatabase)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Create a table for the 'users' table
	createVideosTableSQL := `
	create table videos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_firebase_uid VARCHAR(255) NOT NULL,
		title CHAR(200) NOT NULL,
		author CHAR(200) NOT NULL,
		link CHAR(200) NOT NULL,
		likes INT DEFAULT 0,
		item_categories_id INT NOT NULL,
		explanation TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_firebase_uid) REFERENCES users(firebase_uid),
		FOREIGN KEY (item_categories_id) REFERENCES item_categories(id)
	)`

	// Execute the SQL statement to create the 'users' table
	_, err = db.Exec(createVideosTableSQL)
	if err != nil {
		fmt.Println("Error creating 'Videos' table:", err)
		return
	}

	fmt.Println("'Videos' table created successfully.")
}

// You can define similar functions for creating other tables (blog, book, video, etc.) as needed.
func DropVideosTable() {
	// Define the SQL statement to drop the 'videos' table
	dropVideosTableSQL := "DROP TABLE IF EXISTS videos"

	// Execute the SQL statement to drop the 'videos' table
	_, err := DB.Exec(dropVideosTableSQL)
	if err != nil {
		fmt.Println("Error dropping 'videos' table:", err)
		return
	}

	fmt.Println("'videos' table dropped successfully.")
}

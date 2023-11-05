package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func CreateUsersTable() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

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

	// Create a table for the 'users' table with the updated structure
	CreateUsersTableSQL :=
		`CREATE TABLE users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			firebase_uid VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			password VARCHAR(255) NULL,
			age INT NOT NULL,
			username VARCHAR(255) NOT NULL,
			email VARCHAR(255),
			authority INT(1) DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_firebase_uid (firebase_uid)  -- Add an index to the firebase_uid column
		)
		`

	// Execute the SQL statement to create the 'users' table
	_, err = db.Exec(CreateUsersTableSQL)
	if err != nil {
		fmt.Println("Error creating 'users' table:", err)
		return
	}

	fmt.Println("'users' table created successfully with new columns.")
}

func AddFirebaseUID() {
	// Alter the 'users' table to add the 'firebase_uid' column
	AlterTableSQL := `
        ALTER TABLE users
        ADD COLUMN firebase_uid VARCHAR(255)`

	_, err := DB.Exec(AlterTableSQL)
	if err != nil {
		fmt.Println("Error adding 'firebase_uid' column:", err)
		return
	}

	fmt.Println("'firebase_uid' column added to 'users' table.")

}
func PasswordColumnToNull() {
	// SQL statement to drop the 'password' column
	dropSQL := `ALTER TABLE users DROP COLUMN password;`
	_, err := DB.Exec(dropSQL)
	if err != nil {
		fmt.Println("Error dropping 'password' column:", err)
		return
	}

	// SQL statement to add the 'password' column back as nullable
	addSQL := `ALTER TABLE users ADD COLUMN password VARCHAR(255) NULL;`
	_, err = DB.Exec(addSQL)
	if err != nil {
		fmt.Println("Error adding 'password' column back:", err)
		return
	}

	fmt.Println("'password' column now accepts null in the 'users' table.")
}

func DropUsersTable() {
	// Define the SQL statement to drop the 'users' table
	dropUsersTableSQL := "DROP TABLE IF EXISTS users"

	// Execute the SQL statement to drop the 'users' table
	_, err := DB.Exec(dropUsersTableSQL)
	if err != nil {
		fmt.Println("Error dropping 'users' table:", err)
		return
	}

	fmt.Println("'users' table dropped successfully.")
}

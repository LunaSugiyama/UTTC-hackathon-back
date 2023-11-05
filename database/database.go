package database

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

func InitializeDB() {
	mysqlUser := "test_user"
	mysqlPwd := "password"
	mysqlDatabase := "test_hackathon"

	connStr := fmt.Sprintf("%s:%s@tcp(%s:3307)/%s", mysqlUser, mysqlPwd, "localhost", mysqlDatabase)
	var err error
	DB, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	fmt.Printf("Database connection successful.")
}

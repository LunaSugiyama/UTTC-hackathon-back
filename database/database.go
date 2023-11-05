package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var DB *sql.DB

func InitializeDB() {
	// mysqlUser := "test_user"
	// mysqlPwd := "password"
	// mysqlDatabase := "test_hackathon"

	// connStr := fmt.Sprintf("%s:%s@tcp(%s:3307)/%s", mysqlUser, mysqlPwd, "localhost", mysqlDatabase)
	// DB接続のための準備
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPwd := os.Getenv("MYSQL_PWD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	DB = db

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	fmt.Printf("Database connection successful.")
}

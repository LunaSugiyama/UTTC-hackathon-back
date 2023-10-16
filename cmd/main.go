package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type BlogEntry struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Link      string `json:"link"`
	CreatedAt time.Time
}

func main() {
	r := gin.Default()

	// connect to db
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPwd := os.Getenv("MYSQL_PWD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
	db, err := sql.Open("mysql", connStr)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r.POST("/addBlog", func(c *gin.Context) {
		var blog BlogEntry
		if err := c.ShouldBindJSON(&blog); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insert a new blog entry into the database
		query := "INSERT INTO blog (id, user_id, link, created_at) VALUES (?, ?, ?, ?"
		id := generateUniqueID()
		_, err := db.Exec(query, id, blog.UserID, blog.Link, time.Now())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Blog entry added successfully"})
	})

	r.Run(":8080") // Replace with the port you want to run the server on
}

func generateUniqueID() string {
	// Implement your own logic to generate a unique ID
	return "your_unique_id"
}

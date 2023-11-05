package book

import (
	"log"
	"net/http"
	"time"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func CreateBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert a new book entry into the database
	query := "INSERT INTO books (user_firebase_uid, title, author, link, item_categories_id, explanation, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, dbErr := database.DB.Exec(query, book.UserFirebaseUID, book.Title, book.Author, book.Link, book.ItemCategoriesID, book.Explanation, time.Now())
	if dbErr != nil {
		log.Printf("Error inserting book entry into the database: %v", dbErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert book entry"})
		return
	}

	// Get the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last inserted ID"})
		return
	}

	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	// Query the database to retrieve the inserted book entry
	selectQuery := "SELECT * FROM books WHERE id = ?"
	err = database.DB.QueryRow(selectQuery, lastInsertID).Scan(&book.ID, &book.UserFirebaseUID, &book.Title, &book.Author, &book.Link, &book.Likes, &book.ItemCategoriesID, &book.Explanation, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Error retrieving the created book entry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created book entry"})
		return
	}
	book.CreatedAt = createdAt.Time
	book.UpdatedAt = updatedAt.Time

	// Now, insert rows into the item_curriculums table for each curriculum ID
	for _, curriculumID := range book.CurriculumIDs {
		insertItemCurriculumQuery := "INSERT INTO item_curriculums (item_id, item_categories_id, curriculum_id) VALUES (?, ?, ?)"
		_, insertErr := database.DB.Exec(insertItemCurriculumQuery, book.ID, book.ItemCategoriesID, curriculumID)
		if insertErr != nil {
			log.Printf("Error inserting into item_curriculums table: %v", insertErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert into item_curriculums table"})
			return
		}
	}

	c.JSON(http.StatusOK, book)
}

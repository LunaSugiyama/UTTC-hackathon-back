package book

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func UpdateBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the book with the given ID exists
	checkQuery := "SELECT id FROM books WHERE id = ?"
	var id int
	err := database.DB.QueryRow(checkQuery, book.ID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// The book doesn't exist, return a not found response
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Book with id %d not found", book.ID)})
			return
		} else {
			// Handle other database errors
			log.Printf("Error checking for book: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check book"})
			return
		}
	}

	// Validate that 'title', 'author', 'link', and 'user_firebase_uid' fields are not empty
	if book.Title == "" || book.Author == "" || book.Link == "" || book.UserFirebaseUID == "" || book.ItemCategoriesID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, Author, Link, UserID, and ItemCategoriesID are required fields"})
		return
	}

	// Update the book in the database
	updateQuery := "UPDATE books SET title = ?, author = ?, link = ?, user_firebase_uid = ?, item_categories_id = ?, explanation = ? WHERE id = ?"
	_, dbErr := database.DB.Exec(updateQuery, book.Title, book.Author, book.Link, book.UserFirebaseUID, book.ItemCategoriesID, book.Explanation, book.ID)
	if dbErr != nil {
		log.Printf("Error updating book in the database: %v", dbErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	// Update the curriculum IDs for the book in the item_curriculums table
	// First, delete the existing entries for this book
	deleteQuery := "DELETE FROM item_curriculums WHERE item_id = ? AND item_categories_id = ?"
	_, deleteErr := database.DB.Exec(deleteQuery, book.ID, book.ItemCategoriesID)
	if deleteErr != nil {
		log.Printf("Error deleting existing curriculum entries: %v", deleteErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update curriculum IDs"})
		return
	}

	// Now, insert the updated curriculum IDs
	for _, curriculumID := range book.CurriculumIDs {
		insertItemCurriculumQuery := "INSERT INTO item_curriculums (item_id, item_categories_id, curriculum_id) VALUES (?, ?, ?)"
		_, insertErr := database.DB.Exec(insertItemCurriculumQuery, book.ID, book.ItemCategoriesID, curriculumID)
		if insertErr != nil {
			log.Printf("Error inserting into item_curriculums table: %v", insertErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update curriculum IDs"})
			return
		}
	}

	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	// Query the database to retrieve the updated book entry
	selectQuery := "SELECT * FROM books WHERE id = ?"
	err = database.DB.QueryRow(selectQuery, book.ID).Scan(&book.ID, &book.UserFirebaseUID, &book.Title, &book.Author, &book.Link, &book.Likes, &book.ItemCategoriesID, &book.Explanation, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Error retrieving the updated book entry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated book entry"})
		return
	}

	book.CreatedAt = createdAt.Time
	book.UpdatedAt = updatedAt.Time
	c.JSON(http.StatusOK, book)
}

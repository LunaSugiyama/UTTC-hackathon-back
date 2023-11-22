package book

import (
	"database/sql"
	"log"
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
)

func DeleteBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the book entry exists before attempting to delete it
	checkQuery := "SELECT id FROM books WHERE id = ?"
	var id int
	err := database.DB.QueryRow(checkQuery, book.ID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// The book entry doesn't exist, return an error
			c.JSON(http.StatusNotFound, gin.H{"error": "Book entry not found"})
			return
		} else {
			// Handle other database errors
			log.Printf("Error checking for book entry: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check book entry"})
			return
		}
	}

	// Delete the book entry from the database
	deleteQuery := "DELETE FROM books WHERE id = ?"
	_, dbErr := database.DB.Exec(deleteQuery, book.ID)
	if dbErr != nil {
		log.Printf("Error deleting book entry from the database: %v", dbErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book entry"})
		return
	}

	// Delete related entries from the item_curriculums table
	deleteItemCurriculumQuery := "DELETE FROM item_curriculums WHERE item_id = ? AND item_categories_id = ?"
	_, deleteErr := database.DB.Exec(deleteItemCurriculumQuery, book.ID, book.ItemCategoriesID)
	if deleteErr != nil {
		log.Printf("Error deleting related item_curriculums entries: %v", deleteErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related item_curriculums entries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book entry deleted successfully"})
}

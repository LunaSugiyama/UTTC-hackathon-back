package blog

import (
	"database/sql"
	"log"
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
)

func DeleteBlog(c *gin.Context) {
	var blog model.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the blog entry exists before attempting to delete it
	checkQuery := "SELECT id FROM blogs WHERE id = ?"
	var id int
	err := database.DB.QueryRow(checkQuery, blog.ID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// The blog entry doesn't exist, return an error
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog entry not found"})
			return
		} else {
			// Handle other database errors
			log.Printf("Error checking for blog entry: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check blog entry"})
			return
		}
	}

	// Delete the blog entry from the database
	deleteQuery := "DELETE FROM blogs WHERE id = ?"
	_, dbErr := database.DB.Exec(deleteQuery, blog.ID)
	if dbErr != nil {
		log.Printf("Error deleting blog entry from the database: %v", dbErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog entry"})
		return
	}

	// Delete related entries from the item_curriculums table
	deleteItemCurriculumQuery := "DELETE FROM item_curriculums WHERE item_id = ? AND item_categories_id = ?"
	_, deleteErr := database.DB.Exec(deleteItemCurriculumQuery, blog.ID, 1)
	if deleteErr != nil {
		log.Printf("Error deleting related item_curriculums entries: %v", deleteErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related item_curriculums entries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog entry deleted successfully"})
}

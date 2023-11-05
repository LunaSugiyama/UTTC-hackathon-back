package video

import (
	"database/sql"
	"log"
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
)

func DeleteVideo(c *gin.Context) {
	var video model.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the video entry exists before attempting to delete it
	checkQuery := "SELECT id FROM videos WHERE id = ?"
	var id int
	err := database.DB.QueryRow(checkQuery, video.ID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// The video entry doesn't exist, return an error
			c.JSON(http.StatusNotFound, gin.H{"error": "Video entry not found"})
			return
		} else {
			// Handle other database errors
			log.Printf("Error checking for video entry: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check video entry"})
			return
		}
	}

	// Delete the video entry from the database
	deleteQuery := "DELETE FROM videos WHERE id = ?"
	_, dbErr := database.DB.Exec(deleteQuery, video.ID)
	if dbErr != nil {
		log.Printf("Error deleting video entry from the database: %v", dbErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete video entry"})
		return
	}

	// Delete related entries from the item_curriculums table
	deleteItemCurriculumQuery := "DELETE FROM item_curriculums WHERE item_id = ? AND item_categories_id = ?"
	_, deleteErr := database.DB.Exec(deleteItemCurriculumQuery, video.ID, video.ItemCategoriesID)
	if deleteErr != nil {
		log.Printf("Error deleting related item_curriculums entries: %v", deleteErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related item_curriculums entries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video entry deleted successfully"})
}

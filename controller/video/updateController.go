package video

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

func UpdateVideo(c *gin.Context) {
	var video model.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the video with the given ID exists
	checkQuery := "SELECT id FROM videos WHERE id = ?"
	var id int
	err := database.DB.QueryRow(checkQuery, video.ID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// The video doesn't exist, return a not found response
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Video with id %d not found", video.ID)})
			return
		} else {
			// Handle other database errors
			log.Printf("Error checking for video: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check video"})
			return
		}
	}

	// Validate that 'title', 'link', 'author', 'user_firebase_uid', and 'item_categories_id' fields are not empty
	if video.Title == "" || video.Link == "" || video.Author == "" || video.UserFirebaseUID == "" || video.ItemCategoriesID == 0 || video.Explanation == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, Link, Author, UserID, Explanation and ItemCategoriesID are required fields"})
		return
	}

	// Update the video in the database
	updateQuery := "UPDATE videos SET title = ?, link = ?, author = ?, user_firebase_uid = ?, item_categories_id = ?, explanation = ? WHERE id = ?"
	_, dbErr := database.DB.Exec(updateQuery, video.Title, video.Link, video.Author, video.UserFirebaseUID, video.ItemCategoriesID, video.Explanation, video.ID)
	if dbErr != nil {
		log.Printf("Error updating video in the database: %v", dbErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update video"})
		return
	}

	// Update the curriculum IDs for the video in the item_curriculums table
	// First, delete the existing entries for this video
	deleteQuery := "DELETE FROM item_curriculums WHERE item_id = ? AND item_categories_id = ?"
	_, deleteErr := database.DB.Exec(deleteQuery, video.ID, video.ItemCategoriesID)
	if deleteErr != nil {
		log.Printf("Error deleting existing curriculum entries: %v", deleteErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update curriculum IDs"})
		return
	}

	// Now, insert the updated curriculum IDs
	for _, curriculumID := range video.CurriculumIDs {
		insertItemCurriculumQuery := "INSERT INTO item_curriculums (item_id, item_categories_id, curriculum_id) VALUES (?, ?, ?)"
		_, insertErr := database.DB.Exec(insertItemCurriculumQuery, video.ID, video.ItemCategoriesID, curriculumID)
		if insertErr != nil {
			log.Printf("Error inserting into item_curriculums table: %v", insertErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update curriculum IDs"})
			return
		}
	}

	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	// Query the database to retrieve the updated video entry
	selectQuery := "SELECT * FROM videos WHERE id = ?"
	err = database.DB.QueryRow(selectQuery, video.ID).Scan(&video.ID, &video.UserFirebaseUID, &video.Title, &video.Link, &video.Author, &video.Likes, &video.ItemCategoriesID, &video.Explanation, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Error retrieving the updated video entry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated video entry"})
		return
	}
	video.CreatedAt = createdAt.Time
	video.UpdatedAt = updatedAt.Time
	c.JSON(http.StatusOK, video)
}

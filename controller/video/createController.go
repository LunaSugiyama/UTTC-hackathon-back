package video

import (
	"log"
	"net/http"
	"time"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func CreateVideo(c *gin.Context) {
	var video model.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert a new video entry into the database
	query := "INSERT INTO videos (user_firebase_uid, title, link, author, item_categories_id, explanation, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, dbErr := database.DB.Exec(query, video.UserFirebaseUID, video.Title, video.Link, video.Author, video.ItemCategoriesID, video.Explanation, time.Now())
	if dbErr != nil {
		log.Printf("Error inserting video entry into the database: %v", dbErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert video entry"})
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
	// Query the database to retrieve the inserted video entry
	selectQuery := "SELECT * FROM videos WHERE id = ?"
	err = database.DB.QueryRow(selectQuery, lastInsertID).Scan(&video.ID, &video.UserFirebaseUID, &video.Title, &video.Link, &video.Author, &video.Likes, &video.ItemCategoriesID, &video.Explanation, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Error retrieving the created video entry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created video entry"})
		return
	}

	// Now, insert rows into the item_curriculums table for each curriculum ID
	for _, curriculumID := range video.CurriculumIDs {
		insertItemCurriculumQuery := "INSERT INTO item_curriculums (item_id, item_categories_id, curriculum_id) VALUES (?, ?, ?)"
		_, insertErr := database.DB.Exec(insertItemCurriculumQuery, video.ID, video.ItemCategoriesID, curriculumID)
		if insertErr != nil {
			log.Printf("Error inserting into item_curriculums table: %v", insertErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert into item_curriculums table"})
			return
		}
	}

	video.CreatedAt = createdAt.Time
	video.UpdatedAt = updatedAt.Time
	c.JSON(http.StatusOK, video)
}

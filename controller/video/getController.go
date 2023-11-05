package video

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func GetVideo(c *gin.Context) {
	var video model.Video
	// Get the 'id' parameter from the query string
	videoID := c.Query("id")

	// Check if 'id' is empty or not provided
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is missing or empty"})
		return
	}

	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	var curriculumIDs string

	// Convert the videoID to an integer
	id, err := strconv.Atoi(videoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	// Query the database to retrieve the video entry by ID, including associated curriculum IDs
	query := `
        SELECT v.*, icat.name, GROUP_CONCAT(ic.curriculum_id) AS curriculum_ids
        FROM videos AS v
        LEFT JOIN item_curriculums AS ic ON v.id = ic.item_id AND v.item_categories_id = ic.item_categories_id
        LEFT JOIN item_categories AS icat ON v.item_categories_id = icat.id
		WHERE v.id = ?
        GROUP BY v.id
    `
	err = database.DB.QueryRow(query, id).Scan(
		&video.ID, &video.UserFirebaseUID, &video.Title, &video.Link, &video.Author, &video.Likes, &video.ItemCategoriesID, &video.Explanation,
		&createdAt, &updatedAt, &video.ItemCategoriesName, &curriculumIDs,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// The video entry doesn't exist, return a not found response
			c.JSON(http.StatusNotFound, gin.H{"error": "Video entry not found"})
			return
		} else {
			// Handle other database errors
			log.Printf("Error retrieving video entry: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve video entry"})
			return
		}
	}
	video.CreatedAt = createdAt.Time

	// Split the curriculum IDs into a slice
	curriculumIDSlice := strings.Split(curriculumIDs, ",")
	for _, idStr := range curriculumIDSlice {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			video.CurriculumIDs = append(video.CurriculumIDs, id)
		}
	}

	// Video entry found, return it as JSON
	c.JSON(http.StatusOK, video)
}

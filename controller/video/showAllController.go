package video

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func ShowAllVideos(c *gin.Context) {
	// Query the database to retrieve all video entries with their associated curriculum IDs
	query := `
        SELECT v.*, GROUP_CONCAT(ic.curriculum_id) AS curriculum_ids
        FROM videos AS v
        LEFT JOIN item_curriculums AS ic ON v.id = ic.item_id AND v.item_categories_id = ic.item_categories_id
        GROUP BY v.id
    `
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error querying video entries: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve video entries"})
		return
	}
	defer rows.Close()

	var videos []model.Video

	for rows.Next() {
		var video model.Video
		var createdAt mysql.NullTime // Use mysql.NullTime for MySQL TIMESTAMP columns
		var updatedAt mysql.NullTime // Use mysql.NullTime for MySQL TIMESTAMP columns
		var curriculumIDs string

		if err := rows.Scan(
			&video.ID, &video.UserFirebaseUID, &video.Title, &video.Link, &video.Author, &video.Likes, &video.ItemCategoriesID, &video.Explanation,
			&createdAt, &updatedAt, &curriculumIDs,
		); err != nil {
			log.Printf("Error scanning video entry: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan video entries"})
			return
		}
		video.CreatedAt = createdAt.Time
		video.UpdatedAt = updatedAt.Time

		// Split the curriculum IDs into a slice
		curriculumIDSlice := strings.Split(curriculumIDs, ",")
		for _, idStr := range curriculumIDSlice {
			id, err := strconv.Atoi(idStr)
			if err == nil {
				video.CurriculumIDs = append(video.CurriculumIDs, id)
			}
		}

		videos = append(videos, video)
	}

	// Return the list of videos as JSON
	c.JSON(http.StatusOK, videos)
}

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
	var Images string

	// Convert the videoID to an integer
	id, err := strconv.Atoi(videoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	// Query the database to retrieve the video entry by ID, including associated curriculum IDs
	query := `
	SELECT b.*, icat.name, GROUP_CONCAT(ic.curriculum_id) AS curriculum_ids, GROUP_CONCAT(IFNULL(ii.images, '')) AS images
	FROM videos AS b
	LEFT JOIN item_categories AS icat ON b.item_categories_id = icat.id
	LEFT JOIN (
		SELECT item_id, item_categories_id, GROUP_CONCAT(images) AS images
		FROM item_images
		GROUP BY item_id, item_categories_id
	) AS ii ON b.id = ii.item_id AND b.item_categories_id = ii.item_categories_id
	LEFT JOIN (
		SELECT item_id, item_categories_id, GROUP_CONCAT(curriculum_id) AS curriculum_id
		FROM item_curriculums
		GROUP BY item_id, item_categories_id
	) AS ic ON b.id = ic.item_id AND b.item_categories_id = ic.item_categories_id
	WHERE b.id = ?
	GROUP BY b.id
	
    `
	err = database.DB.QueryRow(query, id).Scan(
		&video.ID, &video.UserFirebaseUID, &video.Title, &video.Link, &video.Author, &video.Likes, &video.ItemCategoriesID, &video.Explanation,
		&createdAt, &updatedAt, &video.ItemCategoriesName, &curriculumIDs, &Images,
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

	// Split the images into a slice
	imagesSlice := strings.Split(Images, ",")
	video.Images = append(video.Images, imagesSlice...)

	// Video entry found, return it as JSON
	c.JSON(http.StatusOK, video)
}

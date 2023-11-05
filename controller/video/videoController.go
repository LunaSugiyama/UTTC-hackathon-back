package video

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"time"
// 	"uttc-hackathon/database"
// 	"uttc-hackathon/model"

// 	"github.com/gin-gonic/gin"
// 	"github.com/go-sql-driver/mysql"
// )

// func CreateVideo(c *gin.Context) {
// 	var video model.Video
// 	if err := c.ShouldBindJSON(&video); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Insert a new video entry into the database
// 	query := "INSERT INTO videos (user_id, title, link, author, item_categories_id, created_at) VALUES (?, ?, ?, ?, ?, ?)"
// 	result, dbErr := database.DB.Exec(query, video.UserID, video.Title, video.Link, video.Author, video.ItemCategoriesID, time.Now())
// 	if dbErr != nil {
// 		log.Printf("Error inserting video entry into the database: %v", dbErr)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert video entry"})
// 		return
// 	}

// 	// Get the last inserted ID
// 	lastInsertID, err := result.LastInsertId()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last inserted ID"})
// 		return
// 	}

// 	var createdAt mysql.NullTime
// 	var updatedAt mysql.NullTime
// 	// Query the database to retrieve the inserted video entry
// 	selectQuery := "SELECT * FROM videos WHERE id = ?"
// 	err = database.DB.QueryRow(selectQuery, lastInsertID).Scan(&video.ID, &video.UserID, &video.Title, &video.Link, &video.Author, &video.Likes, &video.ItemCategoriesID, &createdAt, &updatedAt)
// 	if err != nil {
// 		log.Printf("Error retrieving the created video entry: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created video entry"})
// 		return
// 	}

// 	// Now, insert rows into the item_curriculums table for each curriculum ID
// 	for _, curriculumID := range video.CurriculumIDs {
// 		insertItemCurriculumQuery := "INSERT INTO item_curriculums (item_id, item_categories_id, curriculum_id) VALUES (?, ?, ?)"
// 		_, insertErr := database.DB.Exec(insertItemCurriculumQuery, video.ID, video.ItemCategoriesID, curriculumID)
// 		if insertErr != nil {
// 			log.Printf("Error inserting into item_curriculums table: %v", insertErr)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert into item_curriculums table"})
// 			return
// 		}
// 	}

// 	video.CreatedAt = createdAt.Time
// 	video.UpdatedAt = updatedAt.Time
// 	c.JSON(http.StatusOK, video)
// }

// func UpdateVideo(c *gin.Context) {
// 	var video model.Video
// 	if err := c.ShouldBindJSON(&video); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Check if the video with the given ID exists
// 	checkQuery := "SELECT id FROM videos WHERE id = ?"
// 	var id int
// 	err := database.DB.QueryRow(checkQuery, video.ID).Scan(&id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			// The video doesn't exist, return a not found response
// 			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Video with id %d not found", video.ID)})
// 			return
// 		} else {
// 			// Handle other database errors
// 			log.Printf("Error checking for video: %v", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check video"})
// 			return
// 		}
// 	}

// 	// Validate that 'title', 'link', 'author', 'user_id', and 'item_categories_id' fields are not empty
// 	if video.Title == "" || video.Link == "" || video.Author == "" || video.UserID == 0 || video.ItemCategoriesID == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, Link, Author, UserID, and ItemCategoriesID are required fields"})
// 		return
// 	}

// 	// Update the video in the database
// 	updateQuery := "UPDATE videos SET title = ?, link = ?, author = ?, user_id = ?, item_categories_id = ? WHERE id = ?"
// 	_, dbErr := database.DB.Exec(updateQuery, video.Title, video.Link, video.Author, video.UserID, video.ItemCategoriesID, video.ID)
// 	if dbErr != nil {
// 		log.Printf("Error updating video in the database: %v", dbErr)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update video"})
// 		return
// 	}

// 	// Update the curriculum IDs for the video in the item_curriculums table
// 	// First, delete the existing entries for this video
// 	deleteQuery := "DELETE FROM item_curriculums WHERE item_id = ? AND item_categories_id = ?"
// 	_, deleteErr := database.DB.Exec(deleteQuery, video.ID, video.ItemCategoriesID)
// 	if deleteErr != nil {
// 		log.Printf("Error deleting existing curriculum entries: %v", deleteErr)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update curriculum IDs"})
// 		return
// 	}

// 	// Now, insert the updated curriculum IDs
// 	for _, curriculumID := range video.CurriculumIDs {
// 		insertItemCurriculumQuery := "INSERT INTO item_curriculums (item_id, item_categories_id, curriculum_id) VALUES (?, ?, ?)"
// 		_, insertErr := database.DB.Exec(insertItemCurriculumQuery, video.ID, video.ItemCategoriesID, curriculumID)
// 		if insertErr != nil {
// 			log.Printf("Error inserting into item_curriculums table: %v", insertErr)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update curriculum IDs"})
// 			return
// 		}
// 	}

// 	var createdAt mysql.NullTime
// 	var updatedAt mysql.NullTime
// 	// Query the database to retrieve the updated video entry
// 	selectQuery := "SELECT * FROM videos WHERE id = ?"
// 	err = database.DB.QueryRow(selectQuery, video.ID).Scan(&video.ID, &video.UserID, &video.Title, &video.Link, &video.Author, &video.Likes, &video.ItemCategoriesID, &createdAt, &updatedAt)
// 	if err != nil {
// 		log.Printf("Error retrieving the updated video entry: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated video entry"})
// 		return
// 	}

// 	video.CreatedAt = createdAt.Time
// 	video.UpdatedAt = updatedAt.Time
// 	c.JSON(http.StatusOK, video)
// }

// func DeleteVideo(c *gin.Context) {
// 	var video model.Video
// 	if err := c.ShouldBindJSON(&video); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Check if the video entry exists before attempting to delete it
// 	checkQuery := "SELECT id FROM videos WHERE id = ?"
// 	var id int
// 	err := database.DB.QueryRow(checkQuery, video.ID).Scan(&id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			// The video entry doesn't exist, return an error
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Video entry not found"})
// 			return
// 		} else {
// 			// Handle other database errors
// 			log.Printf("Error checking for video entry: %v", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check video entry"})
// 			return
// 		}
// 	}

// 	// Delete the video entry from the database
// 	deleteQuery := "DELETE FROM videos WHERE id = ?"
// 	_, dbErr := database.DB.Exec(deleteQuery, video.ID)
// 	if dbErr != nil {
// 		log.Printf("Error deleting video entry from the database: %v", dbErr)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete video entry"})
// 		return
// 	}

// 	// Delete related entries from the item_curriculums table
// 	deleteItemCurriculumQuery := "DELETE FROM item_curriculums WHERE item_id = ? AND item_categories_id = ?"
// 	_, deleteErr := database.DB.Exec(deleteItemCurriculumQuery, video.ID, video.ItemCategoriesID)
// 	if deleteErr != nil {
// 		log.Printf("Error deleting related item_curriculums entries: %v", deleteErr)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related item_curriculums entries"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Video entry deleted successfully"})
// }

// func GetVideo(c *gin.Context) {
// 	var video model.Video
// 	// Get the 'id' parameter from the query string
// 	videoID := c.Query("id")

// 	// Check if 'id' is empty or not provided
// 	if videoID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is missing or empty"})
// 		return
// 	}

// 	var createdAt mysql.NullTime
// 	var updatedAt mysql.NullTime
// 	var curriculumIDs string

// 	// Convert the videoID to an integer
// 	id, err := strconv.Atoi(videoID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
// 		return
// 	}

// 	// Query the database to retrieve the video entry by ID, including associated curriculum IDs
// 	query := `
//         SELECT v.*, GROUP_CONCAT(ic.curriculum_id) AS curriculum_ids
//         FROM videos AS v
//         LEFT JOIN item_curriculums AS ic ON v.id = ic.item_id AND v.item_categories_id = ic.item_categories_id
//         WHERE v.id = ?
//         GROUP BY v.id
//     `
// 	err = database.DB.QueryRow(query, id).Scan(
// 		&video.ID, &video.UserID, &video.Title, &video.Link, &video.Author, &video.Likes, &video.ItemCategoriesID,
// 		&createdAt, &updatedAt, &curriculumIDs,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			// The video entry doesn't exist, return a not found response
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Video entry not found"})
// 			return
// 		} else {
// 			// Handle other database errors
// 			log.Printf("Error retrieving video entry: %v", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve video entry"})
// 			return
// 		}
// 	}
// 	video.CreatedAt = createdAt.Time

// 	// Split the curriculum IDs into a slice
// 	curriculumIDSlice := strings.Split(curriculumIDs, ",")
// 	for _, idStr := range curriculumIDSlice {
// 		id, err := strconv.Atoi(idStr)
// 		if err == nil {
// 			video.CurriculumIDs = append(video.CurriculumIDs, id)
// 		}
// 	}

// 	// Video entry found, return it as JSON
// 	c.JSON(http.StatusOK, video)
// }

// func ShowAllVideos(c *gin.Context) {
// 	// Query the database to retrieve all video entries with their associated curriculum IDs
// 	query := `
//         SELECT v.*, GROUP_CONCAT(ic.curriculum_id) AS curriculum_ids
//         FROM videos AS v
//         LEFT JOIN item_curriculums AS ic ON v.id = ic.item_id AND v.item_categories_id = ic.item_categories_id
//         GROUP BY v.id
//     `
// 	rows, err := database.DB.Query(query)
// 	if err != nil {
// 		log.Printf("Error querying video entries: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve video entries"})
// 		return
// 	}
// 	defer rows.Close()

// 	var videos []model.Video

// 	for rows.Next() {
// 		var video model.Video
// 		var createdAt mysql.NullTime // Use mysql.NullTime for MySQL TIMESTAMP columns
// 		var updatedAt mysql.NullTime // Use mysql.NullTime for MySQL TIMESTAMP columns
// 		var curriculumIDs string

// 		if err := rows.Scan(
// 			&video.ID, &video.UserID, &video.Title, &video.Link, &video.Author, &video.Likes, &video.ItemCategoriesID,
// 			&createdAt, &updatedAt, &curriculumIDs,
// 		); err != nil {
// 			log.Printf("Error scanning video entry: %v", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan video entries"})
// 			return
// 		}
// 		video.CreatedAt = createdAt.Time
// 		video.UpdatedAt = updatedAt.Time

// 		// Split the curriculum IDs into a slice
// 		curriculumIDSlice := strings.Split(curriculumIDs, ",")
// 		for _, idStr := range curriculumIDSlice {
// 			id, err := strconv.Atoi(idStr)
// 			if err == nil {
// 				video.CurriculumIDs = append(video.CurriculumIDs, id)
// 			}
// 		}

// 		videos = append(videos, video)
// 	}

// 	// Return the list of videos as JSON
// 	c.JSON(http.StatusOK, videos)
// }

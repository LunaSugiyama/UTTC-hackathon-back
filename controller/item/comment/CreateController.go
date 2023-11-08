package comment

import (
	"net/http"
	"time"
	"uttc-hackathon/model"

	"uttc-hackathon/database" // Import your database package

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// CreateItemComment handles the creation of a new item comment.
func CreateItemComment(c *gin.Context) {
	// Define a struct to bind the request data
	var comment model.Comment

	// Bind the request data to the struct
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// You can perform additional validation and business logic here if needed
	if comment.UserFirebaseUID == "" || comment.ItemID == 0 || comment.ItemCategoriesID == 0 || comment.Comment == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	// Create a new item comment in the database
	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	insertQuery := "INSERT INTO item_comments (user_firebase_uid, item_id, item_categories_id, comment, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := database.DB.Exec(insertQuery, comment.UserFirebaseUID, comment.ItemID, comment.ItemCategoriesID, comment.Comment, time.Now(), time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	selectQuery := "SELECT * FROM item_comments WHERE user_firebase_uid = ? AND item_id = ? AND item_categories_id = ? AND comment = ?"
	err = database.DB.QueryRow(selectQuery, comment.UserFirebaseUID, comment.ItemID, comment.ItemCategoriesID, comment.Comment).Scan(&comment.ID, &comment.UserFirebaseUID, &comment.ItemID, &comment.ItemCategoriesID, &comment.Comment, &createdAt, &updatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Assign the converted time values to the comment struct
	comment.CreatedAt = createdAt.Time
	comment.UpdatedAt = updatedAt.Time

	c.JSON(http.StatusCreated, comment)
}

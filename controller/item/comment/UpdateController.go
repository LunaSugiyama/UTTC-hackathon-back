package comment

import (
	"net/http"
	"uttc-hackathon/model"

	"uttc-hackathon/database" // Import your database package

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// CreateItemComment handles the creation of a new item comment.
func UpdateItemComment(c *gin.Context) {
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
	var updatedAt mysql.NullTime

	updateQuery := "UPDATE item_comments SET comment = ?, updated_at = ? WHERE user_firebase_uid = ? AND item_id = ? AND item_categories_id = ?"
	_, err := database.DB.Exec(updateQuery, comment.Comment, updatedAt, comment.UserFirebaseUID, comment.ItemID, comment.ItemCategoriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Assign the converted time values to the comment struct
	comment.UpdatedAt = updatedAt.Time

	c.JSON(http.StatusOK, comment)
}

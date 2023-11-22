package comment

import (
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
)

func DeleteItemComment(c *gin.Context) {
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if comment.UserFirebaseUID == "" || comment.ItemID == 0 || comment.ItemCategoriesID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	deleteQuery := "DELETE FROM item_comments WHERE user_firebase_uid = ? AND item_id = ? AND item_categories_id = ?"
	result, err := database.DB.Exec(deleteQuery, comment.UserFirebaseUID, comment.ItemID, comment.ItemCategoriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check the number of rows affected by the delete operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found or not authorized to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

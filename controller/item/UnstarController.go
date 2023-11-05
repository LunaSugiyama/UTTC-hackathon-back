package item

import (
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
)

func UnstarItem(c *gin.Context) {
	var starredItem model.StarredItem
	if err := c.ShouldBindJSON(&starredItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the item to unstar exists in the "starred_items" table
	if !recordExists(starredItem.UserFirebaseUID, starredItem.ItemID, starredItem.ItemCategoriesID) {
		// Item does not exist, return an error
		c.JSON(http.StatusNotFound, gin.H{"error": "Item is not starred"})
		return
	}

	// Remove the record from the "starred_items" table
	_, err := database.DB.Exec("DELETE FROM starred_items WHERE user_firebase_uid = ? AND item_id = ? AND item_categories_id = ?", starredItem.UserFirebaseUID, starredItem.ItemID, starredItem.ItemCategoriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item unstarred successfully"})
}

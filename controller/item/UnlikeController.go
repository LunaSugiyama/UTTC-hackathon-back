package item

import (
	"fmt"
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
)

func UnlikeItem(c *gin.Context) {
	var likedItem model.LikedItem
	if err := c.ShouldBindJSON(&likedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("likeditem", likedItem)

	// Check if the item to unlike exists in the "liked_items" table
	if !recordLikedExists(likedItem.UserFirebaseUID, likedItem.ItemID, likedItem.ItemCategoriesID) {
		// Item does not exist, return an error
		c.JSON(http.StatusNotFound, gin.H{"error": "Item is not liked"})
		return
	}

	// Remove the record from the "liked_items" table
	_, err := database.DB.Exec("DELETE FROM liked_items WHERE user_firebase_uid = ? AND item_id = ? AND item_categories_id = ?", likedItem.UserFirebaseUID, likedItem.ItemID, likedItem.ItemCategoriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item unliked successfully"})
}

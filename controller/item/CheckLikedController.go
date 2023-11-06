package item

import (
	"fmt"
	"net/http"
	"strconv"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
)

func CheckLiked(c *gin.Context) {
	var likedItem model.LikedItem

	itemIDStr := c.Query("item_id")
	fmt.Println(itemIDStr)
	itemID, err := strconv.Atoi(itemIDStr)
	fmt.Println(itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item_id"})
		return
	}

	itemCategoriesIDStr := c.Query("item_categories_id")
	itemCategoriesID, err := strconv.Atoi(itemCategoriesIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item_categories_id"})
		return
	}

	likedItem.ItemID = itemID
	likedItem.ItemCategoriesID = itemCategoriesID
	likedItem.UserFirebaseUID = c.Query("user_firebase_uid")

	// if err := c.ShouldBindJSON(&likedItem); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	fmt.Println("likeditem", likedItem)

	var isLiked bool
	query := "SELECT COUNT(*) FROM liked_items WHERE item_id = ? AND item_categories_id = ?"
	err = database.DB.QueryRow(query, likedItem.ItemID, likedItem.ItemCategoriesID).Scan(&isLiked)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if isLiked {
		c.JSON(http.StatusOK, gin.H{"isLiked": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"isLiked": false})
	}
}

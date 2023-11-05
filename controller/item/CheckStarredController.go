package item

import (
	"fmt"
	"net/http"
	"strconv"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
)

// IsItemStarred checks if an item is starred for a user
func IsItemStarred(c *gin.Context) {
	var starredItem model.StarredItem
	starredItem.UserFirebaseUID = c.Query("user_firebase_uid")

	itemIDStr := c.Query("item_id")
	itemID, err := strconv.Atoi(itemIDStr)
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

	starredItem.ItemID = itemID
	starredItem.ItemCategoriesID = itemCategoriesID

	fmt.Println(starredItem)

	// Assuming that you have access to a database connection (e.g., database.DB), you can query the database to check if the item is starred.
	// Here, we use a query similar to what you have in the `recordExists` function.

	var isStarred bool
	query := "SELECT COUNT(*) FROM starred_items WHERE user_firebase_uid = ? AND item_id = ? AND item_categories_id = ?"
	err = database.DB.QueryRow(query, starredItem.UserFirebaseUID, starredItem.ItemID, starredItem.ItemCategoriesID).Scan(&isStarred)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if isStarred {
		c.JSON(http.StatusOK, gin.H{"isStarred": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"isStarred": false})
	}
}

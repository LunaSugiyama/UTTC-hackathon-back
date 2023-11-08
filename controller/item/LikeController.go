package item

import (
	"fmt"
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
)

func LikeItem(c *gin.Context) {
	var likedItem model.LikedItem
	if err := c.ShouldBindJSON(&likedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the name from the "item_categories" table
	itemName, err := getItemCategoryName(likedItem.ItemCategoriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if the name exists
	if itemName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Item category with ID %d does not exist", likedItem.ItemCategoriesID)})
		return
	}

	// Check if the item exists in the table with the retrieved name
	if recordLikedExists(likedItem.UserFirebaseUID, likedItem.ItemID, likedItem.ItemCategoriesID) {
		// Item already exists, return an "already starred" message
		c.JSON(http.StatusConflict, gin.H{"message": "Item is already liked"})
		return
	}

	getItemCategoriesQuery := "SELECT name FROM item_categories WHERE id = ?"
	err = database.DB.QueryRow(getItemCategoriesQuery, likedItem.ItemCategoriesID).Scan(&itemName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(itemName)

	// Update the "likes" column in the "items" table for the liked item
	query := fmt.Sprintf("UPDATE %s SET likes = likes + 1 WHERE id = ?", itemName)
	_, err = database.DB.Exec(query, likedItem.ItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Insert a new record into the "starred_items" table
	_, err = database.DB.Exec("INSERT INTO liked_items (user_firebase_uid, item_id, item_categories_id) VALUES (?, ?, ?)", likedItem.UserFirebaseUID, likedItem.ItemID, likedItem.ItemCategoriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item liked successfully"})

}

func recordLikedExists(userFirebaseUID string, itemID int, itemCategoriesID int) bool {
	var count int
	query := "SELECT COUNT(*) FROM liked_items WHERE user_firebase_uid = ? AND item_id = ? AND item_categories_id = ?"
	err := database.DB.QueryRow(query, userFirebaseUID, itemID, itemCategoriesID).Scan(&count)
	if err != nil {
		// Handle any errors here, such as database connection issues
		return false
	}
	return count > 0
}

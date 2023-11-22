package item

import (
	"fmt"
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
)

func StarItem(c *gin.Context) {
	var starredItem model.StarredItem
	if err := c.ShouldBindJSON(&starredItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the name from the "item_categories" table
	itemName, err := getItemCategoryName(starredItem.ItemCategoriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if the name exists
	if itemName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Item category with ID %d does not exist", starredItem.ItemCategoriesID)})
		return
	}

	// Check if the item exists in the table with the retrieved name
	if recordExists(starredItem.UserFirebaseUID, starredItem.ItemID, starredItem.ItemCategoriesID) {
		// Item already exists, return an "already starred" message
		c.JSON(http.StatusConflict, gin.H{"message": "Item is already starred"})
		return
	}

	// Insert a new record into the "starred_items" table
	_, err = database.DB.Exec("INSERT INTO starred_items (user_firebase_uid, item_id, item_categories_id) VALUES (?, ?, ?)", starredItem.UserFirebaseUID, starredItem.ItemID, starredItem.ItemCategoriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item starred successfully"})
}

// Retrieve the name from the "item_categories" table
func getItemCategoryName(itemCategoriesID int) (string, error) {
	var itemName string
	err := database.DB.QueryRow("SELECT name FROM item_categories WHERE id = ?", itemCategoriesID).Scan(&itemName)
	if err != nil {
		return "", err
	}
	return itemName, nil
}

// Check if the item with a specific name and itemID exists in a table with a dynamic name
func recordExists(UserFirebaseUID string, itemID, itemCategoriesID int) bool {
	var count int
	query := "SELECT COUNT(*) FROM starred_items WHERE user_firebase_uid = ? AND item_id = ? AND item_categories_id = ?"
	err := database.DB.QueryRow(query, UserFirebaseUID, itemID, itemCategoriesID).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

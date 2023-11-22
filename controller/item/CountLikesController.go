package item

import (
	"net/http"
	"strconv"
	"uttc-hackathon/database"

	"github.com/gin-gonic/gin"
)

func CountLikes(c *gin.Context) {
	var itemID int
	var itemCategoriesID int

	itemIDStr := c.Query("item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item_id"})
		return
	}

	itemCategoriesIDStr := c.Query("item_categories_id")
	itemCategoriesID, err = strconv.Atoi(itemCategoriesIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item_categories_id"})
		return
	}

	var count int
	query := "SELECT COUNT(*) FROM liked_items WHERE item_id = ? AND item_categories_id = ?"
	err = database.DB.QueryRow(query, itemID, itemCategoriesID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

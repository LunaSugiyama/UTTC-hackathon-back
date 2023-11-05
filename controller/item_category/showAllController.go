package itemcategory

import (
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// ShowAllItemCategories retrieves all item_categories from the database and returns them as JSON.
func ShowAllItemCategories(c *gin.Context) {
	query := "SELECT * FROM item_categories"
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve item_categories"})
		return
	}
	defer rows.Close()

	var item_categories []model.ItemCategory

	for rows.Next() {
		var item_category model.ItemCategory
		var createdAt, updatedAt mysql.NullTime
		if err := rows.Scan(&item_category.ID, &item_category.Name, &createdAt, &updatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan item_category data"})
			return
		}
		item_category.CreatedAt = createdAt.Time
		item_category.UpdatedAt = updatedAt.Time
		item_categories = append(item_categories, item_category)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error during row iteration"})
		return
	}

	c.JSON(http.StatusOK, item_categories)
}

package comment

import (
	"net/http"
	"strconv"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func ShowItemComment(c *gin.Context) {
	item_id := c.Query("item_id")
	item_categories_id := c.Query("item_categories_id")

	item_id_int, err := strconv.Atoi(item_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ItemID parameter"})
		return
	}

	item_categories_id_int, err := strconv.Atoi(item_categories_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ItemCategoriesID parameter"})
		return
	}

	if item_id_int == 0 || item_categories_id_int == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	selectQuery := "SELECT * FROM item_comments WHERE item_id = ? AND item_categories_id = ?"
	rows, err := database.DB.Query(selectQuery, item_id, item_categories_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close() // Close the rows after we're done with them

	createdAt := mysql.NullTime{}
	updatedAt := mysql.NullTime{}

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(&comment.ID, &comment.UserFirebaseUID, &comment.ItemID, &comment.ItemCategoriesID, &comment.Comment, &createdAt, &updatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		comment.CreatedAt = createdAt.Time
		comment.UpdatedAt = updatedAt.Time
		comments = append(comments, comment)
	}

	c.JSON(http.StatusOK, comments)
}

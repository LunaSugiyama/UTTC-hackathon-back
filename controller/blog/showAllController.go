package blog

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func ShowAllBlogs(c *gin.Context) {
	// Query the database to retrieve all blog entries with their associated curriculum IDs
	query := `
        SELECT b.*, GROUP_CONCAT(ic.curriculum_id) AS curriculum_ids
        FROM blogs AS b
        LEFT JOIN item_curriculums AS ic ON b.id = ic.item_id AND b.item_categories_id = ic.item_categories_id
        GROUP BY b.id
    `
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error querying blog entries: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blog entries"})
		return
	}
	defer rows.Close()

	var blogs []model.Blog

	for rows.Next() {
		var blog model.Blog
		var createdAt mysql.NullTime // Use mysql.NullTime for MySQL TIMESTAMP columns
		var updatedAt mysql.NullTime // Use mysql.NullTime for MySQL TIMESTAMP columns
		var curriculumIDs string

		if err := rows.Scan(
			&blog.ID, &blog.UserFirebaseUID, &blog.Title, &blog.Author, &blog.Link, &blog.Likes, &blog.ItemCategoriesID,
			&blog.Explanation, &createdAt, &updatedAt, &curriculumIDs,
		); err != nil {
			log.Printf("Error scanning blog entry: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan blog entries"})
			return
		}
		blog.CreatedAt = createdAt.Time
		blog.UpdatedAt = updatedAt.Time

		// Split the curriculum IDs into a slice
		curriculumIDSlice := strings.Split(curriculumIDs, ",")
		for _, idStr := range curriculumIDSlice {
			id, err := strconv.Atoi(idStr)
			if err == nil {
				blog.CurriculumIDs = append(blog.CurriculumIDs, id)
			}
		}

		blogs = append(blogs, blog)
	}

	c.JSON(http.StatusOK, blogs)
}

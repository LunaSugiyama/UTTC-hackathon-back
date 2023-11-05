package blog

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func GetBlog(c *gin.Context) {
	var blog model.Blog
	// Get the 'id' parameter from the query string
	blogID := c.Query("id")

	// Check if 'id' is empty or not provided
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is missing or empty"})
		return
	}

	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	var curriculumIDs string

	// Convert the blogID to an integer
	id, err := strconv.Atoi(blogID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	// Query the database to retrieve the blog entry by ID
	query := `
    SELECT blogs.*, item_categories.name, GROUP_CONCAT(item_curriculums.curriculum_id) AS curriculum_ids
    FROM blogs
    LEFT JOIN item_categories ON blogs.item_categories_id = item_categories.id
	LEFT JOIN item_curriculums ON blogs.id = item_curriculums.item_id AND blogs.item_categories_id = item_curriculums.item_categories_id
    WHERE blogs.id = ?
	GROUP BY blogs.id
`
	err = database.DB.QueryRow(query, id).Scan(
		&blog.ID, &blog.UserFirebaseUID, &blog.Title, &blog.Author, &blog.Link, &blog.Likes,
		&blog.ItemCategoriesID, &blog.Explanation, &createdAt, &updatedAt, &blog.ItemCategoriesName, &curriculumIDs,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// The blog entry doesn't exist, return a not found response
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog entry not found"})
			return
		} else {
			// Handle other database errors
			log.Printf("Error retrieving blog entry: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blog entry"})
			return
		}
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

	// Blog entry found, return it as JSON
	c.JSON(http.StatusOK, blog)
}

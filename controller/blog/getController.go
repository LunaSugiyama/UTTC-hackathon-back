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
	var Images string

	// Convert the blogID to an integer
	id, err := strconv.Atoi(blogID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	// Query the database to retrieve the blog entry by ID
	query := `
	SELECT b.*, icat.name, GROUP_CONCAT(ic.curriculum_id) AS curriculum_ids, GROUP_CONCAT(IFNULL(ii.images, '')) AS images
	FROM blogs AS b
	LEFT JOIN item_categories AS icat ON b.item_categories_id = icat.id
	LEFT JOIN (
		SELECT item_id, item_categories_id, GROUP_CONCAT(images) AS images
		FROM item_images
		GROUP BY item_id, item_categories_id
	) AS ii ON b.id = ii.item_id AND b.item_categories_id = ii.item_categories_id
	LEFT JOIN (
		SELECT item_id, item_categories_id, GROUP_CONCAT(curriculum_id) AS curriculum_id
		FROM item_curriculums
		GROUP BY item_id, item_categories_id
	) AS ic ON b.id = ic.item_id AND b.item_categories_id = ic.item_categories_id
	WHERE b.id = ?
	GROUP BY b.id
	
`
	err = database.DB.QueryRow(query, id).Scan(
		&blog.ID, &blog.UserFirebaseUID, &blog.Title, &blog.Author, &blog.Link, &blog.Likes,
		&blog.ItemCategoriesID, &blog.Explanation, &createdAt, &updatedAt, &blog.ItemCategoriesName, &curriculumIDs, &Images,
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

	imageslice := strings.Split(Images, ",")
	blog.Images = append(blog.Images, imageslice...)

	// Blog entry found, return it as JSON
	c.JSON(http.StatusOK, blog)
}

package blog

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func UpdateBlog(c *gin.Context) {
	var blog model.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the blog with the given ID exists
	checkQuery := "SELECT id FROM blogs WHERE id = ?"
	var id int
	err := database.DB.QueryRow(checkQuery, blog.ID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// The blog doesn't exist, return a not found response
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Blog with id %d not found", blog.ID)})
			return
		} else {
			// Handle other database errors
			log.Printf("Error checking for blog: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check blog"})
			return
		}
	}

	// Validate that 'title', 'author', 'link', 'user_firebase_uid', 'explanation', and 'item_categories_id' fields are not empty
	if blog.Title == "" || blog.Author == "" || blog.Link == "" || blog.UserFirebaseUID == "" || blog.ItemCategoriesID == 0 || blog.Explanation == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, Author, Link, UserID, Explanation, and ItemCategoriesID are required fields"})
		return
	}

	// Update the blog in the database
	updateQuery := "UPDATE blogs SET title = ?, author = ?, link = ?, user_firebase_uid = ?, item_categories_id = ?, explanation = ? WHERE id = ?"
	_, dbErr := database.DB.Exec(updateQuery, blog.Title, blog.Author, blog.Link, blog.UserFirebaseUID, blog.ItemCategoriesID, blog.Explanation, blog.ID)
	if dbErr != nil {
		log.Printf("Error updating blog in the database: %v", dbErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog"})
		return
	}

	// Update the curriculum IDs for the blog in the item_curriculums table
	// First, delete the existing entries for this blog
	deleteQuery := "DELETE FROM item_curriculums WHERE item_id = ? AND item_categories_id = ?"
	_, deleteErr := database.DB.Exec(deleteQuery, blog.ID, 1)
	if deleteErr != nil {
		log.Printf("Error deleting existing curriculum entries: %v", deleteErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update curriculum IDs"})
		return
	}

	// Now, insert the updated curriculum IDs
	for _, curriculumID := range blog.CurriculumIDs {
		insertItemCurriculumQuery := "INSERT INTO item_curriculums (item_id, item_categories_id, curriculum_id) VALUES (?, ?, ?)"
		_, insertErr := database.DB.Exec(insertItemCurriculumQuery, blog.ID, blog.ItemCategoriesID, curriculumID)
		if insertErr != nil {
			log.Printf("Error inserting into item_curriculums table: %v", insertErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update curriculum IDs"})
			return
		}
	}
	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	// Query the database to retrieve the updated blog entry
	selectQuery := "SELECT * FROM blogs WHERE id = ?"
	err = database.DB.QueryRow(selectQuery, blog.ID).Scan(&blog.ID, &blog.UserFirebaseUID, &blog.Title, &blog.Author, &blog.Link, &blog.Likes, &blog.ItemCategoriesID, &blog.Explanation, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Error retrieving the updated blog entry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated blog entry"})
		return
	}

	blog.CreatedAt = createdAt.Time
	blog.UpdatedAt = updatedAt.Time
	c.JSON(http.StatusOK, blog)
}

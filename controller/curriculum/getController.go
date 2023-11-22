package curriculum

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// GetCurriculumByID retrieves a curriculum by its ID from the database.
func GetCurriculum(c *gin.Context) {
	var curriculum model.Curriculum

	curriculumID := c.Query("id")

	// Check if 'id' is empty or not provided
	if curriculumID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is missing or empty"})
		return
	}

	// Convert the curriculumID to an integer
	id, err := strconv.Atoi(curriculumID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime

	// Prepare and execute the SQL query to fetch the curriculum by ID
	query := `SELECT * FROM curriculums WHERE id = ?`
	err = database.DB.QueryRow(query, id).Scan(&curriculum.ID, &curriculum.Name, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no curriculum is found, return a not found response
			c.JSON(http.StatusNotFound, gin.H{"error": "Curriculum not found"})
			return
		} else {
			// Handle other database errors
			log.Printf("Error retrieving curriculum: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve curriculum"})
			return
		}
	}

	curriculum.CreatedAt = createdAt.Time
	curriculum.UpdatedAt = updatedAt.Time

	// Curriculum found, return it as JSON
	c.JSON(http.StatusOK, curriculum)
}

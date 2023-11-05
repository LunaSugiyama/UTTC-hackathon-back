package curriculum

import (
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// ShowAllCurriculum retrieves all curriculums from the database and returns them as JSON.
func ShowAllCurriculums(c *gin.Context) {
	// Query all curriculums from the database
	query := "SELECT * FROM curriculums"
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve curriculums"})
		return
	}
	defer rows.Close()

	var curriculums []model.Curriculum

	// Iterate through the rows and scan the results into a slice
	for rows.Next() {
		var curriculum model.Curriculum
		var createdAt, updatedAt mysql.NullTime
		if err := rows.Scan(&curriculum.ID, &curriculum.Name, &createdAt, &updatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan curriculum data"})
			return
		}
		curriculum.CreatedAt = createdAt.Time
		curriculum.UpdatedAt = updatedAt.Time
		curriculums = append(curriculums, curriculum)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error during row iteration"})
		return
	}

	// Return the list of curriculums as JSON
	c.JSON(http.StatusOK, curriculums)
}

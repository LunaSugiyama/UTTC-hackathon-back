package curriculum

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
)

// CreateCurriculum adds a new curriculum to the database and returns it as JSON
func CreateCurriculum(c *gin.Context) {
	var curriculum model.Curriculum
	if err := c.ShouldBindJSON(&curriculum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if curriculum.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	// Create a new record with the current timestamp
	createdAt := time.Now()
	updatedAt := createdAt

	// Execute the SQL INSERT statement
	result, err := database.DB.Exec("INSERT INTO curriculums (name, created_at, updated_at) VALUES (?, ?, ?)", &curriculum.Name, createdAt, updatedAt)
	if err != nil {
		log.Printf("Error creating curriculum: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create curriculum"})
		return
	}

	// Retrieve the auto-generated ID of the newly inserted record
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last inserted ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve curriculum ID"})
		return
	}

	// Set the ID of the curriculum and return it as JSON
	curriculum.ID = int(lastInsertID)
	c.JSON(http.StatusCreated, curriculum)
}

// UpdateCurriculum updates a curriculum by its ID
func UpdateCurriculum(c *gin.Context) {
	// Parse the JSON request body to get the updated curriculum data
	var curriculum model.Curriculum // Define a struct for the request data
	if err := c.ShouldBindJSON(&curriculum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Execute the SQL UPDATE statement
	result, err := database.DB.Exec("UPDATE curriculums SET name=? WHERE id=?", curriculum.Name, curriculum.ID)
	if err != nil {
		log.Printf("Error updating curriculum: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update curriculum"})
		return
	}

	rowsAffected, _ := result.RowsAffected() // Get the number of affected rows

	if rowsAffected == 0 {
		// No rows were updated; the curriculum with the provided ID does not exist.
		c.JSON(http.StatusNotFound, gin.H{"error": "Curriculum not found or no updates necessary"})
		return
	}

	// Return the updated curriculum ID in the response
	c.JSON(http.StatusOK, curriculum)
}

// DeleteCurriculum deletes a curriculum by its ID
func DeleteCurriculum(c *gin.Context) {
	idStr := c.DefaultQuery("id", "0") // Get the ID from the URL parameter as a string with a default value of "0"
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Handle the error by returning a 400 Bad Request response.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Execute the SQL DELETE statement
	result, err := database.DB.Exec("DELETE FROM curriculums WHERE id=?", id)
	if err != nil {
		log.Printf("Error deleting curriculum: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete curriculum"})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		// No rows were deleted; the curriculum with the provided ID does not exist.
		c.JSON(http.StatusNotFound, gin.H{"error": "Curriculum not found"})
		return
	}

	// Return the deleted ID in the response
	c.JSON(http.StatusOK, gin.H{"message": "Curriculum deleted", "deleted_id": id})
}

package user

import (
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
)

func Update(c *gin.Context) {
	// Bind the JSON request body to the updateUserRequest struct
	var updateUserRequest model.User
	if err := c.ShouldBindJSON(&updateUserRequest); err != nil {
		// Handle error, for example:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the required parameters are set
	if updateUserRequest.FirebaseUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firebase UID is required"})
		return
	}

	// Check if the user exists in the "users" table
	if !userExists(updateUserRequest.FirebaseUID) {
		// User does not exist, return an error
		c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
		return
	}

	// Update the user data in the database with the new values
	err := updateUser(updateUserRequest.FirebaseUID, updateUserRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User data updated successfully"})
}

func updateUser(firebaseUID string, request model.User) error {
	query := "UPDATE users SET name=?, email=?, age=?, password=? WHERE firebase_uid=?"
	_, err := database.DB.Exec(query, request.Name, request.Email, request.Age, request.Password, firebaseUID)
	if err != nil {
		return err
	}

	return nil
}

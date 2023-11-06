package user

import (
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
)

func Show(c *gin.Context) {
	var user model.User
	FirebaseUID := c.DefaultQuery("user_firebase_uid", "")

	// Check if the required parameters are set
	if FirebaseUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firebase UID is required"})
		return
	}

	// Check if the user exists in the "users" table
	if !userExists(FirebaseUID) {
		// User does not exist, return an error
		c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
		return
	}

	user, err := getUserData(FirebaseUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	c.JSON(http.StatusOK, user)

}

func userExists(firebaseUID string) bool {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE firebase_uid = ?"
	err := database.DB.QueryRow(query, firebaseUID).Scan(&count)
	if err != nil {
		// Handle the error gracefully, don't panic
		return false
	}

	// Check if count is greater than zero to determine existence
	return count > 0
}

func getUserData(firebaseUID string) (model.User, error) {
	var user model.User
	query := "SELECT firebase_uid, name, email, age FROM users WHERE firebase_uid = ?"
	err := database.DB.QueryRow(query, firebaseUID).Scan(&user.FirebaseUID, &user.Name, &user.Email, &user.Age)
	if err != nil {
		return user, err
	}

	return user, nil
}

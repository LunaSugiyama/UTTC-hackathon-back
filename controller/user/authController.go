package user

import (
	"fmt"
	"log"
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type LoginData struct {
	UID     string `json:"uid"`
	IDToken string `json:"idToken"`
}

func Login(c *gin.Context) {
	var data LoginData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	uid := data.UID
	idToken := data.IDToken
	// fmt.Println(user)
	fmt.Println("token", idToken)
	fmt.Println("uid", uid)

	if idToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": idToken, "uid": uid})
}

func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if the required parameters are set
	if user.FirebaseUID == "" || user.Name == "" || user.Email == "" || user.Age == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	saveUserToSQLDatabase(user.FirebaseUID, user.Email, user.Email, user.Name, user.Age)

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func saveUserToSQLDatabase(firebaseUID, username, email, name string, age int) {
	insertUserSQL := `
        INSERT INTO users (firebase_uid, username, email, name, age, created_at) VALUES (?, ?, ?, ?, ?, NOW())`

	_, err := database.DB.Exec(insertUserSQL, firebaseUID, username, email, name, age)
	if err != nil {
		log.Printf("Error saving user to SQL database: %v", err)
		// You might want to return an error response here, or handle the error according to your application's logic.
		return
	}

	// User successfully inserted into the database
	log.Printf("User inserted into the database.")
}

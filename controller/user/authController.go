package user

import (
	"context"
	"log"
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/firebaseinit"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = []byte("your-secret-key")
var firebaseApp *firebase.App

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// func init() {
// 	opt := option.WithCredentialsFile("/home/denjo/ダウンロード/term4-luna-sugiyama-firebase-adminsdk-1joai-b0f371c4d8.json")
// 	app, err := firebase.NewApp(context.Background(), nil, opt)
// 	if err != nil {
// 		log.Fatalf("error initializing Firebase app: %v\n", err)
// 	}
// 	firebaseApp = app

// 	authClient, err = app.Auth(context.Background())
// 	if err != nil {
// 		log.Fatalf("error creating Firebase Auth client: %v\n", err)
// 	}
// }

func Login(c *gin.Context) {
	// // Parse the token from the client (you need to send the Firebase ID token from the client)
	// authHeader := c.Request.Header.Get("Authorization")
	// if authHeader == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is missing"})
	// 	return
	// }

	// // Check if the header starts with "Bearer " and extract the token part
	// parts := strings.Split(authHeader, " ")
	// if len(parts) != 2 || parts[0] != "Bearer" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
	// 	return
	// }

	// firebaseIDToken := parts[1]

	// // Verify the Firebase ID token on the server.
	// token, err := firebaseinit.AuthClient.VerifyIDToken(context.Background(), firebaseIDToken)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	// 	return
	// }

	// // Obtain the UID from the verified token
	// uid := token.UID

	// The AuthMiddleware will run before this handler, so you can access user information if authenticated
	user, userExists := c.Get("user")
	token, tokenExists := c.Get("token")

	if !userExists || !tokenExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user_id": user})
}

func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if the username already exists
	_, err := firebaseinit.AuthClient.GetUserByEmail(context.Background(), user.Username)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Create a new user in Firebase Authentication
	params := (&auth.UserToCreate{}).
		Email(user.Username).
		Password(user.Password)

	newUser, err := firebaseinit.AuthClient.CreateUser(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User registration failed: " + err.Error()})
		return
	}

	// Save user details to the local SQL database
	saveUserToSQLDatabase(newUser.UID, user.Username, user.Username, user.Name, user.Age)

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func VerifyPassword(password string, userRecord *auth.UserRecord) bool {
	// Implement your password verification logic here
	// You may compare the provided password with userRecord.PasswordHash
	// For security, consider using a library to securely hash and verify passwords
	return userRecord != nil
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

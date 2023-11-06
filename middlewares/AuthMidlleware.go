package middlewares

import (
	"context"
	"net/http"
	"strings"
	"uttc-hackathon/firebaseinit"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the "Authorization" header from the request
		authHeader := c.GetHeader("Authorization")

		// Check if the token is missing
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
			c.Abort()
			return
		}

		// Check if the header is in the correct format (Bearer token)
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Extract the token from the header
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Use the authClient from firebaseinit package
		authClient := firebaseinit.AuthClient

		// Verify the token
		tokenInfo, err := authClient.VerifyIDToken(context.Background(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set the user claims in the Gin context for further use in the handler
		c.Set("user", tokenInfo.UID) // Assuming UID is the user identifier in the claims
		c.Set("token", token)        // Set the token in the context

		c.Next()
	}
}

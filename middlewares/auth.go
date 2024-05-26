package middlewares

import (
	"net/http"
	"strings"

	"github.com/Anjasfedo/go-react-fireauth/configs"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not found"})
			c.Abort()
			return
		}

		// Extract JWT token from the Authorization header
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		userRecord, err := configs.AuthClient.VerifyIDToken(c, tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID Token"})
			c.Abort()
			return
		}

		c.Set("userUID", userRecord.UID)

		c.Next()
	}
}

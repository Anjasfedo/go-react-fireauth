package middlewares

import (
	"net/http"

	"github.com/Anjasfedo/go-react-fireauth/configs"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetHeader("UID")
		if uid == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "UID  header required"})
			c.Abort()
			return
		}

		user, err := configs.AuthClient.GetUser(c, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data"})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}

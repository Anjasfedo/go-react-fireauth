package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"firebase.google.com/go/auth"
	"github.com/Anjasfedo/go-react-fireauth/configs"
)

func main() {
	configs.InitFirebase()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/secure-data", AuthMiddleware(), func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User data not found in context"})
			return
		}
		
		userRecord, ok := user.(*auth.UserRecord)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cast user data to *auth.UserRecord"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"userRecord": userRecord,
		})
	})

	r.Run()
}

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

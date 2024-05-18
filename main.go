package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"google.golang.org/api/iterator"

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

		data, err := allDocs(c, configs.FirestoreClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"userRecord": userRecord,
			"posts":      data,
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

func allDocs(ctx context.Context, client *firestore.Client) ([]map[string]interface{}, error) {
	var data []map[string]interface{}

	iter := client.Collection("posts").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		data = append(data, doc.Data())
	}
	return data, nil
}

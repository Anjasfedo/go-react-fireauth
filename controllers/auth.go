package controllers

import (
	"net/http"
	"time"

	"github.com/Anjasfedo/go-react-fireauth/configs"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func (h AuthController) GenerateJWT(c *gin.Context) {
	var requestBody struct {
		UID string `json:"uid" binding:"required"`
	}

	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := configs.AuthClient.GetUser(c, requestBody.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.UID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		// "iss":     "your-issuer",
		// "aud":     "your-audience",
		// "sub":     "auth",
	})

	// Sign the token with a secret
	tokenString, err := token.SignedString([]byte("anjas gantek"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

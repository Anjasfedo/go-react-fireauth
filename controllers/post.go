package controllers

import (
	"net/http"

	"github.com/Anjasfedo/go-react-fireauth/models"
	"github.com/gin-gonic/gin"
)

type PostController struct{}

var postModel = new(models.Post)

func (u PostController) RetrieveAll(c *gin.Context) {
	ctx := c.Request.Context()

	posts, err := postModel.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to retrieve all posts", "error": err})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Posts found!", "posts": posts})
}

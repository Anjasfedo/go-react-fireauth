package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Anjasfedo/go-react-fireauth/models"
	"github.com/gin-gonic/gin"
)

type PostController struct{}

var postModel = new(models.Post)

func (p PostController) RetrieveAll(c *gin.Context) {
	ctx := c.Request.Context()

	posts, err := postModel.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to retrieve all posts", "error": err})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Posts found!", "posts": posts})
}

func (p PostController) RetrieveById(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		c.Abort()
		return
	}

	post, err := postModel.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error retrieveting post with ID %s: %v\n", id, err)

		if errors.Is(err, models.ErrorDocumentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
			c.Abort()
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to retrieve post", "error": err})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post found!", "post": post})
}

func (p PostController) AddPost(c *gin.Context) {
	var post models.Post

	if err := c.ShouldBindBodyWithJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post received", "post": post})
}

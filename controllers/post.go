package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Anjasfedo/go-react-fireauth/models"
	"github.com/gin-gonic/gin"
)

type PostController struct{}

var postModel = new(models.PostResponse)

func (p PostController) RetrieveAll(c *gin.Context) {
	ctx := c.Request.Context()

	posts, err := postModel.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to retrieve all posts", "error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Posts found", "posts": posts})
}

func (p PostController) RetrieveById(c *gin.Context) {
	ctx := c.Request.Context()
	ID := c.Param("id")

	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		c.Abort()
		return
	}

	post, err := postModel.GetByID(ctx, ID)
	if err != nil {
		log.Printf("Error retrieveting post with ID %s: %v\n", ID, err)

		if errors.Is(err, models.ErrorDocumentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
			c.Abort()
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to retrieve post", "error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post found", "post": post})
}

func (p PostController) AddPost(c *gin.Context) {
	ctx := c.Request.Context()
	// var post models.PostRequest

    err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error parsing multipart form", "error": err.Error()})
		c.Abort()
		return
	}

	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error retrieving image file", "error": err.Error()})
		c.Abort()
		return
	}
	defer file.Close()

	title := c.Request.FormValue("title")
	content := c.Request.FormValue("content")

	ID, err := postModel.Add(ctx, models.PostRequest{Title: title, Content: content}, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to add post", "error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post created", "ID": ID})
}

func (p PostController) UpdatePostByID(c *gin.Context) {
	ctx := c.Request.Context()
	ID := c.Param("id")

	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		c.Abort()
		return
	}

	var post models.PostRequest

	if err := c.ShouldBindBodyWithJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request", "error": err.Error()})
		c.Abort()
		return
	}

	updatedPost, err := postModel.UpdateByID(ctx, ID, post)
	if err != nil {
		if err == models.ErrorDocumentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating post", "error": err.Error()})
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated", "post": updatedPost})
}

func (p PostController) DeletePostById(c *gin.Context) {
	ctx := c.Request.Context()
	ID := c.Param("id")

	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		c.Abort()
		return
	}

	err := postModel.DeleteById(ctx, ID)
	if err != nil {
		if err == models.ErrorDocumentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting post", "error": err.Error()})
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

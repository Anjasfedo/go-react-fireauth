package server

import (
	"github.com/Anjasfedo/go-react-fireauth/controllers"
	"github.com/Anjasfedo/go-react-fireauth/middlewares"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)

	router.Use(middlewares.AuthMiddleware())

	v1 := router.Group("v1")
	{
		postGroup := v1.Group("post")
		{
			post := new(contro)
		}
	}
}
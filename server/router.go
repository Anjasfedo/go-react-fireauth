package server

import (
	"github.com/gin-gonic/gin"

	"github.com/Anjasfedo/go-react-fireauth/controllers"
	"github.com/Anjasfedo/go-react-fireauth/middlewares"
)

func NewRouter() *gin.Engine {
	router := gin.New()

	router.MaxMultipartMemory = 10 << 20
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)
	router.GET("/health", health.Status)

	v1 := router.Group("v1")
	{
		authGroup := v1.Group("auths")
		{
			auth := &controllers.AuthController{}

			authGroup.POST("/token/", auth.GenerateJWT)
		}
		postGroup := v1.Group("posts")
		{
			post := &controllers.PostController{}

			postGroup.Use(middlewares.AuthMiddleware())

			postGroup.GET("/", post.RetrieveAll)
			postGroup.GET("/:id/", post.RetrieveById)
			postGroup.POST("/", post.AddPost)
			postGroup.PUT("/:id/", post.UpdatePostByID)
			postGroup.DELETE("/:id/", post.DeletePostById)
		}
	}

	return router
}

package routes

import (
	"authentication-authorization/controllers"
	"authentication-authorization/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/login", controllers.Login())
	router.POST("/signup", controllers.Signup())

	protected := router.Group("/")
	protected.Use(middleware.Authenticate())
	{
		protected.GET("/users", controllers.GetUsers())
		protected.GET("/user/:id", controllers.GetUser())
	}
}

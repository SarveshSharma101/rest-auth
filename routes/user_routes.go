package routes

import (
	"log"
	"rest-auth/controller"
	"rest-auth/middleware"

	"github.com/gin-gonic/gin"
)

// AddUserRoutes adds user routes to the server
func AddUserRoutes(server *gin.Engine) {
	log.Println("Adding user routes")

	server.POST("/user", controller.RegisterUser)
	server.POST("/login", controller.LoginUser)

	userGroup := server.Group("/user")

	log.Println("Adding authentication middleware")
	userGroup.Use(middleware.AuthMiddleware)
	{

		userGroup.GET("/:emailId", controller.GetUsers)
		userGroup.GET("/", controller.GetUsers)
		userGroup.DELETE("/logout/:emailId", controller.LogoutUser)
	}

	userGroup.Use(middleware.UpdateAuthMiddleware)
	{
		userGroup.PUT("/:emailId", controller.UpdateUser)
		userGroup.PATCH("/:emailId", controller.PatchUser)
		userGroup.DELETE("/:emailId", controller.DeleteUser)
	}
}

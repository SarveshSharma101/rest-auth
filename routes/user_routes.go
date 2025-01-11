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

	server.POST("/users", controller.RegisterUser)
	server.POST("/login", controller.LoginUser)

	userGroup := server.Group("/")

	log.Println("Adding authentication middleware")
	userGroup.Use(middleware.AuthMiddleware)
	{

		userGroup.GET("/users/:id", controller.GetUsers)
		userGroup.PUT("/user", controller.UpdateUser)
		userGroup.PATCH("/user", controller.PatchUser)
		userGroup.DELETE("/user", controller.DeleteUser)
	}
}

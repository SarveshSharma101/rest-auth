package main

import (
	"log"
	"os"

	"rest-auth/routes"
	"rest-auth/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	server := gin.Default()
	routes.AddUserRoutes(server)

	log.Fatal(server.Run(os.Getenv(utils.SERVER_PORT)))
}

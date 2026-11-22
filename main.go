package main

import (
	"authentication-authorization/config"
	"authentication-authorization/helpers"
	"authentication-authorization/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting application...")

	key := config.GenerateRandomKey()
	helpers.SetJWTKey(key)

	fmt.Println("Generated key: %s", key)

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run(":8080")
	log.Println("Server is running on http://localhost:8080")
}

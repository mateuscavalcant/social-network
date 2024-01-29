package main

import (
	"log"
	"os"
	"social-network-go/pkg/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	r.LoadHTMLGlob("client/templates/*")
	r.Static("/client", "./client")

	routes.InitRoutes(r.Group("/"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
	
}

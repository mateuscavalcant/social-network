package main

import (
	"log"
	"net/http"
	"social-network-go/server/database"
	"social-network-go/server/routes"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitializeDB()

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.LoadHTMLGlob("client/templates/*")
	r.Static("/client", "./client")

	routes.InitRoutes(r.Group("/"))
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

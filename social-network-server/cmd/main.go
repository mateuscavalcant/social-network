package main

import (
	"log"
	"net/http"
	"social-network-server/api/routes"
	"social-network-server/pkg/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	database.InitializeDB()

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// Configuração do CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}                             // Permitir todas as origens
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}        // Métodos permitidos
	config.AllowHeaders = []string{"Authorization", "Content-Type"} // Cabeçalhos permitidos

	r.Use(cors.New(config))

	// Inicializar rotas
	routes.InitRoutes(r.Group("/"))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

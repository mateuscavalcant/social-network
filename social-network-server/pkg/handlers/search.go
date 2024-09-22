package handlers

import (
	"log"
	"net/http"
	"social-network-server/config/database"
	"social-network-server/internal/models"
	"social-network-server/pkg/repositories"

	"github.com/gin-gonic/gin"
)

func SearchEngine(c *gin.Context) {
	var req models.SearchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error JSON:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	log.Println("Search term:", req.Search)

	if req.Search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search term is required"})
		return
	}

	db := database.GetDB()
	users, err := repositories.SearchUsers(db, req.Search)
	if err != nil {
		log.Println("Error searching users:", err.Error())
		return
	}

	log.Println("info user:", users)

	resp := map[string]interface{}{
		"users": users,
	}
	c.JSON(http.StatusOK, resp)
}

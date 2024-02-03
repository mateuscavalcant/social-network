package views

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	CON "social-network-go/pkg/database"
	"social-network-go/pkg/models"

	"github.com/gin-gonic/gin"
)

func OtherProfileView(c *gin.Context) {
	var user models.UserProfile
	var targetUserID int
	user.ID = targetUserID

	tmpl := template.Must(template.ParseFiles("client/templates/other_profile.html"))

	username := c.Param("username")
	db := CON.DB()

	queryUserID := "SELECT id FROM user WHERE username = ?"
	errID := db.QueryRow(queryUserID, username).Scan(&targetUserID)
	if errID != nil {
		if errID == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		log.Println("Failed to query target user information:", errID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch target user information",
		})
		return
	}

	query := `SELECT name, username FROM user WHERE id = ?`
	err := db.QueryRow(query, targetUserID).Scan(&user.Name, &user.Username)
	if err != nil {
		log.Println("Failed to scan statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to scan rows",
		})
		return
	}
	
	data := models.UserProfile {
		Name: user.Name,
		Username: user.Username,

	}

	
	tmpl.Execute(c.Writer, data)
}
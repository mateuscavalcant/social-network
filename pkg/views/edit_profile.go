package views

import (
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	CON "social-network-go/pkg/database"
	"social-network-go/pkg/models"
	"social-network-go/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func EditProfileView(c *gin.Context) {
	var user models.UserProfile
	tmpl := template.Must(template.ParseFiles("client/templates/edit_profile.html"))

	idInterface, _ := utils.AllSessions(c)

	id, _ := strconv.Atoi(idInterface.(string))
	db := CON.DB()

	query := `SELECT name, bio, icon FROM user WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&user.Name, &user.Bio, &user.Icon)
	if err != nil {
		log.Println("Failed to scan statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to scan rows",
		})
		return
	}
	
	imageBase64 := base64.StdEncoding.EncodeToString(user.Icon)

	data := models.UserProfile {
		Name: user.Name,
		Bio: user.Bio,
		IconBase64: imageBase64,

	}
	tmpl.Execute(c.Writer, data)
}
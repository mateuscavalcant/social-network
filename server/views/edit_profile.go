package views

import (
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"social-network-go/server/database"
	"social-network-go/server/models"
	"social-network-go/server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// EditProfileView handles rendering the edit profile page.
func EditProfileView(c *gin.Context) {
	// Define a struct to hold user profile data
	var user models.UserProfile
	// Load the HTML template for the edit profile page
	tmpl := template.Must(template.ParseFiles("client/templates/edit_profile.html"))

	// Retrieve the current user's ID from the session
	idInterface, _ := utils.AllSessions(c)
	id, _ := strconv.Atoi(idInterface.(string))
	db := database.GetDB()

	// Query the database to retrieve user profile information
	query := `SELECT name, bio, icon FROM user WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&user.Name, &user.Bio, &user.Icon)
	if err != nil {
		log.Println("Failed to scan statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to scan rows",
		})
		return
	}
	
	// Encode the user's profile image to base64 for displaying in HTML
	imageBase64 := base64.StdEncoding.EncodeToString(user.Icon)

	// Populate the data struct with user profile information
	data := models.UserProfile {
		Name: user.Name,
		Bio: user.Bio,
		IconBase64: imageBase64,
	}

	// Render the edit profile page with the populated data
	tmpl.Execute(c.Writer, data)
}
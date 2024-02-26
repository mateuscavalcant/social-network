package views

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"social-network-go/server/database"
	"social-network-go/server/models"

	"github.com/gin-gonic/gin"
)

// OtherProfileView renders the profile page of another user.
func OtherProfileView(c *gin.Context) {
	// Define a UserProfile struct to hold user profile data
	var user models.UserProfile
	var targetUserID int
	user.ID = targetUserID

	// Load the HTML template for the other profile page
	tmpl := template.Must(template.ParseFiles("client/templates/other_profile.html"))

	// Extract the username parameter from the request URL
	username := c.Param("username")
	db := database.GetDB()

	// Query the database to retrieve the ID of the target user based on their username
	queryUserID := "SELECT id FROM user WHERE username = ?"
	errID := db.QueryRow(queryUserID, username).Scan(&targetUserID)
	if errID != nil {
		if errID == sql.ErrNoRows {
			// Return a 404 Not Found error if the user is not found
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

	// Query the database to retrieve profile information of the target user
	query := `SELECT name, username FROM user WHERE id = ?`
	err := db.QueryRow(query, targetUserID).Scan(&user.Name, &user.Username)
	if err != nil {
		log.Println("Failed to scan statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to scan rows",
		})
		return
	}
	
	// Populate the data struct with user profile information
	data := models.UserProfile {
		Name: user.Name,
		Username: user.Username,
	}

	// Render the other profile page with the populated data
	tmpl.Execute(c.Writer, data)
}
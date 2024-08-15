package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"social-network-server/pkg/database"
	"social-network-server/pkg/models"
	"social-network-server/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// Signup handles user signup requests.
func Signup(c *gin.Context) {
	var user models.User

	// Extract form inputs from the request
	user.Username = strings.TrimSpace(c.PostForm("username"))
	user.Name = strings.TrimSpace(c.PostForm("name"))
	user.Email = strings.TrimSpace(c.PostForm("email"))
	user.Password = strings.TrimSpace(c.PostForm("password"))
	user.ConfirmPassword = strings.TrimSpace(c.PostForm("confirm_password"))

	// Default bio value for new users
	user.Bio = "Your bio"

	// Read default user icon file
	fileBytes, err := ioutil.ReadFile("client/public/images/user-icon.jpg")
	if err != nil {
		log.Println("Error reading file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user icon"})
		return
	}
	user.Icon = fileBytes

	// Connect to the database
	db := database.GetDB()

	// Call the signup service
	errors, err := service.SignupService(db, user)
	if err != nil {
		log.Println("Error in SignupService:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	if errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Successful signup"})
}

package handlers

import (
	"log"
	"social-network-server/internal/models/errs"

	"social-network-server/pkg/validators"
	"strings"

	"github.com/gin-gonic/gin"
)

// ExistEmail checks if the given email already exists in the system.
func ExistEmail(c *gin.Context) {
	// Prepare error response object
	resp := errs.ErrorResponse{
		Error: make(map[string]string),
	}

	// Extract email from the request
	email := strings.TrimSpace(c.PostForm("email"))

	// Check if the email already exists
	existEmail, err := validators.ExistEmail(email)
	if err != nil {
		log.Println("Error checking email existence:", err)
		c.JSON(500, gin.H{"error": "Failed to validate email"})
		return
	}

	// Validate email format
	if email == "" {
		resp.Error["missing"] = "Some values are missing!"
	}
	if validators.ValidateFormatEmail(email) != nil {
		resp.Error["email"] = "Invalid email format!"
	}
	if existEmail {
		resp.Error["email"] = "Email already exists!"
	}

	// If there are errors, return them; otherwise, return success message
	if len(resp.Error) > 0 {
		c.JSON(400, resp)
		return
	}
	c.JSON(200, gin.H{"message": "Email Valid"})
}

// ExistUsername checks if the given username already exists in the system.
func ExistUsername(c *gin.Context) {
	// Prepare error response object
	resp := errs.ErrorResponse{
		Error: make(map[string]string),
	}

	// Extract username from the request
	username := strings.TrimSpace(c.PostForm("username"))

	// Check if the username already exists
	existUsername, err := validators.ExistUsername(username)
	if err != nil {
		log.Println("Error checking email existence:", err)
		c.JSON(500, gin.H{"error": "Failed to validate email"})
		return
	}

	// Validate username format
	if !validators.ValidateFormatUsername(username) {
		resp.Error["username"] = "Invalid username format!"
	}
	if existUsername {
		resp.Error["username"] = "Email already exists!"
	}

	// If there are errors, return them; otherwise, return success message
	if len(resp.Error) > 0 {
		c.JSON(400, resp)
		return
	}
	c.JSON(200, gin.H{"message": "Username Valid"})
}

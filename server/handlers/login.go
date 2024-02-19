package handlers

import (
	"log"
	CON "social-network-go/server/database"
	"social-network-go/server/models"
	"social-network-go/server/models/errs"
	"social-network-go/server/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserLogin handles user login requests.
func UserLogin(c *gin.Context) {
	var user models.User

	// Extract identifier (username or email) and password from the request
	identifier := strings.TrimSpace(c.PostForm("identifier"))
	password := strings.TrimSpace(c.PostForm("password"))

	// Prepare error response object
	resp := errs.ErrorResponse{
		Error: make(map[string]string),
	}

	// Determine if the identifier is an email or username
	isEmail := strings.Contains(identifier, "@")
	var queryField string
	if isEmail {
		queryField = "email"
	} else {
		queryField = "username"
	}

	// Connect to the database
	db := CON.DB()

	// Query the user's ID, email, and password from the database based on the identifier
	err := db.QueryRow("SELECT id, email, password FROM user WHERE "+queryField+"=?", identifier).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		resp.Error["credentials"] = "Invalid credentials"
	}

	// Compare the hashed password from the database with the provided password
	encErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if encErr != nil {
		resp.Error["password"] = "Invalid password"
	}

	// If there are errors, return them; otherwise, log the user in and return success message
	if len(resp.Error) > 0 {
		c.JSON(400, resp)
		return
	}

	// Store user session information
	session := utils.GetSession(c)
	session.Values["id"] = strconv.Itoa(user.ID)
	session.Values["email"] = user.Email
	session.Save(c.Request, c.Writer)

	c.JSON(200, gin.H{"message": "User logged in successfully"})
}

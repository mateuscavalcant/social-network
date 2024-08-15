package handlers

import (
	"log"
	"social-network-server/pkg/database"
	"social-network-server/pkg/models"

	"social-network-server/pkg/models/errs"

	"social-network-server/api/utils"

	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// DeleteAccount handles the logic for deleting a user account.
func DeleteAccount(c *gin.Context) {
	var user models.User

	// Retrieve the current user's ID from the session
	idInterface, _ := utils.AllSessions(c)

	// Convert the user ID to an integer
	id, _ := strconv.Atoi(idInterface.(string))
	user.ID = id

	// Retrieve the identifier (either email or username) and password from the request
	identifier := strings.TrimSpace(c.PostForm("identifier"))
	password := strings.TrimSpace(c.PostForm("password"))
	confirmPassword := strings.TrimSpace(c.PostForm("confirm_password"))

	// Initialize an ErrorResponse object to hold error messages
	resp := errs.ErrorResponse{
		Error: make(map[string]string),
	}

	// Determine whether the identifier is an email or username
	isEmail := strings.Contains(identifier, "@")
	var queryField string
	if isEmail {
		queryField = "email"
	} else {
		queryField = "username"
	}

	// Get a database connection
	db := database.GetDB()

	// Prepare a query to retrieve user data based on the identifier
	stmt, err := db.Prepare("SELECT id, email, password FROM user WHERE " + queryField + " = ?")
	if err != nil {
		log.Println("Error preparing query:", err)
		resp.Error["credentials"] = "Invalid credentials"
		c.JSON(400, resp)
		return
	}
	defer stmt.Close()

	// Execute the query to retrieve user data
	err = stmt.QueryRow(identifier).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		log.Println("Error retrieving user data:", err)
		resp.Error["credentials"] = "Invalid credentials"
		c.JSON(400, resp)
		return
	}

	// Compare the provided password with the hashed password stored in the database
	encErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if encErr != nil {
		resp.Error["password"] = "Invalid password"
		c.JSON(400, resp)
		return
	}

	// Check if the password and confirm password match
	if password != confirmPassword {
		resp.Error["confirm_password"] = "Passwords don't match"
		c.JSON(400, resp)
		return
	}

	// Prepare a query to delete the user account from the database
	stmt, err = db.Prepare("DELETE FROM user WHERE " + queryField + " = ?")
	if err != nil {
		log.Println("Error preparing query:", err)
		resp.Error["delete"] = "Failed to delete user"
		c.JSON(500, resp)
		return
	}
	defer stmt.Close()

	// Execute the query to delete the user account
	_, deleteErr := stmt.Exec(identifier)
	if deleteErr != nil {
		log.Println("Error deleting user:", deleteErr)
		resp.Error["delete"] = "Failed to delete user"
		c.JSON(500, resp)
		return
	}

	// Respond with a success message if the user account deletion is successful
	c.JSON(200, gin.H{"message": "User deleted successfully"})
}
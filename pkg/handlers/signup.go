package handlers

import (
	"io/ioutil"
	"log"
	CON "social-network-go/pkg/database"
	"social-network-go/pkg/models"
	"social-network-go/pkg/models/errs"
	"social-network-go/pkg/validators"
	"strings"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	var user models.User

	bio := "Your bio"
	username := strings.TrimSpace(c.PostForm("username"))
	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.TrimSpace(c.PostForm("email"))
	password := strings.TrimSpace(c.PostForm("password"))
	confirmPassword := strings.TrimSpace(c.PostForm("confirm_password"))

	resp := errs.ErrorResponse{
		Error: make(map[string]string),
	}

	existEmail, err := validators.ExistEmail(email)
	if err != nil {
		log.Println("Error checking email existence:", err)
		c.JSON(500, gin.H{"error": "Failed to validate email"})
		return
	}

	if username == "" || name == "" || email == "" || password == "" || confirmPassword == ""  {
		resp.Error["missing"] = "Some values are missing!"
	}

	if len(username) < 4 || len(username) > 32 {
		resp.Error["username"] = "Username should be between 4 and 32"
	}
	if len(name) < 1 || len(name) > 70 {
		resp.Error["name"] = "Name should be between 1 and 70"
	}
	if name == "" {
		resp.Error["name"] = "Values are missing!"
	}
	if len(bio) > 150 {
		resp.Error["bio"] = "Name should be between 1 and 70"
	}
	if validators.ValidateFormatEmail(email) != nil {
		resp.Error["email"] = "Invalid email format!"
	}
	if existEmail {
		resp.Error["email"] = "Email already exists!"
	}
	if password == "" {
		resp.Error["password"] = "Values are missing!"
	}
	if len(password) < 8 || len(password) > 16 {
		resp.Error["password"] = "Passwords should be between 8 and 16"
	}
	if password != confirmPassword {
		resp.Error["confirm_password"] = "Passwords don't match"
	}
	if len(resp.Error) > 0 {
		c.JSON(400, resp)
		return
	}

	fileBytes, err := ioutil.ReadFile("client/public/images/user-icon-post.png")
	if err != nil {
		log.Println("Error reading file:", err)
	}

	user.Username = username
	user.Email = email
	user.Password = password
	user.Name = name
	user.Bio = bio
	user.Icon = fileBytes

	db := CON.DB()

	query := "INSERT INTO user (username, name, bio, email, password, icon) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(user.Username, user.Name, user.Bio, user.Email, validators.Hash(user.Password), user.Icon)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	
	c.JSON(200, gin.H{"message": "User created successfully"})
}

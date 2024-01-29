package handlers

import (
	"log"
	CON "social-network-go/pkg/database"
	"social-network-go/pkg/models"
	"social-network-go/pkg/models/errs"
	"social-network-go/pkg/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserLogin(c *gin.Context) {
	var user models.User

	identifier := strings.TrimSpace(c.PostForm("identifier"))
	password := strings.TrimSpace(c.PostForm("password"))

	resp := errs.ErrorResponse{
		Error: make(map[string]string),
	}

	isEmail := strings.Contains(identifier, "@")
	var queryField string
	if isEmail {
		queryField = "email"
	} else {
		queryField = "username"
	}

	db := CON.DB()

	err := db.QueryRow("SELECT id, email, password FROM user WHERE "+queryField+"=?", identifier).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		resp.Error["credentials"] = "Invalid credentials"
	}

	encErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if encErr != nil {
		resp.Error["password"] = "Invalid password"
	}

	if len(resp.Error) > 0 {
		c.JSON(400, resp)
		return
	}

	session := utils.GetSession(c)
	session.Values["id"] = strconv.Itoa(user.ID)
	session.Values["email"] = user.Email
	session.Save(c.Request, c.Writer)
	c.JSON(200, gin.H{"message": "User logged in successfully"})
}

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

// Signup lida com solicitações de inscrição do usuário.
func Signup(c *gin.Context) {
	var user models.User

	// Extrair entradas de formulário da solicitação
	user.Username = strings.TrimSpace(c.PostForm("username"))
	user.Name = strings.TrimSpace(c.PostForm("name"))
	user.Email = strings.TrimSpace(c.PostForm("email"))
	user.Password = strings.TrimSpace(c.PostForm("password"))
	user.ConfirmPassword = strings.TrimSpace(c.PostForm("confirm_password"))

	user.Bio = "Your bio"

	fileBytes, err := ioutil.ReadFile("social-network-server/pkg/data/user-icon.jpg")
	if err != nil {
		log.Println("Error reading file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user icon"})
		return
	}
	user.Icon = fileBytes

	// Conectando ao banco de dados
	db := database.GetDB()

	// Chamando o serviço signup
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

	// Returnando mensagem de sucesso
	c.JSON(http.StatusOK, gin.H{"message": "Successful signup"})
}

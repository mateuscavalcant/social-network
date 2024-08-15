package handlers

import (
	"log"
	"net/http"
	"social-network-server/pkg/database"
	"social-network-server/pkg/models/errs"
	"social-network-server/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// UserLogin lida com solicitações de login de usuário.
func UserLogin(c *gin.Context) {
	identifier := strings.TrimSpace(c.PostForm("identifier"))
	password := strings.TrimSpace(c.PostForm("password"))

	// Preparar objeto de resposta de erro
	resp := errs.ErrorResponse{
		Error: make(map[string]string),
	}

	// Conectar ao banco de dados
	db := database.GetDB()

	// Autenticar usuário
	user, err := service.AuthenticateUser(db, identifier, password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			resp.Error["credentials"] = "Credenciais inválidas"
		} else {
			log.Println("Error authenticating user:", err)
			resp.Error["credentials"] = "Erro na autenticação"
		}
		c.JSON(400, resp)
		return
	}

	// Gerar token JWT
	tokenString, err := service.GenerateToken(user)
	if err != nil {
		log.Println("Error generating token:", err)
		c.JSON(500, gin.H{"error": "Falha ao gerar token"})
		return
	}

	// Criar um cookie com o token
	cookie := http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
	}

	// Definir o cookie na resposta
	http.SetCookie(c.Writer, &cookie)

	// Retornar mensagem de sucesso para o cliente
	c.JSON(200, gin.H{
		"message": "Usuário logado com sucesso",
		"token":   tokenString,
	})
}

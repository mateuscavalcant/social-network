package handlers

import (
	"log"
	"social-network-server/config/database"
	"social-network-server/internal/models"

	"social-network-server/internal/models/errs"

	"social-network-server/api/utils"

	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// DeleteAccount lida com a lógica para excluir uma conta de usuário.
func DeleteAccount(c *gin.Context) {
	var user models.User

	// Recupera o ID do usuário atual a partir da sessão
	idInterface, _ := utils.AllSessions(c)

	// Converte o ID do usuário para um inteiro
	id, _ := strconv.Atoi(idInterface.(string))
	user.ID = id

	// Recupera o identificador (seja e-mail ou nome de usuário) e a senha da solicitação
	identifier := strings.TrimSpace(c.PostForm("identifier"))
	password := strings.TrimSpace(c.PostForm("password"))
	confirmPassword := strings.TrimSpace(c.PostForm("confirm_password"))

	// Inicializa um objeto ErrorResponse para armazenar mensagens de erro
	resp := errs.ErrorResponse{
		Error: make(map[string]string),
	}

	// Determina se o identificador é um e-mail ou nome de usuário
	isEmail := strings.Contains(identifier, "@")
	var queryField string
	if isEmail {
		queryField = "email"
	} else {
		queryField = "username"
	}

	// Obtém uma conexão com o banco de dados
	db := database.GetDB()

	// Prepara uma consulta para recuperar os dados do usuário com base no identificador
	stmt, err := db.Prepare("SELECT id, email, password FROM user WHERE " + queryField + " = ?")
	if err != nil {
		log.Println("Erro ao preparar consulta:", err)
		resp.Error["credentials"] = "Credenciais inválidas"
		c.JSON(400, resp)
		return
	}
	defer stmt.Close()

	// Executa a consulta para recuperar os dados do usuário
	err = stmt.QueryRow(identifier).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		log.Println("Erro ao recuperar dados do usuário:", err)
		resp.Error["credentials"] = "Credenciais inválidas"
		c.JSON(400, resp)
		return
	}

	// Compara a senha fornecida com a senha criptografada armazenada no banco de dados
	encErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if encErr != nil {
		resp.Error["password"] = "Senha inválida"
		c.JSON(400, resp)
		return
	}

	// Verifica se a senha e a confirmação da senha correspondem
	if password != confirmPassword {
		resp.Error["confirm_password"] = "As senhas não coincidem"
		c.JSON(400, resp)
		return
	}

	// Prepara uma consulta para excluir a conta do usuário do banco de dados
	stmt, err = db.Prepare("DELETE FROM user WHERE " + queryField + " = ?")
	if err != nil {
		log.Println("Erro ao preparar consulta:", err)
		resp.Error["delete"] = "Falha ao excluir usuário"
		c.JSON(500, resp)
		return
	}
	defer stmt.Close()

	// Executa a consulta para excluir a conta do usuário
	_, deleteErr := stmt.Exec(identifier)
	if deleteErr != nil {
		log.Println("Erro ao excluir usuário:", deleteErr)
		resp.Error["delete"] = "Falha ao excluir usuário"
		c.JSON(500, resp)
		return
	}

	// Responde com uma mensagem de sucesso se a exclusão da conta do usuário for bem-sucedida
	c.JSON(200, gin.H{"message": "Usuário excluído com sucesso"})
}

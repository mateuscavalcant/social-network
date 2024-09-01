package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"social-network-server/config/database"
	"social-network-server/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Follow lida com a lógica para um usuário seguir outro usuário.
func Follow(c *gin.Context) {

	var userfollow models.UserFollow
	var requestBody map[string]string

	// Decodifica o corpo da solicitação para obter o nome de usuário
	if err := c.BindJSON(&requestBody); err != nil {
		log.Println("Falha ao decodificar o corpo da solicitação:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Corpo da solicitação inválido"})
		return
	}

	username, exists := requestBody["username"]
	if !exists || username == "" {
		log.Println("Nome de usuário está vazio")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome de usuário é obrigatório"})
		return
	}

	userId, exists := c.Get("id")
	if !exists {
		log.Println("ID do usuário não encontrado na sessão")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário não encontrado na sessão"})
		return
	}

	id, errId := strconv.Atoi(fmt.Sprintf("%v", userId))
	if errId != nil {
		log.Println("Erro ao converter ID do usuário para int:", errId)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário inválido"})
		c.Abort()
		return
	}

	db := database.GetDB()

	var userID int
	err := db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&userID)
	if err == sql.ErrNoRows {
		log.Println("Usuário não encontrado:", username)
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	} else if err != nil {
		log.Println("Falha ao consultar ID do usuário:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao obter ID do usuário"})
		return
	}

	userfollow.FollowBy = id
	userfollow.FolloTo = userID

	// Prepara uma instrução SQL para inserir uma nova entrada na tabela user_follow
	stmt, err := db.Prepare("INSERT INTO user_follow(followBy, followTo) VALUES(?, ?)")
	if err != nil {
		log.Println("Falha ao preparar instrução", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Falha ao preparar instrução",
		})
		return
	}

	// Executa a instrução SQL para inserir o relacionamento de seguimento
	_, err = stmt.Exec(userfollow.FollowBy, userfollow.FolloTo)
	if err != nil {
		log.Println("Falha ao executar consulta", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Falha ao executar consulta",
		})
		return
	}

	// Responde com uma mensagem JSON indicando a ação de seguir com sucesso
	resp := map[string]interface{}{
		"mssg": "Seguindo ",
	}
	c.JSON(http.StatusOK, resp)
}

func Unfollow(c *gin.Context) {
	var userUnfollow models.UserFollow
	var requestBody map[string]string

	// Decodifica o corpo da solicitação para obter o nome de usuário
	if err := c.BindJSON(&requestBody); err != nil {
		log.Println("Falha ao decodificar o corpo da solicitação:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Corpo da solicitação inválido"})
		return
	}

	username, exists := requestBody["username"]
	if !exists || username == "" {
		log.Println("Nome de usuário está vazio")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome de usuário é obrigatório"})
		return
	}

	userId, exists := c.Get("id")
	if !exists {
		log.Println("ID do usuário não encontrado na sessão")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário não encontrado na sessão"})
		return
	}

	id, errId := strconv.Atoi(fmt.Sprintf("%v", userId))
	if errId != nil {
		log.Println("Erro ao converter ID do usuário para int:", errId)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário inválido"})
		c.Abort()
		return
	}

	db := database.GetDB()

	var userID int
	err := db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&userID)
	if err == sql.ErrNoRows {
		log.Println("Usuário não encontrado:", username)
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	} else if err != nil {
		log.Println("Falha ao consultar ID do usuário:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao obter ID do usuário"})
		return
	}

	userUnfollow.FollowBy = id
	userUnfollow.FolloTo = userID

	stmt, err := db.Prepare("DELETE FROM user_follow WHERE followBy=? AND followTo=?")
	if err != nil {
		log.Println("Falha ao preparar instrução", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao preparar instrução"})
		return
	}

	_, err = stmt.Exec(userUnfollow.FollowBy, userUnfollow.FolloTo)
	if err != nil {
		log.Println("Falha ao executar consulta", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao executar consulta"})
		return
	}

	resp := map[string]interface{}{
		"mssg": "Deixou de seguir",
	}
	c.JSON(http.StatusOK, resp)
}

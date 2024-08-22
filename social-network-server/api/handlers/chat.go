package handlers

import (
	"fmt"
	"log"
	"net/http"
	repo "social-network-server/database"
	"social-network-server/pkg/models/errs"
	"social-network-server/pkg/websockets"
	"social-network-server/service"
	"strings"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Chat é um manipulador HTTP que lida com solicitações de chat.
func Chat(c *gin.Context) {
	userId, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in session"})
		return
	}

	id, err := strconv.Atoi(fmt.Sprintf("%v", userId))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	username := c.Param("username")
	partnerID, err := repo.MessageGetUserIDByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	messages, err := service.GetChatMessages(id, partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}
	currentUsername, err := repo.GetUsernameByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user Username"})
		return
	}

	chatPartnerName, chatPartnerUsername, chatPartnerIcon, err := service.GetChatPartnerInfo(partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat partner info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"currentUsername": gin.H{"username": currentUsername},
		"messages":        messages,
		"chatPartner":     gin.H{"name": chatPartnerName, "username": chatPartnerUsername, "iconBase64": chatPartnerIcon},
	})
}

// WebSocketHandler é um manipulador HTTP para a rota websockets.
func WebSocketHandler(c *gin.Context) {
	ws, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if err != nil {
		log.Println("Erro ao atualizar para WebSocket:", err)
		return
	}
	defer ws.Close()

	userID := websockets.GetUserIDFromContext(c)
	if userID == 0 {
		return
	}

	// Registrar a conexão
	websockets.UserConnections[int64(userID)] = ws
	log.Println("Conexão WebSocket registrada para o usuário:", userID)

	// Iniciar o controle de inatividade
	go websockets.StartInactivityTimer(ws, userID)

	// Iniciar o manuseio de mensagens
	websockets.HandleMessages(ws, userID)
}

func CreateNewMessage(c *gin.Context) {
	var errResp errs.ErrorResponse

	// Parse do corpo da requisição
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.Param("username")
	content := strings.TrimSpace(c.PostForm("content"))
	userId, exists := c.Get("id")
	if !exists {
		log.Println("ID do usuário não encontrado na sessão")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário não encontrado na sessão"})
		return
	}

	id, err := strconv.Atoi(fmt.Sprintf("%v", userId))
	if err != nil {
		log.Println("Erro ao converter ID do usuário para int:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário inválido"})
		return
	}

	// Validação básica
	if content == "" {
		errResp.Error["content"] = "Values are missing!"
	}
	if len(errResp.Error) > 0 {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	// Chama o service para enviar a mensagem
	messageID, err := websockets.SendMessage(id, username, content)
	if err != nil {
		log.Println("Erro ao enviar mensagem:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	resp := map[string]interface{}{
		"messageID": messageID,
		"message":   "Message sent successfully",
	}

	c.JSON(http.StatusOK, resp)
}

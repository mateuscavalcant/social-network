package handlers

import (
	"fmt"
	"log"
	"net/http"
	repo "social-network-server/database"
	"social-network-server/pkg/models"
	"social-network-server/service"
	"strconv"

	"social-network-server/pkg/models/errs"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// userConnections mapeia IDs de usuário para conexões WebSocket
// Canais para enviar mensagens
var (
	userConnections map[int64]*websocket.Conn
	messageQueue    = make(chan models.UserMessage, 100) // Buffer de 100 mensagens

)

func init() {
	userConnections = make(map[int64]*websocket.Conn)
	go handleWebSocketMessages()
}

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
	partnerID, err := repo.MessageGetUserIDByUsername(username) // Supondo uma função no services
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	messages, err := service.GetChatMessages(id, partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	chatPartnerName, chatPartnerUsername, chatPartnerIcon, err := service.GetChatPartnerInfo(partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat partner info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages":    messages,
		"chatPartner": gin.H{"name": chatPartnerName, "username": chatPartnerUsername, "iconBase64": chatPartnerIcon},
	})
}

// Função para enviar mensagens para o canal
func sendMessage(message models.UserMessage) {
	messageQueue <- message
}

// Função para lidar com as mensagens WebSocket
func handleWebSocketMessages() {
	for {
		// Aguarda mensagens no canal
		message := <-messageQueue

		// Verifique se o destinatário está conectado
		destConn, ok := userConnections[int64(message.MessageTo)]
		if !ok {
			log.Println("Recipient is not connected")
			continue
		}

		// Envie a mensagem para o destinatário
		err := destConn.WriteJSON(message) // Use WriteJSON para enviar mensagens JSON via WebSocket
		if err != nil {
			log.Println("Error sending message:", err)
			continue
		}
	}
}

// WebSocketHandler é um manipulador HTTP para a rota WebSocket.
func WebSocketHandler(c *gin.Context) {
	userId, exists := c.Get("id")
	if !exists {
		log.Println("ID do usuário não encontrado na sessão")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário não encontrado na sessão"})
		return
	}

	var id int
	idFloat, ok := userId.(float64)
	if !ok {
		id = int(idFloat)
		return
	} else {
		log.Println("ID do usuário não encontrado no token ou erro de conversão")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário inválido"})
		return
	}

	ws, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if err != nil {
		log.Println("Erro ao atualizar para WebSocket:", err)
		return
	}
	defer ws.Close()

	// Registrar a conexão com o ID convertido corretamente
	userConnections[int64(id)] = ws
	log.Println("Conexão WebSocket registrada para o usuário:", id)

	HandleMessages(ws)
}

// Enviar mensagens para o canal
func HandleMessages(ws *websocket.Conn) {
	defer ws.Close()

	for {
		var msg models.UserMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error receiving message:", err)
			return

		}

		go sendMessage(msg)
	}
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
	messageID, err := service.SendMessage(id, username, content)
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

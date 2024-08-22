package handlers

import (
	"fmt"
	"log"
	"net/http"
	repo "social-network-server/database"
	"social-network-server/pkg/models"
	"social-network-server/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// userConnections mapeia IDs de usuário para conexões WebSocket
// Canais para enviar mensagens
var (
	userChatsConnections map[int64]*websocket.Conn
	chatsQueue           = make(chan models.UserMessage, 100) // Buffer de 100 mensagens
)

func init() {
	userChatsConnections = make(map[int64]*websocket.Conn)
	go handleWebSocketChats()
}

func FeedChats(c *gin.Context) {
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

	chats, err := service.GetUserChats(int64(id))
	if err != nil {
		log.Println("Error in service layer:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	currentUsername, err := repo.GetUsernameByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user Username"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"currentUsername": gin.H{"username": currentUsername},
		"chats":           chats,
	})
}

func WebSocketFeedChats(c *gin.Context) {
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

	ws, err := websocket.Upgrade(c.Writer, c.Request, nil, 4096, 4096)
	if err != nil {
		log.Println("Erro ao atualizar para WebSocket:", err)
		return
	}
	defer ws.Close()

	userChatsConnections[int64(id)] = ws
	HandleMessagesChats(ws)
}

// Função para enviar mensagens para o canal
func sendMessageChats(message models.UserMessage) {
	chatsQueue <- message
}

// Função para lidar com as mensagens WebSocket
func handleWebSocketChats() {
	const workerCount = 5

	for i := 0; i < workerCount; i++ {
		go workerChats()
	}
}

func workerChats() {
	for message := range chatsQueue {
		destConn, ok := userChatsConnections[int64(message.MessageTo)]
		if !ok {
			log.Println("Recipient is not connected")
			continue
		}
		err := destConn.WriteJSON(message)

		if err != nil {
			log.Println("Error sending message: ", err)
		}

	}
}

// Enviar mensagens para o canal
func HandleMessagesChats(ws *websocket.Conn) {
	defer ws.Close()

	for {
		var msg models.UserMessage
		err := ws.ReadJSON(&msg) // Use ReadJSON para ler mensagens JSON do WebSocket
		if err != nil {
			log.Println("Error receiving message:", err)
			return
		}

		// Envie a mensagem para o canal
		go sendMessageChats(msg)
	}
}

/*

 */

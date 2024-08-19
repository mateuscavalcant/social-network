package handlers

import (
	"fmt"
	"log"
	"net/http"
	repo "social-network-server/database"
	"social-network-server/pkg/database"
	"social-network-server/pkg/models"
	"social-network-server/service"
	"strconv"

	"social-network-server/pkg/models/errs"

	"strings"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// userConnections mapeia IDs de usuário para conexões WebSocket
// Canais para enviar mensagens
var (
	userMessageConnections   map[int64]*websocket.Conn
	messageQueue             = make(chan models.UserMessage, 100) // Buffer de 100 mensagens
	connectionMessageMutexes sync.Map
)

func init() {
	userMessageConnections = make(map[int64]*websocket.Conn)
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
	partnerID, err := repo.MessageGetUserIDByUsername(username) // Supondo uma função no service
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
	const workerCount = 5

	for i := 0; i < workerCount; i++ {
		go worker()
	}
}

func worker() {
	for message := range messageQueue {
		destConn, ok := userMessageConnections[int64(message.MessageTo)]
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

// WebSocketHandler é um manipulador HTTP para a rota WebSocket.
func WebSocketHandler(c *gin.Context) {

	cookie, errCookie := c.Request.Cookie("token")

	if errCookie != nil {
		log.Println("Cookie do tokeb não encontrado")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Cookie do token não foi encontrado"})
		return
	}
	tokenString := cookie.Value
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

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte("your_secret_key"), nil
	})

	if err != nil || !token.Valid {
		log.Println("Token JWT inválido:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token JWT inválido"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Println("Token JWT inválido")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token JWT inválido"})
		return
	}

	id, ok = claims["id"].(int)
	if !ok {
		log.Println("ID do usuário não encontrado no token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário não encontrado no token"})
		return
	}

	// Atualizar a conexão para WebSocket
	ws, err := websocket.Upgrade(c.Writer, c.Request, nil, 4096, 4096)
	if err != nil {
		log.Println("Erro ao atualizar para WebSocket:", err)
		return
	}
	defer ws.Close()

	// Registre a conexão com o usuário
	userMessageConnections[int64(id)] = ws

	// Aguardar mensagens do usuário
	HandleMessages(ws)
}

// Enviar mensagens para o canal
func HandleMessages(ws *websocket.Conn) {
	defer ws.Close()

	for {
		var msg models.UserMessage
		err := ws.ReadJSON(&msg) // Use ReadJSON para ler mensagens JSON do WebSocket
		if err != nil {
			log.Println("Error receiving message:", err)
			return
		}

		// Envie a mensagem para o canal
		go sendMessage(msg)
	}
}

func CreateNewMessage(c *gin.Context) {
	var userMessage models.UserMessage
	var errResp errs.ErrorResponse

	// Parse form data
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

	id, errId := strconv.Atoi(fmt.Sprintf("%v", userId))
	if errId != nil {
		log.Println("Erro ao converter ID do usuário para int:", errId)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário inválido"})
		c.Abort()
		return
	}

	if content == "" {
		errResp.Error["content"] = "Values are missing!"
	}

	if len(errResp.Error) > 0 {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	userMessage.Content = content

	db := database.GetDB()

	var userID int
	err := db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&userID)
	if err != nil {
		log.Println("Failed to query user ID", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	userMessage.MessageBy = id
	userMessage.MessageTo = userID

	stmt, err := db.Prepare("INSERT INTO user_message(content, messageBy, messageTo, created_at) VALUES (?, ?, ?, NOW())")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	rs, err := stmt.Exec(userMessage.Content, userMessage.MessageBy, userMessage.MessageTo)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute statement"})
		return
	}

	insertID, _ := rs.LastInsertId()

	resp := map[string]interface{}{
		"messageID": insertID,
		"message":   "Message sent successfully",
	}

	c.JSON(http.StatusOK, resp)
}

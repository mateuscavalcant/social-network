package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"social-network-server/pkg/database"
	"social-network-server/pkg/models"
	"strconv"
	"time"

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
	userConnections   map[int64]*websocket.Conn
	messageQueue      = make(chan models.UserMessage, 100) // Buffer de 100 mensagens
	chatsQueue        = make(chan models.UserMessage, 100) // Buffer de 100 mensagens
	connectionMutexes sync.Map
)

func init() {
	userConnections = make(map[int64]*websocket.Conn)
	go handleWebSocketMessages()
	go handleWebSocketChats()
}

// Chat é um manipulador HTTP que lida com solicitações de chat.
func Chat(c *gin.Context) {
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

	// Extrair o nome de usuário da solicitação
	username := c.Param("username")

	// Obter o banco de dados
	db := database.GetDB()

	// Obter mensagens do banco de dados
	var messages []models.UserMessage
	var chatPartnerName string
	var chatPartnerUsername string
	var chatPartnerIcon []byte

	// Obter o ID do usuário com base no nome de usuário
	var userID int
	err := db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&userID)
	if err != nil {
		log.Println("Failed to query user ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user ID"})
		return
	}

	// Preparar a consulta para buscar mensagens entre usuários
	stmt, err := db.Prepare(`
SELECT user_message.message_id, user_message.messageBy AS message_user_id, user_message.content,
	   user.id AS user_id, user.username, user.name, user.icon, user_message.created_at
FROM user_message
JOIN user ON user.id = user_message.messageBy
WHERE (user_message.messageBy = ? AND user_message.messageTo = ?) OR 
	  (user_message.messageBy = ? AND user_message.messageTo = ?)
ORDER BY user_message.created_at ASC
`)
	if err != nil {
		log.Println("Failed to prepare statement:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	// Executar a consulta SQL
	rows, err := stmt.Query(id, userID, userID, id)
	if err != nil {
		log.Println("Failed to execute query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to execute query"})
		return
	}
	defer rows.Close()

	// Processar as linhas retornadas
	for rows.Next() {
		var message models.UserMessage
		var icon []byte
		var createdAtString string

		err := rows.Scan(&message.MessageID, &message.MessageUserID, &message.Content, &message.UserID, &message.CreatedBy, &message.Name, &icon, &createdAtString)
		if err != nil {
			log.Println("Failed to scan rows:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to scan rows"})
			return
		}

		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtString)
		if err != nil {
			log.Println("Failed to parse created_at:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse created_at"})
			return
		}

		if message.UserID == id {
			message.MessageSession = true
		} else {
			message.MessageSession = false
		}

		// Codificar o ícone em base64, se existir
		var imageBase64 string
		if icon != nil {
			imageBase64 = base64.StdEncoding.EncodeToString(icon)
		}

		// Formatar a data para exibir apenas hora e minutos
		formattedCreatedAt := createdAt.Format("15:04")

		message.CreatedAt = formattedCreatedAt

		// Adicionar mensagem à lista de mensagens
		messages = append(messages, models.UserMessage{
			MessageID:      message.MessageID,
			MessageUserID:  message.MessageUserID,
			Content:        message.Content,
			UserID:         message.UserID,
			CreatedBy:      message.CreatedBy,
			Name:           message.Name,
			IconBase64:     imageBase64,
			MessageSession: message.MessageSession,
			CreatedAt:      message.CreatedAt, // Usar a data formatada aqui
		})
	}

	// Obter o nome e o ícone do usuário com quem está conversando
	err = db.QueryRow("SELECT name, icon FROM user WHERE id = ?", userID).Scan(&chatPartnerName, &chatPartnerIcon)
	if err != nil {
		log.Println("Failed to query chat partner details:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get chat partner details"})
		return
	}
	var chatPartnerIconBase64 string
	if chatPartnerIcon != nil {
		chatPartnerIconBase64 = base64.StdEncoding.EncodeToString(chatPartnerIcon)
	}

	// Verificar se o usuário possui uma conexão WebSocket ativa
	conn, ok := userConnections[int64(id)]
	if !ok {
		log.Println("User does not have an active WebSocket connection")
	} else {
		connMutex, _ := connectionMutexes.LoadOrStore(int64(id), &sync.Mutex{})
		connMutex.(*sync.Mutex).Lock()
		defer connMutex.(*sync.Mutex).Unlock()
		// Enviar as mensagens via WebSocket
		err = conn.WriteJSON(gin.H{"messages": messages, "chatPartner": gin.H{"name": chatPartnerName, "iconBase64": chatPartnerIconBase64}})
		if err != nil {
			log.Println("Error sending messages via WebSocket:", err)
			// Tratar o erro de forma apropriada
		}
	}

	err = db.QueryRow("SELECT username FROM user WHERE id = ?", id).Scan(&chatPartnerUsername)
	if err != nil {
		log.Println("Failed to query chat partner details:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get chat partner details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages":            messages,
		"chatPartner":         gin.H{"name": chatPartnerName, "iconBase64": chatPartnerIconBase64},
		"chatPartnerUsername": gin.H{"username": chatPartnerUsername},
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
	// Extrair o token JWT dos parâmetros de consulta
	tokenString := c.Query("token")
	if tokenString == "" {
		log.Println("Token JWT não fornecido")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token JWT não fornecido"})
		return
	}

	// Validar o token JWT
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

	// Obter o ID do usuário a partir das claims do token
	userId, ok := claims["id"].(float64)
	if !ok {
		log.Println("ID do usuário não encontrado no token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário não encontrado no token"})
		return
	}

	// Atualizar a conexão para WebSocket
	ws, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if err != nil {
		log.Println("Erro ao atualizar para WebSocket:", err)
		return
	}
	defer ws.Close()

	// Registre a conexão com o usuário
	userConnections[int64(userId)] = ws

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
		sendMessage(msg)
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

	db := database.GetDB()

	query := `
    SELECT 
        user.id AS user_id, 
        user.username, 
        user.name, 
        user.icon, 
        user_message.content, 
        user_message.created_at
    FROM user_message
    JOIN user ON (user.id = user_message.messageTo OR user.id = user_message.messageBy)
    WHERE (user_message.messageBy = ? OR user_message.messageTo = ?)
    AND user.id != ?
    AND user_message.created_at = (
        SELECT MAX(user_message2.created_at)
        FROM user_message AS user_message2
        WHERE (
            (user_message2.messageBy = user_message.messageBy AND user_message2.messageTo = user_message.messageTo) 
            OR (user_message2.messageBy = user_message.messageTo AND user_message2.messageTo = user_message.messageBy)
        )
    )
    ORDER BY user_message.created_at DESC
`

	rows, err := db.Query(query, id, id, id)
	if err != nil {
		log.Println("Failed to query statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to execute query",
		})
		return
	}
	defer rows.Close()

	chats := []models.UserMessage{}

	for rows.Next() {
		var chat models.UserMessage
		var icon []byte
		var createdAtString string

		err := rows.Scan(&chat.UserID, &chat.CreatedBy, &chat.Name, &icon, &chat.Content, &createdAtString)
		if err != nil {
			log.Println("Failed to scan statement", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to scan rows",
			})
			return
		}

		var imageBase64 string
		if icon != nil {
			imageBase64 = base64.StdEncoding.EncodeToString(icon)
		}

		chats = append(chats, models.UserMessage{
			UserID:     chat.UserID,
			CreatedBy:  chat.CreatedBy,
			Name:       chat.Name,
			IconBase64: imageBase64,
			Content:    chat.Content,
			CreatedAt:  chat.CreatedAt,
		})
	}
	log.Println("Number of posts retrieved:", len(chats)) // Log para verificar o número de posts

	// Verificar se o usuário possui uma conexão WebSocket ativa
	conn, ok := userConnections[int64(id)]
	if !ok {
		log.Println("User does not have an active WebSocket connection")
	} else {
		connMutex, _ := connectionMutexes.LoadOrStore(int64(id), &sync.Mutex{})
		connMutex.(*sync.Mutex).Lock()
		defer connMutex.(*sync.Mutex).Unlock()
		// Enviar as mensagens via WebSocket
		err = conn.WriteJSON(gin.H{"chats": chats})
		if err != nil {
			log.Println("Error sending messages via WebSocket:", err)
			// Tratar o erro de forma apropriada
		}
	}

	// Obter o user e o ícone do usuário

	var chatPartnerUsername string
	var chatPartnerIcon []byte

	err = db.QueryRow("SELECT username, icon FROM user WHERE id = ?", id).Scan(&chatPartnerUsername, &chatPartnerIcon)
	if err != nil {
		log.Println("Failed to query chat partner details:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get chat partner details"})
		return
	}
	var chatPartnerIconBase64 string
	if chatPartnerIcon != nil {
		chatPartnerIconBase64 = base64.StdEncoding.EncodeToString(chatPartnerIcon)
	}

	c.JSON(http.StatusOK, gin.H{
		"chats":       chats,
		"chatPartner": gin.H{"username": chatPartnerUsername, "iconBase64": chatPartnerIconBase64},
	})
}

func WebSocketFeedChats(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		log.Println("Token JWT não fornecido")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token JWT não fornecido"})
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

	userId, ok := claims["id"].(float64)
	if !ok {
		log.Println("ID do usuário não encontrado no token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário não encontrado no token"})
		return
	}

	ws, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if err != nil {
		log.Println("Erro ao atualizar para WebSocket:", err)
		return
	}
	defer ws.Close()

	userConnections[int64(userId)] = ws
	HandleMessagesChats(ws)
}

// Função para enviar mensagens para o canal
func sendMessageChats(message models.UserMessage) {
	chatsQueue <- message
}

// Função para lidar com as mensagens WebSocket
func handleWebSocketChats() {
	for {
		// Aguarda mensagens no canal
		message := <-chatsQueue

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
		sendMessageChats(msg)
	}
}

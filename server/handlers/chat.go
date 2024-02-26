package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	"social-network-go/server/database"
	"social-network-go/server/models"
	"social-network-go/server/models/errs"
	"social-network-go/server/utils"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// userConnections mapeia IDs de usuário para conexões WebSocket.
var userConnections map[int64]*websocket.Conn

// Canal para enviar mensagens
var messageQueue = make(chan models.UserMessage, 100) // Buffer de 100 mensagens

var connectionMutexes sync.Map


func init() {
    userConnections = make(map[int64]*websocket.Conn)
    go handleWebSocketMessages()
}

// Chat é um manipulador HTTP que lida com solicitações de chat.
func Chat(c *gin.Context) {
    // Extrair o ID do usuário da sessão
    idInterface, _ := utils.AllSessions(c)
    id, _ := strconv.Atoi(idInterface.(string))

    // Extrair o nome de usuário da solicitação
    username := c.Param("username")

    // Obter o banco de dados
	db := database.GetDB()

    // Obter mensagens do banco de dados
    var messages []models.UserMessage

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
        SELECT user_message.message_id, user_message.id AS message_user_id, user_message.content,
               user.id AS user_id, user.username, user.name, user.icon
        FROM user_message
        JOIN user ON user.id = user_message.id
        WHERE (user_message.id = ? AND user_message.messageTo = ?) OR 
              (user_message.id = ? AND user_message.messageTo = ?)
        ORDER BY user_message.created_at DESC
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

        err := rows.Scan(&message.MessageID, &message.MessageUserID, &message.Content, &message.UserID, &message.CreatedBy, &message.Name, &icon)
        if err != nil {
            log.Println("Failed to scan rows:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to scan rows"})
            return
        }
        // Verificar se a mensagem é da sessão atual ou de outra pessoa
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
        })
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
        err = conn.WriteJSON(messages)
        if err != nil {
            log.Println("Error sending messages via WebSocket:", err)
            // Tratar o erro de forma apropriada
        }
    }

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
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
	ws, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer ws.Close()

	// Registre a conexão com o usuário
	idInterface, _ := utils.AllSessions(c)
	id, _ := strconv.Atoi(idInterface.(string))
	userConnections[int64(id)] = ws // Armazene o ponteiro ws, que é do tipo *websocket.Conn

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
    idInterface, _ := utils.AllSessions(c)
	id, _ := strconv.Atoi(idInterface.(string))


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

    stmt, err := db.Prepare("INSERT INTO user_message(content, messageBy, messageTo, id, created_at) VALUES (?, ?, ?, ?, NOW())")
    if err != nil {
        log.Println("Error preparing SQL statement:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
        return
    }
    defer stmt.Close()

    rs, err := stmt.Exec(userMessage.Content, userMessage.MessageBy, userMessage.MessageTo, id)
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
package websockets

import (
	"fmt"
	"log"
	"net/http"
	repo "social-network-server/database"
	"social-network-server/pkg/models"
	"social-network-server/pkg/models/errs"
	"social-network-server/service"
	"strings"

	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Mapeamento de conexões WebSocket por ID de usuário
var (
	UserConnections   map[int64]*websocket.Conn
	connectionMutexes sync.Map
	workerPool        *WorkerPool
)

func init() {
	UserConnections = make(map[int64]*websocket.Conn)
	workerPool = NewWorkerPool(10) // Pool com 10 workers
	go handleWebSocketMessages()
}

// Pool de workers para processar mensagens
type WorkerPool struct {
	workers  int
	jobQueue chan models.UserMessage
	wg       sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	pool := &WorkerPool{
		workers:  numWorkers,
		jobQueue: make(chan models.UserMessage, 100), // Buffer com 100 mensagens
	}
	pool.startWorkers()
	return pool
}

func (pool *WorkerPool) startWorkers() {
	for i := 0; i < pool.workers; i++ {
		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done()
			for job := range pool.jobQueue {
				processMessage(job)
			}
		}()
	}
}

func (pool *WorkerPool) Submit(job models.UserMessage) {
	select {
	case pool.jobQueue <- job:
		// Mensagem enviada para o pool com sucesso
	default:
		// Buffer cheio, decidir como lidar
		log.Println("Buffer de mensagens cheio, mensagem será descartada ou reprocessada.")
	}
}

func (pool *WorkerPool) Shutdown() {
	close(pool.jobQueue)
	pool.wg.Wait()
}

func processMessage(message models.UserMessage) {
	// Processar a mensagem e enviar via WebSocket
	if conn, ok := UserConnections[int64(message.MessageTo)]; ok {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Println("Erro ao enviar mensagem:", err)
		}
	} else {
		log.Println("Destinatário não está conectado")
	}
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

	userInfosName, userInfosUsername, userInfosIcon, err := service.GetuserInfosInfo(partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat partner info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"currentUsername": currentUsername,
		"messages":        messages,
		"userInfos":     gin.H{"name": userInfosName, "username": userInfosUsername, "iconBase64": userInfosIcon},
	})
}

// Função para enviar mensagens para o pool
func sendMessage(message models.UserMessage) {
	workerPool.Submit(message)
}

// Função para lidar com as mensagens WebSocket de forma eficiente em lote
func handleWebSocketMessages() {
	batch := make([]models.UserMessage, 0, 10) // Processar lotes de 10 mensagens
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case message := <-workerPool.jobQueue:
			batch = append(batch, message)
			if len(batch) >= 10 {
				flushMessages(batch)
				batch = batch[:0] // Limpar o batch
			}
		case <-ticker.C:
			if len(batch) > 0 {
				flushMessages(batch)
				batch = batch[:0]
			}
		}
	}
}

func flushMessages(batch []models.UserMessage) {
	for _, message := range batch {
		if conn, ok := UserConnections[int64(message.MessageTo)]; ok {
			// Enviar todas as mensagens em um único payload JSON
			err := conn.WriteJSON(batch)
			if err != nil {
				log.Println("Erro ao enviar mensagens:", err)
			}
		} else {
			log.Println("Destinatário não está conectado")
		}
	}
}

// Função para iniciar o controle de inatividade
func StartInactivityTimer(ws *websocket.Conn, userID int) {
	inactivityDuration := 30 * time.Second
	inactivityTimer := time.NewTimer(inactivityDuration)

	for {
		select {
		case <-inactivityTimer.C:
			// Fechar a conexão após 30 segundos de inatividade
			log.Println("Fechando conexão por inatividade:", userID)
			ws.Close()
			delete(UserConnections, int64(userID))
			return
		case <-time.After(1 * time.Second): // Checa a cada segundo se a conexão ainda está ativa
			if _, isConnected := UserConnections[int64(userID)]; !isConnected {
				inactivityTimer.Stop()
				return
			}
		}
	}
}

// Função para gerenciar o timeout de conexão ociosa com o PongHandler e redefinir o timer de inatividade
func HandleMessages(ws *websocket.Conn, userID int) {
	defer ws.Close()

	ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	ws.SetPongHandler(func(appData string) error {
		ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg models.UserMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Erro ao receber mensagem:", err)
			return
		}

		// Redefinir o timer de inatividade sempre que uma mensagem for recebida
		go ResetInactivityTimer(userID)

		sendMessage(msg)
	}
}

// Função para redefinir o timer de inatividade
func ResetInactivityTimer(userID int) {
	if _, ok := UserConnections[int64(userID)]; ok {
		log.Println("Redefinindo timer de inatividade para o usuário:", userID)
	}
}

// Helper para extrair o ID do usuário do contexto
func GetUserIDFromContext(c *gin.Context) int {
	userId, exists := c.Get("id")
	if !exists {
		log.Println("User ID not found in session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in session"})
		return 0
	}

	var id int
	idFloat, ok := userId.(float64)
	if !ok {
		id = int(idFloat)

	} else {
		log.Println("ID do usuário não encontrado no token ou erro de conversão")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
	}

	return id
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
		log.Println("User ID not found in session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in session"})
		return
	}

	id, err := strconv.Atoi(fmt.Sprintf("%v", userId))
	if err != nil {
		log.Println("Erro ao converter ID do usuário para int:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
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
	messageID, err := SendMessage(id, username, content)
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

func SendMessage(senderID int, receiverUsername, content string) (int64, error) {
	// Obtém o ID do usuário destinatário
	receiverID, err := repo.MessageGetUserIDByUsername(receiverUsername)
	if err != nil {
		log.Println("Erro ao obter ID do destinatário:", err)
		return 0, err
	}

	// Cria a mensagem
	message := models.UserMessage{
		MessageBy: senderID,
		MessageTo: receiverID,
		Content:   content,
	}

	// Salva a mensagem no banco de dados
	messageID, err := repo.SaveMessage(message)
	if err != nil {
		log.Println("Erro ao salvar a mensagem:", err)
		return 0, err
	}

	// Verifica se o destinatário está online
	conn, isOnline := UserConnections[int64(receiverID)]
	if isOnline {
		// Envia a mensagem via WebSocket
		go func() {
			if err := conn.WriteJSON(message); err != nil {
				log.Printf("Erro ao enviar mensagem via WebSocket para o usuário %d: %v", receiverID, err)
			}
		}()
	} else {
		log.Printf("O destinatário %d não está online. A mensagem foi apenas armazenada no banco de dados.", receiverID)
	}

	return messageID, nil
}

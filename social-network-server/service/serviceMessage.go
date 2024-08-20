package service

import (
	"encoding/base64"
	"fmt"
	"log"
	"social-network-server/database"
	"social-network-server/pkg/models"
	"time"
)

// Obter mensagens entre usuários e processá-las
func GetChatMessages(user1ID, user2ID int) ([]models.UserMessage, error) {
	messages, err := database.GetUserMessages(user1ID, user2ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving messages: %w", err)
	}

	for i, message := range messages {
		createdAt, err := time.Parse("2006-01-02 15:04:05", message.CreatedAt)
		if err != nil {
			log.Println("Failed to parse created_at:", err)
			continue
		}
		messages[i].CreatedAt = createdAt.Format("15:04")
		messages[i].MessageSession = message.UserID == user1ID

		// Codificar o ícone em base64
		if message.Icon != nil {
			messages[i].IconBase64 = base64.StdEncoding.EncodeToString(message.Icon)
		}
	}
	return messages, nil
}

func SendMessage(senderID int, receiverUsername, content string) (int64, error) {
	// Obtém o ID do usuário destinatário
	receiverID, err := database.MessageGetUserIDByUsername(receiverUsername)
	if err != nil {
		log.Println("Erro ao obter ID do destinatário:", err)
		return 0, nil
	}

	// Cria a mensagem
	message := models.UserMessage{
		MessageBy: senderID,
		MessageTo: receiverID,
		Content:   content,
	}

	// Salva a mensagem no banco de dados
	messageID, err := database.SaveMessage(message)
	if err != nil {
		log.Println("Erro ao salvar a mensagem:", err)
		return 0, nil
	}

	return messageID, nil
}

// Obter informações de parceiro de chat
func GetChatPartnerInfo(userID int) (string, string, string, error) {
	name, username, icon, err := database.GetUserInfo(userID)
	if err != nil {
		return "", "", "", fmt.Errorf("error retrieving chat partner info: %w", err)
	}
	iconBase64 := ""
	if icon != nil {
		iconBase64 = base64.StdEncoding.EncodeToString(icon)
	}
	return name, username, iconBase64, nil
}

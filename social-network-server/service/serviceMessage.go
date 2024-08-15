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

// Salvar nova mensagem
func SendMessage(message models.UserMessage) (int64, error) {
	return database.SaveMessage(message)
}

// Obter informações de parceiro de chat
func GetChatPartnerInfo(userID int) (string, string, error) {
	name, icon, err := database.GetUserInfo(userID)
	if err != nil {
		return "", "", fmt.Errorf("error retrieving chat partner info: %w", err)
	}
	iconBase64 := ""
	if icon != nil {
		iconBase64 = base64.StdEncoding.EncodeToString(icon)
	}
	return name, iconBase64, nil
}

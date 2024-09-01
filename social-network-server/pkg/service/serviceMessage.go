package service

import (
	"encoding/base64"
	"fmt"
	"log"
	"social-network-server/internal/models"
	"social-network-server/pkg/repositories"

	"time"
)

// Obter mensagens entre usuários e processá-las
func GetChatMessages(user1ID, user2ID int) ([]models.UserMessage, error) {
	messages, err := repositories.GetUserMessages(user1ID, user2ID, 30, 0)
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

// Obter informações de parceiro de chat
func GetuserInfosInfo(userID int) (string, string, string, error) {
	name, username, icon, err := repositories.GetUserInfo(userID)
	if err != nil {
		return "", "", "", fmt.Errorf("error retrieving chat partner info: %w", err)
	}
	iconBase64 := ""
	if icon != nil {
		iconBase64 = base64.StdEncoding.EncodeToString(icon)
	}
	return name, username, iconBase64, nil
}

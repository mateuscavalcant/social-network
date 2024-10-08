package service

import (
	"fmt"
	"social-network-server/config/database"
	"social-network-server/internal/models"
	"social-network-server/pkg/repositories"
)

func GetUserChats(userID int64) ([]models.UserMessage, error) {
	db := database.GetDB()

	chats, err := repositories.FetchUserChats(db, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user chats: %w", err)
	}

	return chats, nil
}

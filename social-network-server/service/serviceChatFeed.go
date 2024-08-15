package service

import (
	"fmt"
	repo "social-network-server/database"
	"social-network-server/pkg/database"
	"social-network-server/pkg/models"
)

func GetUserChats(userID int64) ([]models.UserMessage, error) {
	db := database.GetDB()

	chats, err := repo.FetchUserChats(db, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user chats: %w", err)
	}

	return chats, nil
}

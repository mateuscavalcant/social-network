package service

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"social-network-server/config/database"
	"social-network-server/internal/models"
)

type UserService struct{}

func (us *UserService) GetUserIcon(userID int) (*models.UserIconResponse, error) {
	db := database.GetDB()

	var username string
	var icon []byte

	err := db.QueryRow("SELECT username, icon FROM user WHERE id = ?", userID).Scan(&username, &icon)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuário não encontrado")
		}
		log.Println("Erro ao consultar detalhes do usuário:", err)
		return nil, errors.New("falha ao obter detalhes do usuário")
	}

	var iconBase64 string
	if icon != nil {
		iconBase64 = base64.StdEncoding.EncodeToString(icon)
	}

	return &models.UserIconResponse{
		Username:   username,
		IconBase64: iconBase64,
	}, nil
}

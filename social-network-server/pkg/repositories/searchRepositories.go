package repositories

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"social-network-server/internal/models"
)

func SearchUsers(db *sql.DB, searchTerm string) ([]models.User, error) {

	var icon []byte

	query := `
	SELECT id, name, username, icon, bio FROM user WHERE username LIKE ? OR name LIKE ?
	`

	rows, err := db.Query(query, "%"+searchTerm+"%", "%"+searchTerm+"%")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Name, &user.Username, &icon, &user.Bio)

		if icon != nil {
			user.IconBase64 = base64.StdEncoding.EncodeToString(icon)
		}
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("usuário não encontrado")
			}
			log.Println("Erro ao consultar perfil do usuário:", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

package repositories

import (
	"database/sql"
	"log"
	"social-network-server/internal/models"
)

// GetUserByIdentifier obtém o usuário do banco de dados com base no identificador (email ou nome de usuário).
func GetUserByIdentifier(db *sql.DB, identifier string, isEmail bool) (models.User, error) {
	var user models.User
	queryField := "username"
	if isEmail {
		queryField = "email"
	}

	query := "SELECT id, email, password FROM user WHERE " + queryField + "=?"
	err := db.QueryRow(query, identifier).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		return user, err
	}
	return user, nil
}

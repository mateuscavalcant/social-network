package service

import (
	"database/sql"
	"os"
	"social-network-server/database"
	"social-network-server/pkg/models"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// AuthenticateUser autentica o usuário com base no identificador e senha fornecidos.
func AuthenticateUser(db *sql.DB, identifier, password string) (models.User, error) {
	isEmail := strings.Contains(identifier, "@")
	user, err := database.GetUserByIdentifier(db, identifier, isEmail)
	if err != nil {
		return user, err
	}

	// Comparar a senha criptografada do banco de dados com a senha fornecida
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

// GenerateToken gera um token JWT para o usuário autenticado.
func GenerateToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
	})

	// Assinar e obter o token completo codificado como uma string
	tokenString, err := token.SignedString([]byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

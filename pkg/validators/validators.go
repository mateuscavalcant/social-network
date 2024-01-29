package validators

import (
	CON "social-network-go/pkg/database"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte(err.Error())
	}
	return hash
}

func ValidateFormatEmail(email string) error {
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return err
	}
	return nil
}

func ExistEmail(email string) (bool, error) {
	db := CON.DB()

	var emailCount int

	err := db.QueryRow("SELECT COUNT(id) AS emailCount FROM user WHERE email=?", email).Scan(&emailCount)
	if err != nil {
		return false, err
	}

	return emailCount > 0, nil
}

func ExistUsername(username string) (bool, error) {
	db := CON.DB()

	var userCount int

	err := db.QueryRow("SELECT COUNT(id) AS userCount FROM users1 WHERE username=?", username).Scan(&userCount)
	if err != nil {
		return false, err
	}
	return userCount > 0, nil
}

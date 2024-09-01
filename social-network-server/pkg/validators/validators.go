package validators

import (
	"encoding/base64"
	"regexp"
	"social-network-server/config/database"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

// Hash hashes the given password using bcrypt and returns it as a string.
func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidateFormatEmail(email string) error {
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return err
	}
	return nil
}

func ValidateFormatUsername(username string) bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return regex.MatchString(username)
}

func ExistEmail(email string) (bool, error) {
	db := database.GetDB()

	var emailCount int

	err := db.QueryRow("SELECT COUNT(id) AS emailCount FROM user WHERE email=?", email).Scan(&emailCount)
	if err != nil {
		return false, err
	}

	return emailCount > 0, nil
}

func ExistUsername(username string) (bool, error) {
	db := database.GetDB()

	var userCount int

	err := db.QueryRow("SELECT COUNT(id) AS userCount FROM users WHERE username=?", username).Scan(&userCount)
	if err != nil {
		return false, err
	}
	return userCount > 0, nil
}

func ConvertByteToBase64(x []byte) string {
	iconBase64 := base64.StdEncoding.EncodeToString(x)
	return iconBase64
}

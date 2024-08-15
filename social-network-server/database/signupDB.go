package database

import (
	"database/sql"
	"log"
	"social-network-server/pkg/models"
)

// CreateUser inserts a new user into the database.
func CreateUser(db *sql.DB, user models.User) error {
	query := "INSERT INTO user (username, name, bio, email, password, icon) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Name, user.Bio, user.Email, user.Password, user.Icon)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		return err
	}

	return nil
}

// CheckEmailExistence checks if the email already exists in the database.
func CheckEmailExistence(db *sql.DB, email string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM user WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&count)
	if err != nil {
		log.Println("Error checking email existence:", err)
		return false, err
	}
	return count > 0, nil
}

package database

import (
	"database/sql"
	"log"
	"social-network-server/pkg/models"
)

// CreateUser inserts a new user into the database.
func CreatePost(db *sql.DB, userPost models.UserPost, id int) error {

	err := db.QueryRow("SELECT username FROM user WHERE id = ?", id).Scan(&userPost.CreatedBy)
	if err != nil {
		log.Println("Error querying username:", err)
		return nil
	}

	query := "INSERT INTO user_post(content, createdBy, id, created_at) VALUES (?, ?, ?, NOW())"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userPost.Content, userPost.CreatedBy, id)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		return err
	}

	return nil
}

package database

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"social-network-server/pkg/database"
	"social-network-server/pkg/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.GetDB(),
	}
}

func (ur *UserRepository) GetUserIDByUsername(username string) (int, error) {
	var id int
	err := ur.db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("usuário não encontrado")
		}
		log.Println("Erro ao consultar ID do usuário:", err)
		return 0, err
	}
	return id, nil
}

func (ur *UserRepository) GetUserProfile(userID int) (*models.UserProfile, error) {
	var user models.UserProfile
	var icon []byte
	query := `
		SELECT user.username, user.name, user.icon, user.bio,
		       IFNULL(follower_counts.follower_count, 0) AS follower_count,
		       IFNULL(followed_counts.following_count, 0) AS following_count
		FROM user
		LEFT JOIN (
			SELECT followTo, COUNT(followBy) AS follower_count
			FROM user_follow
			GROUP BY followTo
		) AS follower_counts ON follower_counts.followTo = user.id
		LEFT JOIN (
			SELECT followBy, COUNT(followTo) AS following_count
			FROM user_follow
			GROUP BY followBy
		) AS followed_counts ON followed_counts.followBy = user.id
		WHERE user.id = ?
	`
	err := ur.db.QueryRow(query, userID).Scan(&user.Username, &user.Name, &icon, &user.Bio, &user.FollowByCount, &user.FollowToCount)
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

	return &user, nil
}

func (ur *UserRepository) GetUserPosts(userID int) ([]models.UserPost, error) {
	posts := []models.UserPost{}
	query := `
		SELECT user_post.postID, user_post.id, user_post.content, user.id, user.username, user.name, user.icon
		FROM user_post
		JOIN user ON user.id = user_post.id
		WHERE user.id = ?
		ORDER BY user_post.created_at DESC
	`
	rows, err := ur.db.Query(query, userID)
	if err != nil {
		log.Println("Erro ao consultar posts:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.UserPost
		var icon []byte
		err := rows.Scan(&post.PostID, &post.PostUserID, &post.Content, &post.UserID, &post.CreatedBy, &post.Name, &icon)
		if err != nil {
			log.Println("Erro ao escanear post:", err)
			return nil, err
		}
		if icon != nil {
			post.IconBase64 = base64.StdEncoding.EncodeToString(icon)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (ur *UserRepository) IsFollowing(followByID, followToID int) (bool, error) {
	var count int
	err := ur.db.QueryRow("SELECT COUNT(*) FROM user_follow WHERE followBy = ? AND followTo = ?", followByID, followToID).Scan(&count)
	if err != nil {
		log.Println("Erro ao verificar status de seguimento:", err)
		return false, err
	}
	return count > 0, nil
}

func (ur *UserRepository) GetUsernameByID(userID int) (string, error) {
	var username string
	err := ur.db.QueryRow("SELECT username FROM user WHERE id = ?", userID).Scan(&username)
	if err != nil {
		log.Println("Erro ao consultar username:", err)
		return "", err
	}
	return username, nil
}

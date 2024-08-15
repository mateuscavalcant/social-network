package service

import (
	"encoding/base64"
	"errors"
	"log"
	"social-network-server/pkg/database"
	"social-network-server/pkg/models"
)

type FeedService struct{}

func (fs *FeedService) GetFeed(userID int) ([]models.UserPost, error) {
	db := database.GetDB()

	var posts []models.UserPost

	query := `
        SELECT user_post.postID, user_post.id AS post_user_id, user_post.content,
               user.id AS user_id, user.username, user.name, user.icon
        FROM user_post
        JOIN user ON user.id = user_post.id
        WHERE user.id = ? OR user.id IN (
            SELECT user_follow.followTo
            FROM user_follow
            WHERE user_follow.followBy = ?
        )
        ORDER BY user_post.created_at DESC`

	rows, err := db.Query(query, userID, userID)
	if err != nil {
		return nil, errors.New("failed to execute query")
	}
	defer rows.Close()

	for rows.Next() {
		var post models.UserPost
		var icon []byte

		if err := rows.Scan(&post.PostID, &post.PostUserID, &post.Content, &post.UserID, &post.CreatedBy, &post.Name, &icon); err != nil {
			log.Println("Failed to scan statement:", err)
			return nil, err
		}

		if icon != nil {
			post.IconBase64 = base64.StdEncoding.EncodeToString(icon)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (fs *FeedService) CreatePost(userPost *models.UserPost, userID int) error {
	db := database.GetDB()

	stmt, err := db.Prepare("INSERT INTO user_post(content, createdBy, id, created_at) VALUES (?, ?, ?, NOW())")
	if err != nil {
		return errors.New("failed to prepare statement")
	}
	defer stmt.Close()

	_, err = stmt.Exec(userPost.Content, userPost.CreatedBy, userID)
	if err != nil {
		return errors.New("failed to execute statement")
	}

	return nil
}

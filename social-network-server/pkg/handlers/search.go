package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"social-network-server/config/database"
	"social-network-server/internal/models"

	"github.com/gin-gonic/gin"
)

type SearchRequest struct {
	Search string `json:"search"`
}

func Search(c *gin.Context) {
	var req SearchRequest

	// Ler o corpo da requisição JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Erro ao vincular JSON:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de post inválidos"})
		return
	}

	db := database.GetDB()

	searchTerm := req.Search

	users, err := searchUsers(db, searchTerm)
	if err != nil {
		log.Println("Erro ao buscar usuários:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuários"})
		return
	}

	posts, err := searchPosts(db, searchTerm)
	if err != nil {
		log.Println("Erro ao buscar posts:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar posts"})
		return
	}

	resp := map[string]interface{}{
		"users": users,
		"posts": posts,
	}
	c.JSON(http.StatusOK, resp)
}

func searchUsers(db *sql.DB, searchTerm string) ([]models.User, error) {
	query := `
		SELECT id, username, name
		FROM user
		WHERE username LIKE ? OR name LIKE ?
	`
	rows, err := db.Query(query, "%"+searchTerm+"%", "%"+searchTerm+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func searchPosts(db *sql.DB, searchTerm string) ([]models.UserPost, error) {
	query := `
		SELECT id, content
		FROM user_post
		WHERE content LIKE ?
	`
	rows, err := db.Query(query, "%"+searchTerm+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.UserPost
	for rows.Next() {
		var post models.UserPost
		if err := rows.Scan(&post.PostID, &post.Content); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

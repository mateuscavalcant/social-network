package handlers

import (
	"fmt"
	"log"
	"net/http"
	"social-network-server/pkg/database"
	"social-network-server/pkg/models"
	"social-network-server/pkg/models/errs"
	"social-network-server/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Feed(c *gin.Context) {
	// Obter o ID do usuário da sessão JWT
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in session"})
		return
	}

	id, err := strconv.Atoi(fmt.Sprintf("%v", userID))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	feedService := service.FeedService{}
	posts, err := feedService.GetFeed(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Obter o ícone e o username do usuário
	userService := service.UserService{}
	userDetails, err := userService.GetUserIcon(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"chatPartner": gin.H{
			"username":   userDetails.Username,
			"iconBase64": userDetails.IconBase64,
		},
	})
}
func CreateNewPost(c *gin.Context) {
	var userPost models.UserPost
	errresp := errs.ErrorResponse{Error: make(map[string]string)}

	// Obter o ID do usuário da sessão JWT
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in session"})
		return
	}

	id, err := strconv.Atoi(fmt.Sprintf("%v", userID))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.ShouldBindJSON(&userPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de post inválidos"})
		return
	}

	if userPost.Content == "" {
		errresp.Error["content"] = "O conteúdo do post não pode estar vazio!"
	}

	if len(errresp.Error) > 0 {
		c.JSON(http.StatusBadRequest, errresp)
		return
	}

	feedService := service.FeedService{}
	if err := feedService.CreatePost(&userPost, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post criado com sucesso!"})
}

// DeletePost deletes a post based on its ID.
func DeletePost(c *gin.Context) {
	postID := c.PostForm("post")
	// Obter o ID do usuário da sessão JWT
	id, exists := c.Get("id")
	if !exists {
		// Lidar com o caso em que o ID do usuário não está disponível
		log.Println("User ID not found in session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in session"})
		return
	}

	db := database.GetDB()

	var postAuthorID int
	err := db.QueryRow("SELECT id FROM user_post WHERE postID=?", postID).Scan(&postAuthorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch post details",
		})
		return
	}

	if postAuthorID != id {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to delete this post",
		})
		return
	}

	_, err = db.Exec("DELETE FROM user_post WHERE postID=?", postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete post",
		})
		return
	}

	resp := map[string]interface{}{
		"mssg": "Post Deleted!",
	}
	c.JSON(http.StatusOK, resp)
}

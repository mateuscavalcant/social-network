package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"social-network-server/pkg/database"
	"social-network-server/pkg/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Follow handles the logic for a user to follow another user.
func Follow(c *gin.Context) {

	var userfollow models.UserFollow
	var requestBody map[string]string

	// Decode the request body to get the username
	if err := c.BindJSON(&requestBody); err != nil {
		log.Println("Failed to bind request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	username, exists := requestBody["username"]
	if !exists || username == "" {
		log.Println("Username is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	userId, exists := c.Get("id")
	if !exists {
		log.Println("ID do usuário não encontrado na sessão")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário não encontrado na sessão"})
		return
	}

	id, errId := strconv.Atoi(fmt.Sprintf("%v", userId))
	if errId != nil {
		log.Println("Erro ao converter ID do usuário para int:", errId)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário inválido"})
		c.Abort()
		return
	}

	db := database.GetDB()

	var userID int
	err := db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&userID)
	if err == sql.ErrNoRows {
		log.Println("User not found:", username)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		log.Println("Failed to query user ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	userfollow.FollowBy = id
	userfollow.FolloTo = userID

	// Prepare an SQL statement to insert a new entry into the user_follow table
	stmt, err := db.Prepare("INSERT INTO user_follow(followBy, followTo) VALUES(?, ?)")
	if err != nil {
		log.Println("Failed to prepare statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to prepare statement",
		})
		return
	}

	// Execute the SQL statement to insert the follow relationship
	_, err = stmt.Exec(userfollow.FollowBy, userfollow.FolloTo)
	if err != nil {
		log.Println("Failed to execute query", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to execute query",
		})
		return
	}

	// Respond with a JSON message indicating the successful follow action
	resp := map[string]interface{}{
		"mssg": "Followed ",
	}
	c.JSON(http.StatusOK, resp)
}

func Unfollow(c *gin.Context) {
	var userUnfollow models.UserFollow
	var requestBody map[string]string

	// Decode the request body to get the username
	if err := c.BindJSON(&requestBody); err != nil {
		log.Println("Failed to bind request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	username, exists := requestBody["username"]
	if !exists || username == "" {
		log.Println("Username is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	userId, exists := c.Get("id")
	if !exists {
		log.Println("ID do usuário não encontrado na sessão")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário não encontrado na sessão"})
		return
	}

	id, errId := strconv.Atoi(fmt.Sprintf("%v", userId))
	if errId != nil {
		log.Println("Erro ao converter ID do usuário para int:", errId)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário inválido"})
		c.Abort()
		return
	}

	db := database.GetDB()

	var userID int
	err := db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&userID)
	if err == sql.ErrNoRows {
		log.Println("User not found:", username)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		log.Println("Failed to query user ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	userUnfollow.FollowBy = id
	userUnfollow.FolloTo = userID

	stmt, err := db.Prepare("DELETE FROM user_follow WHERE followBy=? AND followTo=?")
	if err != nil {
		log.Println("Failed to prepare statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}

	_, err = stmt.Exec(userUnfollow.FollowBy, userUnfollow.FolloTo)
	if err != nil {
		log.Println("Failed to execute query", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		return
	}

	resp := map[string]interface{}{
		"mssg": "Unfollowed",
	}
	c.JSON(http.StatusOK, resp)
}

package handlers

import (
	"log"
	"net/http"
	CON "social-network-go/server/database"
	"social-network-go/server/utils"

	"github.com/gin-gonic/gin"
)

func Follow(c *gin.Context) {
	id, _ := utils.AllSessions(c)
	username := c.PostForm("username")

	db := CON.DB()

	var userID int
	err := db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&userID)
	if err != nil {
		log.Println("Failed to query user ID", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user ID",
		})
		return
	}

	stmt, err := db.Prepare("INSERT INTO user_follow(followBy, followTo) VALUES(?, ?)")
	if err != nil {
		log.Println("Failed to prepare statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to prepare statement",
		})
		return
	}

	_, err = stmt.Exec(id, userID)
	if err != nil {
		log.Println("Failed to execute query", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to execute query",
		})
		return
	}

	resp := map[string]interface{}{
		"mssg": "Followed ",
	}
	c.JSON(http.StatusOK, resp)
}

func Unfollow(c *gin.Context) {
	id, _ := utils.AllSessions(c)
	username := c.PostForm("username")
	db := CON.DB()

	var userID int
	err := db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&userID)
	if err != nil {
		log.Println("Failed to query user ID", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user ID",
		})
		return
	}

	stmt, err := db.Prepare("DELETE FROM user_follow WHERE followBy=? AND followTo=?")
	if err != nil {
		log.Println("Failed to prepare statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to prepare statement",
		})
		return
	}
	
	_, err = stmt.Exec(id, userID)
	if err != nil {
		log.Println("Failed to query statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to execute query",
		})
		return
	}

	resp := map[string]interface{}{
		"mssg": "Unfollowed ",
	}
	c.JSON(http.StatusOK, resp)
}

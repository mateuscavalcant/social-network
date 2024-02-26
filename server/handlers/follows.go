package handlers

import (
	"log"
	"net/http"
	"social-network-go/server/database"
	"social-network-go/server/utils"

	"github.com/gin-gonic/gin"
)

// Follow handles the logic for a user to follow another user.
func Follow(c *gin.Context) {
	// Retrieve the current user's ID from the session
	id, _ := utils.AllSessions(c)
	// Retrieve the username of the user to be followed from the request body
	username := c.PostForm("username")

	db := database.GetDB()

	var userID int
	// Query the database to get the ID of the user to be followed
	err := db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&userID)
	if err != nil {
		log.Println("Failed to query user ID", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user ID",
		})
		return
	}

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
	_, err = stmt.Exec(id, userID)
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

// Unfollow handles the logic for a user to unfollow another user.
func Unfollow(c *gin.Context) {
	// Retrieve the current user's ID from the session
	id, _ := utils.AllSessions(c)
	// Retrieve the username of the user to be unfollowed from the request body
	username := c.PostForm("username")
	db := database.GetDB()

	var userID int
	// Query the database to get the ID of the user to be unfollowed
	err := db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&userID)
	if err != nil {
		log.Println("Failed to query user ID", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user ID",
		})
		return
	}

	// Prepare an SQL statement to delete the follow relationship from the user_follow table
	stmt, err := db.Prepare("DELETE FROM user_follow WHERE followBy=? AND followTo=?")
	if err != nil {
		log.Println("Failed to prepare statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to prepare statement",
		})
		return
	}
	
	// Execute the SQL statement to delete the follow relationship
	_, err = stmt.Exec(id, userID)
	if err != nil {
		log.Println("Failed to query statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to execute query",
		})
		return
	}

	// Respond with a JSON message indicating the successful unfollow action
	resp := map[string]interface{}{
		"mssg": "Unfollowed ",
	}
	c.JSON(http.StatusOK, resp)
}
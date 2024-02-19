package handlers

import (
	"database/sql"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	CON "social-network-go/server/database"
	"social-network-go/server/models"
	"social-network-go/server/models/errs"
	"social-network-go/server/utils"
	"social-network-go/server/views"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)


    var user models.UserProfile
	var post models.UserPost



// AnotherUserProfile handles requests to view another user's profile.
func AnotherUserProfile(c *gin.Context) {
	// Extract username from request parameters
	username := c.Param("username")
	db := CON.DB()
	var targetUserID int
	post.UserID = targetUserID

	// Query the database to get the ID of the target user
	queryUserID := "SELECT id FROM user WHERE username = ?"
	err := db.QueryRow(queryUserID, username).Scan(&targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		log.Println("Failed to query target user information:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch target user information",
		})
		return
	}

	// Query user profile information including follower and following counts
	queryUser := `
	SELECT
		user.username, user.name, user.bio,
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
	err = db.QueryRow(queryUser, targetUserID).Scan(&user.Username, &user.Name, &user.Bio, &user.FollowByCount, &user.FollowToCount)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		log.Println("Failed to query user information:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user information",
		})
		return
	}

	// Query user's posts
	posts := []models.UserPost{}
	query := "SELECT user_post.post_id, user_post.id AS user_post_id, user_post.content, user.id AS user_id, user.username, user.name FROM user_post JOIN user ON user.id = user_post.id WHERE user.id = ?"
	rows, err := db.Query(query, targetUserID)
	if err != nil {
		log.Println("Failed to query statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to execute query",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&post.PostID, &post.PostUserID, &post.Content, &post.UserID, &post.CreatedBy, &post.Name)
		if err != nil {
			log.Println("Failed to scan statement", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to scan rows",
			})
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println("Failed 3", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error occurred while iterating rows",
		})
		return
	}
	countPosts := len(posts)
	user.Posts = countPosts

	// Query user's icon
	queryIcon := `SELECT icon FROM user WHERE id = ?`
	err = db.QueryRow(queryIcon, targetUserID).Scan(&user.Icon)
	if err != nil {
		log.Println("Failed to scan statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to scan rows",
		})
		return
	}

	// Encode user's icon to base64
	imageBase64 := base64.StdEncoding.EncodeToString(user.Icon)

	// Retrieve current user ID from session
	idInterface, _ := utils.AllSessions(c)
	id, _ := strconv.Atoi(idInterface.(string))

	// Check if the current user is following the target user
	queryFollow := "SELECT COUNT(*) FROM user_follow WHERE followBy = ? AND followTo = ?"
	var followCount int
	errFollow := db.QueryRow(queryFollow, id, targetUserID).Scan(&followCount)
	if errFollow != nil {
		log.Println("Failed to check follow status:", errFollow)
	}

	// Set FollowBy field based on follow status
	user.FollowBy = followCount > 0

	// Return the target user's public profile with their public posts
	c.JSON(http.StatusOK, gin.H{
		"profile": user,
		"posts":   posts,
		"icon":    imageBase64, // Send the base64 encoded image to the client
	})
}

// Profile handles requests to view the user's own profile.
func Profile(c *gin.Context) {
	// Retrieve user ID from session
	idInterface, _ := utils.AllSessions(c)
	id, _ := strconv.Atoi(idInterface.(string))
	db := CON.DB()

	post.UserID = id

	// Query user's profile information
	queryUser := `
	SELECT
		user.username, user.name, user.bio,
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
	err := db.QueryRow(queryUser, id).Scan(&user.Username, &user.Name, &user.Bio, &user.FollowByCount, &user.FollowToCount)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		log.Println("Failed to query user information:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user information",
		})
		return
	}

	// Query user's posts
	posts := []models.UserPost{}
	query := "SELECT user_post.post_id, user_post.id AS user_post_id, user_post.content, user.id AS user_id, user.username, user.name FROM user_post JOIN user ON user.id = user_post.id WHERE user.id = ?"
	rows, err := db.Query(query, id)
	if err != nil {
		log.Println("Failed to query statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to execute query",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&post.PostID, &post.PostUserID, &post.Content, &post.UserID, &post.CreatedBy, &post.Name)
		if err != nil {
			log.Println("Failed to scan statement", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to scan rows",
			})
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println("Failed 3", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error occurred while iterating rows",
		})
		return
	}
	countPosts := len(posts)
	user.Posts = countPosts

	// Query user's icon
	queryIcon := `SELECT icon FROM user WHERE id = ?`
	err = db.QueryRow(queryIcon, id).Scan(&user.Icon)
	if err != nil {
		log.Println("Failed to scan statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to scan rows",
		})
		return
	}

	// Encode user's icon to base64
	imageBase64 := base64.StdEncoding.EncodeToString(user.Icon)

	// Return user's profile with their posts
	c.JSON(http.StatusOK, gin.H{
		"profile": user,
		"posts":   posts,
		"icon":    imageBase64,
	})
}

// RenderProfileTemplate renders the user's profile template based on the request.
func RenderProfileTemplate(c *gin.Context) {
	// Retrieve user ID from session
	idInterface, _ := utils.AllSessions(c)
	id, _ := strconv.Atoi(idInterface.(string))

	// Retrieve username from request parameters
	username := c.Param("username")

	db := CON.DB()

	// Check if the user with the given username exists
	queryExist := "SELECT COUNT(*) FROM user WHERE username = ?"
	var count int
	err := db.QueryRow(queryExist, username).Scan(&count)
	if err != nil {
		log.Println("Failed to query user existence:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check user existence",
		})
		return
	}

	// If user does not exist, render a not found page
	if count == 0 {
		c.HTML(http.StatusOK, "notfounduser.html", gin.H{})
		return
	}

	// Retrieve user session information
	var userSession models.User
	queryUserSession := "SELECT id, username FROM user WHERE id = ?"
	err = db.QueryRow(queryUserSession, id).Scan(&userSession.ID, &userSession.Username)
	if err != nil {
		log.Println("Failed to query user session information:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user session information",
		})
		return
	}

	// If user session's username does not match the requested username, render the other user's profile template
	if userSession.Username != username {
		views.RenderProfile(c, "other_profile.html", gin.H{
			"username": username,
		})
		return
	}

	// Render the user's own profile template
	views.RenderProfile(c, "profile.html", gin.H{
		"username": username,
	})
}

// EditProfile handles requests to edit the user's profile.
func EditProfile(c *gin.Context) {
	// Retrieve user ID from session
	idInterface, exists := utils.AllSessions(c)
	if exists == false || idInterface == nil {
		c.Redirect(http.StatusUnauthorized, "/login")
		return
	}

	id, _ := strconv.Atoi(idInterface.(string))

	var fileBytes []byte
	file, _, err := c.Request.FormFile("icon")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting image from form"})
		return
	} else if err == nil {
		defer file.Close()

		fileBytes, err = ioutil.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading image"})
			return
		}
	}

	name := strings.TrimSpace(c.PostForm("name"))
	bio := strings.TrimSpace(c.PostForm("bio"))

	resp := errs.ErrorResponse{
		Error: make(map[string]string),
	}

	if len(name) < 1 || len(name) > 70 {
		resp.Error["name"] = "Name should be between 1 and 70"
	}
	if len(bio) > 150 {
		resp.Error["bio"] = "Bio should be between 1 and 150"
	}

	if len(resp.Error) > 0 {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	db := CON.DB()

	stmt, err := db.Prepare("UPDATE user SET name=?, bio=? WHERE id=?")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	defer stmt.Close()

	if fileBytes != nil {
		stmt, err = db.Prepare("UPDATE user SET name=?, bio=?, icon=? WHERE id=?")
		if err != nil {
			log.Println("Error preparing SQL statement:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, bio, fileBytes, id)
	} else {
		stmt, err = db.Prepare("UPDATE user SET name=?, bio=? WHERE id=?")
		if err != nil {
			log.Println("Error preparing SQL statement:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, bio, id)
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
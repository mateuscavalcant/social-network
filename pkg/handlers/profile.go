package handlers

import (
	"database/sql"
	"log"
	"net/http"
	CON "social-network-go/pkg/database"
	"social-network-go/pkg/models"
	"social-network-go/pkg/models/errs"
	"social-network-go/pkg/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)


func CreateProfile(c *gin.Context) {
	var user models.UserProfile
	name := strings.TrimSpace(c.PostForm("name"))
	bio := strings.TrimSpace(c.PostForm("bio"))

	resp := errs.ErrorResponse{
		Error: make(map[string]string),
	}

	if name == "" {
		resp.Error["name"] = "Some values are missing!"
	}
	if len(name) < 1 || len(name) > 64 {
		resp.Error["name"] = "Name should be between 1 and 64"
	}
	if len(bio) > 120 {
		resp.Error["bio"] = "bio should be 120"
	}

	db := CON.DB()

	user.Name = name
	user.Bio = bio

	query := "INSERT INTO user (name, bio) VALUES (?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(user.Name, user.Bio)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(200, gin.H{"message": "Profile created successfully"})
}

func AnotherUserProfile(c *gin.Context) {
	utils.LoggedIn(c, "/welcome")

	username := c.Param("username")
	var post models.UserPost

	// Obtenha o ID do usuário alvo usando o nome de usuário
	db := CON.DB()
	var targetUserID int
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

	idInterface, _ := utils.AllSessions(c)
	id, _ := strconv.Atoi(idInterface.(string))

	if id == targetUserID {
		var user models.UserProfile
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

		errUser := db.QueryRow(queryUser, id).Scan(&user.Username, &user.Name, &user.Bio, &user.FollowByCount, &user.FollowToCount)
		if errUser != nil {
			if errUser == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "User not found",
				})
				return
			}
			log.Println("Failed to query user information:", errUser)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch user information",
			})
			return
		}

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

		log.Println("Count Posts:", user.Posts)

		// Consulta para verificar se o usuário atual está seguindo o usuário-alvo
		queryFollow := "SELECT COUNT(*) FROM user_follow WHERE followBy = ? AND followTo = ?"
		var followCount int
		errFollow := db.QueryRow(queryFollow, id, targetUserID).Scan(&followCount)
		if errFollow != nil {
			log.Println("Failed to check follow status:", errFollow)
			// Trate o erro, se necessário
		}

		user.FollowBy = followCount > 0

		c.JSON(http.StatusOK, gin.H{
			"profile": user,
			"posts":   posts,
		})

	} else {
		var user models.UserProfile
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
		errUser := db.QueryRow(queryUser, targetUserID).Scan(&user.Username, &user.Name, &user.Bio, &user.FollowByCount, &user.FollowToCount)
		if errUser != nil {
			if errUser == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "User not found",
				})
				return
			}
			log.Println("Failed to query target user information:", errUser)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch target user information",
			})
			return
		}
		targetUserPosts := []models.UserPost{}

		targetUserPostsQuery := "SELECT user_post.post_id, user_post.id AS user_post_id, user_post.content, user.id AS user_id, user.username, user.name FROM user_post JOIN user ON user.id = user_post.id WHERE user.id = ?"

		rowsTargetUserPosts, err := db.Query(targetUserPostsQuery, targetUserID)
		if err != nil {
			log.Println("Failed to query target user posts", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to execute query for target user posts",
			})
			return
		}
		defer rowsTargetUserPosts.Close()

		for rowsTargetUserPosts.Next() {
			err := rowsTargetUserPosts.Scan(&post.PostID, &post.PostUserID, &post.Content, &post.UserID, &post.CreatedBy, &post.Name)
			if err != nil {
				log.Println("Failed to scan target user post", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to scan rows for target user posts",
				})
				return
			}
			log.Println("@", post.CreatedBy)
			log.Println("Name:", post.Name)

			targetUserPosts = append(targetUserPosts, post)
		}

		if err := rowsTargetUserPosts.Err(); err != nil {
			log.Println("Failed while iterating target user posts", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error occurred while iterating target user posts",
			})
			return
		}
		countPosts := len(targetUserPosts)
		user.Posts = countPosts

		log.Println("Count Posts:", user.Posts)

		// Consulta para verificar se o usuário atual está seguindo o usuário-alvo
		queryFollow := "SELECT COUNT(*) FROM user_follow WHERE followBy = ? AND followTo = ?"
		var followCount int
		errFollow := db.QueryRow(queryFollow, id, targetUserID).Scan(&followCount)
		if errFollow != nil {
			log.Println("Failed to check follow status:", errFollow)
			// Trate o erro, se necessário
		}

		// Se followCount for maior que 0, o usuário atual está seguindo o usuário-alvo
		user.FollowBy = followCount > 0

		// Retorne o perfil público do usuário alvo com seus posts públicos
		c.JSON(http.StatusOK, gin.H{
			"profile": user,
			"posts":   targetUserPosts,
		})
	}
}

func Profile(c *gin.Context) {
	var user models.UserProfile
	var post models.UserPost

	utils.LoggedIn(c, "/welcome")

	idInterface, _ := utils.AllSessions(c)
	id, _ := strconv.Atoi(idInterface.(string))
	db := CON.DB()

	post.UserID = id

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
		log.Println("@", post.CreatedBy)
		log.Println("Name:", post.Name)

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

	log.Println("Count Posts:", user.Posts)

	c.JSON(http.StatusOK, gin.H{
		"profile": user,
		"posts":   posts,
	})
}

func RenderProfileTemplate(c *gin.Context) {
	username := c.Param("username")
	idInterface, _ := utils.AllSessions(c)
	id, _ := strconv.Atoi(idInterface.(string))

	db := CON.DB()

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

	if count == 0 {
		c.HTML(http.StatusOK, "notfounduser.html", gin.H{})
		return
	}

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

	if userSession.Username != username {
		c.HTML(http.StatusOK, "other_profile.html", gin.H{
			"username": username,
		})
		return
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"username": username,
	})
}

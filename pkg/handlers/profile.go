package handlers

import (
	"database/sql"
	"encoding/base64"
	"io/ioutil"
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


    var user models.UserProfile
	var post models.UserPost


func AnotherUserProfile(c *gin.Context) {
	username := c.Param("username")
	// Obtenha o ID do usuário alvo usando o nome de usuário
	db := CON.DB()
	var targetUserID int
	post.UserID = targetUserID

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

	err1 := db.QueryRow(queryUser, targetUserID).Scan(&user.Username, &user.Name, &user.Bio, &user.FollowByCount, &user.FollowToCount)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		log.Println("Failed to query user information:", err1)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user information",
		})
		return
	}

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

	queryIcon := `SELECT icon FROM user WHERE id = ?`
	errIcon := db.QueryRow(queryIcon, targetUserID).Scan(&user.Icon)
	if errIcon != nil {
		log.Println("Failed to scan statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to scan rows",
		})
		return
	}
	
	imageBase64 := base64.StdEncoding.EncodeToString(user.Icon)

	idInterface, _ := utils.AllSessions(c)
	id, _ := strconv.Atoi(idInterface.(string))
	// Consulta para verificar se o usuário atual está seguindo o usuário-alvo
	queryFollow := "SELECT COUNT(*) FROM user_follow WHERE followBy = ? AND followTo = ?"
	var followCount int
	errFollow := db.QueryRow(queryFollow, id, targetUserID).Scan(&followCount)
	if errFollow != nil {
		log.Println("Failed to check follow status:", errFollow)
	}

	// Se followCount for maior que 0, o usuário atual está seguindo o usuário-alvo
	user.FollowBy = followCount > 0

	// Retorne o perfil público do usuário alvo com seus posts públicos
	c.JSON(http.StatusOK, gin.H{
		"profile": user,
		"posts":   posts,
		"icon":    imageBase64, // Envie a imagem codificada em base64 para o cliente
	})
}
	


func Profile(c *gin.Context) {
	idInterface, exists := utils.AllSessions(c)
	if exists == false || idInterface == nil {
        c.Redirect(http.StatusUnauthorized, "/login")
        return
    }
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

	queryIcon := `SELECT icon FROM user WHERE id = ?`
	errIcon := db.QueryRow(queryIcon, id).Scan(&user.Icon)
	if errIcon != nil {
		log.Println("Failed to scan statement", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to scan rows",
		})
		return
	}
	
	imageBase64 := base64.StdEncoding.EncodeToString(user.Icon)
	
	c.JSON(http.StatusOK, gin.H{
		"profile": user,
		"posts":   posts,
		"icon":    imageBase64, 
	})
}

func RenderProfileTemplate(c *gin.Context) {
	idInterface, exists := utils.AllSessions(c)
	if exists == false || idInterface == nil {
        c.Redirect(http.StatusUnauthorized, "/login")
        return
    }
	id, _ := strconv.Atoi(idInterface.(string))

	username := c.Param("username")

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

func EditProfile(c *gin.Context) {
	idInterface, exists := utils.AllSessions(c)
	if exists == false || idInterface == nil {
        c.Redirect(http.StatusUnauthorized, "/login")
        return
    }
	
	id, _ := strconv.Atoi(idInterface.(string))

	file, _, err := c.Request.FormFile("icon")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao obter a imagem do formulário"})
        return
    }
    defer file.Close()

    // Lê o conteúdo do arquivo
    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler a imagem"})
        return
    }

    // Obtém outros dados do formulário
    username := strings.TrimSpace(c.PostForm("username"))
    name := strings.TrimSpace(c.PostForm("name"))
    bio := strings.TrimSpace(c.PostForm("bio"))

    // Validação dos dados do formulário
    resp := errs.ErrorResponse{
        Error: make(map[string]string),
    }

    if len(username) < 4 || len(username) > 32 {
        resp.Error["username"] = "Username should be between 4 and 32"
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

    user.Username = username
    user.Name = name
    user.Bio = bio
	user.Icon = fileBytes

    db := CON.DB()

	stmt, err := db.Prepare("UPDATE user SET username=?, name=?, bio=?, icon=? WHERE id=?")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	defer stmt.Close()
	
	_, err = stmt.Exec(user.Username, user.Name, user.Bio, user.Icon, id)
	if err != nil {
		log.Println("Error executing prepared SQL statement:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
package handlers

import (
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


func Feed(c *gin.Context) {
	utils.LoggedIn(c, "/welcome")

	idInterface, _ := utils.AllSessions(c)
	id, _ := strconv.Atoi(idInterface.(string))
	db := CON.DB()

	var post models.UserPost
	post.UserID = id

	posts := []models.UserPost{}

	query := `
		SELECT user_post.post_id, user_post.id AS post_user_id, user_post.content,
		       user.id AS user_id, user.username, user.name
		FROM user_post
		JOIN user ON user.id = user_post.id
		WHERE user.id = ? OR user.id IN (
		    SELECT user_follow.followTo
		    FROM user_follow
		    WHERE user_follow.followBy = ?
		)
	`

	rows, err := db.Query(query, id, id)
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
		log.Println("CreatedBy:", post.CreatedBy)

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println("Failed 3", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error occurred while iterating rows",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func CreateNewPost(c *gin.Context) {
    var userPost models.UserPost
    errresp := errs.ErrorResponse{
        Error: make(map[string]string),
    }

    content := strings.TrimSpace(c.PostForm("content"))
    idInterface, _ := utils.AllSessions(c)
    if idInterface == nil {
        // Se o usuário não estiver logado, retornar um erro de autenticação
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Unauthorized",
        })
        return
    }

    if content == "" {
        errresp.Error["content"] = "Values are missing!"
    }

    if len(errresp.Error) > 0 {
        c.JSON(400, errresp)
        return
    }

    id, _ := strconv.Atoi(idInterface.(string))
    userPost.Content = content

    db := CON.DB()

    // Recuperar o username com base no id do usuário
    var username string
    err := db.QueryRow("SELECT username FROM user WHERE id = ?", id).Scan(&username)
    if err != nil {
        log.Println("Error querying username:", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to query username",
        })
        return
    }

    // Atribuir o username recuperado ao campo createdBy
    userPost.CreatedBy = username

    stmt, err := db.Prepare("INSERT INTO user_post(content, createdBy, id) VALUES (?, ?, ?)")
    if err != nil {
        log.Println("Error preparing SQL statement:", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to prepare statement",
        })
        return
    }

    rs, err := stmt.Exec(userPost.Content, userPost.CreatedBy, id)
    if err != nil {
        log.Println("Error executing SQL statement:", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to execute statement",
        })
        return
    }

    insertID, _ := rs.LastInsertId()

    resp := map[string]interface{}{
        "postID": insertID,
        "mssg":   "Post Created!!",
    }
    c.JSON(http.StatusOK, resp)
}

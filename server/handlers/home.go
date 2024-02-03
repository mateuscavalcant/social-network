package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	CON "social-network-go/server/database"
	"social-network-go/server/models"
	"social-network-go/server/models/errs"
	"social-network-go/server/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Feed(c *gin.Context) {
    idInterface, exists := utils.AllSessions(c)
    if exists == false || idInterface == nil {
        c.Redirect(http.StatusUnauthorized, "/login")
        return
    }

    idString, ok := idInterface.(string)
    if !ok {
        c.String(http.StatusInternalServerError, "Internal Server Error")
        return
    }

    id, err := strconv.Atoi(idString)
    if err != nil {
        c.String(http.StatusInternalServerError, "Internal Server Error")
        return
    }

	db := CON.DB()

	var post models.UserPost
	post.UserID = id

	posts := []models.UserPost{}

	query := `
    SELECT user_post.post_id, user_post.id AS post_user_id, user_post.content,
           user.id AS user_id, user.username, user.name, user.icon
    FROM user_post
    JOIN user ON user.id = user_post.id
    WHERE user.id = ? OR user.id IN (
        SELECT user_follow.followTo
        FROM user_follow
        WHERE user_follow.followBy = ?
    )
    ORDER BY user_post.created_at ASC
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
    var icon []byte 

    err := rows.Scan(&post.PostID, &post.PostUserID, &post.Content, &post.UserID, &post.CreatedBy, &post.Name, &icon)
    if err != nil {
        log.Println("Failed to scan statement", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to scan rows",
        })
        return
    }

    var imageBase64 string
    if icon != nil {
        imageBase64 = base64.StdEncoding.EncodeToString(icon)
    }

    posts = append(posts, models.UserPost{
        PostID:      post.PostID,
        PostUserID:  post.PostUserID,
        Content:     post.Content,
        UserID:      post.UserID,
        CreatedBy:   post.CreatedBy,
        Name:        post.Name,
        IconBase64:        imageBase64,
    })
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

    var username string
    err := db.QueryRow("SELECT username FROM user WHERE id = ?", id).Scan(&username)
    if err != nil {
        log.Println("Error querying username:", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to query username",
        })
        return
    }

    userPost.CreatedBy = username

    stmt, err := db.Prepare("INSERT INTO user_post(content, createdBy, id, created_at) VALUES (?, ?, ?, NOW())")

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

package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"social-network-server/pkg/database"
	"social-network-server/pkg/models/errs"
	"social-network-server/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context) {
	username := c.Param("username")
	userService := service.NewUserServiceProfile()

	// Obter o ID do usuário da sessão JWT
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário não encontrado na sessão"})
		return
	}

	id, err := strconv.Atoi(fmt.Sprintf("%v", userID))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário inválido"})
		return
	}

	// Obter informações do perfil usando o serviço
	profile, posts, isCurrentUser, chatPartner, err := userService.GetUserProfileAndPosts(username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile":       profile,
		"posts":         posts,
		"icon":          profile.IconBase64,
		"isCurrentUser": isCurrentUser,
		"chatPartner":   chatPartner,
	})
}

// EditProfile handles requests to edit the user's profile.
func EditProfile(c *gin.Context) {
	// Obter o ID do usuário da sessão JWT
	userId, exists := c.Get("id")
	if !exists {
		// Lidar com o caso em que o ID do usuário não está disponível
		log.Println("ID do usuário não encontrado na sessão")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário não encontrado na sessão"})
		return
	}
	id, _ := userId.(int)

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

	db := database.GetDB()

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

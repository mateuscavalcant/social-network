package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"social-network-server/config/database"
	"social-network-server/internal/models/errs"
	"social-network-server/pkg/service"
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in session"})
		return
	}

	id, err := strconv.Atoi(fmt.Sprintf("%v", userID))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Obter informações do perfil usando o serviço
	profile, posts, isCurrentUser, userInfos, err := userService.GetUserProfileAndPosts(username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile":       profile,
		"posts":         posts,
		"icon":          profile.IconBase64,
		"isCurrentUser": isCurrentUser,
		"userInfos":     userInfos,
	})
}

// EditProfile lida com solicitações para editar o perfil do usuário.
func EditProfile(c *gin.Context) {
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

	if len(name) > 70 {
		resp.Error["name"] = "Name should be up to 70 characters"
	}
	if len(bio) > 150 {
		resp.Error["bio"] = "Bio should be up to 150 characters"
	}

	if len(resp.Error) > 0 {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	db := database.GetDB()

	// Buscar detalhes do usuário da sessão
	var currentName, currentBio string
	err = db.QueryRow("SELECT name, bio FROM user WHERE id=?", id).Scan(&currentName, &currentBio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current user data"})
		return
	}

	// Se nenhum novo nome for fornecido, mantenha o atual
	if name == "" {
		name = currentName
	}

	// Criar consulta SQL com base nos campos fornecidos
	var query string
	var args []interface{}

	if fileBytes != nil {
		query = "UPDATE user SET name=?, bio=?, icon=? WHERE id=?"
		args = []interface{}{name, bio, fileBytes, id}
	} else {
		query = "UPDATE user SET name=?, bio=? WHERE id=?"
		args = []interface{}{name, bio, id}
	}

	// Exexcuta a consulta
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare SQL statement"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

package handlers

import (
	"log"
	CON "social-network-go/pkg/database"
	"social-network-go/pkg/models"
	"social-network-go/pkg/models/errs"
	"social-network-go/pkg/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


func DeleteAccount(c *gin.Context) {
    var user models.User

    idInterface, _ := utils.AllSessions(c)

    id, _ := strconv.Atoi(idInterface.(string))
    user.ID = id

    identifier := strings.TrimSpace(c.PostForm("identifier"))
    password := strings.TrimSpace(c.PostForm("password"))
    confirmPassword := strings.TrimSpace(c.PostForm("confirm_password"))

    resp := errs.ErrorResponse{
        Error: make(map[string]string),
    }

    isEmail := strings.Contains(identifier, "@")
    var queryField string
    if isEmail {
        queryField = "email"
    } else {
        queryField = "username"
    }

    db := CON.DB()

    stmt, err := db.Prepare("SELECT id, email, password FROM user WHERE " + queryField + " = ?")
    if err != nil {
        log.Println("Erro ao preparar a consulta:", err)
        resp.Error["credentials"] = "Invalid credentials"
        c.JSON(400, resp)
        return
    }
    defer stmt.Close()

    err = stmt.QueryRow(identifier).Scan(&user.ID, &user.Email, &user.Password)
    if err != nil {
        log.Println("Erro ao recuperar os dados do usuário:", err)
        resp.Error["credentials"] = "Invalid credentials"
        c.JSON(400, resp)
        return
    }

    encErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if encErr != nil {
        resp.Error["password"] = "Invalid password"
        c.JSON(400, resp)
        return
    }

    if password != confirmPassword {
        resp.Error["confirm_password"] = "Passwords don't match"
        c.JSON(400, resp)
        return
    }

    stmt, err = db.Prepare("DELETE FROM user WHERE " + queryField + " = ?")
    if err != nil {
        log.Println("Erro ao preparar a consulta:", err)
        resp.Error["delete"] = "Failed to delete user"
        c.JSON(500, resp)
        return
    }
    defer stmt.Close()

    _, deleteErr := stmt.Exec(identifier)
    if deleteErr != nil {
        log.Println("Erro ao excluir o usuário:", deleteErr)
        resp.Error["delete"] = "Failed to delete user"
        c.JSON(500, resp)
        return
    }
    c.JSON(200, gin.H{"message": "User deleted successfully"})
}

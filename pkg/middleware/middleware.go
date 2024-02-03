package middleware

import (
    "net/http"
    "social-network-go/pkg/utils"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        session := utils.GetSession(c)
        userID := session.Values["id"]

        if userID == nil {
            c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        c.Next()
    }
}

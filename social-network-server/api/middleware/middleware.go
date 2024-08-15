package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware é um middleware para verificar se o token JWT é válido e relacionado a um usuário autenticado.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtenha o token JWT do cabeçalho de autorização
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		// O token deve estar no formato "Bearer {token}"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		// Extrair o token JWT
		tokenString := tokenParts[1]

		// Parse e verifique o token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verifique o método de assinatura
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Token inválido")
			}
			return []byte(os.Getenv("SESSION_SECRET")), nil
		})
		if err != nil {
			log.Println("Erro ao analisar o token JWT:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Verifique se o token é válido
		if !token.Valid {
			log.Println("Token Inválido")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Obtenha o ID do usuário das reivindicações do token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Reivindicação inválida")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Reivindicações inválidas"})
			c.Abort()
			return
		}

		// Extrair o ID do usuário como um valor genérico
		userID, ok := claims["id"]
		if !ok {
			log.Println("Erro ao converter ID do usuário para int")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário inválido"})
			c.Abort()
			return
		}

		// Converter o userID para int
		idInt, err := strconv.Atoi(fmt.Sprintf("%v", userID))
		if err != nil {
			log.Println("Erro ao converter ID do usuário para int:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário inválido"})
			c.Abort()
			return
		}

		// Definir o ID do usuário no contexto da requisição
		c.Set("id", idInt)

		// Continuar com a solicitação
		c.Next()
	}
}

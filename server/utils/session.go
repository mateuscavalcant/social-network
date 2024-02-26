package utils

import (
	"crypto/rand"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(generateRandomKey(32)))

func generateRandomKey(length int) []byte {
    key := make([]byte, length)
    _, err := rand.Read(key)
    if err != nil {
        log.Fatal("Erro ao gerar a chave aleatória:", err)
    }
    return key
}

func GetSession(c *gin.Context) *sessions.Session {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		log.Println("Error: ", err)

	}
	Err(err)
	return session
}

func AllSessions(c *gin.Context) (interface{}, interface{}) {
	session := GetSession(c)
	id := session.Values["id"]
	email := session.Values["email"]
	return id, email
}


// GetSession retrieves the session from the request.
func GetSessionHTTP(r *http.Request) *sessions.Session {
    session, err := store.Get(r, "session")
    if err != nil {
        log.Println("Erro ao recuperar a sessão:", err)
    }
    return session
}

// AllSessions retrieves all session data.
func AllSessionsHTTP(r *http.Request) (interface{}, interface{}) {
    session := GetSessionHTTP(r)
    id := session.Values["id"]
    email := session.Values["email"]
    return id, email
}

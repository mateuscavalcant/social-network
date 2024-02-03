package utils

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

func GetSession(c *gin.Context) *sessions.Session {
	session, err := store.Get(c.Request, "session")
	Err(err)
	return session
}

func AllSessions(c *gin.Context) (interface{}, interface{}) {
	session := GetSession(c)
	id := session.Values["id"]
	email := session.Values["email"]
	return id, email
}

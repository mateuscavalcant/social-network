package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Err(err interface{}) {
	if err != nil {
		log.Fatal(err)
	}
}

func Json(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func Ses(c *gin.Context) interface{} {
	id, username := AllSessions(c)
	return map[string]interface{}{
		"id":       id,
		"username": username,
	}
}

func LoggedIn(c *gin.Context, urlRedirect string) {
	var URL string
	if urlRedirect == "" {
		URL = "/login"
	} else {
		URL = urlRedirect
	}
	id, _ := AllSessions(c)
	if id == nil {
		c.Redirect(http.StatusFound, URL)
	}
}

func NotLoggedIn(c *gin.Context) {
	id, _ := AllSessions(c)
	if id != nil {
		c.Redirect(http.StatusFound, "/")
	}
}

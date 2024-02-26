package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Err(err interface{}) {
	if err != nil {
		log.Println("Error: ", err)
	}
}


func Ses(c *gin.Context) interface{} {
	id, username := AllSessions(c)
	return map[string]interface{}{
		"id":       id,
		"username": username,
	}
}

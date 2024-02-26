package views

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

func ChatView(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("client/templates/chat.html"))
	tmpl.Execute(c.Writer, nil)
}
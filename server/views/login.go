package views

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

func LoginView(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("client/templates/login.html"))
	tmpl.Execute(c.Writer, nil)
}
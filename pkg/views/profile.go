package views

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

func ProfileView(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("client/templates/profile.html"))

	tmpl.Execute(c.Writer, nil)
}
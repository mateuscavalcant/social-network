package views

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

func HomeView(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("client/templates/home.html"))
	tmpl.Execute(c.Writer, nil)
}
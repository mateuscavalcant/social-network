package views

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

func OtherProfileView(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("client/templates/other_profile.html"))

	tmpl.Execute(c.Writer, nil)
}
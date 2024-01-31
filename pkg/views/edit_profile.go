package views

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

func EditProfileView(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("client/templates/edit_profile.html"))
	tmpl.Execute(c.Writer, nil)
}
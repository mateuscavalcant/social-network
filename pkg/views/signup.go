package views

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

func SignupView(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("client/templates/signup.html"))
	tmpl.Execute(c.Writer, nil)
}
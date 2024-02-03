package views

import (
	"html/template"

	"github.com/gin-gonic/gin"
)


func RenderProfile(c *gin.Context, tmpl string, data interface{}) {
	t := template.Must(template.ParseFiles("client/templates/" + tmpl))
	t.Execute(c.Writer, data)
}

package views

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

func DeleteAccountView(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("client/templates/delete_account.html"))
	tmpl.Execute(c.Writer, nil)
}
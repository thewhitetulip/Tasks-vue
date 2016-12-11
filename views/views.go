package views

/*Holds the fetch related view handlers*/

import (
	"html/template"
	"net/http"
)

var homeTemplate *template.Template
var templates *template.Template

var message string //message will store the message to be shown as notification
var err error

func Home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, nil)
}

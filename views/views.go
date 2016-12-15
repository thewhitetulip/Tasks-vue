package views

import (
	"html/template"
	"log"
	"net/http"
)

var homeTemplate *template.Template
var templates *template.Template

var message string //message will store the message to be shown as notification
var err error

// Home renders the basic html page to the Vue front end, the Vue front end will then
// update the html content according to the user interactions by making AJAX calls
// we have used Go language's templating mechanism to split templates into logical parts
func Home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, nil)
}

// PopulateTemplates parses all templates present in
// the templates folder and returns a pointer in templates
// we use templates variable to locate other templates
func PopulateTemplates() {
	templates, err = template.ParseGlob("./templates/*.html")
	if err != nil {
		log.Println(err)
	}
	homeTemplate = templates.Lookup("tasks.html")

}

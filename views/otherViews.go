package views

/*
Holds the non insert/update/delete related view handlers
*/

import (
	"html/template"
	"log"
)

//PopulateTemplates is used to parse all templates present in
//the templates folder
func PopulateTemplates() {
	templates, err = template.ParseGlob("./templates/*.html")
	if err != nil {
		log.Println(err)
	}
	homeTemplate = templates.Lookup("tasks.html")

}

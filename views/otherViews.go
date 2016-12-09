package views

/*
Holds the non insert/update/delete related view handlers
*/

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/thewhitetulip/Tasks-vue/db"
	"github.com/thewhitetulip/Tasks-vue/sessions"
	"github.com/thewhitetulip/Tasks-vue/utils"
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

//CompleteTaskFunc is used to show the complete tasks, handles "/completed/" url
func CompleteTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	redirectURL := utils.GetRedirectUrl(r.Referer())
	id, err := strconv.Atoi(r.URL.Path[len("/complete/"):])
	if err != nil {
		log.Println(err)
	} else {
		username := sessions.GetCurrentUserName(r)
		err = db.CompleteTask(username, id)
		if err != nil {
			message = "Complete task failed"
		} else {
			message = "Task marked complete"
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}

}

//SearchTaskFunc is used to handle the /search/ url, handles the search function
func SearchTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	r.ParseForm()
	query := r.Form.Get("query")

	username := sessions.GetCurrentUserName(r)
	context, err := db.SearchTask(username, query)
	if err != nil {
		log.Println("error fetching search results")
	}

	categories := db.GetCategories(username)
	context.Categories = categories

	// searchTemplate.Execute(w, context)

}

//UpdateTaskFunc is used to update a task, handes "/update/" URL
func UpdateTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	r.ParseForm()
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		log.Println(err)
	}
	category := r.Form.Get("category")
	title := r.Form.Get("title")
	content := r.Form.Get("content")
	priority, err := strconv.Atoi(r.Form.Get("priority"))
	if err != nil {
		log.Println(err)
	}
	username := sessions.GetCurrentUserName(r)

	var hidden int
	hideTimeline := r.FormValue("hide")
	if hideTimeline != "" {
		hidden = 1
	} else {
		hidden = 0
	}

	err = db.UpdateTask(id, title, content, category, priority, username, hidden)
	if err != nil {
		message = "Error updating task"
	} else {
		message = "Task updated"
		log.Println(message)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//UpdateCategoryFunc is used to update a task, handes "/upd-category/" URL
func UpdateCategoryFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	var redirectURL string
	r.ParseForm()
	oldName := r.URL.Path[len("/upd-category/"):]
	newName := r.Form.Get("catname")
	username := sessions.GetCurrentUserName(r)
	err := db.UpdateCategoryByName(username, oldName, newName)
	if err != nil {
		message = "error updating category"
		log.Println("not updated category " + oldName)
		redirectURL = "/category/" + oldName
	} else {
		message = "cat " + oldName + " -> " + newName
		redirectURL = "/category/" + newName
	}
	log.Println("redirecting to " + redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusFound)

}

//SignUpFunc will enable new users to sign up to our service
func SignUpFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")
	email := r.Form.Get("email")

	log.Println(username, password, email)

	err := db.CreateUser(username, password, email)
	if err != nil {
		http.Error(w, "Unable to sign user up", http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, "/login/", 302)
	}
}

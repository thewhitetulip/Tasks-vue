package views

/*Holds the fetch related view handlers*/

import (
	"html/template"
	"net/http"
	"time"

	"github.com/thewhitetulip/Tasks-vue/db"
	"github.com/thewhitetulip/Tasks-vue/sessions"

	md "github.com/shurcooL/github_flavored_markdown"
)

var homeTemplate *template.Template
var templates *template.Template

var message string //message will store the message to be shown as notification
var err error

func Home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, nil)
}

// ShowAllTasksFunc is used to handle the "/" URL which is the default ons
// TODO add http404 error
func ShowAllTasksFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	username := sessions.GetCurrentUserName(r)
	context, err := db.GetTasks(username, "pending", "")
	for i := 0; i < len(context.Tasks); i++ {
		context.Tasks[i].Content = string(md.Markdown([]byte(context.Tasks[i].Content)))
	}
	categories := db.GetCategories(username)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	} else {
		if message != "" {
			context.Message = message
		}
		context.CSRFToken = "abcd"
		context.Categories = categories
		message = ""
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "csrftoken", Value: "abcd", Expires: expiration}
		http.SetCookie(w, &cookie)
		homeTemplate.Execute(w, context)
	}
}

// ShowTrashTaskFunc is used to handle the "/trash" URL which is used to show the deleted tasks
func ShowTrashTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	username := sessions.GetCurrentUserName(r)
	categories := db.GetCategories(username)
	context, err := db.GetTasks(username, "deleted", "")
	for i := 0; i < len(context.Tasks); i++ {
		context.Tasks[i].Content = string(md.Markdown([]byte(context.Tasks[i].Content)))
	}
	context.Categories = categories
	if err != nil {
		http.Redirect(w, r, "/trash", http.StatusInternalServerError)
	}
	if message != "" {
		context.Message = message
		message = ""
	}
	// deletedTemplate.Execute(w, context)
}

// ShowCompleteTasksFunc is used to populate the "/completed/" URL
func ShowCompleteTasksFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return

	}

	username := sessions.GetCurrentUserName(r)
	categories := db.GetCategories(username)
	context, err := db.GetTasks(username, "completed", "")
	for i := 0; i < len(context.Tasks); i++ {
		context.Tasks[i].Content = string(md.Markdown([]byte(context.Tasks[i].Content)))
	}
	context.Categories = categories
	if err != nil {
		http.Redirect(w, r, "/completed", http.StatusInternalServerError)
	}
	// completedTemplate.Execute(w, context)
}

// ShowCategoryFunc will populate the /category/<id> URL which shows all the tasks related
// to that particular category
func ShowCategoryFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	category := r.URL.Path[len("/category/"):]
	username := sessions.GetCurrentUserName(r)
	context, err := db.GetTasks(username, "", category)
	for i := 0; i < len(context.Tasks); i++ {
		context.Tasks[i].Content = string(md.Markdown([]byte(context.Tasks[i].Content)))
	}
	categories := db.GetCategories(username)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}
	if message != "" {
		context.Message = message
	}
	context.CSRFToken = "abcd"
	context.Categories = categories
	message = ""
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "csrftoken", Value: "abcd", Expires: expiration}
	http.SetCookie(w, &cookie)
	homeTemplate.Execute(w, context)
}

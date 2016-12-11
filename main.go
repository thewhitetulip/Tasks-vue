package main

/**
 * This is the main file for the Task application
 * License: MIT
 **/
import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/thewhitetulip/Tasks-vue/config"
	"github.com/thewhitetulip/Tasks-vue/views"

	"github.com/gorilla/mux"
)

func main() {
	values, err := config.ReadConfig("config.json")
	var port *string

	if err != nil {
		port = flag.String("port", "", "IP address")
		flag.Parse()

		// User is expected to give :8080 like input, if they give 8080
		// we'll append the required ':'.
		if !strings.HasPrefix(*port, ":") {
			*port = ":" + *port
			log.Println("port is " + *port)
		}

		values.ServerPort = *port
	}

	views.PopulateTemplates()

	r := mux.NewRouter()

	http.Handle("/static/", http.FileServer(http.Dir("public")))

	r.HandleFunc("/api/task/", views.GetTasksFuncAPI).Methods("GET")
	r.HandleFunc("/api/task/", views.AddTaskFuncAPI).Methods("PUT")
	r.HandleFunc("/api/task/", views.UpdateTaskFuncAPI).Methods("POST")
	r.HandleFunc("/api/task/{id}", views.TrashTaskFuncAPI).Methods("DELETE")
	r.HandleFunc("/api/deleted/", views.GetDeletedTaskFuncAPI).Methods("GET")

	r.HandleFunc("/api/completed/", views.GetCompletedTaskFuncAPI).Methods("GET")
	r.HandleFunc("/api/categories/", views.GetCategoryFuncAPI).Methods("GET")

	r.HandleFunc("/api/category/{category}", views.ShowCategoryFuncAPI).Methods("GET")
	r.HandleFunc("/api/category/{category}", views.DeleteCategoryFuncAPI).Methods("DELETE")
	r.HandleFunc("/api/complete-task/{id}", views.CompleteTaskFuncAPI).Methods("GET")
	r.HandleFunc("/api/incomplete-task/{id}", views.RestoreFromCompleteFuncAPI).Methods("GET")
	r.HandleFunc("/api/restore-task/{id}", views.RestoreTaskFuncAPI).Methods("GET")
	r.HandleFunc("/", views.Home).Methods("GET")

	r.HandleFunc("/api/comment/", views.AddCommentFuncAPI).Methods("PUT")
	r.HandleFunc("/api/comment/{id}", views.DeleteCommentFuncAPI).Methods("DELETE")
	r.HandleFunc("/api/category/", views.AddCategoryFuncAPI).Methods("PUT")
	r.HandleFunc("/api/category/{category}", views.UpdateCategoryFuncAPI).Methods("POST")

	//Login logout
	r.HandleFunc("/api/login/", views.LoginFuncAPI).Methods("POST", "GET")
	r.HandleFunc("/api/logout/", views.RequiresLogin(views.LogoutFuncAPI)).Methods("GET")
	r.HandleFunc("/api/signup/", views.SignUpFuncAPI).Methods("POST")

	//these handlers perform action like delete, mark as complete etc
	// http.HandleFunc("/files/", views.RequiresLogin(views.UploadedFileHandler))
	// http.HandleFunc("/search/", views.RequiresLogin(views.SearchTaskFunc))
	http.Handle("/", r)
	log.Println("running server on ", values.ServerPort)
	log.Fatal(http.ListenAndServe(values.ServerPort, nil))
}

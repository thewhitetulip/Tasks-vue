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
	r.HandleFunc("/api/task/", views.DeleteTaskFuncAPI).Methods("DELETE")
	r.HandleFunc("/api/deleted/", views.GetDeletedTaskFuncAPI).Methods("GET")

	r.HandleFunc("/api/completed/", views.GetCompletedTaskFuncAPI).Methods("GET")
	r.HandleFunc("/api/categories/", views.GetCategoryFuncAPI).Methods("GET")

	r.HandleFunc("/api/category/{category}", views.ShowCategoryFuncAPI).Methods("GET")
	r.HandleFunc("/api/category/{category}", views.DeleteCategoryFuncAPI).Methods("DELETE")
	r.HandleFunc("/api/trash-task/{id}", views.TrashTaskFuncAPI).Methods("GET")
	r.HandleFunc("/api/complete-task/{id}", views.CompleteTaskFuncAPI).Methods("GET")
	r.HandleFunc("/api/incomplete-task/{id}", views.RestoreFromCompleteFuncAPI).Methods("GET")
	r.HandleFunc("/api/restore-task/{id}", views.RestoreTaskFuncAPI).Methods("GET")
	r.HandleFunc("/", views.Home).Methods("GET")

	r.HandleFunc("/api/comment/", views.AddCommentFunc).Methods("PUT")
	// r.HandleFunc("/api/get-token/", views.GetTokenHandler).Methods("POST")
	r.HandleFunc("/api/category/", views.AddCategoryFuncAPI).Methods("PUT")
	r.HandleFunc("/api/category/{category}", views.UpdateCategoryFuncAPI).Methods("POST")
	// r.HandleFunc("/api/delete-category/", views.DeleteCategoryFuncAPI).Methods("DELETE")

	//Login logout
	// http.HandleFunc("/login/", views.LoginFunc)
	// http.HandleFunc("/logout/", views.RequiresLogin(views.LogoutFunc))
	// http.HandleFunc("/signup/", views.SignUpFunc)

	// http.HandleFunc("/add-category/", views.RequiresLogin(views.AddCategoryFunc))
	// http.HandleFunc("/add/", views.RequiresLogin(views.AddTaskFunc))

	// //these handlers are used to delete
	// http.HandleFunc("/del-comment/", views.RequiresLogin(views.DeleteCommentFunc))
	// http.HandleFunc("/del-category/", views.RequiresLogin(views.DeleteCategoryFunc))
	// http.HandleFunc("/delete/", views.RequiresLogin(views.DeleteTaskFunc))

	// //these handlers update
	// http.HandleFunc("/upd-category/", views.RequiresLogin(views.UpdateCategoryFunc))
	// http.HandleFunc("/update/", views.RequiresLogin(views.UpdateTaskFunc))

	// //these handlers are used for restoring tasks

	// //these handlers fetch set of tasks
	// http.HandleFunc("/", views.RequiresLogin(views.ShowAllTasksFunc))
	// http.HandleFunc("/category/", views.RequiresLogin(views.ShowCategoryFunc))
	// http.HandleFunc("/deleted/", views.RequiresLogin(views.ShowTrashTaskFunc))
	// http.HandleFunc("/completed/", views.RequiresLogin(views.ShowCompleteTasksFunc))

	// //these handlers perform action like delete, mark as complete etc
	// http.HandleFunc("/files/", views.RequiresLogin(views.UploadedFileHandler))
	// http.HandleFunc("/trash/", views.RequiresLogin(views.TrashTaskFunc))
	// http.HandleFunc("/edit/", views.RequiresLogin(views.EditTaskFunc))
	// http.HandleFunc("/search/", views.RequiresLogin(views.SearchTaskFunc))
	http.Handle("/", r)
	log.Println("running server on ", values.ServerPort)
	log.Fatal(http.ListenAndServe(values.ServerPort, nil))
}

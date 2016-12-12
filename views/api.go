package views

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/thewhitetulip/Tasks-vue/db"
	"github.com/thewhitetulip/Tasks-vue/sessions"
	"github.com/thewhitetulip/Tasks-vue/types"
)

const unableToProcess = "Something went wrong"

//GetTasksFuncAPI fetches tasks depending on the request, the authorization will be taken care by our middleare
//in this function we will return all the tasks to the user or tasks per category
//GET /tasks/
func GetTasksFuncAPI(w http.ResponseWriter, r *http.Request) {
	var strTaskID string
	var err error
	var task types.Task
	var tasks types.Tasks
	var status types.Status
	isError := false

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	username := sessions.GetCurrentUserName(r)
	log.Println("token is valid " + username + " is logged in")

	strTaskID = r.URL.Path[len("/task/"):]
	//this is when we get a request for all the tasks for that user
	if strTaskID == "" {
		context, err := db.GetTasks(username, "pending", "")
		if err != nil {
			isError = true
		} else {
			tasks = context.Tasks
			for _, tasks := range tasks {
				tasks.ShowComment = false
			}
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(tasks)

			if err != nil {
				panic(err)
			}
			return
		}
	} else {
		//this is when we get a request for a particular task
		taskID, err := strconv.Atoi(strTaskID)
		if err != nil {
			isError = true
		} else {
			ctx, err := db.GetTaskByID(username, taskID)
			task = ctx.Tasks[0]
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(task)
			if err != nil {
				panic(err)
			}
			return
		}
	}

	if isError {
		log.Println("GetTasksFuncAPI: api.go: taskID")
		status = types.Status{http.StatusInternalServerError, unableToProcess}
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(status)

		if err != nil {
			panic(err)
		}
	}
}

//AddTaskFuncAPI will add the tasks for the user
func AddTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	var hidden = 0 // assume that it is not hidden, if form element is true, set this to 1

	statusCode := http.StatusOK
	message := "Task added to db"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	username := sessions.GetCurrentUserName(r)

	r.ParseForm()
	category := r.FormValue("category")
	title := template.HTMLEscapeString(r.Form.Get("title"))
	content := template.HTMLEscapeString(r.Form.Get("content"))
	taskPriority, priorityErr := strconv.Atoi(r.FormValue("priority"))
	hideTimeline := r.FormValue("ishidden")

	if priorityErr != nil {
		log.Print(priorityErr)
		message = "Bad task priority"
	}
	priorityList := []int{1, 2, 3}
	found := false
	for _, priority := range priorityList {
		if taskPriority == priority {
			found = true
		}
	}
	//If someone gives us incorrect priority number, we give the priority
	//to that task as 1 i.e. Low
	if !found {
		taskPriority = 1
	}

	if hideTimeline == "true" {
		hidden = 1
	}

	if title != "" && content != "" {
		taskTruth := db.AddTask(title, content, category, taskPriority, username, hidden)
		if taskTruth != nil {
			statusCode = http.StatusInternalServerError
			message = unableToProcess
		}
	}

	status := types.Status{statusCode, message}
	json.NewEncoder(w).Encode(status)
}

//UpdateTaskFuncAPI will add the tasks for the user
func UpdateTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	taskErr := false
	statusCode := http.StatusOK
	var hidden = 0
	message := "updated task id "

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	username := sessions.GetCurrentUserName(r)

	r.ParseForm()
	strID := r.Form.Get("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		log.Println(err)
		taskErr = true
	}

	category := r.Form.Get("category")
	title := r.Form.Get("title")
	content := r.Form.Get("content")
	priority, err := strconv.Atoi(r.Form.Get("priority"))
	hideTimeline := r.FormValue("ishidden")

	if err != nil {
		log.Println(err)
		priority = 1
	}

	if hideTimeline == "true" {
		hidden = 1
	}

	if strID != "" && title != "" && content != "" {
		err = db.UpdateTask(id, title, content, category, priority, username, hidden)
		if err != nil {
			taskErr = true
		}
	} else {
		taskErr = true
	}

	if taskErr {
		statusCode = http.StatusBadRequest
		message = unableToProcess
	}

	status := types.Status{statusCode, message}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
}

//GetCompletedTaskFuncAPI will get the deleted tasks for the user
func GetCompletedTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	var tasks types.Tasks
	var status types.Status
	var statusCode = http.StatusOK

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	username := sessions.GetCurrentUserName(r)

	//this is when we get a request for all the deleted tasks for that user
	context, err := db.GetTasks(username, "completed", "")

	if err != nil {
		log.Println("GetTasksFuncAPI: api.go: Server error")
		statusCode = http.StatusInternalServerError
		status = types.Status{statusCode, unableToProcess}
		w.WriteHeader(statusCode)
		err = json.NewEncoder(w).Encode(status)

		if err != nil {
			panic(err)
		}
		return
	}

	tasks = context.Tasks
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(tasks)

	if err != nil {
		panic(err)
	}
	return
}

//GetDeletedTaskFuncAPI will get the deleted tasks for the user
func GetDeletedTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	var tasks types.Tasks
	var status types.Status
	var statusCode = http.StatusOK

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	username := sessions.GetCurrentUserName(r)
	log.Println("token is valid " + username + " is logged in")

	// this is when we get a request for all the deleted tasks for that user
	context, err := db.GetTasks(username, "deleted", "")
	if err != nil {
		log.Println("GetTasksFuncAPI: api.go: Server error")
		statusCode = http.StatusInternalServerError
		status = types.Status{statusCode, unableToProcess}
		w.WriteHeader(statusCode)
		err = json.NewEncoder(w).Encode(status)

		if err != nil {
			panic(err)
		}
		return
	}

	tasks = context.Tasks
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(tasks)

	if err != nil {
		panic(err)
	}
	return
}

//GetCategoryFuncAPI will return the categories for the user
//depends on the ID that we get, if we get all, then return all categories of the user as a JSON.
func GetCategoryFuncAPI(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusOK

	username := sessions.GetCurrentUserName(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	categories, err := db.GetCategories(username)
	if err != nil {
		statusCode = http.StatusInternalServerError
		status := types.Status{statusCode, unableToProcess}
		w.WriteHeader(statusCode)
		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}
		return
	}

	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(categories)
	if err != nil {
		panic(err)
	}

}

//AddCategoryFuncAPI will add the category for the user
func AddCategoryFuncAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	var message = "Added Category"
	var statusCode = http.StatusOK

	r.ParseForm()

	category := r.Form.Get("categoryName")
	if strings.Trim(category, " ") != "" {
		username := sessions.GetCurrentUserName(r)
		log.Println("adding category")
		err := db.AddCategory(username, category)
		if err != nil {
			statusCode = http.StatusInternalServerError
			message = unableToProcess
		}
	} else {
		statusCode = http.StatusInternalServerError
		message = unableToProcess
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	status := types.Status{statusCode, message}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
}

//UpdateCategoryFuncAPI will update the category for the user
func UpdateCategoryFuncAPI(w http.ResponseWriter, r *http.Request) {
	message := "Category Updated"
	statusCode := http.StatusOK
	var err error

	r.ParseForm()

	oldName := r.URL.Path[len("/category/"):]
	oldName = strings.Trim(oldName, "/")

	newName := r.Form.Get("newCategoryName")

	if strings.Trim(newName, " ") == "" {
		statusCode = http.StatusBadRequest
		message = unableToProcess
	} else {
		username := sessions.GetCurrentUserName(r)
		err = db.UpdateCategoryByName(username, oldName, newName)
		if err != nil {
			message = unableToProcess
			statusCode = http.StatusInternalServerError
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	status := types.Status{statusCode, message}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
}

//DeleteCategoryFuncAPI will delete the category for the user
func DeleteCategoryFuncAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	message := "Category Deleted"
	var statusCode = http.StatusOK

	categoryName := r.URL.Path[len("/category/"):]
	categoryName = strings.Trim(categoryName, "/")
	categoryName = strings.Trim(categoryName, " ")

	if categoryName == "" {
		message = unableToProcess
	} else {
		username := sessions.GetCurrentUserName(r)
		err = db.DeleteCategoryByName(username, categoryName)
		if err != nil {
			message = unableToProcess
			statusCode = http.StatusInternalServerError
			log.Println(err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	status := types.Status{statusCode, message}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
}

// ShowCategoryFuncAPI will return all the tasks of a particular category
// we will be returning a status internal server error in case we do not find the
// tasks of that category, the url it will handle is GET /categories/<value>; if value is nil it'll return a JSON error
func ShowCategoryFuncAPI(w http.ResponseWriter, r *http.Request) {
	message := "Success"
	statusCode := http.StatusOK

	category := r.URL.Path[len("/category/"):]

	if category == "" {
		statusCode = http.StatusBadRequest
		message = unableToProcess
	} else {
		username := sessions.GetCurrentUserName(r)
		log.Println("fetching tasks for " + category)
		context, err := db.GetTasks(username, "", category)
		if err != nil {
			log.Println("ShowCategoryFuncAPI: api.go: Server error")
			message = unableToProcess
			statusCode = http.StatusInternalServerError
			log.Println(message)
		} else {
			err = json.NewEncoder(w).Encode(context.Tasks)
			if err != nil {
				panic(err)
			}
			return
		}

		/*for i := 0; i < len(context.Tasks); i++ {
			context.Tasks[i].Content = string(md.Markdown([]byte(context.Tasks[i].Content)))
		}*/

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	status := types.Status{statusCode, message}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}

}

//CompleteTaskFuncAPI will delete a task which is passed as an ID
func CompleteTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusOK
	message := "Task completed"

	id, err := strconv.Atoi(r.URL.Path[len("/complete-task/"):])
	if err != nil {
		log.Println("CompleteTaskFuncAPI", err)
		message = unableToProcess
		statusCode = http.StatusBadRequest
	} else {
		username := sessions.GetCurrentUserName(r)
		err = db.CompleteTask(username, id)
		if err != nil {
			message = unableToProcess
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	status := types.Status{statusCode, message}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
	return
}

func RestoreTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	var statusCode = http.StatusOK
	var message = "Task restored from Trash"

	id, err := strconv.Atoi(r.URL.Path[len("/restore-task/"):])

	if err != nil {
		log.Println("RestoreTaskFunAPI", err)
		message = unableToProcess
		statusCode = http.StatusBadRequest
	} else {
		username := sessions.GetCurrentUserName(r)
		err = db.RestoreTask(username, id)
		if err != nil {
			message = unableToProcess
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	status := types.Status{statusCode, message}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
	return

}

// TrashTaskFuncAPI handles the GET /trash-task/ and trashes the ID passed in the URL.
// return JSON {http status code, output of operation}.
func TrashTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	var statusCode = http.StatusOK
	var message = "Task Trashed"

	id, err := strconv.Atoi(r.URL.Path[len("/task/"):])
	if err != nil {
		log.Println("TrashTaskFunc", err)
		message = unableToProcess
		statusCode = http.StatusBadRequest
	} else {
		username := sessions.GetCurrentUserName(r)
		err = db.TrashTask(username, id)
		if err != nil {
			message = unableToProcess
			statusCode = http.StatusInternalServerError
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	status := types.Status{statusCode, message}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
	return
}

// RestoreFromCompleteFuncAPI handles the GET /incomplete-task/ and restores the status of task
// from complete to pending.
// Returns a JSON {http status code, output} like {400, "Not Deleted"} or {200, "Comment deleted"}.
// The status code is also written in the HTTP header of the response.
func RestoreFromCompleteFuncAPI(w http.ResponseWriter, r *http.Request) {
	var statusCode = http.StatusOK
	var message = "Marked Incomplete"

	id, err := strconv.Atoi(r.URL.Path[len("/incomplete-task/"):])

	if err != nil {
		log.Println("api.go: RestoreFromComplete", err)
		message = unableToProcess
		statusCode = http.StatusBadRequest
	} else {
		username := sessions.GetCurrentUserName(r)
		err = db.RestoreTaskFromComplete(username, id)
		if err != nil {
			message = unableToProcess
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	status := types.Status{statusCode, message}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
	return
}

// AddCommentFuncAPI handles the PUT /comment/ and adds a new comment to the database.
// Each comment has a parent task, comments can't be parents of comments.
// Returns a JSON {http status code, output} like {400, "Unable to add comment"} or {200, "Comment Added"}.
// The status code is also written in the HTTP header of the response.
func AddCommentFuncAPI(w http.ResponseWriter, r *http.Request) {
	var statusCode = http.StatusOK
	var message = "Comment Added"

	r.ParseForm()
	text := r.Form.Get("content")
	taskID := r.Form.Get("taskID")

	taskIDInt, err := strconv.Atoi(taskID)

	if err != nil {
		message = unableToProcess
		statusCode = http.StatusBadRequest
		log.Println(err)
	} else {
		username := sessions.GetCurrentUserName(r)
		err = db.AddComments(username, taskIDInt, text)
		statusCode = http.StatusOK
		if err != nil {
			log.Println("unable to insert into db")
			message = unableToProcess
			statusCode = http.StatusInternalServerError
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	status := types.Status{statusCode, message}
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
	return
}

// DeleteCommentFuncAPI handles the DELETE /comment/12 and deletes a comment with ID 12.
// Returns a JSON {http status code, output} like {400, "Not Deleted"} or {200, "Comment deleted"}.
// The status code is also written in the HTTP header of the response.
func DeleteCommentFuncAPI(w http.ResponseWriter, r *http.Request) {
	var statusCode = http.StatusOK
	var message = "Comment deleted"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Path[len("/comment/"):]

	commentID, err := strconv.Atoi(id)

	if err != nil {
		message = unableToProcess
	} else {
		username := sessions.GetCurrentUserName(r)
		err = db.DeleteCommentByID(username, commentID)
		if err != nil {
			message = unableToProcess
			statusCode = http.StatusInternalServerError
		}
	}

	w.WriteHeader(statusCode)
	status := types.Status{statusCode, message}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
	return
}

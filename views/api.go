package views

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/thewhitetulip/Tasks-vue/db"
	"github.com/thewhitetulip/Tasks-vue/types"
)

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var mySigningKey = []byte("secret")

//GetTokenHandler will get a token for the username and password
func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Write([]byte("Method not allowed"))
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	log.Println(username, " ", password)
	if username == "" || password == "" {
		w.Write([]byte("Invalid Username or password"))
		return
	}
	if db.ValidUser(username, password) {
		/* Set token claims */

		// Create the Claims
		claims := MyCustomClaims{
			username,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		/* Sign the token with our secret */
		tokenString, err := token.SignedString(mySigningKey)
		if err != nil {
			log.Println("Something went wrong with signing token")
			w.Write([]byte("Authentication failed"))
			return
		}

		/* Finally, write the token to the browser window */
		w.Write([]byte(tokenString))
	} else {
		w.Write([]byte("Authentication failed"))
	}
}

//ValidateToken will validate the token
func ValidateToken(myToken string) (bool, string) {
	token, err := jwt.ParseWithClaims(myToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})

	if err != nil {
		return false, ""
	}

	claims := token.Claims.(*MyCustomClaims)
	return token.Valid, claims.Username
}

//GetTasksFuncAPI fetches tasks depending on the request, the authorization will be taken care by our middleare
//in this function we will return all the tasks to the user or tasks per category
//GET /api/get-tasks/
func GetTasksFuncAPI(w http.ResponseWriter, r *http.Request) {
	var strTaskID string
	var err error
	var message string
	var task types.Task
	var tasks types.Tasks
	var status types.Status

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	username := "suraj"

	/*token := r.Header["Token"][0]

	IsTokenValid, username := ValidateToken(token)
	//When the token is not valid show the default error JSON document
	if !IsTokenValid {
		status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(status)

		if err != nil {
			panic(err)
		}
		return
	}*/

	log.Println("token is valid " + username + " is logged in")

	strTaskID = r.URL.Path[len("/api/task/"):]
	//this is when we get a request for all the tasks for that user
	if strTaskID == "" {
		context, err := db.GetTasks(username, "pending", "")
		if err != nil {
			message = "GetTasksFuncAPI: api.go: Server error"
			log.Println(message)
			status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

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
	//this is when we get a request for a particular task

	taskID, err := strconv.Atoi(strTaskID)
	if err != nil {
		message = "GetTasksFuncAPI: api.go: Invalid taskID " + strTaskID
		log.Println(message)

		status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(status)

		if err != nil {
			panic(err)
		}
		return
	}
	ctx, err := db.GetTaskByID(username, taskID)
	task = ctx.Tasks[0]

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		panic(err)
	}
}

//AddTaskFuncAPI will add the tasks for the user
func AddTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	// token := r.Header["Token"][0]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	username := "suraj"
	// IsTokenValid, username := ValidateToken(token)
	//When the token is not valid show the default error JSON document
	// if !IsTokenValid {
	// 	status := types.Status{StatusCode: http.StatusInternalServerError, Message: message}
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	err = json.NewEncoder(w).Encode(status)
	//
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		return
	// 	}

	log.Println("token is valid " + username + " is logged in")

	r.ParseForm()
	category := r.FormValue("category")
	title := template.HTMLEscapeString(r.Form.Get("title"))
	content := template.HTMLEscapeString(r.Form.Get("content"))
	taskPriority, priorityErr := strconv.Atoi(r.FormValue("priority"))

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
	var hidden int
	hideTimeline := r.FormValue("ishidden")
	log.Println(hideTimeline)
	if hideTimeline == "true" {
		hidden = 1
	} else {
		hidden = 0
	}
	var taskErr bool

	if title != "" && content != "" {
		taskTruth := db.AddTask(title, content, category, taskPriority, username, hidden)
		if taskTruth != nil {
			taskErr = true
		}
	}

	var statusCode int
	var message string

	if !taskErr {
		statusCode = http.StatusInternalServerError
		message = "Error adding task to db"
	} else {
		statusCode = http.StatusOK
		message = "Task added to db"
	}
	status := types.Status{StatusCode: statusCode, Message: message}
	json.NewEncoder(w).Encode(status)
	log.Println("added task to db")
}

//UpdateTaskFuncAPI will add the tasks for the user
func UpdateTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	var taskErr bool
	// token := r.Header["Token"][0]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	username := "suraj"
	// IsTokenValid, username := ValidateToken(token)
	// When the token is not valid show the default error JSON document
	// if !IsTokenValid {
	// status := types.Status{StatusCode: http.StatusInternalServerError, Message: message}
	// w.WriteHeader(http.StatusInternalServerError)
	// err = json.NewEncoder(w).Encode(status)
	//
	// if err != nil {
	// panic(err)
	// }
	// return
	// }

	log.Println("update func token is valid " + username + " is logged in")

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

	if err != nil {
		log.Println(err)
		priority = 1
	}

	var hidden int
	hideTimeline := r.FormValue("hide")
	if hideTimeline != "" {
		hidden = 1
	} else {
		hidden = 0
	}

	if strID != "" && title != "" && content != "" {
		err = db.UpdateTask(id, title, content, category, priority, username, hidden)
		if err != nil {
			taskErr = true
		}
		taskErr = false
	} else {
		taskErr = true
	}

	var statusCode int
	var message string

	if taskErr {
		statusCode = http.StatusBadRequest
		message = "unable to update task id "
	} else {
		statusCode = http.StatusOK
		message = "updated task id "
	}

	status := types.Status{StatusCode: statusCode, Message: message}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
}

//DeleteTaskFuncAPI will add the tasks for the user
func DeleteTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		token := r.Header["Token"][0]
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		IsTokenValid, username := ValidateToken(token)
		//When the token is not valid show the default error JSON document
		if !IsTokenValid {
			status := types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		}

		log.Println("token is valid " + username + " is logged in")
		var strtaskID string
		strtaskID = r.URL.Path[len("/api/delete-task/"):]
		var statusCode int
		var message string

		taskID, err := strconv.Atoi(strtaskID)
		if err != nil {
			log.Println("invalid task id")
			statusCode = http.StatusBadRequest
			message = "invalid TaskID"
		} else {
			db.TrashTask(username, taskID)
			if err != nil {
				statusCode = http.StatusOK
				message = "deleted task id " + strtaskID
			} else {
				statusCode = http.StatusOK
				message = "deleted task id " + strtaskID
			}
		}

		status := types.Status{StatusCode: statusCode, Message: message}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}
	}
}

//GetCompletedTaskFuncAPI will get the deleted tasks for the user
func GetCompletedTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var err error
		var message string
		var tasks types.Tasks
		var status types.Status

		//token := r.Header["Token"][0]

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		username := "suraj"
		//IsTokenValid, username := ValidateToken(token)
		//When the token is not valid show the default error JSON document
		//if !IsTokenValid {
		//	status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
		//	w.WriteHeader(http.StatusInternalServerError)
		//	err = json.NewEncoder(w).Encode(status)
		//
		//			if err != nil {
		//				panic(err)
		//			}
		//			return
		//		}

		log.Println("token is valid " + username + " is logged in")

		//this is when we get a request for all the deleted tasks for that user
		context, err := db.GetTasks(username, "completed", "")
		if err != nil {
			message = "GetTasksFuncAPI: api.go: Server error"
			log.Println(message)
			status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		}

		tasks = context.Tasks
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(tasks)

		if err != nil {
			panic(err)
		}
		return
	}
}

//GetDeletedTaskFuncAPI will get the deleted tasks for the user
func GetDeletedTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var err error
		var message string
		var tasks types.Tasks
		var status types.Status

		//token := r.Header["Token"][0]

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		username := "suraj"
		//IsTokenValid, username := ValidateToken(token)
		//When the token is not valid show the default error JSON document
		//if !IsTokenValid {
		//	status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
		//	w.WriteHeader(http.StatusInternalServerError)
		//	err = json.NewEncoder(w).Encode(status)
		//
		//			if err != nil {
		//				panic(err)
		//			}
		//			return
		//		}

		log.Println("token is valid " + username + " is logged in")

		//this is when we get a request for all the deleted tasks for that user
		context, err := db.GetTasks(username, "deleted", "")
		if err != nil {
			message = "GetTasksFuncAPI: api.go: Server error"
			log.Println(message)
			status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		}

		tasks = context.Tasks
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(tasks)

		if err != nil {
			panic(err)
		}
		return
	}
}

//GetCategoryFuncAPI will return the categories for the user
//depends on the ID that we get, if we get all, then return all categories of the user as a JSON.
func GetCategoryFuncAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	//var message string
	//var status types.Status

	//	token := r.Header["Token"][0]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	username := "suraj"

	/*	IsTokenValid, username := ValidateToken(token)
		//When the token is not valid show the default error JSON document
		if !IsTokenValid {
			status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		} */

	log.Println("token is valid " + username + " is logged in")
	categories, _ := db.GetCategories(username)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(categories)
	if err != nil {
		panic(err)
	}

}

//AddCategoryFuncAPI will add the category for the user
func AddCategoryFuncAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	var message string
	var status types.Status
	var statusCode int
	var categoryErr bool

	// token := r.Header["Token"][0]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	username := "suraj"

	//IsTokenValid, username := ValidateToken(token)
	//When the token is not valid show the default error JSON document
	//if !IsTokenValid {
	//	status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
	//	w.WriteHeader(http.StatusInternalServerError)
	//	err = json.NewEncoder(w).Encode(status)

	//	if err != nil {
	//		panic(err)
	//	}
	//	return
	//}

	log.Println("token is valid " + username + " is logged in")
	r.ParseForm()

	category := r.Form.Get("categoryName")
	if strings.Trim(category, " ") != "" {
		log.Println("adding category")
		err := db.AddCategory(username, category)
		if err != nil {
			categoryErr = true
		} else {
			categoryErr = false
		}
	} else {
		categoryErr = true
	}

	if categoryErr {
		statusCode = http.StatusInternalServerError
		message = "error adding category" + category
	} else {
		statusCode = http.StatusOK
		message = "added category " + category
	}

	status = types.Status{StatusCode: statusCode, Message: message}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
}

//UpdateCategoryFuncAPI will update the category for the user
func UpdateCategoryFuncAPI(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var err error
	var message string
	var catErr bool
	var status types.Status

	//token := r.Header["Token"][0]
	username := "suraj"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	/*	IsTokenValid, username := ValidateToken(token)
		//When the token is not valid show the default error JSON document
		if !IsTokenValid {
			status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		}
	*/
	log.Println("token is valid " + username + " is logged in")

	r.ParseForm()

	oldName := r.URL.Path[len("/api/category/"):]
	oldName = strings.Trim(oldName, "/")

	newName := r.Form.Get("newCategoryName")

	if strings.Trim(newName, " ") != "" {
		err = db.UpdateCategoryByName(username, oldName, newName)
		if err != nil {
			catErr = true
		}
		catErr = false
	} else {
		catErr = true
	}

	if catErr {
		statusCode = http.StatusInternalServerError
		message = "unable to update category from " + oldName + " to " + newName
	} else {
		statusCode = http.StatusOK
		message = "updated category from " + oldName + " to " + newName

	}

	status = types.Status{StatusCode: statusCode, Message: message}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
}

//DeleteCategoryFuncAPI will delete the category for the user
func DeleteCategoryFuncAPI(w http.ResponseWriter, r *http.Request) {
	categoryName := r.URL.Path[len("/api/category/"):]
	categoryName = strings.Trim(categoryName, "/")
	var statusCode int
	var err error
	var message string
	var catErr bool

	// token := r.Header["Token"][0]
	username := "suraj"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//	IsTokenValid, username := ValidateToken(token)
	//When the token is not valid show the default error JSON document
	//	if !IsTokenValid {
	//		status := types.Status{StatusCode: http.StatusInternalServerError, Message: message}
	//		w.WriteHeader(http.StatusInternalServerError)
	//		err = json.NewEncoder(w).Encode(status)

	//		if err != nil {
	//			panic(err)
	//		}
	//		return
	//	}

	log.Println("token is valid " + username + " is logged in")

	categoryName = strings.Trim(categoryName, " ")

	if categoryName != "" {
		catErr = true
	}
	err = db.DeleteCategoryByName(username, categoryName)
	if err != nil {
		catErr = true
		log.Println(err)
	} else {
		catErr = false
	}

	if catErr {
		statusCode = http.StatusBadRequest
		message = "error deleting category" + categoryName
	} else {
		statusCode = http.StatusOK
		message = "deleted category " + categoryName
	}

	status := types.Status{StatusCode: statusCode, Message: message}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
}

// ShowCategoryFuncAPI will return all the tasks of a particular category
// we will be returning a status internal server error in case we do not find the
// tasks of that category, the url it will handle is GET /api/categories/<value>; if value is nil it'll return a JSON error
func ShowCategoryFuncAPI(w http.ResponseWriter, r *http.Request) {
	var status types.Status
	category := r.URL.Path[len("/api/category/"):]
	username := "suraj"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if category == "" {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(types.Status{StatusCode: http.StatusInternalServerError, Message: "Invalid Request"})
		return
	}

	log.Println("fetching tasks for " + category)

	context, err := db.GetTasks(username, "", category)

	if err != nil {
		message = "ShowCategoryFuncAPI: api.go: Server error"
		log.Println(message)
		status = types.Status{StatusCode: http.StatusInternalServerError, Message: "error fetching categories"}
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}
		return
	}

	//for i := 0; i < len(context.Tasks); i++ {
	//	context.Tasks[i].Content = string(md.Markdown([]byte(context.Tasks[i].Content)))
	//}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(context.Tasks)
	if err != nil {
		panic(err)
	}

}

//CompleteTaskFuncAPI will delete a task which is passed as an ID
func CompleteTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	var status types.Status
	id, err := strconv.Atoi(r.URL.Path[len("/api/complete-task/"):])
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		log.Println("TrashTaskFunc", err)
		message = "Incorrect command"
		status = types.Status{StatusCode: http.StatusBadRequest, Message: "Invalid Task"}
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}
		return

	}
	username := "suraj"
	err = db.CompleteTask(username, id)
	if err != nil {
		message = "Error completing task"
	} else {
		message = "Task completed"
	}
	w.WriteHeader(http.StatusOK)
	status = types.Status{http.StatusOK, message}
	err = json.NewEncoder(w).Encode(status)

	if err != nil {
		panic(err)
	}
	return

}
func RestoreTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	var status types.Status
	id, err := strconv.Atoi(r.URL.Path[len("/api/restore-task/"):])
	log.Println("restore" + string(id))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		log.Println("TrashTaskFunc", err)
		message = "Incorrect command"
		status = types.Status{StatusCode: http.StatusBadRequest, Message: "Invalid Task"}
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}
		return

	}
	username := "suraj"
	err = db.RestoreTask(username, id)
	if err != nil {
		message = "Error trashing task"
	} else {
		message = "Task trashed"
	}
	w.WriteHeader(http.StatusOK)
	status = types.Status{http.StatusOK, "task trashed"}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
	return

}

// TrashTaskFuncAPI will delete a task which is passed as an ID
func TrashTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	var status types.Status
	id, err := strconv.Atoi(r.URL.Path[len("/api/trash-task/"):])
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		log.Println("TrashTaskFunc", err)
		message = "Incorrect command"
		status = types.Status{StatusCode: http.StatusBadRequest, Message: "Invalid Task"}
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}
		return

	}
	username := "suraj"
	err = db.TrashTask(username, id)
	if err != nil {
		message = "Error trashing task"
	} else {
		message = "Task trashed"
	}
	w.WriteHeader(http.StatusOK)
	status = types.Status{http.StatusOK, "task trashed"}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
	return
}

// RestoreFromCompleteFuncAPI will delete a task which is passed as an ID
func RestoreFromCompleteFuncAPI(w http.ResponseWriter, r *http.Request) {
	var status types.Status
	id, err := strconv.Atoi(r.URL.Path[len("/api/incomplete-task/"):])
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		log.Println("TrashTaskFunc", err)
		message = "Incorrect command"
		status = types.Status{StatusCode: http.StatusBadRequest, Message: "Invalid Task"}
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}
		return

	}
	username := "suraj"
	err = db.RestoreTaskFromComplete(username, id)
	if err != nil {
		message = "Error trashing task"
	} else {
		message = "Task trashed"
	}
	w.WriteHeader(http.StatusOK)
	status = types.Status{http.StatusOK, "task trashed"}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
	return
}

//AddCommentFuncAPI will be used
func AddCommentFuncAPI(w http.ResponseWriter, r *http.Request) {
	var status types.Status
	var message string

	r.ParseForm()
	text := r.Form.Get("content")
	id := r.Form.Get("taskID")

	idInt, err := strconv.Atoi(id)
	log.Println(r.Form)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		if text == "" {
			message = "blank comment"
		} else {
			message = "Error adding comment"

		}
		log.Println(err)
		status = types.Status{http.StatusBadRequest, message}
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(status)
		return
	}
	//username := sessions.GetCurrentUserName(r)
	username := "suraj"
	err = db.AddComments(username, idInt, text)
	htstatus := http.StatusOK
	if err != nil {
		log.Println("unable to insert into db")
		message = "Comment not added"
		htstatus = http.StatusInternalServerError
	} else {
		message = "Comment added"
	}
	status = types.Status{htstatus, message}
	w.WriteHeader(htstatus)
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
	return
}

//DeleteCommentFuncAPI will delete any category
func DeleteCommentFuncAPI(w http.ResponseWriter, r *http.Request) {
	var status types.Status

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Path[len("/api/comment/"):]
	commentID, err := strconv.Atoi(id)
	if err != nil {
		status = types.Status{http.StatusBadRequest, "Invalid comment"}
		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}
		return
	}
	//username := sessions.GetCurrentUserName(r)
	username := "suraj"

	err = db.DeleteCommentByID(username, commentID)
	htStatus := http.StatusOK
	if err != nil {
		message = "comment not deleted"
		htStatus = http.StatusInternalServerError
	} else {
		message = "comment deleted"
	}
	status = types.Status{htStatus, message}
	w.WriteHeader(htStatus)
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		panic(err)
	}
	return
}

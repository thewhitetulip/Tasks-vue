package views

/*
Holds the delete related view handlers
*/

import (
	"log"
	"net/http"
	"strconv"

	"github.com/thewhitetulip/Tasks-vue/db"
	"github.com/thewhitetulip/Tasks-vue/sessions"
	"github.com/thewhitetulip/Tasks-vue/utils"
)

//TrashTaskFunc is used to populate the trash tasks
func TrashTaskFunc(w http.ResponseWriter, r *http.Request) {
	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	redirectURL := utils.GetRedirectUrl(r.Referer())

	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.URL.Path[len("/trash/"):])
	if err != nil {
		log.Println("TrashTaskFunc", err)
		message = "Incorrect command"
		http.Redirect(w, r, redirectURL, http.StatusFound)
	} else {
		username := sessions.GetCurrentUserName(r)
		err = db.TrashTask(username, id)
		if err != nil {
			message = "Error trashing task"
		} else {
			message = "Task trashed"
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)

	}
}

//RestoreTaskFunc is used to restore task from trash, handles "/restore/" URL
func RestoreTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.URL.Path[len("/restore/"):])
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/deleted", http.StatusBadRequest)
	} else {
		username := sessions.GetCurrentUserName(r)
		err = db.RestoreTask(username, id)
		if err != nil {
			message = "Restore failed"
		} else {
			message = "Task restored"
		}
		http.Redirect(w, r, "/deleted/", http.StatusFound)
	}

}

//DeleteTaskFunc is used to delete a task, trash = move to recycle bin, delete = permanent delete
func DeleteTaskFunc(w http.ResponseWriter, r *http.Request) {
	username := sessions.GetCurrentUserName(r)
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	id := r.URL.Path[len("/delete/"):]
	if id == "all" {
		err := db.DeleteAll(username)
		if err != nil {
			message = "Error deleting tasks"
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		id, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
		} else {
			err = db.DeleteTask(username, id)
			if err != nil {
				message = "Error deleting task"
			} else {
				message = "Task deleted"
			}
			http.Redirect(w, r, "/deleted", http.StatusFound)
		}
	}

}

//RestoreFromCompleteFunc restores the task from complete to pending
func RestoreFromCompleteFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.URL.Path[len("/incomplete/"):])
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/completed", http.StatusBadRequest)
	} else {
		username := sessions.GetCurrentUserName(r)
		err = db.RestoreTaskFromComplete(username, id)
		if err != nil {
			message = "Restore failed"
		} else {
			message = "Task restored"
		}
		http.Redirect(w, r, "/completed", http.StatusFound)
	}

}

//DeleteCategoryFunc will delete any category
func DeleteCategoryFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	categoryName := r.URL.Path[len("/del-category/"):]
	username := sessions.GetCurrentUserName(r)
	err := db.DeleteCategoryByName(username, categoryName)
	if err != nil {
		message = "error deleting category"
	} else {
		message = "Category " + categoryName + " deleted"
	}

	http.Redirect(w, r, "/", http.StatusFound)

}

//DeleteCommentFunc will delete any category
func DeleteCommentFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	id := r.URL.Path[len("/del-comment/"):]
	commentID, err := strconv.Atoi(id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	username := sessions.GetCurrentUserName(r)

	err = db.DeleteCommentByID(username, commentID)

	if err != nil {
		message = "comment not deleted"
	} else {
		message = "comment deleted"
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

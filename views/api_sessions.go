package views

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/thewhitetulip/Tasks-vue/db"
	"github.com/thewhitetulip/Tasks-vue/sessions"
	"github.com/thewhitetulip/Tasks-vue/types"
)

//RequiresLogin is a middleware which will be used for each httpHandler to check if there is any active session
func RequiresLogin(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessions.IsLoggedIn(r) {
			http.Redirect(w, r, "/login/", 302)
			return
		}
		handler(w, r)
	}
}

//LogoutFuncAPI Implements the logout functionality. WIll delete the session information from the cookie store
func LogoutFuncAPI(w http.ResponseWriter, r *http.Request) {
	var status types.Status
	var message = "Login successful"
	htStatus := http.StatusOK

	session, err := sessions.Store.Get(r, "session")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		htStatus = http.StatusInternalServerError
		message = "Could not logout"
	} else {
		if session.Values["loggedin"] != "false" {
			session.Values["loggedin"] = "false"
			session.Save(r, w)
		}
	}
	w.WriteHeader(htStatus)
	err = json.NewEncoder(w).Encode(status)

	status = types.Status{htStatus, message}
	if err != nil {
		panic(err)
	}
}

//LoginFuncAPI implements the login functionality, will add a cookie to the cookie store for managing authentication
func LoginFuncAPI(w http.ResponseWriter, r *http.Request) {
	var status types.Status
	var message string
	var htStatus = http.StatusOK

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	session, err := sessions.Store.Get(r, "session")

	if r.Method == "GET" {
		value := (session.Values["loggedin"] == "true")
		w.WriteHeader(htStatus)
		type loggedIn struct {
			Loggedin bool `json:"loggedin"`
		}
		logged := loggedIn{Loggedin: value}
		json.NewEncoder(w).Encode(logged)
		return
	}

	if err != nil {
		log.Println("error identifying session")
		htStatus = http.StatusInternalServerError
		message = "Could not login"
	} else {

		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if (username != "" && password != "") && db.ValidUser(username, password) {
			session.Values["loggedin"] = "true"
			session.Values["username"] = username
			session.Save(r, w)
			log.Print("user ", username, " is authenticated")
			message = "Logged in successfully"
		} else {
			htStatus = http.StatusBadRequest
			message = "Invalid user name or password;"
		}
	}

	w.WriteHeader(htStatus)
	status = types.Status{htStatus, message}
	err = json.NewEncoder(w).Encode(status)

	if err != nil {
		panic(err)
	}
}

//SignUpFuncAPI will enable new users to sign up to our service
func SignUpFuncAPI(w http.ResponseWriter, r *http.Request) {
	var status types.Status
	var message = "Sign up success"
	var statusCode = http.StatusOK
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	email := r.Form.Get("email")

	if username != "" && password != "" && email != "" {
		log.Println(username, password, email)

		err := db.CreateUser(username, password, email)
		if err != nil {
			statusCode = http.StatusInternalServerError
			message = "Something went wront"
		}
	} else {
		statusCode = http.StatusBadRequest
		message = "Invalid input"
	}

	w.WriteHeader(statusCode)
	status = types.Status{statusCode, message}
	err = json.NewEncoder(w).Encode(status)

	if err != nil {
		panic(err)
	}
}

package controllers

import (
	"html/template"
	"net/http"

	"forum-image-upload/db"
	"forum-image-upload/model"
)

// PageNotFound controller
func PageNotFound(w http.ResponseWriter, r *http.Request) {
	var response model.Response
	var user model.User
	username, authorized := IsAuthorized(r)
	if authorized {
		var err error
		user, err = db.GetUserByEmail(username)
		if err != nil {
			user = db.GetUserByName(username)
		}
	}
	response.User = user
	response.LoggedIn = authorized

	t := template.Must(template.New("notfound").ParseFiles("static/404.html", "static/header.html", "static/footer.html"))
	w.WriteHeader(404)
	t.Execute(w, response)
}

// InternalServerError internal server error
func InternalServerError(w http.ResponseWriter, r *http.Request) {
	var response model.Response
	var user model.User
	username, authorized := IsAuthorized(r)
	if authorized {
		var err error
		user, err = db.GetUserByEmail(username)
		if err != nil {
			user = db.GetUserByName(username)
		}
	}
	response.User = user
	response.LoggedIn = authorized

	t := template.Must(template.New("serverError").ParseFiles("static/500.html", "static/header.html", "static/footer.html"))
	w.WriteHeader(500)
	t.Execute(w, response)
}

// BadRequest controller
func BadRequest(w http.ResponseWriter, r *http.Request, err string) {
	var response model.Response
	response.Errors = append(response.Errors, err)
	var user model.User
	username, authorized := IsAuthorized(r)
	if authorized {
		var err error
		user, err = db.GetUserByEmail(username)
		if err != nil {
			user = db.GetUserByName(username)
		}
	}
	response.User = user
	response.LoggedIn = authorized

	t := template.Must(template.New("badrequest").ParseFiles("static/400.html", "static/header.html", "static/footer.html"))
	w.WriteHeader(400)
	t.Execute(w, response)
}

// UnauthorizedAccess controller
func UnauthorizedAccess(w http.ResponseWriter, r *http.Request) {
	var response model.Response
	var user model.User
	username, authorized := IsAuthorized(r)
	if authorized {
		var err error
		user, err = db.GetUserByEmail(username)
		if err != nil {
			user = db.GetUserByName(username)
		}
	}
	response.User = user
	response.LoggedIn = authorized

	t := template.Must(template.New("badrequest").ParseFiles("static/401.html", "static/header.html", "static/footer.html"))
	w.WriteHeader(http.StatusUnauthorized)
	t.Execute(w, response)
}

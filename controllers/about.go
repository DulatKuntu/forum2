package controllers

import (
	"html/template"
	"net/http"

	"forum-image-upload/db"
	"forum-image-upload/model"
)

//About just about page
func About(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/about" {
		if r.Method == "GET" {
			var response model.Response
			uname, loggedIn := IsAuthorized(r)
			if loggedIn {
				user, err := db.GetUserByEmail(uname)
				if err != nil {
					user = db.GetUserByName(uname)
				}
				response.User = user
			}
			response.LoggedIn = loggedIn
			t := template.Must(template.New("about").ParseFiles("static/about.html", "static/header.html", "static/footer.html"))
			t.Execute(w, response)
		} else {
			BadRequest(w, r, r.Method+" is not allowed")
		}
	} else {
		PageNotFound(w, r)
	}
}

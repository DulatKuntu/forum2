package controllers

import (
	"html/template"
	"net/http"

	"forum-image-upload/db"
)

//ProfilePage returns data about user
func ProfilePage(w http.ResponseWriter, r *http.Request) {
	username, auth := IsAuthorized(r)

	if !auth || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.URL.Path == "/profile" {
		if r.Method == "GET" {
			response := GetResponse(r)
			response.Posts = db.GetPostsByUserID(response.User.UserID)
			t := template.Must(template.New("profile").ParseFiles("static/profile.html", "static/header.html", "static/footer.html"))
			t.Execute(w, response)
		} else {
			BadRequest(w, r, r.Method+" is not allowed")
		}
	} else {
		PageNotFound(w, r)
	}
}

package controllers

import (
	"html/template"
	"net/http"
	"time"

	"forum-image-upload/db"
	"forum-image-upload/model"
	uuid "github.com/satori/go.uuid"
)

var cookies map[string]*http.Cookie

// Login controller
func Login(w http.ResponseWriter, r *http.Request) {
	if cookies == nil {
		cookies = map[string]*http.Cookie{}
	}
	_, auth := IsAuthorized(r)
	if auth {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
		if r.Method == "GET" {
			var response model.Response
			t := template.Must(template.New("login").ParseFiles("static/login.html", "static/header.html", "static/footer.html"))
			t.Execute(w, response)
		} else if r.Method == "POST" {
			//check credentials and login the user
			r.ParseForm()
			var response model.Response
			response.User.Username = r.FormValue("login")
			response.User.Password = r.FormValue("password")
			if db.CheckCredentials(response.User.Username, response.User.Password) {
				user, err := db.GetUserByEmail(response.User.Username)
				if err != nil {
					response.User = db.GetUserByName(response.User.Username)
				} else {
					response.User = user
				}
				u := uuid.NewV4()
				sessionToken := u.String()

				cookie := &http.Cookie{
					Name:    "session_token",
					Value:   sessionToken, // Some encoded value
					Path:    "/",          // Otherwise it defaults to the /login if you create this on /login (standard cookie behaviour)
					Expires: time.Now().Add(7200 * time.Second),
				}
				cookies[response.User.Username] = cookie
				http.SetCookie(w, cookie)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				// invalid credentials
				response.Errors = append(response.Errors, "invalid username or password")
				response.LoggedIn = false
				t := template.Must(template.New("login").ParseFiles("static/login.html", "static/header.html", "static/footer.html"))
				t.Execute(w, response)
			}
		} else {
			BadRequest(w, r, r.Method+" is not allowed")
		}
	

}

//IsAuthorized simple middlwear
func IsAuthorized(r *http.Request) (string, bool) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", false
		}
		return "", false
	}
	for username, cookie := range cookies {
		if cookie.Value == c.Value {

			if time.Until(cookie.Expires) <= 0 {

				return "", false
			}
			return username, true
		}
	}
	return "", false
}

package controllers

import (
	"html/template"
	"net/http"
	"forum-image-upload/db"
	"forum-image-upload/model"
)
// Signup controller
func Signup(w http.ResponseWriter, r *http.Request) {
	_, auth := IsAuthorized(r)
	if auth {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.URL.Path == "/signup" {
		if r.Method == "GET" {
			response := GetResponse(r)
			t := template.Must(template.New("signup").ParseFiles("static/signup.html", "static/header.html", "static/footer.html"))
			t.Execute(w, response)
		} else if r.Method == "POST" {
			//insert new user
			var response model.Response
			var user model.User
			r.ParseForm()
			user.Username = r.FormValue("username")
			user.Email = r.FormValue("email")
			user.Password = r.FormValue("password")
			retypedPassword := r.FormValue("repassword")

			//check password, username, email
			if user.IsValidPassword() != "" {
				response.Errors = append(response.Errors, user.IsValidPassword())
			}
			if retypedPassword != user.Password {
				response.Errors = append(response.Errors, "passwords are not same!!!")
			}

			if user.IsValidUsername() {
				usernameExists, err := db.UsernameExists(user.Username)
				if err != nil {
					InternalServerError(w, r)
					return
				}

				if usernameExists {
					response.Errors = append(response.Errors, "username is already in use")
				}
			} else {
				response.Errors = append(response.Errors, "username should be an alphanumeric string with length 3 to 20")
			}

			if user.IsValidEmail() {
				emailExists, err := db.EmailExists(user.Email)
				if err != nil {
					InternalServerError(w, r)
					return
				}

				if emailExists {
					response.Errors = append(response.Errors, "email is already in use")
				}
			} else {
				response.Errors = append(response.Errors, "invalid email address")

			}

			// if there are no errors then signup and make status registered
			if len(response.Errors) == 0 {
				registered := db.SignupUser(user)
				if !registered {
					InternalServerError(w, r)
					return
				}
				response.Status = "REGISTERED"
			} else {
				response.Status = "ERROR"
			}

			response.User = user
			t := template.Must(template.New("signup").ParseFiles("static/signup.html", "static/header.html", "static/footer.html"))
			t.Execute(w, response)

		} else {
			BadRequest(w, r, r.Method+" is not allowed")
		}
	} else {
		PageNotFound(w, r)
	}

}

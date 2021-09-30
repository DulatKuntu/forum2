package controllers

import (
	"html/template"
	"net/http"
	"strings"

	"forum-image-upload/db"
	"forum-image-upload/model"
)

// MainPage controller
func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
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
			allPosts := db.GetAllPosts()

			r.ParseForm()
			title := r.FormValue("title")
			category := r.FormValue("category")
			for i := 0; i < len(allPosts); i++ {
				if !strings.HasPrefix(strings.ToLower(allPosts[i].Title), strings.ToLower(title)) {
					allPosts = append(allPosts[:i], allPosts[i+1:]...)
					i--
				}
			}

			if loggedIn {
				likedbyme := r.FormValue("liked")
				myposts := r.FormValue("mypost")
				if likedbyme == "on" {
					for i := 0; i < len(allPosts); i++ {
						liked := false
						for _, likeid := range allPosts[i].Likes {
							if likeid == response.User.UserID {
								liked = true
								break
							}
						}
						if !liked {
							allPosts = append(allPosts[:i], allPosts[i+1:]...)
							i--
						}
					}
				}

				if myposts == "on" {
					for i := 0; i < len(allPosts); i++ {
						if allPosts[i].UserID != response.User.UserID {
							allPosts = append(allPosts[:i], allPosts[i+1:]...)
							i--
						}
					}
				}
			}

			if category != "" {
				for i := 0; i < len(allPosts); i++ {
					if strings.ToLower(allPosts[i].Category) != strings.ToLower(category) {
						allPosts = append(allPosts[:i], allPosts[i+1:]...)
						i--
					}
				}
			}

			response.Posts = allPosts

			t := template.Must(template.New("index").ParseFiles("static/index.html", "static/header.html", "static/footer.html"))
			t.Execute(w, response)
		} else {
			BadRequest(w, r, r.Method+" is not allowed")
		}
	} else {
		PageNotFound(w, r)
	}

}

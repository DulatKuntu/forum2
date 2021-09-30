package controllers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"forum-image-upload/db"
	"forum-image-upload/model"
)

//GetPost show post details
func GetPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/post" {
		if r.Method == "GET" {
			r.ParseForm()
			postid, _ := strconv.Atoi(r.FormValue("id"))
			post := db.GetPostByPostID(postid)
			post.Comments = db.GetCommentsByPostID(postid)
			post.Likes = db.GetLikesByPostID(postid)
			post.Dislikes = db.GetDislikesByPostID(postid)
			response := GetResponse(r)
			response.Post = post
			t := template.Must(template.New("post").ParseFiles("static/post.html", "static/header.html", "static/footer.html"))
			t.Execute(w, response)
		} else {
			BadRequest(w, r, r.Method+" is not allowed")
		}
	} else {
		PageNotFound(w, r)
	}
}

//CreatePost creates post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	username, auth := IsAuthorized(r)
	if !auth || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.URL.Path == "/create" {
		if r.Method == "GET" {
			// return template
			response := GetResponse(r)
			t := template.Must(template.New("create").ParseFiles("static/create.html", "static/header.html", "static/footer.html"))
			t.Execute(w, response)
		} else if r.Method == "POST" {
			// create  post
			var response model.Response
			r.ParseMultipartForm(0)
			title := r.FormValue("title")
			text := r.FormValue("text")
			category := r.FormValue("category")
			var post model.Post

			post.Title = title
			post.Text = text
			post.Category = category
			post.Date = time.Now().String()

			errors := post.Validate()
			if len(errors) != 0 {
				response.Post = post
				response.Errors = errors
				t := template.Must(template.New("create").ParseFiles("static/create.html", "static/header.html", "static/footer.html"))
				t.Execute(w, response)
				return
			}

			///uploading file
			err := r.ParseMultipartForm(20 << 20)

			if err != nil {
				fmt.Println("ERROR IS ", err.Error())
				InvalidPostHandler(w, r, "size of the image is too big! It should be max of 20mb!", post)
				return
			}
			file, fileHeader, fileExists := r.FormFile("file")

			if fileExists == nil {
				//there is some file

				if fileHeader.Size > 20*1024*1024 {
					InvalidPostHandler(w, r, "image size is too big, max 20mb!", post)
					return
				}

				defer file.Close()
				buff := make([]byte, 512)
				temp := file
				_, err := temp.Read(buff)

				if err != nil {
					InvalidPostHandler(w, r, "invalid image", post)
					return
				}

				ext := http.DetectContentType(buff)

				if ext != "image/jpeg" && ext != "image/svg" && ext != "image/png" && ext != "image/gif" && ext != "image/jpg" {
					InvalidPostHandler(w, r, ext+" is invalid image format!", post)
					return
				}

				_, err = file.Seek(0, io.SeekStart) // makes sure the file is read from the start
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				err = os.MkdirAll("./images", os.ModePerm) // makes uploads directory if it doesnt exist
				if err != nil {
					InternalServerError(w, r)
					return
				}

				strPath := time.Now().UnixNano()

				dst, err := os.Create(fmt.Sprintf("./images/%d.%s", strPath, strings.Split(ext, "/")[1]))

				if err != nil {
					InternalServerError(w, r)
					return
				}

				defer dst.Close()

				_, err = io.Copy(dst, file) // copies the file content into the new one in the uploads folder
				if err != nil {
					InternalServerError(w, r)
					return
				}

				post.Image = fmt.Sprintf("./images/%d.%s", strPath, strings.Split(ext, "/")[1])
				defer file.Close()

			}

			user, err := db.GetUserByEmail(username)
			if err != nil {
				user = db.GetUserByName(username)
			}
			post.UserID = user.UserID
			err = db.InsertPost(post)
			if err != nil {
				BadRequest(w, r, r.Method+" is not allowed") //TODO review
				return
			}
			response = GetResponse(r)
			response.Status = "SUCCESS"
			response.Post = post
			w.WriteHeader(200)
			t := template.Must(template.New("create").ParseFiles("static/create.html", "static/header.html", "static/footer.html"))
			t.Execute(w, response)
		} else {
			BadRequest(w, r, r.Method+" is not allowed")
		}
	} else {
		PageNotFound(w, r)
	}

}

//InvalidPostHandler handler
func InvalidPostHandler(w http.ResponseWriter, r *http.Request, err string, post model.Post) {
	var response model.Response
	response.Errors = append(response.Errors, err)
	response.Status = "FAIL"
	response.Post = post
	w.WriteHeader(400)
	t := template.Must(template.New("create").ParseFiles("static/create.html", "static/header.html", "static/footer.html"))
	t.Execute(w, response)
}

//DeletePost deletes post
func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/delete" {

	} else {
		PageNotFound(w, r)
	}
}

//LikePost like post
func LikePost(w http.ResponseWriter, r *http.Request) {
	username, auth := IsAuthorized(r)
	if !auth || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.URL.Path == "/like" {
		if r.Method == "GET" {
			r.ParseForm()
			postid, _ := strconv.Atoi(r.FormValue("id"))
			user, err := db.GetUserByEmail(username)
			if err != nil {
				user = db.GetUserByName(username)
			}

			db.LikePost(user.UserID, postid, true)

			redir := r.FormValue("redir")
			http.Redirect(w, r, redir, http.StatusSeeOther)
		} else {
			BadRequest(w, r, r.Method+" is not allowed")
		}
	} else {
		PageNotFound(w, r)
	}
}

//DislikePost dislike post
func DislikePost(w http.ResponseWriter, r *http.Request) {
	username, auth := IsAuthorized(r)
	if !auth || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return

	}
	if r.URL.Path == "/dislike" {
		if r.Method == "GET" {
			r.ParseForm()
			postid, _ := strconv.Atoi(r.FormValue("id"))
			user, err := db.GetUserByEmail(username)
			if err != nil {
				user = db.GetUserByName(username)
			}
			db.LikePost(user.UserID, postid, false)

			redir := r.FormValue("redir")
			http.Redirect(w, r, redir, http.StatusSeeOther)
		} else {
			BadRequest(w, r, r.Method+" is not allowed")
		}
	} else {
		PageNotFound(w, r)
	}
}

//LikeCom like comment
func LikeCom(w http.ResponseWriter, r *http.Request) {
	username, auth := IsAuthorized(r)
	if !auth || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		r.ParseForm()
		commentid, _ := strconv.Atoi(r.FormValue("id"))
		user, err := db.GetUserByEmail(username)
		if err != nil {
			user = db.GetUserByName(username)
		}

		db.LikeCom(user.UserID, commentid, true)
		postid := db.GetPostID(commentid)
		redir := "/post?id=" + strconv.Itoa(postid) + "#comment" + strconv.Itoa(commentid)
		http.Redirect(w, r, redir, http.StatusSeeOther)
	} else {
		BadRequest(w, r, r.Method+" is not allowed")
	}

}

//DislikeCom dislike comment
func DislikeCom(w http.ResponseWriter, r *http.Request) {
	username, auth := IsAuthorized(r)
	if !auth || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return

	}
	if r.Method == "GET" {
		r.ParseForm()
		commentid, _ := strconv.Atoi(r.FormValue("id"))
		user, err := db.GetUserByEmail(username)
		if err != nil {
			user = db.GetUserByName(username)
		}
		db.LikeCom(user.UserID, commentid, false)

		postid := db.GetPostID(commentid)
		redir := "/post?id=" + strconv.Itoa(postid) + "#comment" + strconv.Itoa(commentid)
		fmt.Print(redir)
		http.Redirect(w, r, redir, http.StatusSeeOther)
	} else {
		BadRequest(w, r, r.Method+" is not allowed")
	}

}

//CommentPost creates comment to the post
func CommentPost(w http.ResponseWriter, r *http.Request) {
	username, auth := IsAuthorized(r)
	if !auth || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return

	}
	if r.URL.Path == "/comment" {
		if r.Method == "POST" {
			r.ParseForm()
			postid, err := strconv.Atoi(r.FormValue("id"))
			if err != nil {
				PageNotFound(w, r)
				return
			}
			text := r.FormValue("text")
			user, err := db.GetUserByEmail(username)
			if err != nil {
				fmt.Println(err)
				user = db.GetUserByName(username)
			}
			err = db.CommentPost(user.UserID, postid, text)
			if err != nil {
				InternalServerError(w, r)
				return
			}
			http.Redirect(w, r, "/post?id="+r.FormValue("id"), http.StatusSeeOther)

		} else {
			BadRequest(w, r, r.Method+" is not allowed")
		}
	} else {
		PageNotFound(w, r)
	}
}

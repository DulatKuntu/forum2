package main

import (
	"fmt"
	"net/http"

	controllers "forum-image-upload/controllers"
	db "forum-image-upload/db"
)

func main() {
	port := "5555"
	_, err := db.ConnectDb()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Couldn't start the server database creation error")
		return
	}
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", controllers.MainPage)
	http.HandleFunc("/about", controllers.About)
	http.HandleFunc("/signup", controllers.Signup)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/create", controllers.CreatePost)
	http.HandleFunc("/cleanup", controllers.Cleanup)
	http.HandleFunc("/like", controllers.LikePost)
	http.HandleFunc("/dislike", controllers.DislikePost)
	http.HandleFunc("/likecom", controllers.LikeCom)
	http.HandleFunc("/dislikecom", controllers.DislikeCom)
	http.HandleFunc("/comment", controllers.CommentPost)
	http.HandleFunc("/post", controllers.GetPost)
	http.HandleFunc("/profile", controllers.ProfilePage)

	fmt.Println("running on port ", port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println("occured error while running the server!")
		return
	}

}

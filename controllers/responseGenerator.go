package controllers

import (
	"net/http"

	"forum-image-upload/db"
	"forum-image-upload/model"
)

//GetResponse adds user details to the response
func GetResponse(r *http.Request) model.Response {
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
	return response
}

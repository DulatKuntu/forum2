package controllers

import (
	"net/http"

	"forum-image-upload/db"
)

//Cleanup used to drop all tables
func Cleanup(w http.ResponseWriter, r *http.Request) {
	db.DropTables()
}

package model

import (
	"strings"
)

// Post just post
type Post struct {
	PostID   int
	UserID   int
	Title    string
	Text     string
	Category string
	Image    string
	Date     string
	Comments []Comment
	Likes    []int
	Dislikes []int
}

//Validate used to check if post is valid
func (post Post) Validate() []string {
	var errors []string
	strTitle := strings.ReplaceAll(post.Title, " ", "")
	if strTitle == "" {
		errors = append(errors, "title cannot be empty!")
	}

	if len(strTitle) > 100 {
		errors = append(errors, "length of the title cannot be longer than 100 characters!")
	}

	strText := strings.ReplaceAll(post.Text, " ", "")

	if strText == "" {
		errors = append(errors, "text of the post cannot be empty")
	}

	cat := strings.ReplaceAll(post.Category, " ", "")

	if cat == "" {

		errors = append(errors, "select valid category!")

	}

	return errors
}

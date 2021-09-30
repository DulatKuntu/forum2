package model

//Response Custom response for any request
type Response struct {
	Status   string
	LoggedIn bool
	User     User
	Post     Post
	Comment  Comment
	Posts    []Post
	Errors   []string
}

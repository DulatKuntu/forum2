package model

// Comment just model
type Comment struct {
	CommentID int
	PostID    int
	User      User
	Comment   string
	Likes     int
	Dislikes  int
}

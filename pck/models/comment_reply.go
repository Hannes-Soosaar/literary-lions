package models

type CommentReply struct {
	ID        int
	Body      string
	UserId    int
	CommentId int
	PostId    int
	CreatedAt string
	Active    int
}

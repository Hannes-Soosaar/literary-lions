package models

type Comment struct {
	ID         int
	Body       string
	UserId     int
	Likes      int
	Dislikes   int
	PostID     int
	CreatedAt  string
	ModifiedAt string
	Active     int
}

package models

type Post struct {
	ID         int
	Title      string
	Body       string
	Likes      int
	Dislikes   int
	UserId     int
	CategoryID int
	CreatedAt  string
	ModifiedAt string
	Active     int
}

type Posts struct {
	AllPosts []Post
}

type FilteredPosts struct {
	FilteredPosts []Post
}

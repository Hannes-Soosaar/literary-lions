package models

// import "gitea.kood.tech/hannessoosaar/literary-lions/pck/models"

type TemplateData struct {
	Username                   string
	RegistrationSuccessMessage string
	ErrorMessage               string
	Title                      string
	Uuid                       string
	CategoryPage               bool
	IsLoggedIn                 bool
	ProfilePage                bool
	MainPage                   bool
	UserPostsPage              bool
	LikedPostsPage             bool
	DislikedPostsPage          bool
	AllPosts                   Posts
	Categories                 []Category
	PostComments               []PostComment
	Comments                   []Comment
	CommentReplies             []CommentReply
	User                       User
	StaticURL                  string
	CreatePostPage             bool
	ShowComments               bool
	PostCreatedMessage         string
	DisplayCatID               int
	QueryNoResult              string
	SearchQuery                string
	FilterType                 string
	EmptyMessage               string
	Message                    Message
}

func DefaultTemplateData() TemplateData {
	return TemplateData{
		Username:                   "",
		RegistrationSuccessMessage: "",
		ErrorMessage:               "",
		Title:                      "Lions",
		Uuid:                       "",
		CategoryPage:               false,
		IsLoggedIn:                 false,
		ProfilePage:                false,
		MainPage:                   false,
		UserPostsPage:              false,
		LikedPostsPage:             false,
		DislikedPostsPage:          false,
		CreatePostPage:             false,
		AllPosts:                   Posts{},
		Categories:                 []Category{},
		PostComments:               []PostComment{},
		Comments:                   []Comment{},
		CommentReplies:             []CommentReply{},
		User:                       User{},
		StaticURL:                  "http://localhost:8080/static",
		ShowComments:               false,
		PostCreatedMessage:         "",
		DisplayCatID:               0,
		QueryNoResult:              "",
		SearchQuery:                "",
		FilterType:                 "",
		EmptyMessage:               "",
		Message:                    *GetInstance(),
	}
}

package models

// import "gitea.kood.tech/hannessoosaar/literary-lions/pck/models"

type TemplateData struct {
	Username                   string
	RegistrationSuccessMessage string
	ErrorMessage               string
	Title                      string
	Uuid                       string
	IsLoggedIn                 bool
	ProfilePage                bool
	MainPage                   bool
	AllPosts                   Posts
	Categories                 []Category
	PostComments               []PostComment
	Comments                   []Comment
	User                       User
	FilteredPosts              FilteredPosts
	StaticURL                  string
	CreatePostPage             bool
	ShowComments               bool
}

// TODO migrate all user data fields to the user struct
func DefaultTemplateData() TemplateData {
	return TemplateData{
		Username:                   "",
		RegistrationSuccessMessage: "",
		ErrorMessage:               "",
		Title:                      "Lions",
		Uuid:                       "",
		IsLoggedIn:                 false,
		ProfilePage:                false,
		MainPage:                   false,
		CreatePostPage:             false,
		AllPosts:                   Posts{},
		Categories:                 []Category{},
		PostComments:               []PostComment{},
		Comments:                   []Comment{},
		User:                       User{},
		FilteredPosts:              FilteredPosts{},
		StaticURL:                  "http://localhost:8082/static",
		ShowComments:               false,
	}
}

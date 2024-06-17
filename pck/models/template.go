package models

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
}

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
		AllPosts:                   Posts{},
	}
}

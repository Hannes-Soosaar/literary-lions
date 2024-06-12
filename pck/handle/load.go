package handle

import (
	"fmt"
	"net/http"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/render"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func LandingPageHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		SuccessMessage             string
		RegistrationSuccessMessage string
		ErrorMessage               string
		Title                      string
		Uuid                       string
	}{
		SuccessMessage:             "",
		Title:                      "Lions",
		RegistrationSuccessMessage: "",
		ErrorMessage:               "",
		Uuid:                       "",
	}
	render.RenderLandingPage(w, "index.html", data)
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var errorMessage string
	var successMessage string
	var user models.User

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request method", http.StatusMethodNotAllowed)
		return
	} else {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := utils.HashString(r.FormValue("password"))
	err := utils.AddNewUser(username, email, password)
	if err != nil {
		errorMessage = err.Error()
		fmt.Println(errorMessage)
	} else {
		successMessage = fmt.Sprintf("%s was added with the email %s", username, email)
		user = utils.FindUserByUserName(username)
	}
	fmt.Println(successMessage)
	data := struct {
		SuccessMessage             string
		ErrorMessage               string
		RegistrationSuccessMessage string
		Title                      string
		Uuid                       string
	}{
		SuccessMessage:             successMessage,
		ErrorMessage:               errorMessage,
		RegistrationSuccessMessage: "",
		Title:                      "Login page",
		Uuid:                       user.UUID,
	}
	render.RenderLandingPage(w, "index.html", data)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request method", http.StatusMethodNotAllowed)
		return
	} else {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		}
	}

	var errorMessage string
	var successMessage string
	email := r.FormValue("email")
	password := r.FormValue("password")
	uuid, isActiveUser, err := utils.ValidateUser(email, password)

	if err != nil {
		errorMessage = err.Error()
	}

	if isActiveUser {
		successMessage = "Logged on as " + email
	} else {
		successMessage = "Not a active user "
		uuid = ""
	}

	data := struct {
		SuccessMessage             string
		ErrorMessage               string
		RegistrationSuccessMessage string
		Title                      string
		Uuid                       string
	}{
		SuccessMessage:             successMessage,
		ErrorMessage:               errorMessage,
		RegistrationSuccessMessage: "",
		Title:                      "Login page",
		Uuid:                       uuid,
	}

	render.RenderLandingPage(w, "index.html", data)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	data := struct {
		SuccessMessage             string
		RegistrationSuccessMessage string
		ErrorMessage               string
		Title                      string
		Uuid                       string
	}{
		SuccessMessage:             "You are now logged out!",
		Title:                      "Lions",
		RegistrationSuccessMessage: "",
		ErrorMessage:               "",
		Uuid:                       "",
	}
	render.RenderLandingPage(w, "index.html", data)
}

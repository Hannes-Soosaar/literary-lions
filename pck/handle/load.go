package handle

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/render"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func LandingPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Loading Page")
	utils.FindUserByUserName("bob")
	data := struct {
		RegistrationSuccessMessage string
		ErrorMessage               string
		Title                      string
	}{
		Title:                      "Lions",
		RegistrationSuccessMessage: "",
		ErrorMessage:               "",
	}

	render.RenderLandingPage(w, "index.html", data)
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}
	successMessage := r.URL.Query().Get("success")

	if successMessage != "" {
		data := struct {
			RegistrationSuccessMessage string
			ErrorMessage               string
			Title                      string
		}{
			Title:                      "Lions",
			RegistrationSuccessMessage: successMessage,
			ErrorMessage:               "",
		}
		render.RenderLandingPage(w, "index.html", data)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	switch { //TODO make database lookup functions to check if a username or email has already been used for another account and make the corresponding switch cases!
	case username == "" || email == "" || password == "":
		errorMessage := "None of the fields can be empty!"
		data := struct {
			RegistrationSuccessMessage string
			ErrorMessage               string
			Title                      string
		}{
			Title:                      "Lions",
			RegistrationSuccessMessage: "",
			ErrorMessage:               errorMessage,
		}
		render.RenderLandingPage(w, "index.html", data)
		return
	case !strings.Contains(email, "@"):
		errorMessage := "Invalid email format! Please enter a valid email address."
		data := struct {
			RegistrationSuccessMessage string
			ErrorMessage               string
			Title                      string
		}{
			Title:                      "Lions",
			RegistrationSuccessMessage: "",
			ErrorMessage:               errorMessage,
		}
		render.RenderLandingPage(w, "index.html", data)
	default:
		utils.AddUserTest(username, email, password)

		successMessage := "Registration successful!"
		encodedMessage := url.QueryEscape(successMessage)

		redirectURL := "/register?success=" + encodedMessage
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
	}

	data := struct {
		ErrorMessage               string
		RegistrationSuccessMessage string
		Title                      string
	}{
		ErrorMessage:               "",
		RegistrationSuccessMessage: "",
		Title:                      "Login page",
	}
	render.RenderLandingPage(w, "index.html", data)
}

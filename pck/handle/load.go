package handle

import (
	"fmt"
	"net/http"
	"net/url"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/render"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func LandingPageHandler(w http.ResponseWriter, r *http.Request) {
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

	if successMessage != "" { // why are we loading the  page if we get a happy path ?
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

	fmt.Println("We made it here")

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
		//! move to front end 
	// case !strings.Contains(email, "@"):
	// 	errorMessage := "Invalid email format! Please enter a valid email address."
	// 	data := struct {
	// 		RegistrationSuccessMessage string
	// 		ErrorMessage               string
	// 		Title                      string
	// 	}{
	// 		Title:                      "Lions",
	// 		RegistrationSuccessMessage: "",
	// 		ErrorMessage:               errorMessage,
	// 	}

		render.RenderLandingPage(w, "index.html", data)

	default:
		utils.AddNewUser(username, email, password)

		successMessage := "Registration successful!"
		encodedMessage := url.QueryEscape(successMessage)

		redirectURL := "/register?success=" + encodedMessage
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost{
		http.Error(w, "Invalid Request method", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password :=r.FormValue("password")
	var errorMessage string
	//? simplify to only send back an error and uuid if uuid is not present it means no user exists
	uuid,isActiveUser,err := utils.ValidateUser(email,password)
	
	if err != nil {
		errorMessage = err.Error()
	}

	// err := r.ParseForm()
	// if err != nil {
	// 	http.Error(w, "Failed to parse form data", http.StatusBadRequest)
	// }
	data := struct {
		ErrorMessage               string
		RegistrationSuccessMessage string
		Title                      string
		Uuid					   string
	}{
		ErrorMessage:               errorMessage,
		RegistrationSuccessMessage: "",
		Title:                      "Login page",
		Uuid: 						uuid,
	}
		if isActiveUser{
			render.RenderLandingPage(w, "index.html", data)
		}else {
			fmt.Println(" NOT A USER REGISTER")
		}
}

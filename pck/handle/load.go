package handle

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/render"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

type contextKey string

const userContextKey = contextKey("username")

var sessionStore = map[string]string{}
var loggedInUsername string

func LandingPageHandler(w http.ResponseWriter, r *http.Request) {

	sessionToken, err := r.Cookie("session_token")
	isLoggedIn := err == nil && isValidSession(sessionToken.Value)

	data := struct {
		Username                   string
		RegistrationSuccessMessage string
		ErrorMessage               string
		Title                      string
		Uuid                       string
		IsLoggedIn                 bool
		ProfilePage                bool
		MainPage                   bool
	}{
		Username:                   "",
		Title:                      "Lions",
		RegistrationSuccessMessage: "",
		ErrorMessage:               "",
		Uuid:                       "",
		IsLoggedIn:                 isLoggedIn,
		ProfilePage:                false,
		MainPage:                   true,
	}
	if isLoggedIn {
		if data.Username == "" {
			data.Username = loggedInUsername
		}
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
		IsLoggedIn                 bool
		ProfilePage                bool
		MainPage                   bool
	}{
		SuccessMessage:             successMessage,
		ErrorMessage:               errorMessage,
		RegistrationSuccessMessage: "",
		Title:                      "Login page",
		Uuid:                       user.UUID,
		IsLoggedIn:                 false,
		ProfilePage:                false,
		MainPage:                   false,
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
	username := r.FormValue("username")
	password := r.FormValue("password")
	uuid, isActiveUser, err := utils.ValidateUser(username, password)

	if err != nil {
		errorMessage = err.Error()
	}

	if isActiveUser {

		sessionToken, err := utils.GenerateUUID()
		if err != nil {
			errorMessage = errorMessage + " Failed to generate UUID"
		}

		sessionStore[sessionToken] = username
		fmt.Println("Session token/UUID and username:", sessionStore)

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(30 * 24 * time.Hour),
			Path:    "/",
		})
	} else {
		uuid = ""
	}

	data := struct {
		Username                   string
		ErrorMessage               string
		RegistrationSuccessMessage string
		Title                      string
		Uuid                       string
		IsLoggedIn                 bool
		ProfilePage                bool
		MainPage                   bool
	}{
		Username:                   username,
		ErrorMessage:               errorMessage,
		RegistrationSuccessMessage: "",
		Title:                      "Login page",
		Uuid:                       uuid,
		IsLoggedIn:                 false,
		ProfilePage:                false,
		MainPage:                   false,
	}

	if isActiveUser {
		data.IsLoggedIn = true
		loggedInUsername = username
	}
	render.RenderLandingPage(w, "index.html", data)

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})

	data := struct {
		Username                   string
		RegistrationSuccessMessage string
		ErrorMessage               string
		Title                      string
		Uuid                       string
		IsLoggedIn                 bool
		ProfilePage                bool
		MainPage                   bool
	}{
		Username:                   "",
		Title:                      "Lions",
		RegistrationSuccessMessage: "",
		ErrorMessage:               "",
		Uuid:                       "",
		IsLoggedIn:                 false,
		ProfilePage:                false,
		MainPage:                   false,
	}
	for key := range sessionStore {
		delete(sessionStore, key)
		break
	}
	render.RenderLandingPage(w, "index.html", data)
}

func isValidSession(sessionToken string) bool {
	_, isValidSession := sessionStore[sessionToken]
	return isValidSession
}

func AuthSessionToken(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("TEST")
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		sessionToken := cookie.Value
		username, exists := sessionStore[sessionToken]
		if !exists {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
		fmt.Println("CTX", ctx)
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Check session token
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	sessionToken := cookie.Value
	username, exists := sessionStore[sessionToken]
	if !exists {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Retrieve username from session token
	ctx := context.WithValue(r.Context(), userContextKey, username)

	// Proceed with handling the request
	ctxUsername, ok := ctx.Value(userContextKey).(string)
	if !ok {
		http.Error(w, "Unable to retrieve username from context", http.StatusInternalServerError)
		return
	}

	isLoggedIn := true
	data := struct {
		Username                   string
		RegistrationSuccessMessage string
		ErrorMessage               string
		Title                      string
		Uuid                       string
		IsLoggedIn                 bool
		ProfilePage                bool
		MainPage                   bool
	}{
		Username:                   ctxUsername,
		Title:                      "Your Profile",
		RegistrationSuccessMessage: "",
		ErrorMessage:               "",
		Uuid:                       "",
		IsLoggedIn:                 isLoggedIn,
		ProfilePage:                true,
		MainPage:                   false,
	}

	if isLoggedIn {
		render.RenderProfile(w, "index.html", data)
	} else {
		data.ErrorMessage = "You need to be logged in to access your profile!"
		render.RenderLandingPage(w, "index.html", data)
	}
}

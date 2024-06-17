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
	allPosts := utils.RetrieveAllPosts()

	data := struct {
		Username                   string
		RegistrationSuccessMessage string
		ErrorMessage               string
		Title                      string
		Uuid                       string
		IsLoggedIn                 bool
		ProfilePage                bool
		MainPage                   bool
		AllPosts                   models.Posts
	}{
		Username:                   "",
		Title:                      "Lions",
		RegistrationSuccessMessage: "",
		ErrorMessage:               "",
		Uuid:                       "",
		IsLoggedIn:                 isLoggedIn,
		ProfilePage:                false,
		MainPage:                   true,
		AllPosts:                   allPosts,
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
	}
	fmt.Println(successMessage)
	allPosts := utils.RetrieveAllPosts()
	data := struct {
		Username                   string
		ErrorMessage               string
		RegistrationSuccessMessage string
		Title                      string
		Uuid                       string
		IsLoggedIn                 bool
		ProfilePage                bool
		MainPage                   bool
		AllPosts                   models.Posts
	}{
		Username:                   "",
		ErrorMessage:               errorMessage,
		RegistrationSuccessMessage: "Account created successfully! You can now log in.",
		Title:                      "Login page",
		Uuid:                       "",
		IsLoggedIn:                 false,
		ProfilePage:                false,
		MainPage:                   false,
		AllPosts:                   allPosts,
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
		for key, storedUsername := range sessionStore {
			if storedUsername == username {
				delete(sessionStore, key)
				break
			}
		}

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
		errorMessage = "Not a valid user!"
	}

	allPosts := utils.RetrieveAllPosts()

	data := struct {
		Username                   string
		ErrorMessage               string
		RegistrationSuccessMessage string
		Title                      string
		Uuid                       string
		IsLoggedIn                 bool
		ProfilePage                bool
		MainPage                   bool
		AllPosts                   models.Posts
	}{
		Username:                   username,
		ErrorMessage:               errorMessage,
		RegistrationSuccessMessage: "",
		Title:                      "Login page",
		Uuid:                       uuid,
		IsLoggedIn:                 false,
		ProfilePage:                false,
		MainPage:                   false,
		AllPosts:                   allPosts,
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

	allPosts := utils.RetrieveAllPosts()

	data := struct {
		Username                   string
		RegistrationSuccessMessage string
		ErrorMessage               string
		Title                      string
		Uuid                       string
		IsLoggedIn                 bool
		ProfilePage                bool
		MainPage                   bool
		AllPosts                   models.Posts
	}{
		Username:                   "",
		Title:                      "Lions",
		RegistrationSuccessMessage: "",
		ErrorMessage:               "",
		Uuid:                       "",
		IsLoggedIn:                 false,
		ProfilePage:                false,
		MainPage:                   false,
		AllPosts:                   allPosts,
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

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle like logic here (update database, etc.)
	// Example: Updating database
	// Replace this with your actual database update logic
	// _, err := db.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", postID)

	// Redirect back to the previous page (referer)
	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle like logic here (update database, etc.)
	// Example: Updating database
	// Replace this with your actual database update logic
	// _, err := db.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?", postID)

	// Redirect back to the previous page (referer)
	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

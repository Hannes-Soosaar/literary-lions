package handle

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/render"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

type contextKey string

const userContextKey = contextKey("username")

var sessionStore = map[string]string{}

func LandingPageHandler(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session_token")
	isLoggedIn := err == nil && isValidSession(sessionToken.Value)
	allPosts := utils.GetAllPosts()
	categories := utils.GetActiveCategories()
	comments := utils.GetActiveComments()
	data := models.DefaultTemplateData()
	replies := utils.GetAllReplies()
	data.CommentReplies = replies
	data.IsLoggedIn = isLoggedIn
	data.MainPage = true
	data.ProfilePage = false
	data.AllPosts = allPosts
	data.Comments = comments
	data.Categories = categories
	if isLoggedIn {
		if data.Username == "" {
			data.Username = GetUsernameFromCookie(r)
			if data.Username == "" {
				render.RenderLandingPage(w, "index.html", data)
			}
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
		fmt.Println("We are in the error path of the Registration handler")
		errorMessage = err.Error()
		fmt.Println(errorMessage)
	} else {
		successMessage = fmt.Sprintf("%s was added with the email %s", username, email)
	}
	fmt.Println(successMessage)
	allPosts := utils.GetAllPosts()
	data := models.DefaultTemplateData()
	categories := utils.GetActiveCategories()
	replies := utils.GetAllReplies()
	data.CommentReplies=replies
	data.Categories = categories
	data.AllPosts = allPosts
	data.ErrorMessage = errorMessage
	data.RegistrationSuccessMessage = "Account created successfully! You can now log in."
	data.Title = "Registration"
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
		sessionToken := uuid // Made it so the session has the same uuid as the user.
		if err != nil {
			errorMessage = errorMessage + " Failed to generate UUID"
		}
		sessionStore[sessionToken] = username
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

	allPosts := utils.GetAllPosts()
	data := models.DefaultTemplateData()
	categories := utils.GetActiveCategories()
	comments := utils.GetActiveComments()
	replies := utils.GetAllReplies()
	data.Comments = comments
	data.Username = username
	data.User = utils.FindUserByUserName(username)
	data.Categories = categories
	data.ErrorMessage = errorMessage
	data.Uuid = uuid
	data.AllPosts = allPosts
	data.Title = "Login"
	data.CommentReplies=replies
	if isActiveUser {
		data.IsLoggedIn = true
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
	sessionStore = map[string]string{} //! empties the sessions storage variable on logout.
	allPosts := utils.GetAllPosts()
	data := models.DefaultTemplateData()
	data.Title = "Logout"
	data.AllPosts = allPosts
	categories := utils.GetActiveCategories()
	data.Categories = categories
	for key := range sessionStore {
		delete(sessionStore, key)
		break
	}
	// render.RenderLandingPage(w, "index.html", data)
	LandingPageHandler(w,r)
}

func isValidSession(sessionToken string) bool {
	_, isValidSession := sessionStore[sessionToken]
	return isValidSession
}

func AuthSessionToken(next http.HandlerFunc) http.HandlerFunc {
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
	username, exists := sessionStore[sessionToken] //! Only works in this file as sessionStore is stored as a global variable.
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
	data := models.DefaultTemplateData()
	user := utils.FindUserByUserName(username)
	categories := utils.GetActiveCategories()
	comments := utils.GetActiveComments()
	replies := utils.GetAllReplies()
	data.Comments = comments
	data.Categories = categories
	data.User = user
	data.Categories = categories
	data.Username = ctxUsername
	data.ProfilePage = true
	data.Title = "Your Profile"
	data.IsLoggedIn = true
	data.CommentReplies=replies
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
	var postIDstr string
	path := r.URL.Path
	parts := strings.Split(path, "/")
	for _, part := range parts {
		// Check if the part is a numeric string
		_, err := strconv.Atoi(part)
		if err == nil { // Found the numeric part which is the postID
			postIDstr = part
			break
		}
	}
	if postIDstr == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
	}
	postID, _ := strconv.Atoi(postIDstr)
	HasLiked, HasDisliked := CheckUserActivity(postID, r)
	username := GetUsernameFromCookie(r)
	user := utils.FindUserByUserName(username)
	if !HasDisliked {
		if !HasLiked {
			db, err := sql.Open("sqlite3", config.LION_DB)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			defer db.Close()
			_, err = db.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", postIDstr)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			MarkPostAsLiked(user.ID, postID)
		} else {
			db, err := sql.Open("sqlite3", config.LION_DB)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			defer db.Close()
			_, err = db.Exec("UPDATE posts SET likes = likes - 1 WHERE id = ?", postIDstr)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			MarkPostAsUnliked(user.ID, postID)
		}
	}
	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var postIDstr string
	path := r.URL.Path
	parts := strings.Split(path, "/")
	for _, part := range parts {
		// Check if the part is a numeric string
		_, err := strconv.Atoi(part)
		if err == nil { // Found the numeric part which is the postID
			postIDstr = part
			break
		}
	}
	if postIDstr == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
	}
	postID, _ := strconv.Atoi(postIDstr)
	HasLiked, HasDisliked := CheckUserActivity(postID, r)
	username := GetUsernameFromCookie(r)
	user := utils.FindUserByUserName(username)
	if !HasLiked {
		if !HasDisliked {
			db, err := sql.Open("sqlite3", config.LION_DB)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			defer db.Close()
			_, err = db.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?", postIDstr)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			MarkPostAsDisliked(user.ID, postID)
		} else {
			db, err := sql.Open("sqlite3", config.LION_DB)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			defer db.Close()
			_, err = db.Exec("UPDATE posts SET dislikes = dislikes - 1 WHERE id = ?", postIDstr)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			MarkPostAsUndisliked(user.ID, postID)
		}
	}
	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func GetGetUserPostHistoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get user activity activated")
	LandingPageHandler(w, r)
}

func UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if !verifyPostMethod(w, r) {
		return
	}
	verifiedUserName := verifySession(r)
	if verifiedUserName == "" {
		fmt.Printf("not a user log in")
		LandingPageHandler(w, r)
	}
	sessionUser := utils.FindUserByUserName(verifiedUserName)
	var updatedUser models.User
	updatedUser.Password = r.FormValue("newPassword")
	passwordAgain := r.FormValue("newPasswordAgain")
	if  (updatedUser.Password == passwordAgain ) || updatedUser.Password ==""{
	} else {
		updatedUser.Password="0"
	}
	if (r.FormValue("email")) == "" {
		updatedUser.Email = sessionUser.Email
	} else {
		updatedUser.Email = r.FormValue("email")
	}
	if (r.FormValue("username")) == "" {
		updatedUser.Username = sessionUser.Username
	} else {
		updatedUser.Username = r.FormValue("username")
	}
	if (r.FormValue("role")) == "" {
		updatedUser.Role = sessionUser.Role
	} else {
		updatedUser.Role = r.FormValue("role")
	}
	userId := r.FormValue("ID")
	parsedInt, err := strconv.Atoi(userId)
	if err != nil {
		fmt.Println("Unable to get user ID")
	}
	updatedUser.ID = int(parsedInt)
	if sessionUser.ID == updatedUser.ID {
		utils.UpdateUserProfile(updatedUser)
	} else {
		fmt.Println("Not a user")
	}
	LogoutHandler(w,r)
}

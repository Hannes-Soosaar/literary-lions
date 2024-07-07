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
	models.GetInstance().SetSuccess("")
	models.GetInstance().SetError(nil)
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
		models.GetInstance().SetError(err)
	} else {
		successMessage = fmt.Sprintf("%s was added with the email %s", username, email)

	}
	models.GetInstance().SetSuccess(successMessage)
	allPosts := utils.GetAllPosts()
	data := models.DefaultTemplateData()
	categories := utils.GetActiveCategories()
	replies := utils.GetAllReplies()
	data.CommentReplies = replies
	data.Categories = categories
	data.AllPosts = allPosts
	data.ErrorMessage = errorMessage
	data.RegistrationSuccessMessage = "Account created successfully! You can now log in."
	data.Title = "Registration"
	render.RenderLandingPage(w, "index.html", data)
	models.GetInstance().SetSuccess("")
	models.GetInstance().SetError(nil)
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
		sessionToken := uuid
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
		errorMessage = "The session is not Invalid"
	}

	if models.GetInstance().GetError() != nil {
		errorMessage = models.GetInstance().GetError().Error()
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
	data.CommentReplies = replies
	if isActiveUser {
		data.IsLoggedIn = true
	}
	render.RenderLandingPage(w, "index.html", data)
	models.GetInstance().SetSuccess("")
	models.GetInstance().SetError(nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})
	sessionStore = map[string]string{}
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
	LandingPageHandler(w, r)
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
		// fmt.Println("CTX", ctx)
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
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
	data.CommentReplies = replies
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
		_, err := strconv.Atoi(part)
		if err == nil {
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
		_, err := strconv.Atoi(part)
		if err == nil {
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
	models.GetInstance().SetSuccess("")
	models.GetInstance().SetError(nil)
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

	formUsr := r.FormValue("username")
	formEmail := r.FormValue("email")
	formPwdNew := r.FormValue("newPassword")
	formPwdNewRepeat := r.FormValue("newPasswordAgain")

	userId := r.FormValue("ID")
	parsedInt, err := strconv.Atoi(userId)
	if err != nil {
		fmt.Println("Unable to get user ID")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if (formPwdNew != "") && (formPwdNewRepeat != "") {
		updatedUser.Password = formPwdNew

		if updatedUser.Password != formPwdNewRepeat {
			data := models.DefaultTemplateData()
			data.ProfileErrorMessage = "New passwords don't match!"
			render.RenderProfile(w, "index.html", data)
			return
		}
	}

	if formEmail == "" {
		updatedUser.Email = sessionUser.Email
	} else {
		fmt.Println("Email:", formEmail)
		if utils.UserWithEmailExists(formEmail) {
			data := models.DefaultTemplateData()
			data.ProfileErrorMessage = "This email is already in use!"
			if sessionUser.Email == formEmail {
				data.ProfileErrorMessage = "This is already your email!"
			}
			render.RenderProfile(w, "index.html", data)
			return
		}
		updatedUser.Email = formEmail
	}
	if formUsr == "" {
		updatedUser.Username = sessionUser.Username
	} else {
		fmt.Println("Username:", formUsr)
		if utils.UserWithUserNameExists(formUsr) {
			data := models.DefaultTemplateData()
			data.ProfileErrorMessage = "This username is already in use!"
			if sessionUser.Username == formUsr {
				data.ProfileErrorMessage = "This is already your username!"
			}
			render.RenderProfile(w, "index.html", data)
			return
		}
		updatedUser.Username = formUsr
	}

	if (r.FormValue("role")) == "" {
		updatedUser.Role = sessionUser.Role
	} else {
		updatedUser.Role = r.FormValue("role")
	}

	updatedUser.ID = int(parsedInt)

	if sessionUser.ID == updatedUser.ID {
		utils.UpdateUserProfile(updatedUser)
	} else {
		fmt.Println("Not a user")
	}

	LogoutHandler(w, r)
}

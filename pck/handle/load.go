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
	allPosts := utils.RetrieveAllPosts()

	data := models.DefaultTemplateData()
	data.IsLoggedIn = isLoggedIn
	data.MainPage = true
	data.ProfilePage = false
	data.AllPosts = allPosts

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
	allPosts := utils.RetrieveAllPosts()

	data := models.DefaultTemplateData()
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

	data := models.DefaultTemplateData()
	data.Username = username
	data.ErrorMessage = errorMessage
	data.Uuid = uuid
	data.AllPosts = allPosts
	data.Title = "Login"

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

	allPosts := utils.RetrieveAllPosts()
	data := models.DefaultTemplateData()
	data.Title = "Logout"
	data.AllPosts = allPosts

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
	data := models.DefaultTemplateData()
	data.Username = ctxUsername
	data.Title = "Your Profile"
	data.ProfilePage = true
	data.IsLoggedIn = true

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

func MarkPostAsLiked(userID, postID int) error {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		return err
	}
	defer db.Close()

	// Check if there's already an entry for this user and post
	var likeActivity bool
	err = db.QueryRow("SELECT like_activity FROM user_activity WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&likeActivity)
	switch {
	case err == sql.ErrNoRows:
		// No existing entry, insert a new one
		_, err := db.Exec("INSERT INTO user_activity (user_id, post_id, like_activity) VALUES (?, ?, 1)", userID, postID)
		if err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		// There's an existing entry, update it
		if !likeActivity {
			_, err := db.Exec("UPDATE user_activity SET like_activity = 1 WHERE user_id = ? AND post_id = ?", userID, postID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func MarkPostAsUnliked(userID, postID int) error {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		return err
	}
	defer db.Close()
	// Check if there's already an entry for this user and post
	var likeActivity bool
	err = db.QueryRow("SELECT like_activity FROM user_activity WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&likeActivity)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No existing entry?")
		// No existing entry, insert a new one with dislike_activity set to 1
		_, err := db.Exec("INSERT INTO user_activity (user_id, post_id, like_activity) VALUES (?, ?, 0)", userID, postID)
		if err != nil {
			fmt.Println("error", err)
			return err
		}
	case err != nil:
		return err
	default:
		// There's an existing entry, update it
		if likeActivity {
			_, err := db.Exec("UPDATE user_activity SET like_activity = 0 WHERE user_id = ? AND post_id = ?", userID, postID)
			if err != nil {
				fmt.Println("Liking set to 0")
				return err
			}
		}
	}

	return nil
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

func MarkPostAsDisliked(userID, postID int) error {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		return err
	}
	defer db.Close()

	// Check if there's already an entry for this user and post
	var dislikeActivity bool
	err = db.QueryRow("SELECT dislike_activity FROM user_activity WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&dislikeActivity)
	switch {
	case err == sql.ErrNoRows:
		// No existing entry, insert a new one
		_, err := db.Exec("INSERT INTO user_activity (user_id, post_id, dislike_activity) VALUES (?, ?, 1)", userID, postID)
		if err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		// There's an existing entry, update it
		if !dislikeActivity {
			_, err := db.Exec("UPDATE user_activity SET dislike_activity = 1 WHERE user_id = ? AND post_id = ?", userID, postID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func MarkPostAsUndisliked(userID, postID int) error {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		return err
	}
	defer db.Close()
	// Check if there's already an entry for this user and post
	var dislikeActivity bool
	err = db.QueryRow("SELECT dislike_activity FROM user_activity WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&dislikeActivity)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No existing entry?")
		// No existing entry, insert a new one with dislike_activity set to 1
		_, err := db.Exec("INSERT INTO user_activity (user_id, post_id, dislike_activity) VALUES (?, ?, 0)", userID, postID)
		if err != nil {
			fmt.Println("error", err)
			return err
		}
	case err != nil:
		return err
	default:
		// There's an existing entry, update it
		if dislikeActivity {
			_, err := db.Exec("UPDATE user_activity SET dislike_activity = 0 WHERE user_id = ? AND post_id = ?", userID, postID)
			if err != nil {
				fmt.Println("Disiking set to 0")
				return err
			}
		}
	}

	return nil
}

func CheckUserActivity(postID int, r *http.Request) (bool, bool) {

	type UserActivity struct {
		HasLiked    bool
		HasDisliked bool
	}

	username := GetUsernameFromCookie(r)
	user := utils.FindUserByUserName(username)

	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("Database error:", err)
		return false, false
	}
	defer db.Close()

	query := "SELECT like_activity, dislike_activity FROM user_activity WHERE user_id = ? AND post_id = ?"
	row := db.QueryRow(query, user.ID, postID)

	var activity UserActivity
	err = row.Scan(&activity.HasLiked, &activity.HasDisliked)
	if err != nil {
		if err == sql.ErrNoRows {
			// No activity found for user and post
			return false, false
		}
		fmt.Println("Database error:", err)
		return false, false
	}

	// Update HasLiked and HasDisliked based on user activity
	return activity.HasLiked, activity.HasDisliked
}

func GetUsernameFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return ""
	}
	sessionToken := cookie.Value
	username, exists := sessionStore[sessionToken]
	if !exists {
		return ""
	}
	return username
}

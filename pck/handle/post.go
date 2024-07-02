package handle

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/render"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
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
	user := utils.FindUserByUserName(username)
	categories := utils.GetActiveCategories()
	data.Categories = categories
	data.User = user
	data.Username = ctxUsername
	data.CreatePostPage = true
	data.IsLoggedIn = true
	data.Title = "Create Post"

	if isLoggedIn {
		render.RenderPostPage(w, "index.html", data)
	} else {
		data.ErrorMessage = "You need to be logged in to access your profile!"
		render.RenderLandingPage(w, "index.html", data)
	}
}

func SubmitPostHandler(w http.ResponseWriter, r *http.Request) {

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

	user := utils.FindUserByUserName(username)

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

	userID := user.ID
	categoryIDstr := r.Form.Get("category")
	categoryID, err := strconv.Atoi(categoryIDstr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	title := r.Form.Get("title")
	body := r.Form.Get("body")
	err = utils.AddNewPost(categoryID, title, body, userID)

	if err != nil {
		http.Error(w, "Error creating post", http.StatusBadRequest)
	}

	data := models.DefaultTemplateData()
	data.PostCreatedMessage = "Post was created successfully!"
	categories := utils.GetActiveCategories()
	data.Categories = categories
	data.User = user
	data.Username = user.Username
	data.CreatePostPage = true
	data.IsLoggedIn = true
	render.RenderPostPage(w, "index.html", data)
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

func UserPostsHandler(w http.ResponseWriter, r *http.Request) {
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
	user := utils.FindUserByUserName(username)
	categories := utils.GetActiveCategories()
	comments := utils.GetActiveComments()
	data.Comments = comments
	data.Categories = categories
	data.User = user
	allPosts := utils.RetrieveAllPosts()
	data.AllPosts = utils.UserPostsFinder(allPosts, data.User.ID)
	data.Categories = categories
	data.Username = ctxUsername
	data.UserPostsPage = true
	data.IsLoggedIn = true
	if isLoggedIn {
		render.RenderUserPosts(w, "posts-by-user.html", data)
	} else {
		data.ErrorMessage = "You need to be logged in to access your profile!"
		render.RenderLandingPage(w, "index.html", data)
	}
}

func LikedAndDislikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	var liked, disliked bool = false, false
	path := r.URL.Path
	switch path {
	case "/liked-posts":
		{
			liked = true
		}
	case "/disliked-posts":
		{
			disliked = true
		}

	}
	fmt.Println(path)
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
	user := utils.FindUserByUserName(username)
	categories := utils.GetActiveCategories()
	comments := utils.GetActiveComments()
	data.Comments = comments
	data.Categories = categories
	data.User = user
	allPosts := utils.RetrieveAllPosts()
	data.Categories = categories
	data.Username = ctxUsername
	data.IsLoggedIn = true
	if liked || disliked {
		if isLoggedIn {
			switch {
			case liked:
				{
					data.AllPosts = utils.FindUserLikedPosts(allPosts, data.User.ID)
					data.LikedPostsPage = true
				}
			case disliked:
				{
					data.AllPosts = utils.FindUserDislikedPosts(allPosts, data.User.ID)
					data.DislikedPostsPage = true
				}
			}
			render.RenderUserPosts(w, "posts-by-user.html", data)
		} else {
			data.ErrorMessage = "You need to be logged in to access your profile!"
			render.RenderLandingPage(w, "index.html", data)
		}
	} else {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
}

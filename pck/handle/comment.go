package handle

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	fmt.Printf("the full content of r is %v \n", r)
	commentIdString := r.FormValue("commentID")
	commentId, _ := strconv.Atoi(commentIdString)
	postIdString := r.FormValue("postID")
	postId, _ := strconv.Atoi(postIdString)

	comment := r.FormValue("comment")
	if !verifyPostMethod(w, r) {
		return
	}
	verifiedUserName := verifySession(r)
	if verifiedUserName == "" {
		fmt.Printf("not a user log in")
		LandingPageHandler(w, r)
	}
	sessionUser := utils.FindUserByUserName(verifiedUserName)

	if postIdString == "" {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	fmt.Printf("the Comment we want to reply to is: %d \n ", commentId)
	fmt.Printf("the Post we want to comment on is: %d \n ", postId)
	fmt.Printf("the comment we want to add to is: %s \n ", comment)

	if commentIdString == "" {
		fmt.Printf("Posting comment \n")
		utils.PostComment(sessionUser.ID, comment, postId)
	} else {
		fmt.Printf("Replying to comment \n")
		utils.CommentReply(comment, sessionUser.ID, commentId, postId)
	}

	http.Redirect(w, r, referer, http.StatusFound)
}

func CommentLikeHandler(w http.ResponseWriter, r *http.Request) {
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

func CommentDislikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// TODO: migrate to r.FromValue ?

	var commentIDstr string

	path := r.URL.Path
	parts := strings.Split(path, "/")
	for _, part := range parts {
		// Check if the part is a numeric string
		_, err := strconv.Atoi(part)
		if err == nil { // Found the numeric part which is the postID
			commentIDstr = part
			break
		}
	}
	if commentIDstr == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
	}
	commentID, _ := strconv.Atoi(commentIDstr)

	HasLiked, HasDisliked := CheckUserActivity(commentID, r) // ? OK

	username := GetUsernameFromCookie(r)       // ? OK
	user := utils.FindUserByUserName(username) // ? OK
	if !HasLiked {
		if !HasDisliked {
			db, err := sql.Open("sqlite3", config.LION_DB)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			defer db.Close()
			_, err = db.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?", commentIDstr)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			MarkPostAsDisliked(user.ID, commentID) // TODO mod this function for comments
		} else {
			db, err := sql.Open("sqlite3", config.LION_DB)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			defer db.Close()
			_, err = db.Exec("UPDATE posts SET dislikes = dislikes - 1 WHERE id = ?", commentIDstr)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			MarkPostAsUndisliked(user.ID, commentID)
		}
	}
	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

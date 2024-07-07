package handle

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
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
		LandingPageHandler(w, r)
	}
	sessionUser := utils.FindUserByUserName(verifiedUserName)
	if postIdString == "" {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	if commentIdString == "" {
		utils.PostComment(sessionUser.ID, comment, postId)
	} else {
		utils.CommentReply(comment, sessionUser.ID, commentId, postId)
	}
	http.Redirect(w, r, referer, http.StatusFound)
}

func CommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var commentIDstr string
	path := r.URL.Path
	parts := strings.Split(path, "/")
	for _, part := range parts {
		_, err := strconv.Atoi(part)
		if err == nil { 
			commentIDstr = part
			break
		}
	}
	if commentIDstr == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
	}
	commentID, _ := strconv.Atoi(commentIDstr)
	HasLiked, HasDisliked := CheckUserReplyActivity(commentID, r)
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
			_, err = db.Exec("UPDATE comments SET likes = likes + 1 WHERE id = ?", commentIDstr)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			MarkCommentAsLiked(user.ID, commentID)
		} else {
			db, err := sql.Open("sqlite3", config.LION_DB)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			defer db.Close()
			_, err = db.Exec("UPDATE comments SET likes = likes - 1 WHERE id = ?", commentIDstr)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			MarkCommentAsUnliked(user.ID, commentID)
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
	var commentIDstr string
	path := r.URL.Path
	parts := strings.Split(path, "/")
	for _, part := range parts {
		_, err := strconv.Atoi(part)
		if err == nil {
			commentIDstr = part
			break
		}
	}
	if commentIDstr == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
	}
	commentID, _ := strconv.Atoi(commentIDstr)
	HasLiked, HasDisliked := CheckUserReplyActivity(commentID, r)
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
			_, err = db.Exec("UPDATE comments SET dislikes = dislikes + 1 WHERE id = ?", commentIDstr)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			MarkCommentAsDisliked(user.ID, commentID)
		} else {
			db, err := sql.Open("sqlite3", config.LION_DB)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			defer db.Close()
			_, err = db.Exec("UPDATE comments SET dislikes = dislikes - 1 WHERE id = ?", commentIDstr)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			MarkCommentAsUndisliked(user.ID, commentID)
		}
	}
	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

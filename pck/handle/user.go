package handle

import (
	"database/sql"
	"fmt"
	"net/http"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

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

func CheckUserReplyActivity(commentID int, r *http.Request) (bool, bool) {
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
	query := "SELECT like_activity, dislike_activity FROM user_reply_activity WHERE user_id = ? AND comment_id = ?"
	row := db.QueryRow(query, user.ID, commentID)
	var activity UserActivity
	err = row.Scan(&activity.HasLiked, &activity.HasDisliked)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, false
		}
		fmt.Println("Database error:", err)
		return false, false
	}
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

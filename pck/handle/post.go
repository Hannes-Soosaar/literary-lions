package handle

import (
	"database/sql"
	"fmt"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
)

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

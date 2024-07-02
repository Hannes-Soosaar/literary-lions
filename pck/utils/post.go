package utils

import (
	"database/sql"
	"fmt"

	// "gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
)

func AddNewPost(categoryID int, title string, body string, userID int) error {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		return err
	}
	defer db.Close()
	statement, err := db.Prepare("INSERT INTO posts(title, body, likes, dislikes, user_id, category_id, created_at, modified_at, active) VALUES (?, ?, 0, 0, ?, ?, datetime('now'), datetime('now'), 1)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(title, body, userID, categoryID)
	if err != nil {
		return err
	}

	return nil
}

func RetrieveAllPosts() models.Posts {
	var posts models.Posts
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		return posts
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, title, body, likes, dislikes, user_id, category_id, created_at, modified_at, active FROM posts ORDER BY created_at DESC")
	if err != nil {
		fmt.Println("Error querying DB:", err)
		return posts
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.Likes, &post.Dislikes, &post.UserId, &post.CategoryID, &post.CreatedAt, &post.ModifiedAt, &post.Active)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		posts.AllPosts = append(posts.AllPosts, post)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}
	return posts
}

func FilterPostsByCategoryID(posts models.Posts, categoryID int) models.FilteredPosts {
	var filteredPosts models.FilteredPosts
	for _, post := range posts.AllPosts {
		if post.CategoryID == categoryID {
			filteredPosts.FilteredPosts = append(filteredPosts.FilteredPosts, post)
		}
	}
	return filteredPosts
}

func FilterPostByID(posts models.Posts, postID int) models.FilteredPosts {
	var filteredPost models.FilteredPosts
	for _, post := range posts.AllPosts {
		if post.ID == postID {
			filteredPost.FilteredPosts = append(filteredPost.FilteredPosts, post)
		}
	}
	return filteredPost
}

func UserPostsFinder(posts models.Posts, userID int) models.Posts {
	var userPosts models.Posts
	for _, post := range posts.AllPosts {
		if post.UserId == userID {
			userPosts.AllPosts = append(userPosts.AllPosts, post)
		}
	}
	return userPosts
}

func FindUserLikedPosts(posts models.Posts, userID int) models.Posts {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		return posts
	}
	defer db.Close()
	rows, err := db.Query("SELECT user_id, post_id, like_activity FROM user_activity")
	if err != nil {
		fmt.Println("Error querying DB:", err)
		return posts
	}
	defer rows.Close()

	var likedPostIDs []int
	for rows.Next() {
		var userIDfromDB int
		var postID int
		var likeActivity bool
		err := rows.Scan(&userIDfromDB, &postID, &likeActivity)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		if userIDfromDB == userID && likeActivity {
			likedPostIDs = append(likedPostIDs, postID)
		}
	}
	if err = rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}

	var likedPosts models.Posts
	for _, post := range posts.AllPosts {
		for _, postID := range likedPostIDs {
			if postID == post.ID {
				likedPosts.AllPosts = append(likedPosts.AllPosts, post)
				continue
			}
		}
	}

	return likedPosts
}

func FindUserDislikedPosts(posts models.Posts, userID int) models.Posts {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		return posts
	}
	defer db.Close()
	rows, err := db.Query("SELECT user_id, post_id, dislike_activity FROM user_activity")
	if err != nil {
		fmt.Println("Error querying DB:", err)
		return posts
	}
	defer rows.Close()

	var dislikedPostIDs []int
	for rows.Next() {
		var userIDfromDB int
		var postID int
		var dislikeActivity bool
		err := rows.Scan(&userIDfromDB, &postID, &dislikeActivity)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		if userIDfromDB == userID && dislikeActivity {
			dislikedPostIDs = append(dislikedPostIDs, postID)
		}
	}
	if err = rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}

	var dislikedPosts models.Posts
	for _, post := range posts.AllPosts {
		for _, postID := range dislikedPostIDs {
			if postID == post.ID {
				dislikedPosts.AllPosts = append(dislikedPosts.AllPosts, post)
				continue
			}
		}
	}

	return dislikedPosts
}

func UserPostsFinder(posts models.Posts, userID int) models.Posts {
	var userPosts models.Posts
	for _, post := range posts.AllPosts {
		if post.UserId == userID {
			userPosts.AllPosts = append(userPosts.AllPosts, post)
		}
	}
	return userPosts
}

func FindUserLikedPosts(posts models.Posts, userID int) models.Posts {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		return posts
	}
	defer db.Close()
	rows, err := db.Query("SELECT user_id, post_id, like_activity FROM user_activity")
	if err != nil {
		fmt.Println("Error querying DB:", err)
		return posts
	}
	defer rows.Close()

	var likedPostIDs []int
	for rows.Next() {
		var userIDfromDB int
		var postID int
		var likeActivity bool
		err := rows.Scan(&userIDfromDB, &postID, &likeActivity)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		if userIDfromDB == userID && likeActivity {
			likedPostIDs = append(likedPostIDs, postID)
		}
	}
	if err = rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}

	var likedPosts models.Posts
	for _, post := range posts.AllPosts {
		for _, postID := range likedPostIDs {
			if postID == post.ID {
				likedPosts.AllPosts = append(likedPosts.AllPosts, post)
				continue
			}
		}
	}

	return likedPosts
}

func FindUserDislikedPosts(posts models.Posts, userID int) models.Posts {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		return posts
	}
	defer db.Close()
	rows, err := db.Query("SELECT user_id, post_id, dislike_activity FROM user_activity")
	if err != nil {
		fmt.Println("Error querying DB:", err)
		return posts
	}
	defer rows.Close()

	var dislikedPostIDs []int
	for rows.Next() {
		var userIDfromDB int
		var postID int
		var dislikeActivity bool
		err := rows.Scan(&userIDfromDB, &postID, &dislikeActivity)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		if userIDfromDB == userID && dislikeActivity {
			dislikedPostIDs = append(dislikedPostIDs, postID)
		}
	}
	if err = rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}

	var dislikedPosts models.Posts
	for _, post := range posts.AllPosts {
		for _, postID := range dislikedPostIDs {
			if postID == post.ID {
				dislikedPosts.AllPosts = append(dislikedPosts.AllPosts, post)
				continue
			}
		}
	}

	return dislikedPosts
}

func FindPostsByUserName(userName string) models.Posts {
	var posts models.Posts
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		return posts
	}
	defer db.Close()
		statement, err := db.Prepare("SELECT id, title, body, likes, dislikes, user_id, category_id, created_at, modified_at, active FROM posts WHERE user_id = ? ORDER BY created_at DESC")
	if err != nil {
		return posts
	}
	defer statement.Close()
	rows, err := statement.Query(userName)
	if err != nil {
		fmt.Println("Error querying DB:", err)
		return posts
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.Likes, &post.Dislikes, &post.UserId, &post.CategoryID, &post.CreatedAt, &post.ModifiedAt, &post.Active)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		posts.AllPosts = append(posts.AllPosts, post)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}
	return models.Posts{}
}

func FindPostByCategory(categoryName string) models.Posts {
	//TODO: GetPostFromCategory

	return models.Posts{}
}

func UpdatedComment(commentId int) {

}

func UpdateEmotes(emote string) {
	//TODO: GetAllPostLikes
	//TODO: GetAllPostDislikes no need to implement, if there is a counter you can just
	//TODO: AddLike and RemoveLike
	//TODO: Add dislike and RemoveDisLike

}
func CreateNewPost(post string, userName string) {
	// TODO: AddUserPost
}


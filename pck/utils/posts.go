package utils

import (
	"database/sql"
	"fmt"

	// "gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
)

func RetrieveAllPosts() models.Posts {
	var posts models.Posts

	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		return posts
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, body, likes, dislikes, user_id, category_id, created_at, modified_at, active FROM posts")
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

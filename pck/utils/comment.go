package utils

import (
	"database/sql"
	"fmt"
	"log"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
)

//TODO: GetCommentsForPost
//TODO: GetCommentsContaining

// ? HANDLED also with html template logic for all comments.
func GetActivePostComments(postId int) []models.Comment {
	var activeComments []models.Comment
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query := "SELECT * FROM comments WHERE active = ? AND post_id = ?"
	rows, err := db.Query(query, config.ACTIVE, postId)
	if err != nil {
		fmt.Printf("there is an error getting rows %v \n", err)
		return []models.Comment{}
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.Body, &comment.UserId, &comment.Likes, &comment.Dislikes, &comment.PostID, &comment.CreatedAt, &comment.ModifiedAt, &comment.Active)
		if err != nil {
			fmt.Printf("error reading from a row %v  \n", err)
			return activeComments
		}
		activeComments = append(activeComments, comment)
	}
	err = rows.Err()
	if err != nil {
		fmt.Printf("error occurred during rows iteration %v \n", err)
		return activeComments
	}
	return activeComments
}
func GetActiveComments() []models.Comment {
	var activeComments []models.Comment
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query := "SELECT * FROM comments WHERE active = ?"
	rows, err := db.Query(query, config.ACTIVE)
	if err != nil {
		fmt.Printf("there is an error getting rows %v \n", err)
		return []models.Comment{}
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.Body, &comment.UserId, &comment.Likes, &comment.Dislikes, &comment.PostID, &comment.CreatedAt, &comment.ModifiedAt, &comment.Active)
		if err != nil {
			fmt.Printf("error reading from a row %v  \n", err)
			return activeComments
		}
		activeComments = append(activeComments, comment)
	}
	err = rows.Err()
	if err != nil {
		fmt.Printf("error occurred during rows iteration %v \n", err)
		return activeComments
	}
	return activeComments
}

func GetActiveUserComments(userId int) []models.Comment {
	var activeComments []models.Comment
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query := "SELECT * FROM comments WHERE active = ? AND user_id = ?"
	rows, err := db.Query(query, config.ACTIVE, userId)
	if err != nil {
		fmt.Printf("there is an error getting rows %v \n", err)
		return []models.Comment{}
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.Body, &comment.UserId, &comment.Likes, &comment.Dislikes, &comment.PostID, &comment.CreatedAt, &comment.ModifiedAt, &comment.Active)
		if err != nil {
			fmt.Printf("error reading from a row %v  \n", err)
			return activeComments
		}
		activeComments = append(activeComments, comment)
	}
	err = rows.Err()
	if err != nil {
		fmt.Printf("error occurred during rows iteration %v \n", err)
		return activeComments
	}
	return activeComments
}

func PostComment(userId int, comment string, postId string) error {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		log.Fatal(err)
		return err
	}
	 query := "INSERT INTO comments(body,user_id,likes,dislikes,post_id,created_at,modified_at,active) VALUES (?,?,0,0,?,datetime('now'),datetime('now'),1)"
	_,err = db.Exec(query,comment,userId,postId)
	if err != nil {
	log.Fatal(err)
		return err
	}

	return nil
}

func LikeComment(commentID int) {

}

func DislikeComment(commentID int) {

}

func RemoveComment(commentID int) {

}

package utils

import (
	"database/sql"
	"fmt"
	"log"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
)

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



func PostComment(userId int, comment string, postId int) error {
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
//TODO: start work on liking comments
func LikeComment(commentID int) {

}
//TODO: start work on liking comments
func DislikeComment(commentID int) {

}

//! discontinued logic.
// func PostChildComment(parentCommentId int, childCommentId int) error {
// 	db, err := sql.Open("sqlite3", config.LION_DB)
// 	if err != nil {
// 		log.Fatal(err)
// 		return err
// 	}
// 	query := "INSERT INTO comment_relations (parent_comment_id,child_comment_id)"
// 	_,err = db.Exec(query,parentCommentId,childCommentId)
// 	if err != nil {
// 	log.Fatal(err)
// 		return err
// 	}
// 	return nil
// }



//! discontinued logic 
// func RemoveCommentById(commentID int) (string, error) {
// 	var successMessage string
// 	db, err := sql.Open("sqlite3", config.LION_DB)
// 	if err != nil {
// 		log.Fatal(err)
// 		return successMessage, err
// 	}
// 	query := "DELETE FROM comments WHERE id = ?"
// 	_,err = db.Exec(query,commentID)
// 	if err != nil {
// 	log.Fatal(err)
// 		return successMessage, err
// 	}
// 	return successMessage, err
// }

func CommentReply(reply string, userId, CommentId, postId int) error {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		log.Fatal(err)
		return err
	}
	query := "INSERT INTO comment_replies(body,user_id,comment_id,post_id,created_at, active) VALUES (?,?,?,?,datetime('now'),1)"
	_,err = db.Exec(query,reply,userId,CommentId,postId)
	if err != nil {
	log.Fatal(err)
		return err
	} 
	return nil
}
// NOT USED ON THE SITE YET
func GetCommentRepliesById (CommentId int) []models.CommentReply {
	var commentReplies []models.CommentReply
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query := "SELECT * FROM comment_replies WHERE  comment_id = ?"
	rows, err := db.Query(query,CommentId)
	if err != nil {
		fmt.Printf("there is an error getting rows %v \n", err)
		return []models.CommentReply{}
	}
	defer rows.Close()
	for rows.Next() {
		var CommentReply models.CommentReply
		err := rows.Scan(&CommentReply.ID,&CommentReply.Body,&CommentReply.UserId,&CommentReply.CommentId, &CommentReply.PostId, &CommentReply.CreatedAt, &CommentReply.CreatedAt )
		if err != nil {
			fmt.Printf("error reading from a row %v  \n", err)
			return commentReplies
		}
		commentReplies = append(commentReplies, CommentReply)
	}
	err = rows.Err()
	if err != nil {
		fmt.Printf("error occurred during rows iteration %v \n", err)
		return commentReplies
	}
	return commentReplies
}

// Gets all replies 
func GetAllReplies () []models.CommentReply {
	var commentReplies []models.CommentReply
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query := "SELECT * FROM comment_replies"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("there is an error getting rows %v \n", err)
		return []models.CommentReply{}
	}
	defer rows.Close()
	for rows.Next() {
		var CommentReply models.CommentReply
		err := rows.Scan(&CommentReply.ID,&CommentReply.Body,&CommentReply.UserId,&CommentReply.CommentId, &CommentReply.PostId, &CommentReply.CreatedAt, &CommentReply.CreatedAt )
		if err != nil {
			fmt.Printf("error reading from a row %v  \n", err)
			return commentReplies
		}
		commentReplies = append(commentReplies, CommentReply)
	}
	err = rows.Err()
	if err != nil {
		fmt.Printf("error occurred during rows iteration %v \n", err)
		return commentReplies
	}
	return commentReplies
}

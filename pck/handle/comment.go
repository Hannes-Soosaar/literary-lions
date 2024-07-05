package handle

import (
	"fmt"
	"net/http"
	"strconv"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)



func CommentHandler(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	fmt.Printf("the full content of r is %v \n", r)
	commentIdString := r.FormValue("commentID")
	commentId,_ := strconv.Atoi(commentIdString)
	postIdString := r.FormValue("postID")
	postId,_ := strconv.Atoi( postIdString )

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

	fmt.Printf("the Comment we want to reply to is: %d \n ",commentId)
	fmt.Printf("the Post we want to comment on is: %d \n ",postId)
	fmt.Printf("the comment we want to add to is: %s \n ",comment)

	if commentIdString == "" {
		fmt.Printf("Posting comment \n")
	utils.PostComment(sessionUser.ID, comment, postId)
	} else {
		fmt.Printf("Replying to comment \n")
		utils.CommentReply(comment,sessionUser.ID,commentId,postId)
	}
	
	http.Redirect(w, r, referer, http.StatusFound)
}

package handle

import (
	"fmt"
	"net/http"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	if !verifyPostMethod(w, r) {
		return
	}
	verifiedUserName := verifySession(r)
	if verifiedUserName == "" {
		fmt.Printf("not a user log in")
		LandingPageHandler(w, r)
	}
	sessionUser := utils.FindUserByUserName(verifiedUserName)
	postId := r.FormValue("postID")
	fmt.Println("POST ID", postId)
	if postId == "" {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	comment := r.FormValue("comment")
	utils.PostComment(sessionUser.ID, comment, postId)
	http.Redirect(w, r, referer, http.StatusFound)
}

package main

import (
	"fmt"
	"net/http"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/handle"
)

func main() {

	fs := http.FileServer(http.Dir("../../static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Printf("server is running and listening on  Port %s \n", config.PORT)
	http.HandleFunc("/", handle.LandingPageHandler)
	http.HandleFunc("/register", handle.RegistrationHandler)
	http.HandleFunc("/login", handle.LoginHandler)
	http.HandleFunc("/logout", handle.LogoutHandler)
	http.HandleFunc("/profile", handle.ProfileHandler)
	http.HandleFunc("/like/", handle.LikeHandler)
	http.HandleFunc("/dislike/", handle.DislikeHandler)
	http.HandleFunc("/category/", handle.CategoryHandler)
	http.HandleFunc("/postComment/", handle.CommentHandler)
	http.HandleFunc("/create-post", handle.CreatePostHandler)
	http.HandleFunc("/submit-post", handle.SubmitPostHandler)
	http.HandleFunc("/search", handle.SearchHandler)
	http.HandleFunc("/userPostHistory", handle.GetGetUserPostHistoryHandler)
	http.HandleFunc("/updateUserProfile", handle.UpdateUserProfileHandler) 


	err := http.ListenAndServe(config.PORT, nil)
	if err != nil {
		fmt.Printf("Error:%s", err)
	}

	fmt.Printf("Server started on Port: %s \n", config.PORT)
}

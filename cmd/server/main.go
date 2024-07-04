package main

import (
	"fmt"
	"net/http"
	"strings"

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
	http.HandleFunc("/postComment/", handle.CommentHandler)
	http.HandleFunc("/create-post", handle.CreatePostHandler)
	http.HandleFunc("/submit-post", handle.SubmitPostHandler)
	http.HandleFunc("/category/", handleCategoryOrSearch)
	http.HandleFunc("/search", handle.SearchHandler)
	http.HandleFunc("/your-posts", handle.UserPostsHandler)
	http.HandleFunc("/liked-posts", handle.LikedAndDislikedPostsHandler) //TODO: make copy for comments
	http.HandleFunc("/disliked-posts", handle.LikedAndDislikedPostsHandler) // TODO: make copy for comments
	http.HandleFunc("/userPostHistory", handle.GetGetUserPostHistoryHandler)
	http.HandleFunc("/updateUserProfile", handle.UpdateUserProfileHandler)

	err := http.ListenAndServe(config.PORT, nil)
	if err != nil {
		fmt.Printf("Error:%s", err)
	}

	fmt.Printf("Server started on Port: %s \n", config.PORT)
}

func handleCategoryOrSearch(w http.ResponseWriter, r *http.Request) {
	// Split the URL path
	parts := strings.Split(r.URL.Path, "/")
	// Check if the URL path matches /category/any_number/search
	if len(parts) >= 4 && parts[3] == "search" {
		// Call SearchHandler for /category/any_number/search
		handle.SearchHandler(w, r)
		return
	}

	// Call CategoryHandler for other /category/ URLs
	handle.CategoryHandler(w, r)
}

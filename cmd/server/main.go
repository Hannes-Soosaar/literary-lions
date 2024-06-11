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
	fmt.Printf("server is running and listening on  Port %s", config.PORT)
	http.HandleFunc("/", handle.LandingPageHandler)
	http.HandleFunc("/register", handle.RegistrationHandler)
	http.HandleFunc("/login", handle.LoginHandler)
	err := http.ListenAndServe(config.PORT, nil)
	if err != nil {
		fmt.Printf("Error:%s", err)
	}

	fmt.Printf("Server started on Port: %s \n", config.PORT)
}

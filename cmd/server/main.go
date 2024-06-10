package main

import (
	"fmt"
	"log"
	"net/http"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/handle"
	utils "gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)


func main() {
	
	fs := http.FileServer(http.Dir("../../static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Printf("server is running and listening on  Port %s", config.PORT)
	http.HandleFunc("/", handle.LandingPageHandler)
	err := http.ListenAndServe(config.PORT, nil)
	if err != nil {
		fmt.Printf("Error:%s", err)
	}

	fmt.Printf("Server started on Port: %s \n", config.PORT)	

	fmt.Println("Hello Lions!")
	for i := 0; i < 5; i++ {
		Id, err := utils.GenerateUUID()
		if err != nil {
			log.Printf("Error 1! %v", err)
		}
		fmt.Printf("Id %d, is %s", i, Id)
	}

}

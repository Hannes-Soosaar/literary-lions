package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"log"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	utils "gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

// move out from here!
func RenderLandingPage(w http.ResponseWriter, tmpl string, data  interface{}){
	template := template.Must(template.ParseFiles(
		filepath.Join("../../template","base.html"),
		filepath.Join("../../template","head.html"),
		filepath.Join("../../template",tmpl),
	))
	err := template.ExecuteTemplate(w,"base.html",data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


// move to handler
func LandingPageHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "Test page",
	}
RenderLandingPage(w, "index.html", data)
}

func main() {
		fs := http.FileServer(http.Dir("../../static"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		fmt.Printf("server is running and listening on  Port %s", config.PORT)
		http.HandleFunc("/", LandingPageHandler)
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

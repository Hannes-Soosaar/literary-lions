package handle

import (
	"fmt"
	"net/http"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/render"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func LandingPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Loading Page")
	utils.FindUserByUserName("bob")
	data := struct {
		Title string
	}{
		Title: "Lions",
	}

	render.RenderLandingPage(w, "index.html", data)
}

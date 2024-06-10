package handle

import (
	"net/http"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/render"
)


func LandingPageHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "Lions",
	}
	render.RenderLandingPage(w, "index.html", data)
}
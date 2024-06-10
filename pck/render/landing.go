package render

import (
	"html/template"
	"net/http"
	"path/filepath"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func RenderLandingPage(w http.ResponseWriter, tmpl string, data interface{}) {
	utils.FindUserByUserName("bob")
	template := template.Must(template.ParseFiles(
		filepath.Join("../../template", "index.html"),
		filepath.Join("../../template", "head.html"),
		filepath.Join("../../template", "navbar.html"),
		filepath.Join("../../template", "sidebar.html"),
		filepath.Join("../../template", "forum-content.html"),
	))
	err := template.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

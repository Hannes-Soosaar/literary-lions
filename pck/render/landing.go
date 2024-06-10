package render

import (
	"html/template"
	"net/http"
	"path/filepath"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)



func RenderLandingPage(w http.ResponseWriter, tmpl string, data interface{}) {
	utils.GetUserByUserName("bob")
	template := template.Must(template.ParseFiles(
		filepath.Join("../../template", "base.html"),
		filepath.Join("../../template", "head.html"),
		filepath.Join("../../template", tmpl),
	))
	err := template.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
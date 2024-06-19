package render

import (
	"html/template"
	"net/http"
	"path/filepath"
)


func RenderProfile(w http.ResponseWriter, tmpl string, data interface{}) {
	template := template.Must(template.ParseFiles(
		filepath.Join("../../template", "index.html"),
		filepath.Join("../../template", "head.html"),
		filepath.Join("../../template", "navbar.html"),
		filepath.Join("../../template", "sidebar.html"),
		filepath.Join("../../template", "profile.html"),
	))
	err := template.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

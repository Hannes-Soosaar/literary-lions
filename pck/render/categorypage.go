package render

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func RenderCategoryPage(w http.ResponseWriter, tmpl string, data interface{}) {
	template := template.Must(template.ParseFiles(
		filepath.Join("../../template", "categories.html"),
		filepath.Join("../../template", "head.html"),
		filepath.Join("../../template", "navbar.html"),
		filepath.Join("../../template", "sidebar.html"),
		filepath.Join("../../template", "filtered-posts.html"),
	))
	err := template.ExecuteTemplate(w, "categories.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

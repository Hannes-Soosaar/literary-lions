package render

import (
	"html/template"
	"net/http"
	"path/filepath"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func RenderCategoryPage(w http.ResponseWriter, tmpl string, data interface{}) {
	funcMap := template.FuncMap{
        "add": utils.Add,
    }
	template := template.Must(template.New("categories.html").Funcs(funcMap).ParseFiles(
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

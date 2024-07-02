package render

import (
	"html/template"
	"net/http"
	"path/filepath"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func RenderUserPosts(w http.ResponseWriter, tmpl string, data interface{}) {
	funcMap := template.FuncMap{
		"add": utils.Add,
	}

	template := template.Must(template.New("categories.html").Funcs(funcMap).ParseFiles(
		filepath.Join("../../template", "index.html"),
		filepath.Join("../../template", "head.html"),
		filepath.Join("../../template", "navbar.html"),
		filepath.Join("../../template", "sidebar.html"),
		filepath.Join("../../template", tmpl),
	))
	err := template.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

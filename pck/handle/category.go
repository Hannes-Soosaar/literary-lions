package handle

import (
	"net/http"
	"strconv"
	"strings"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/render"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}
	categoryIDstr := parts[1]
	categoryID, err := strconv.Atoi(categoryIDstr)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}

	sessionToken, err := r.Cookie("session_token")
	isLoggedIn := err == nil && isValidSession(sessionToken.Value)
	allPosts := utils.RetrieveAllPosts()
	categories := utils.GetActiveCategories()

	if len(parts) == 2 {
		data := models.DefaultTemplateData()
		data.IsLoggedIn = isLoggedIn
		data.ProfilePage = false
		data.FilteredPosts = utils.FilterPostsByCategoryID(allPosts, categoryID)
		data.Categories = categories
		if isLoggedIn {
			if data.Username == "" {
				data.Username = GetUsernameFromCookie(r)
				if data.Username == "" {
					render.RenderCategoryPage(w, "categories.html", data)
				}
			}
		}
		render.RenderCategoryPage(w, "categories.html", data)
	}

	if len(parts) == 3 {
		//Send back a template with one post + comments
	}
}

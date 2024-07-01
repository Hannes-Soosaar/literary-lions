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
	if len(categoryIDstr) != 1 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	categoryID, err := strconv.Atoi(categoryIDstr)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}

	sessionToken, err := r.Cookie("session_token")
	isLoggedIn := err == nil && isValidSession(sessionToken.Value)
	allPosts := utils.RetrieveAllPosts()
	categories := utils.GetActiveCategories()
	data := models.DefaultTemplateData()
	comments := utils.GetActiveComments()
	data.Comments = comments
	data.IsLoggedIn = isLoggedIn
	data.ProfilePage = false
	data.Categories = categories
	data.ShowComments = false

	if len(parts) == 2 {

		data.FilteredPosts = utils.FilterPostsByCategoryID(allPosts, categoryID)
		for _, cat := range data.Categories {
			if cat.ID == categoryID {
				data.Title = cat.Category
			}
		}
		data.ShowComments = false
		data.DisplayCatID = categoryID
		if isLoggedIn {
			if data.Username == "" {
				data.Username = GetUsernameFromCookie(r)
				if data.Username == "" {
					render.RenderCategoryPage(w, "category-filtered-posts.html", data)
				}
			}
		}
		render.RenderCategoryPage(w, "category-filtered-posts.html", data)
	}

	if len(parts) == 3 {
		data.DisplayCatID = categoryID
		postIDstr := parts[2]
		postID, err := strconv.Atoi(postIDstr)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
		}
		data.FilteredPosts = utils.FilterPostByID(allPosts, postID)
		data.ShowComments = true
		if isLoggedIn {
			if data.Username == "" {
				data.Username = GetUsernameFromCookie(r)
				if data.Username == "" {
					render.RenderCategoryPage(w, "filtered-posts.html", data)
				}
			}
		}
		render.RenderCategoryPage(w, "filtered-posts.html", data)
	}
	if len(parts) > 3 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

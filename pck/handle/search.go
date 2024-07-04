package handle

import (
	"net/http"
	"strconv"
	"strings"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/render"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if !verifyGetMethod(w, r) {
		return
	}
	r.ParseForm()
	FilterType := r.Form.Get("filter-type")
	SearchQuery := r.Form.Get("search-query")
	sessionToken, err := r.Cookie("session_token")
	isLoggedIn := err == nil && isValidSession(sessionToken.Value)
	categories := utils.GetActiveCategories()
	comments := utils.GetActiveComments()
	data := models.DefaultTemplateData()
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 3 && parts[2] == "search" {
		categoryIDstr := parts[1]
		if len(categoryIDstr) != 1 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		categoryID, err := strconv.Atoi(categoryIDstr)

		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
		}
		data.CategoryPage = true
		data.DisplayCatID = categoryID
		allPosts := utils.FilterPostForSearch(FilterType, SearchQuery, categoryID)
		data.AllPosts = allPosts
	} else {
		allPosts := utils.FilterPostForSearch(FilterType, SearchQuery, 0)
		data.AllPosts = allPosts
		data.MainPage = true
	}

	data.IsLoggedIn = isLoggedIn
	data.Title = "Search results"
	data.ProfilePage = false
	data.SearchQuery = SearchQuery
	data.FilterType = FilterType
	if (len(data.AllPosts.AllPosts)) == 0 {
		var message string = "No results for keyword \"" + SearchQuery + "\""
		switch FilterType {
		case "likes":
			{
				message += " filtering by likes"
			}
		case "dislikes":
			{
				message += " filtering by dislikes"
			}
		case "time_new":
			{
				message += " filtering by newest first"
			}
		case "time_old":
			{
				message += " filtering by oldest first"
			}
		}
		data.QueryNoResult = message
	}
	data.Comments = comments
	data.Categories = categories
	if isLoggedIn {
		if data.Username == "" {
			data.Username = GetUsernameFromCookie(r)
			if data.Username == "" {
				if len(parts) >= 3 && parts[2] == "search" {
					render.RenderCategoryPage(w, "category-filtered-posts.html", data)
					return

				} else {
					render.RenderLandingPage(w, "index.html", data)
					return
				}

			}
		}
	}
	if len(parts) >= 3 && parts[2] == "search" {
		render.RenderCategoryPage(w, "category-filtered-posts.html", data)
		return
	} else {
		render.RenderLandingPage(w, "index.html", data)
		return
	}
}

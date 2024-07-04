package handle

import (
	"net/http"

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
	allPosts := utils.FilterPostForSearch(FilterType, SearchQuery)
	categories := utils.GetActiveCategories()
	comments := utils.GetActiveComments()
	data := models.DefaultTemplateData()
	data.IsLoggedIn = isLoggedIn
	data.MainPage = true
	data.Title = "Search results"
	data.ProfilePage = false
	data.AllPosts = allPosts
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
				render.RenderLandingPage(w, "index.html", data)
			}
		}
	}
	render.RenderLandingPage(w, "index.html", data)
}

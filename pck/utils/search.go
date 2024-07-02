package utils

import (
	"sort"
	"strings"
	"time"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
)

func FilterPostForSearch(FilterType string, SearchQuery string) models.Posts {
	var SortedPosts models.Posts
	var allPosts []models.Post
	allPosts = RetrieveAllPosts().AllPosts
	if SearchQuery != "" {
		allPosts = FilterPostByKeyword(SearchQuery, allPosts)
	}
	switch FilterType {
	case "likes":
		{
			sort.Slice(allPosts, func(i, j int) bool {
				return allPosts[i].Likes > allPosts[j].Likes
			})
		}
	case "dislikes":
		{
			sort.Slice(allPosts, func(i, j int) bool {
				return allPosts[i].Dislikes > allPosts[j].Dislikes
			})
		}
	case "time_new":
		{
			sort.Slice(allPosts, func(i, j int) bool {
				timeI, _ := time.Parse("2006-01-02 15:04:05", allPosts[i].CreatedAt)
				timeJ, _ := time.Parse("2006-01-02 15:04:05", allPosts[j].CreatedAt)
				return timeI.After(timeJ)
			})
		}
	case "time_old":
		{
			sort.Slice(allPosts, func(i, j int) bool {
				timeI, _ := time.Parse("2006-01-02 15:04:05", allPosts[i].CreatedAt)
				timeJ, _ := time.Parse("2006-01-02 15:04:05", allPosts[j].CreatedAt)
				return timeJ.After(timeI)
			})
		}
	}
	SortedPosts.AllPosts = allPosts
	return SortedPosts
}

func FilterPostByKeyword(SearchQuery string, allPosts []models.Post) []models.Post {
	var filteredPosts []models.Post
	for _, post := range allPosts {
		if strings.Contains(post.Title, SearchQuery) || strings.Contains(post.Body, SearchQuery) {
			filteredPosts = append(filteredPosts, post)
		}
	}
	return filteredPosts
}

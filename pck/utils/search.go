package utils

import (
	"sort"
	"strings"
	"time"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
)

func FilterPostForSearch(FilterType string, SearchQuery string, catID int) models.Posts {
	var SortedPosts models.Posts
	var allPosts []models.Post
	var postsToSort []models.Post
	if catID != 0 {
		allPosts = GetAllPosts().AllPosts
		for _, post := range allPosts {
			if post.CategoryID == catID {
				postsToSort = append(postsToSort, post)
			}
		}
	} else {
		postsToSort = GetAllPosts().AllPosts
	}
	if SearchQuery != "" {
		postsToSort = FilterPostByKeyword(SearchQuery, postsToSort)
	}
	switch FilterType {
	case "likes":
		{
			sort.Slice(postsToSort, func(i, j int) bool {
				return postsToSort[i].Likes > postsToSort[j].Likes
			})
		}
	case "dislikes":
		{
			sort.Slice(postsToSort, func(i, j int) bool {
				return postsToSort[i].Dislikes > postsToSort[j].Dislikes
			})
		}
	case "time_new":
		{
			sort.Slice(postsToSort, func(i, j int) bool {
				timeI, _ := time.Parse("2006-01-02 15:04:05", postsToSort[i].CreatedAt)
				timeJ, _ := time.Parse("2006-01-02 15:04:05", postsToSort[j].CreatedAt)
				return timeI.After(timeJ)
			})
		}
	case "time_old":
		{
			sort.Slice(postsToSort, func(i, j int) bool {
				timeI, _ := time.Parse("2006-01-02 15:04:05", postsToSort[i].CreatedAt)
				timeJ, _ := time.Parse("2006-01-02 15:04:05", postsToSort[j].CreatedAt)
				return timeJ.After(timeI)
			})
		}
	}
	SortedPosts.AllPosts = postsToSort
	return SortedPosts
}

func FilterPostByKeyword(SearchQuery string, allPosts []models.Post) []models.Post {

	var filteredPosts []models.Post
	lowerCaseSearchQuery := strings.ToLower(SearchQuery)

	if strings.Contains(lowerCaseSearchQuery, "title:") {
		_, lowerCaseSearchQuery, _ = strings.Cut(lowerCaseSearchQuery, "title:")
		for _, post := range allPosts {
			if strings.Contains(strings.ToLower(post.Title), lowerCaseSearchQuery) {
				filteredPosts = append(filteredPosts, post)
			}
		}
		return filteredPosts
	}

	if strings.Contains(lowerCaseSearchQuery, "body:") {
		_, lowerCaseSearchQuery, _ = strings.Cut(lowerCaseSearchQuery, "body:")
		for _, post := range allPosts {
			if strings.Contains(strings.ToLower(post.Body), lowerCaseSearchQuery) {
				filteredPosts = append(filteredPosts, post)
			}
		}
		return filteredPosts
	}

	for _, post := range allPosts {
		if strings.Contains(strings.ToLower(post.Title), lowerCaseSearchQuery) || strings.Contains(strings.ToLower(post.Body), lowerCaseSearchQuery) {
			filteredPosts = append(filteredPosts, post)
		}
	}

	return filteredPosts
}

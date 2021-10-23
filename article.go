package main

import (
	"sort"
	"time"
)

type Article struct {
	Title     string    `json:"title"`
	URL       *JSONURL  `json:"url,string"`
	FirstSeen time.Time `json:"first_seen,string"`
}

type sortableArticleSlice []*Article

func (articles sortableArticleSlice) Len() int {
	return len(articles)
}

func (articles sortableArticleSlice) Less(i, j int) bool {
	// We want the newest articles early in the list, so "Less" to us means "After"
	return (articles[i].FirstSeen.After(articles[j].FirstSeen))
}

func (articles sortableArticleSlice) Swap(i, j int) {
	articles[i], articles[j] = articles[j], articles[i]
}

// ArticleList is the struct that gets marshaled into a file by writeBlob
type ArticleList struct {
	Articles []*Article `json:"articles"`
}

// NewArticleList returns an ArticleList with the given Articles.
//
// The Articles will be ordered descending by Article.FirstSeen and limited to the first limit
// elements. If limit is 0, all elements are included.
func NewArticleList(articles []*Article, limit int) *ArticleList {
	if limit == 0 {
		limit = len(articles)
	}

	// sort
	sorted := make([]*Article, len(articles))
	copy(sorted, articles)
	sort.Sort(sortableArticleSlice(sorted))

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	// trim
	list := new(ArticleList)
	list.Articles = make([]*Article, min(len(articles), limit))
	copy(list.Articles, sorted)

	return list
}

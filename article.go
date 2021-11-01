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
	Limit    int
}

// Merge adds newArticles to the ArticleList.
//
// The Articles will be deduped (by URL, ordered descending by Article.FirstSeen, and limited to the
// first list.Limit elements. If list.Limit is 0, all elements are included.)
func (list *ArticleList) Merge(newArticles []*Article) {
	list.Articles = append(list.Articles, newArticles...)

	// dedupe
	urls := make(map[string]bool)
	for i := 0; i < len(list.Articles); i++ {
		a := list.Articles[i]
		urlString := jsonURLToURL(a.URL).String()
		if _, ok := urls[urlString]; ok {
			list.Articles[i] = list.Articles[len(list.Articles)-1]
			list.Articles = list.Articles[:len(list.Articles)-1]
		} else {
			urls[urlString] = true
		}
	}

	// sort
	sorted := make([]*Article, len(list.Articles))
	copy(sorted, list.Articles)
	sort.Sort(sortableArticleSlice(sorted))

	// trim
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	limit := min(len(list.Articles), list.Limit)
	if limit == 0 {
		limit = len(list.Articles)
	}
	list.Articles = make([]*Article, limit)
	copy(list.Articles, sorted)
}

// NewArticleList returns a new, empty ArticleList.
func NewArticleList(limit int) *ArticleList {
	return &ArticleList{Limit: limit}
}

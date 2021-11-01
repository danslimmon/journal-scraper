package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gocolly/colly/v2"
)

// elementToArticle takes an HTML element and returns a corresponding *Article.
//
// e corresponds to an <a> tag with the ".issue-item__title" class. baseURL is the base URL of the
// scrape (e.g. "https://onlinelibrary.wiley.com/"). t is the time we should set as the Article's
// FirstSeen value.
//
// If an *Article can't be inferred from the element, we return nil.
func elementToArticle(e *colly.HTMLElement, baseURL string, t time.Time) (*Article, error) {
	a := new(Article)

	href := e.Attr("href")
	if href == "" {
		return nil, fmt.Errorf("element has no or empty href value")
	}
	baseU, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	hrefU, err := url.Parse(href)
	if err != nil {
		return nil, err
	}
	a.URL = urlToJSONURL(baseU.ResolveReference(hrefU))

	// On some pages, <h2> is used for titles. On other pages, <h3> is used.
	title := e.ChildText("h2")
	if title == "" {
		title = e.ChildText("h3")
	}
	if title == "" {
		return nil, fmt.Errorf("no article title could be derived from element")
	}
	a.Title = title

	a.FirstSeen = t

	return a, nil
}

func main() {
	c := colly.NewCollector()
	articles := make([]*Article, 0)
	now := time.Now().UTC()

	// Each article has an <a>
	c.OnHTML("a.issue-item__title", func(e *colly.HTMLElement) {
		article, _ := elementToArticle(e, "https://onlinelibrary.wiley.com/", now)
		if article == nil {
			// This is screen scraping, which is always chaotic. We can't afford to be panicking on
			// every error.
			return
		}
		articles = append(articles, article)
	})

	err := c.Visit("https://onlinelibrary.wiley.com/toc/14764431/0/0")
	if err != nil {
		panic(err.Error())
	}
	if len(articles) < 5 {
		panic("Ended up with an article list that's too short to be correct")
	}

	as := NewDiskArticleStore("/tmp/articles.json", 1000)
	if err := as.Load(); err != nil {
		panic(err.Error())
	}
	as.Add(articles)
	if err := as.Save(); err != nil {
		panic(err.Error())
	}
}

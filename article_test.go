package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewArticleList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	foo := &Article{
		Title:     "foo",
		URL:       urlToJSONURL(mustURLParse("http://example.com/foo")),
		FirstSeen: mustTimeParse("2006-01-02T15:04:05Z00:00", "2021-10-15T00:00:00Z00:00"),
	}
	bar := &Article{
		Title:     "bar",
		URL:       urlToJSONURL(mustURLParse("http://example.com/bar")),
		FirstSeen: mustTimeParse("2006-01-02T15:04:05Z00:00", "2021-10-14T00:00:00Z00:00"),
	}

	list := NewArticleList([]*Article{bar, foo}, 0)
	assert.Equal(2, len(list.Articles))
	// make sure the list got sorted correctly (newest first)
	assert.Equal(foo, list.Articles[0])
	assert.Equal(bar, list.Articles[1])
}

func Test_NewArticleList_Limit(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	articles := make([]*Article, 10)
	t0 := mustTimeParse("2006-01-02T15:04:05Z00:00", "2021-10-14T00:00:00Z00:00")
	for i := 0; i < len(articles); i++ {
		articles[i] = &Article{
			Title:     "blah",
			URL:       new(JSONURL),
			FirstSeen: t0.Add(time.Duration(i*24) * time.Hour),
		}
	}

	list := NewArticleList(articles, 3)
	assert.Equal(3, len(list.Articles))

	// make sure trimming happens after sorting
	refTime := t0.Add(time.Duration(6*24) * time.Hour)
	for i := range list.Articles {
		assert.True(list.Articles[i].FirstSeen.After(refTime))
	}
}

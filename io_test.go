package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_JSONURL_MarshalJSON(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := "http://example.com/foo"
	u := mustURLParse(s)
	ju := urlToJSONURL(u)

	b, err := json.Marshal(ju)
	assert.Nil(err)
	assert.Equal("\""+s+"\"", string(b))
}

func Test_JSONURL_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := "http://example.com/foo"
	ju := new(JSONURL)
	u := mustURLParse(s)

	err := json.Unmarshal([]byte("\""+s+"\""), ju)
	assert.Nil(err)

	rslt := jsonURLToURL(ju)
	assert.Nil(err)
	assert.Equal(u, rslt)
}

func Test_DiskArticleStore(t *testing.T) {
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
	articles := []*Article{foo, bar}

	f, err := ioutil.TempFile("", "journal-scraper-test-")
	assert.Nil(err)
	f.Close()

	as := NewDiskArticleStore(f.Name(), 3)
	// This initial Load() is called on a nonexistent file
	assert.Nil(as.Load())
	as.Add(articles)
	assert.Nil(as.Save())

	// Now read back what we wrote
	as.Load()
	das, ok := as.(*DiskArticleStore)
	assert.True(ok)
	list := das.articleList
	assert.Equal(2, len(list.Articles))
	for i := range list.Articles {
		assert.Equal(articles[i], list.Articles[i])
	}
}

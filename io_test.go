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

func Test_WriteBlob(t *testing.T) {
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
	err = writeBlob(f.Name(), articles, 0)
	assert.Nil(err)

	b, err := ioutil.ReadFile(f.Name())
	assert.Nil(err)

	list := new(ArticleList)
	err = json.Unmarshal(b, list)
	assert.Nil(err)
	assert.Equal(2, len(list.Articles))
	for i := range list.Articles {
		assert.Equal(articles[i], list.Articles[i])
	}
}

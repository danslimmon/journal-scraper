package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
)

// dummyHandler implements http.Handler, always returning Response with 200 status.
type dummyHandler struct {
	Response []byte
}

func (h *dummyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write(h.Response)
}

func Test_ElementToArticle(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	h := &dummyHandler{
		[]byte(`<a href="/blah/blah"><h3>Blah Blah</h2></a>`),
	}
	server := httptest.NewServer(h)
	defer server.Close()
	client := server.Client()

	var err error
	var article *Article
	c := colly.NewCollector()
	c.SetClient(client)
	c.OnHTML("a", func(e *colly.HTMLElement) {
		article, err = elementToArticle(e, "http://example.com")
	})

	visitErr := c.Visit(server.URL)
	assert.Nil(visitErr)

	assert.Nil(err)
	assert.NotNil(article)
	if article == nil {
		return
	}
	assert.Equal("http://example.com/blah/blah", article.URL.String())
	assert.Equal("Blah Blah", article.Title)
}

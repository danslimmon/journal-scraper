package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
)

const (
	MaxArticles = 1000
)

// Custom JSONURL type to make *url.URL marshal/unmarshal as a string instead of an object
type JSONURL url.URL

func (ju *JSONURL) MarshalJSON() ([]byte, error) {
	u := jsonURLToURL(ju)
	return []byte(fmt.Sprintf("\"%s\"", u.String())), nil
}

func (ju *JSONURL) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	*ju = JSONURL(*u)
	return nil
}

func urlToJSONURL(u *url.URL) *JSONURL {
	ju := JSONURL(*u)
	return &ju
}

func jsonURLToURL(ju *JSONURL) *url.URL {
	u := url.URL(*ju)
	return &u
}

// ArticleStore is the interface through which we store articles that have been scraped.
//
// Code using ArticleStore should call Load(), then Add(), then Save().
type ArticleStore interface {
	// Load prepares the ArticleStore by retrieving its current contents from persistent storage.
	Load() error
	// Add adds the given articles to the ArticleStore.
	Add([]*Article)
	// Save writes the ArticleStore's contents to persistent storage.
	Save() error
}

// DiskArticleStore is an ArticleStore implementation that uses files on disk.
//
// It's used for testing and development.
type DiskArticleStore struct {
	articleList *ArticleList

	FilePath string
	Limit    int
}

// Load unmarshals the contents of the target file into the DiskArticleStore.
//
// If the file is nonexistent or empty, we end up with an empty list.
func (as *DiskArticleStore) Load() error {
	as.articleList = NewArticleList(MaxArticles)

	if _, err := os.Stat(as.FilePath); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	b, err := ioutil.ReadFile(as.FilePath)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}

	err = json.Unmarshal(b, as.articleList)
	return err
}

func (as *DiskArticleStore) Add(articles []*Article) {
	as.articleList.Merge(articles)
}

func (as *DiskArticleStore) Save() error {
	b, err := json.Marshal(as.articleList)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(as.FilePath, b, 0600)
}

func NewDiskArticleStore(filePath string, limit int) ArticleStore {
	return &DiskArticleStore{FilePath: filePath, Limit: limit}
}

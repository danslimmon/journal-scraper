package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
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

// writeBlob writes the article list to the given file.
func writeBlob(filePath string, articles []*Article, limit int) error {
	list := NewArticleList(articles, limit)
	b, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, b, 0600)
}

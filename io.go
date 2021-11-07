package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
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

// S3ArticleStore is an ArticleStore implementation that uses S3.
type S3ArticleStore struct {
	articleList *ArticleList
	s3          s3iface.S3API

	Bucket string
	Key    string
	Limit  int
}

// Load unmarshals the contents of the target S3 object into the S3ArticleStore.
//
// If the object is nonexistent or empty, we end up with an empty list.
func (as *S3ArticleStore) Load() error {
	as.articleList = NewArticleList(MaxArticles)

	rslt, err := as.s3.GetObject(&s3.GetObjectInput{Bucket: &as.Bucket, Key: &as.Key})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == s3.ErrCodeNoSuchKey {
				return nil
			}
			return err
		} else {
			return err
		}
	}

	b, err := ioutil.ReadAll(rslt.Body)
	rslt.Body.Close()
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}

	err = json.Unmarshal(b, as.articleList)
	return err
}

func (as *S3ArticleStore) Add(articles []*Article) {
	as.articleList.Merge(articles)
}

func (as *S3ArticleStore) Save() error {
	b, err := json.Marshal(as.articleList)
	if err != nil {
		return err
	}

	// Write body to temp file because file implements io.ReadCloser
	f, err := ioutil.TempFile("", "journal-scraper-")
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		panic("test failed to write to temp file: " + err.Error())
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		panic("test failed to seek to beginning of file: " + err.Error())
	}

	_, err = as.s3.PutObject(&s3.PutObjectInput{
		Body:   f,
		Bucket: aws.String(as.Bucket),
		Key:    aws.String(as.Key),
	})
	return err
}

func NewS3ArticleStore(bucket, key string, limit int) ArticleStore {
	sess := session.Must(session.NewSession())

	return &S3ArticleStore{
		s3:     s3.New(sess),
		Bucket: bucket,
		Key:    key,
		Limit:  limit,
	}
}

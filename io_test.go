package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
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

type testS3ArticleStore_s3 struct {
	s3iface.S3API

	body io.ReadCloser

	getObjectCalls []*s3.GetObjectInput
	putObjectCalls []*s3.PutObjectInput
}

func (s *testS3ArticleStore_s3) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	s.getObjectCalls = append(s.getObjectCalls, input)
	if s.body == nil {
		return nil, awserr.New(s3.ErrCodeNoSuchKey, "No such key", nil)
	}

	return &s3.GetObjectOutput{Body: s.body}, nil
}

func (s *testS3ArticleStore_s3) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	s.putObjectCalls = append(s.putObjectCalls, input)

	f, err := ioutil.TempFile("", "journal-scraper-")
	if err != nil {
		panic("test failed to create temp file: " + err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(input.Body)
	if err != nil {
		panic("test failed to read S3 object body: " + err.Error())
	}

	_, err = f.Write(bodyBytes)
	if err != nil {
		panic("test failed to write to temp file: " + err.Error())
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		panic("test failed to seek to beginning of file: " + err.Error())
	}

	s.body = f
	return &s3.PutObjectOutput{}, nil
}

func (s *testS3ArticleStore_s3) Cleanup() {
	if s.body == nil {
		return
	}
	if f, ok := s.body.(*os.File); ok {
		os.Remove(f.Name())
	}
}

func Test_S3ArticleStore(t *testing.T) {
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

	s3mock := &testS3ArticleStore_s3{
		getObjectCalls: make([]*s3.GetObjectInput, 0),
		putObjectCalls: make([]*s3.PutObjectInput, 0),
	}
	defer s3mock.Cleanup()
	sas := &S3ArticleStore{
		Bucket: "foo_bucket",
		Key:    "article_store.json",
		Limit:  3,
		s3:     s3mock,
	}
	as := ArticleStore(sas)

	// This initial Load() is called on a nonexistent object
	assert.Nil(as.Load())
	as.Add(articles)
	assert.Nil(as.Save())

	// Now read back what we wrote
	as.Load()
	list := sas.articleList
	assert.Equal(2, len(list.Articles))
	for i := range list.Articles {
		assert.Equal(articles[i], list.Articles[i])
	}
}

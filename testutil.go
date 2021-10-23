package main

import (
	"net/url"
	"time"
)

func mustURLParse(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err.Error())
	}
	return u
}

func mustTimeParse(layout, s string) time.Time {
	t, err := time.Parse(layout, s)
	if err != nil {
		panic(err.Error())
	}
	return t
}

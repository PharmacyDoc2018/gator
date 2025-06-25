package main

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
)

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	}
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var fetchedRSSFeed RSSFeed

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	err = xml.Unmarshal(resData, fetchedRSSFeed)
	if err != nil {
		return nil, err
	}
}

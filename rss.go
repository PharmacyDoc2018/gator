package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/PharmacyDoc2018/gator/internal/database"
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
	} `xml:"channel"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var fetchedRSSFeed RSSFeed

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
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

	err = xml.Unmarshal(resData, &fetchedRSSFeed)
	if err != nil {
		return nil, err
	}

	fetchedRSSFeed.Channel.Title = html.UnescapeString(fetchedRSSFeed.Channel.Title)
	fetchedRSSFeed.Channel.Description = html.UnescapeString(fetchedRSSFeed.Channel.Description)
	for i := range fetchedRSSFeed.Channel.Item {
		fetchedRSSFeed.Channel.Item[i].Title = html.UnescapeString(fetchedRSSFeed.Channel.Item[i].Title)
		fetchedRSSFeed.Channel.Item[i].Description = html.UnescapeString(fetchedRSSFeed.Channel.Item[i].Description)
	}

	return &fetchedRSSFeed, nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	params := database.MarkFeedFetchedParams{
		ID:        feed.ID,
		UpdatedAt: time.Now(),
	}

	err = s.db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		return err
	}

	fetchedFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	for _, item := range fetchedFeed.Channel.Item {
		fmt.Println(item.Title)
		fmt.Println()
	}

	return nil
}

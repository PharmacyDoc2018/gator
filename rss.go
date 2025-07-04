package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/PharmacyDoc2018/gator/internal/database"
	"github.com/google/uuid"
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

var layouts = []string{
	time.RFC3339,
	"2006-01-02 15:04:05",   // MySQL DATETIME
	"2006-01-02",            // ISO Date
	"01/02/2006",            // US format
	"02 Jan 2006",           // e.g. 02 Jan 2006
	"02 Jan 2006 15:04",     // e.g. 02 Jan 2006 15:04
	"02 Jan 2006 15:04:05",  // e.g. 02 Jan 2006 15:04:05
	"Jan 2, 2006 at 3:04pm", // informal format
	"January 2, 2006",       // full month name
	"2006/01/02",            // slash-separated
	"Mon, 02 Jan 2006 15:04:05 EDT",
}

// TryParseDate attempts to parse a string using multiple layouts
func TryParseDate(value string) (time.Time, error) {
	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not parse time: unsupported format\n%s", value)
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
		pubTime, err := TryParseDate(item.PubDate)
		if err != nil {
			return err
		}
		params := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: sql.NullTime{
				Time:  pubTime,
				Valid: item.PubDate != "",
			},
			FeedID: feed.ID,
		}

		post, err := s.db.CreatePost(context.Background(), params)
		if err != nil {
			if fmt.Sprint(err) == "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
				// Ignore error. Expected duplicate.
				continue
			} else {
				return err
			}
		}

		fmt.Println("New post added:", post.Title)
	}

	return nil
}

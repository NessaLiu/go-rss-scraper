package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/NessaLiu/go-rss-scraper/internal/database"
	"github.com/google/uuid"
)

// Scraper is a long-running job - it will run in the background as our server runs
func startScraping(
	db *database.Queries, // connection to DB
	concurrency int, // # of concurrent units - how many go routines we want to do the scraping on
	timeBetweenRequest time.Duration, // how much time in between each request to go scrape a new RSS feed
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	// ticker.C is a channel where every [timeBetweenRequest], a value will be sent across the channel
	// so, we run this for loop every [timeBetweenRequest] (and starts when it hits this line)
	for ; ; <-ticker.C {
		// Grab the next batch of feeds to fetch
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(), // context.Background is like the "global context", what you use if you don't have access to a scoped context
			int32(concurrency),
		)
		if err != nil {
			log.Println("Error fetching feeds:", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done() // deferring this means it will always be called at the end of the function

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed from URL", feed.Url)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{} // If the item description is empty, we will set it to null in the database
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		// Parse the published date to get the time - RFC1123Z is a layout
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Could not parse published date %v with err %v", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: publishedAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			// The duplication error is expected since we don't want to create posts that already exist
			// We only log an error if it is another unexpected error with a post creation
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Failed to created post:", err)
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

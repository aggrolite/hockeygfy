package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	dev := flag.Bool("dev", false, "Read from developer config")
	flag.Parse()

	cfg, err := LoadConfig(dev)
	if err != nil {
		log.Fatalf("Problem loading config: %s\n", err)
	}

	r, err := NewRedditBot(cfg.Reddit)
	if err != nil {
		log.Fatalf("Problem creating new reddit bot: %s\n", err)
	}

	t := NewTwitterBot(cfg.Twitter)

	// First wave of scraping and tweeting.
	run(t, r)

	// Now tick over the same function every hour.
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		if r.Client.TokenExpiry.Before(time.Now()) {
			log.Print("Token is expired! Requesting new token...")
			err = r.Reauthorize(cfg.Reddit)
			if err != nil {
				log.Printf("Problem creating new token: %s. Will try again.\n", err)
				continue
			}
		}
		run(t, r)
	}

}

func run(t *TwitterBot, r *RedditBot) {
	log.Print("Finding new links...")
	links, err := r.FetchNewLinks()
	if err != nil {
		log.Printf("Problem fetching links: %s", err)
		return
	}
	if len(links) == 0 {
		log.Print("No new links found.")
		return
	}
	log.Printf("Found %d new links. Queueing as new tweets.\n", len(links))
	t.QueueNewTweets(links)
}

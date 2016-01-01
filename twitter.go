package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

type TwitterBot struct {
	client *anaconda.TwitterApi
	queue  chan *RedditLink
}

type TwitterConfig struct {
	AccessToken       string
	AccessTokenSecret string
	ConsumerKey       string
	ConsumerSecret    string
}

func NewTwitterBot(cfg *TwitterConfig) *TwitterBot {
	anaconda.SetConsumerKey(cfg.ConsumerKey)
	anaconda.SetConsumerSecret(cfg.ConsumerSecret)
	t := &TwitterBot{
		client: anaconda.NewTwitterApi(cfg.AccessToken, cfg.AccessTokenSecret),
		queue:  make(chan *RedditLink, 100),
	}
	t.run()
	return t
}

func (t *TwitterBot) QueueNewTweets(links []*RedditLink) {
	for _, l := range links {
		log.Printf("Queueing new link to be tweeted. URL: %s, ID: %s\n", l.URL, l.ID)
		t.queue <- l
	}
}

func (t *TwitterBot) publishTweet(link *RedditLink) (int64, error) {

	body := fmt.Sprintf("%s %s #gfycat #nhl", link.Title, link.URL)
	tweet, err := t.client.PostTweet(body, url.Values{})
	if err != nil {
		return 0, err
	}
	return tweet.Id, nil
}

func (t *TwitterBot) run() {
	go func() {
		log.Print("running twitter bot")
		ticker := time.NewTicker(30 * time.Minute)
		for range ticker.C {
			log.Print("fetching link from queue")
			link := <-t.queue

			log.Print("Publishing new tweet.")
			tweet, err := t.publishTweet(link)
			if err != nil {
				log.Printf("Problem publishing tweet: %s\n", err)
				continue
			}
			log.Printf("New tweet published! ID=%d\n", tweet)
		}
	}()
}

package main

import (
	"database/sql"
	"github.com/jzelinskie/geddit"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type RedditBot struct {
	Client *geddit.OAuthSession
}

type RedditConfig struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

// For debugging, it may be nice to track Reddit IDs for future issues.
type RedditLink struct {
	Author string
	ID     string
	Title  string
	URL    string
}

func NewRedditBot(cfg *RedditConfig) (*RedditBot, error) {
	o, err := geddit.NewOAuthSession(
		cfg.ClientID,
		cfg.ClientSecret,
		"hockeygfy.com bot by u/aggrolite v0.1 see source https://github.com/jzelinskie/geddit",
		"http://hockeygfy.com",
	)
	if err != nil {
		return nil, err
	}

	err = o.LoginAuth(cfg.Username, cfg.Password)
	if err != nil {
		return nil, err
	}

	r := &RedditBot{o}
	return r, nil
}

func (r *RedditBot) Reauthorize(cfg *RedditConfig) error {
	return r.Client.LoginAuth(cfg.Username, cfg.Password)
}

func (r *RedditBot) linkExists(db *sql.DB, id string) (bool, error) {
	stmt, err := db.Prepare("select count(1) from links where reddit_id = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(id).Scan(&count)
	if err != nil {
		return false, nil
	}
	log.Printf("count=%v\n", count)
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (r *RedditBot) updateLinks(db *sql.DB, l *geddit.Submission) error {

	stmt, err := db.Prepare("insert into links(reddit_id, title, url) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(l.FullID, l.Title, l.URL)
	if err != nil {
		return err
	}

	return nil
}

func (r *RedditBot) FetchNewLinks() ([]*RedditLink, error) {
	opts := geddit.ListingOptions{
		Limit: 100,
	}
	links, err := r.Client.SubredditSubmissions("hockey", geddit.HotSubmissions, opts)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", "db/hockeygfy.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var newLinks []*RedditLink
	for _, l := range links {
		if l.Domain != "gfycat.com" {
			continue
		}
		exists, err := r.linkExists(db, l.FullID)
		if err != nil {
			log.Printf("Problem checking for link in DB: %v\n", err)
			continue
		}
		if exists {
			log.Println("Link exists in DB. Skipping.")
			continue
		}
		err = r.updateLinks(db, l)
		if err != nil {
			log.Println("Problem updating DB: %v", err)
			continue
		}
		r := &RedditLink{
			Author: l.Author,
			ID:     l.FullID,
			Title:  l.Title,
			URL:    l.URL,
		}
		newLinks = append(newLinks, r)
	}
	return newLinks, nil
}

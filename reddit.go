package main

import (
	"github.com/golang/groupcache/lru"
	"github.com/jzelinskie/geddit"
)

type RedditBot struct {
	Cache  *lru.Cache
	client *geddit.OAuthSession
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

	r := &RedditBot{
		Cache:  lru.New(500),
		client: o,
	}
	return r, nil
}

func (r *RedditBot) FetchNewLinks() ([]*RedditLink, error) {
	opts := geddit.ListingOptions{
		Limit: 100,
	}
	links, err := r.client.SubredditSubmissions("hockey", geddit.HotSubmissions, opts)
	if err != nil {
		return nil, err
	}

	var newLinks []*RedditLink
	for _, l := range links {
		if l.Domain != "gfycat.com" {
			continue
		}
		_, ok := r.Cache.Get(l.FullID)
		if ok {
			continue
		}
		r.Cache.Add(l.FullID, 1)
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

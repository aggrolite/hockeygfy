package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Reddit  *RedditConfig
	Twitter *TwitterConfig
}

func LoadConfig(dev *bool) (*Config, error) {
	path := "config/prod.json"
	if *dev {
		path = "config/dev.json"
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(file)
	c := &Config{}
	err = d.Decode(c)
	return c, err
}

package main

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Master MasterConfig
	Server ServerConfig
}

type MasterConfig struct {
	AuthorName      string
	AuthorEmail     string
	OriginURL       string
	TargetDir       string
	FeedTitle       string
	FeedDescription string
}

type ServerConfig struct {
	Timeout int
}

func GetConfig() Config {
	var config Config

	// defaults when ENV is empty
	config.Master.AuthorName = "username"
	config.Master.OriginURL = "https://rss.example.com"
	config.Master.TargetDir = "~/rss-notes"
	config.Master.FeedTitle = ""
	config.Master.FeedDescription = ""
	config.Master.AuthorEmail = ""
	config.Server.Timeout = 10

	applyEnvOverrides(&config)

	// feed title/description fallbacks
	if config.Master.FeedTitle == "" {
		config.Master.FeedTitle = config.Master.AuthorName + " notes"
	}
	if config.Master.FeedDescription == "" {
		config.Master.FeedDescription = "Notes by " + config.Master.AuthorName
	}

	return config
}

func applyEnvOverrides(config *Config) {
	if v := os.Getenv("RSS_AUTHOR_NAME"); v != "" {
		config.Master.AuthorName = v
	}
	if v := os.Getenv("RSS_AUTHOR_EMAIL"); v != "" {
		config.Master.AuthorEmail = v
	}
	if v := os.Getenv("RSS_ORIGIN_URL"); v != "" {
		config.Master.OriginURL = v
	}
	if v := os.Getenv("RSS_TARGET_DIR"); v != "" {
		config.Master.TargetDir = v
	}
	if v := os.Getenv("RSS_FEED_TITLE"); v != "" {
		config.Master.FeedTitle = v
	}
	if v := os.Getenv("RSS_FEED_DESCRIPTION"); v != "" {
		config.Master.FeedDescription = v
	}
	if v := os.Getenv("RSS_TIMEOUT"); v != "" {
		timeout, err := strconv.Atoi(v)
		if err != nil || timeout < 0 {
			log.Fatalf("invalid RSS_TIMEOUT %q: must be >= 0", v)
		}
		config.Server.Timeout = timeout
	}
}

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
	AuthorUsername string
	OriginURL      string
	TargetDir      string
}

type ServerConfig struct {
	Host    string
	Timeout int
}

func GetConfig() Config {
	var config Config

	// Defaults when ENV is empty
	config.Master.AuthorUsername = "username"
	config.Master.OriginURL = "https://rss.example.com"
	config.Master.TargetDir = "~/rss-notes"
	config.Server.Host = "127.0.0.1"
	config.Server.Timeout = 10

	applyEnvOverrides(&config)

	return config
}

func applyEnvOverrides(config *Config) {
	if v := os.Getenv("RSS_AUTHOR_USERNAME"); v != "" {
		config.Master.AuthorUsername = v
	}
	if v := os.Getenv("RSS_ORIGIN_URL"); v != "" {
		config.Master.OriginURL = v
	}
	if v := os.Getenv("RSS_TARGET_DIR"); v != "" {
		config.Master.TargetDir = v
	}
	if v := os.Getenv("RSS_HOST"); v != "" {
		config.Server.Host = v
	}
	if v := os.Getenv("RSS_TIMEOUT"); v != "" {
		timeout, err := strconv.Atoi(v)
		if err != nil || timeout < 0 {
			log.Fatalf("invalid RSS_TIMEOUT %q: must be >= 0", v)
		}
		config.Server.Timeout = timeout
	}
}

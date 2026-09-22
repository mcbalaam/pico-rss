package main

import (
	"log"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Master MasterConfig `toml:"master"`
	Server ServerConfig `toml:"server"`
}

type MasterConfig struct {
	AuthorUsername string `toml:"author_username"`
	OriginURL      string `toml:"origin_url"`
	TargetDir      string `toml:"target_dir"`
}

type ServerConfig struct {
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	Timeout int    `toml:"timeout"`
}

func GetConfig() Config {
	var config Config

	// Defaults for local run and for pure-ENV (docker) mode.
	config.Master.AuthorUsername = "mcbalaam"
	config.Master.OriginURL = "https://rss.mcblm.xyz"
	config.Master.TargetDir = "~/rss-notes"
	config.Server.Host = "127.0.0.1"
	config.Server.Port = 8066
	config.Server.Timeout = 10

	// CONFIG_PATH lets docker/compose point at a mounted file.
	// Default keeps backward compat with ./config.toml.
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.toml"
	}

	if _, err := os.Stat(configPath); err == nil {
		if _, err := toml.DecodeFile(configPath, &config); err != nil {
			log.Fatal(err)
		}
	} else if configPath != "config.toml" || !os.IsNotExist(err) {
		// Explicit path that is missing/unreadable -> fail loud.
		// Default path missing -> fine, fall through to defaults+ENV.
		if os.Getenv("CONFIG_PATH") != "" {
			log.Fatalf("config file %q: %v", configPath, err)
		}
	}

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
	if v := os.Getenv("RSS_PORT"); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil || port <= 0 || port > 65535 {
			log.Fatalf("invalid RSS_PORT %q: must be 1-65535", v)
		}
		config.Server.Port = port
	}
	if v := os.Getenv("RSS_TIMEOUT"); v != "" {
		timeout, err := strconv.Atoi(v)
		if err != nil || timeout < 0 {
			log.Fatalf("invalid RSS_TIMEOUT %q: must be >= 0", v)
		}
		config.Server.Timeout = timeout
	}
}

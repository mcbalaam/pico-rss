package main

import (
	"log"

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

	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

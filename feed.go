package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/feeds"
)

func CollectNotesFromDir(dir string, config *Config) ([]*feeds.Item, error) {
	expanded, err := expandHome(dir)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(expanded)
	if err != nil {
		return nil, err
	}

	var items []*feeds.Item
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(e.Name()), ".md") {
			continue
		}
		item, err := GetNote(filepath.Join(expanded, e.Name()), config)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Created.After(items[j].Created)
	})

	return items, nil
}

func BuildFeed(config *Config, items []*feeds.Item) *feeds.Feed {
	origin := strings.TrimRight(config.Master.OriginURL, "/")

	author := &feeds.Author{Name: config.Master.AuthorName, Email: config.Master.AuthorEmail}
	if author.Name == "" && author.Email == "" {
		author = nil
	}

	feed := &feeds.Feed{
		Title:       config.Master.FeedTitle,
		Link:        &feeds.Link{Href: origin},
		Description: config.Master.FeedDescription,
		Author:      author,
		Created:     time.Now(),
	}

	if len(items) > 0 {
		feed.Items = items
		latest := items[0].Created
		for _, it := range items {
			if it.Created.After(latest) {
				latest = it.Created
			}
			if it.Updated.After(latest) {
				latest = it.Updated
			}
		}
		feed.Created = latest
		feed.Updated = latest
	}

	return feed
}

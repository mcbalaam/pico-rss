package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Link struct {
	Href, Rel, Type, Length string
}

type Author struct {
	Name, Email string
}

type Item struct {
	Title       string
	Link        *Link
	Description string
	Content     string
	Created     time.Time
	Updated     time.Time
	Id          string
	Author      *Author
}

type Feed struct {
	Title       string
	Link        *Link
	Description string
	Author      *Author
	Created     time.Time
	Updated     time.Time
	Items       []*Item
}

func CollectNotesFromDir(dir string, config *Config) ([]*Item, error) {
	expanded, err := expandHome(dir)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(expanded)
	if err != nil {
		return nil, err
	}

	var items []*Item
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

func BuildFeed(config *Config, items []*Item) *Feed {
	origin := strings.TrimRight(config.Master.OriginURL, "/")

	var author *Author
	if config.Master.AuthorName != "" || config.Master.AuthorEmail != "" {
		author = &Author{Name: config.Master.AuthorName, Email: config.Master.AuthorEmail}
	}

	feed := &Feed{
		Title:       config.Master.FeedTitle,
		Link:        &Link{Href: origin},
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

func anyTimeFormat(format string, times ...time.Time) string {
	for _, t := range times {
		if !t.IsZero() {
			return t.Format(format)
		}
	}
	return ""
}

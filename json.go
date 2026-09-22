package main

import (
	"encoding/json"
	"time"
)

const jsonFeedVersion = "https://jsonfeed.org/version/1.1"

type JSONAuthor struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type JSONItem struct {
	Id            string      `json:"id"`
	Url           string      `json:"url,omitempty"`
	Title         string      `json:"title,omitempty"`
	ContentHTML   string      `json:"content_html,omitempty"`
	Summary       string      `json:"summary,omitempty"`
	PublishedDate *time.Time  `json:"date_published,omitempty"`
	ModifiedDate  *time.Time  `json:"date_modified,omitempty"`
	Author        *JSONAuthor `json:"author,omitempty"`
	Authors       []*JSONAuthor `json:"authors,omitempty"`
}

type JSONFeed struct {
	Version     string      `json:"version"`
	Title       string      `json:"title"`
	HomePageUrl string      `json:"home_page_url,omitempty"`
	Description string      `json:"description,omitempty"`
	Author      *JSONAuthor   `json:"author,omitempty"`
	Authors     []*JSONAuthor `json:"authors,omitempty"`
	Items       []*JSONItem   `json:"items,omitempty"`
}

type JSON struct {
	*Feed
}

func (f *JSON) ToJSON() (string, error) {
	return f.JSONFeed().ToJSON()
}

func (f *JSONFeed) ToJSON() (string, error) {
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (f *JSON) JSONFeed() *JSONFeed {
	feed := &JSONFeed{
		Version:     jsonFeedVersion,
		Title:       f.Title,
		Description: f.Description,
	}
	if f.Link != nil {
		feed.HomePageUrl = f.Link.Href
	}
	if f.Author != nil {
		author := &JSONAuthor{Name: f.Author.Name}
		feed.Author = author
		feed.Authors = []*JSONAuthor{author}
	}
	for _, e := range f.Items {
		feed.Items = append(feed.Items, newJSONItem(e))
	}
	return feed
}

func newJSONItem(i *Item) *JSONItem {
	id := i.Id
	if id == "" && i.Link != nil {
		id = i.Link.Href
	}
	item := &JSONItem{
		Id:          id,
		Title:       i.Title,
		Summary:     i.Description,
		ContentHTML: i.Content,
	}
	if i.Link != nil {
		item.Url = i.Link.Href
	}
	if i.Author != nil && (i.Author.Name != "" || i.Author.Email != "") {
		author := &JSONAuthor{Name: i.Author.Name}
		item.Author = author
		item.Authors = []*JSONAuthor{author}
	}
	if !i.Created.IsZero() {
		t := i.Created
		item.PublishedDate = &t
	}
	if !i.Updated.IsZero() {
		t := i.Updated
		item.ModifiedDate = &t
	}
	return item
}

func (f *Feed) ToJSON() (string, error) {
	j := &JSON{f}
	return j.ToJSON()
}

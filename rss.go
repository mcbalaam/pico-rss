package main

import (
	"encoding/xml"
	"fmt"
	"time"
)

type RssFeedXml struct {
	XMLName          xml.Name `xml:"rss"`
	Version          string   `xml:"version,attr"`
	ContentNamespace string   `xml:"xmlns:content,attr"`
	Channel          *RssFeed
}

type RssContent struct {
	XMLName xml.Name `xml:"content:encoded"`
	Content string   `xml:",cdata"`
}

type RssImage struct {
	XMLName xml.Name `xml:"image"`
	Url     string   `xml:"url"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Width   int      `xml:"width,omitempty"`
	Height  int      `xml:"height,omitempty"`
}

type RssFeed struct {
	XMLName        xml.Name `xml:"channel"`
	Title          string   `xml:"title"`
	Link           string   `xml:"link"`
	Description    string   `xml:"description"`
	Language       string   `xml:"language,omitempty"`
	Copyright      string   `xml:"copyright,omitempty"`
	ManagingEditor string   `xml:"managingEditor,omitempty"`
	PubDate        string   `xml:"pubDate,omitempty"`
	LastBuildDate  string   `xml:"lastBuildDate,omitempty"`
	Image          *RssImage
	Items          []*RssItem `xml:"item"`
}

type RssItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Content     *RssContent
	Author      string `xml:"author,omitempty"`
	Guid        *RssGuid
	PubDate     string `xml:"pubDate,omitempty"`
}

type RssGuid struct {
	XMLName     xml.Name `xml:"guid"`
	Id          string   `xml:",chardata"`
	IsPermaLink string   `xml:"isPermaLink,attr,omitempty"`
}

type Rss struct {
	*Feed
}

func newRssItem(i *Item) *RssItem {
	item := &RssItem{
		Title:       i.Title,
		Description: i.Description,
		PubDate:     anyTimeFormat(time.RFC1123Z, i.Created, i.Updated),
	}
	if i.Id != "" {
		item.Guid = &RssGuid{Id: i.Id}
	} else if i.Link != nil && i.Link.Href != "" {
		item.Guid = &RssGuid{Id: i.Link.Href, IsPermaLink: "true"}
	}
	if i.Link != nil {
		item.Link = i.Link.Href
	}
	if len(i.Content) > 0 {
		item.Content = &RssContent{Content: i.Content}
	}
	if i.Author != nil && i.Author.Name != "" {
		item.Author = i.Author.Name
	}
	return item
}

func (r *Rss) RssFeed() *RssFeed {
	pub := anyTimeFormat(time.RFC1123Z, r.Created, r.Updated)
	build := anyTimeFormat(time.RFC1123Z, r.Updated)
	author := ""
	if r.Author != nil {
		author = r.Author.Email
		if len(r.Author.Name) > 0 {
			if author != "" {
				author = fmt.Sprintf("%s (%s)", author, r.Author.Name)
			} else {
				author = r.Author.Name
			}
		}
	}
	var href string
	if r.Link != nil {
		href = r.Link.Href
	}
	channel := &RssFeed{
		Title:          r.Title,
		Link:           href,
		Description:    r.Description,
		ManagingEditor: author,
		PubDate:        pub,
		LastBuildDate:  build,
	}
	for _, i := range r.Items {
		channel.Items = append(channel.Items, newRssItem(i))
	}
	return channel
}

func (r *Rss) FeedXml() interface{} {
	return r.RssFeed().FeedXml()
}

func (r *RssFeed) FeedXml() interface{} {
	return &RssFeedXml{
		Version:          "2.0",
		Channel:          r,
		ContentNamespace: "http://purl.org/rss/1.0/modules/content/",
	}
}

func (f *Feed) ToRss() (string, error) {
	r := &Rss{f}
	return ToXML(r)
}

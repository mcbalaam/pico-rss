package main

import (
	"crypto/rand"
	"encoding/xml"
	"fmt"
	"net/url"
	"time"
)

const atomNS = "http://www.w3.org/2005/Atom"

type AtomPerson struct {
	Name  string `xml:"name,omitempty"`
	Uri   string `xml:"uri,omitempty"`
	Email string `xml:"email,omitempty"`
}

type AtomSummary struct {
	XMLName xml.Name `xml:"summary"`
	Content string   `xml:",chardata"`
	Type    string   `xml:"type,attr"`
}

type AtomContent struct {
	XMLName xml.Name `xml:"content"`
	Content string   `xml:",chardata"`
	Type    string   `xml:"type,attr"`
}

type AtomAuthor struct {
	XMLName xml.Name `xml:"author"`
	AtomPerson
}

type AtomEntry struct {
	XMLName   xml.Name `xml:"entry"`
	Title     string   `xml:"title"`
	Updated   string   `xml:"updated"`
	Id        string   `xml:"id"`
	Content   *AtomContent
	Links     []AtomLink   `xml:"link"`
	Summary   *AtomSummary `xml:"summary,omitempty"`
	Author    *AtomAuthor  `xml:"author,omitempty"`
	Published string       `xml:"published,omitempty"`
}

type AtomLink struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Length  string   `xml:"length,attr,omitempty"`
}

type AtomFeed struct {
	XMLName  xml.Name `xml:"feed"`
	Xmlns    string   `xml:"xmlns,attr"`
	Title    string   `xml:"title"`
	Id       string   `xml:"id"`
	Updated  string   `xml:"updated"`
	Subtitle string   `xml:"subtitle,omitempty"`
	Link     *AtomLink
	Author   *AtomAuthor  `xml:"author,omitempty"`
	Entries  []*AtomEntry `xml:"entry"`
}

type Atom struct {
	*Feed
}

type UUID [16]byte

func NewUUID() *UUID {
	u := &UUID{}
	_, err := rand.Read(u[:16])
	if err != nil {
		panic(err)
	}
	u[8] = (u[8] | 0x80) & 0xBf
	u[6] = (u[6] | 0x40) & 0x4f
	return u
}

func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func newAtomEntry(i *Item) *AtomEntry {
	id := i.Id
	link := i.Link
	if link == nil {
		link = &Link{}
	}
	if len(id) == 0 {
		if len(link.Href) > 0 && (!i.Created.IsZero() || !i.Updated.IsZero()) {
			dateStr := anyTimeFormat("2006-01-02", i.Updated, i.Created)
			host, path := link.Href, "/invalid.html"
			if u, err := url.Parse(link.Href); err == nil {
				host, path = u.Host, u.Path
			}
			id = fmt.Sprintf("tag:%s,%s:%s", host, dateStr, path)
		} else {
			id = "urn:uuid:" + NewUUID().String()
		}
	}
	var name, email string
	if i.Author != nil {
		name, email = i.Author.Name, i.Author.Email
	}
	linkRel := link.Rel
	if linkRel == "" {
		linkRel = "alternate"
	}
	x := &AtomEntry{
		Title:   i.Title,
		Links:   []AtomLink{{Href: link.Href, Rel: linkRel, Type: link.Type}},
		Id:      id,
		Updated: anyTimeFormat(time.RFC3339, i.Updated, i.Created),
	}
	if !i.Created.IsZero() {
		x.Published = anyTimeFormat(time.RFC3339, i.Created)
	}
	if len(i.Description) > 0 {
		x.Summary = &AtomSummary{Content: i.Description, Type: "html"}
	}
	if len(i.Content) > 0 {
		x.Content = &AtomContent{Content: i.Content, Type: "html"}
	}
	if len(name) > 0 || len(email) > 0 {
		x.Author = &AtomAuthor{AtomPerson: AtomPerson{Name: name, Email: email}}
	}
	return x
}

func (a *Atom) AtomFeed() *AtomFeed {
	updated := anyTimeFormat(time.RFC3339, a.Updated, a.Created)
	link := a.Link
	if link == nil {
		link = &Link{}
	}
	// ensure updated is not empty
	if updated == "" {
		updated = time.Now().Format(time.RFC3339)
	}
	feed := &AtomFeed{
		Xmlns:    atomNS,
		Title:    a.Title,
		Link:     &AtomLink{Href: link.Href, Rel: link.Rel},
		Subtitle: a.Description,
		Id:       link.Href,
		Updated:  updated,
	}
	if feed.Id == "" {
		feed.Id = "urn:uuid:" + NewUUID().String()
	}
	if a.Author != nil {
		feed.Author = &AtomAuthor{AtomPerson: AtomPerson{Name: a.Author.Name, Email: a.Author.Email}}
	}
	for _, e := range a.Items {
		feed.Entries = append(feed.Entries, newAtomEntry(e))
	}
	return feed
}

func (a *Atom) FeedXml() interface{} {
	return a.AtomFeed()
}

func (a *AtomFeed) FeedXml() interface{} {
	return a
}

func (f *Feed) ToAtom() (string, error) {
	a := &Atom{f}
	return ToXML(a)
}

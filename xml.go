package main

import (
	"encoding/xml"
	"io"
)

type XmlFeed interface {
	FeedXml() interface{}
}

func ToXML(feed XmlFeed) (string, error) {
	x := feed.FeedXml()
	data, err := xml.MarshalIndent(x, "", "  ")
	if err != nil {
		return "", err
	}
	s := xml.Header[:len(xml.Header)-1] + string(data)
	return s, nil
}

func WriteXML(feed XmlFeed, w io.Writer) error {
	x := feed.FeedXml()
	if _, err := w.Write([]byte(xml.Header[:len(xml.Header)-1])); err != nil {
		return err
	}
	e := xml.NewEncoder(w)
	e.Indent("", "  ")
	return e.Encode(x)
}

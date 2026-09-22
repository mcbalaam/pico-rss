package main

import (
	"fmt"
	"net/http"
)

func writeFeed(w http.ResponseWriter, contentType, body string) {
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	_, _ = w.Write([]byte(body))
}

func registerRoutes(mux *http.ServeMux, config *Config) {
	mux.HandleFunc("/rss.xml", func(w http.ResponseWriter, r *http.Request) {
		items, err := CollectNotesFromDir(config.Master.TargetDir, config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rss, err := BuildFeed(config, items).ToRss()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeFeed(w, "application/rss+xml", rss)
	})

	mux.HandleFunc("/atom.xml", func(w http.ResponseWriter, r *http.Request) {
		items, err := CollectNotesFromDir(config.Master.TargetDir, config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		atom, err := BuildFeed(config, items).ToAtom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeFeed(w, "application/atom+xml", atom)
	})

	mux.HandleFunc("/feed.json", func(w http.ResponseWriter, r *http.Request) {
		items, err := CollectNotesFromDir(config.Master.TargetDir, config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		js, err := BuildFeed(config, items).ToJSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeFeed(w, "application/feed+json", js)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!doctype html><html><head><meta charset="utf-8"><title>%s notes</title>`+
			`<link rel="alternate" type="application/rss+xml" title="RSS" href="/rss.xml">`+
			`<link rel="alternate" type="application/atom+xml" title="Atom" href="/atom.xml">`+
			`<link rel="alternate" type="application/feed+json" title="JSON Feed" href="/feed.json">`+
			`</head><body><h1>%s notes</h1><ul>`+
			`<li><a href="/rss.xml">RSS</a></li>`+
			`<li><a href="/atom.xml">Atom</a></li>`+
			`<li><a href="/feed.json">JSON Feed</a></li>`+
			`</ul></body></html>`,
			config.Master.AuthorUsername, config.Master.AuthorUsername)
	})
}

package main

import (
	"fmt"
	"html"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func writeFeed(w http.ResponseWriter, contentType, body string) {
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	_, _ = w.Write([]byte(body))
}

// findNotePath resolves a URL slug to a .md file in dir.
// Slug must match the filename without extension exactly.
func findNotePath(dir, slug string) (string, error) {
	expanded, err := expandHome(dir)
	if err != nil {
		return "", err
	}
	entries, err := os.ReadDir(expanded)
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".md") {
			continue
		}
		if strings.TrimSuffix(name, filepath.Ext(name)) == slug {
			return filepath.Join(expanded, name), nil
		}
	}
	return "", os.ErrNotExist
}

func servePost(w http.ResponseWriter, r *http.Request, config *Config, slug string) {
	if slug == "" || slug == "." || slug == ".." ||
		strings.Contains(slug, "/") || strings.Contains(slug, "\\") {
		http.NotFound(w, r)
		return
	}
	path, err := findNotePath(config.Master.TargetDir, slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	item, err := GetNote(path, config)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!doctype html><html><head><meta charset="utf-8"><title>%s</title>`+
		`<link rel="alternate" type="application/rss+xml" title="RSS" href="/rss.xml">`+
		`<link rel="alternate" type="application/atom+xml" title="Atom" href="/atom.xml">`+
		`<link rel="alternate" type="application/feed+json" title="JSON Feed" href="/feed.json">`+
		`</head><body><p><a href="/">← all posts</a></p>`+
		`<h1>%s</h1><p><time>%s</time></p><article>%s</article></body></html>`,
		html.EscapeString(item.Title),
		html.EscapeString(item.Title),
		html.EscapeString(item.Created.Format("02.01.2006")),
		item.Content)
}

func serveIndex(w http.ResponseWriter, r *http.Request, config *Config) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!doctype html><html><head><meta charset="utf-8"><title>%s notes</title>`+
		`<link rel="alternate" type="application/rss+xml" title="RSS" href="/rss.xml">`+
		`<link rel="alternate" type="application/atom+xml" title="Atom" href="/atom.xml">`+
		`<link rel="alternate" type="application/feed+json" title="JSON Feed" href="/feed.json">`+
		`</head><body><h1>%s notes</h1><ul>`+
		`<li><a href="/rss.xml">RSS</a></li>`+
		`<li><a href="/atom.xml">Atom</a></li>`+
		`<li><a href="/feed.json">JSON Feed</a></li>`+
		`</ul>`,
		html.EscapeString(config.Master.AuthorUsername),
		html.EscapeString(config.Master.AuthorUsername))

	items, err := CollectNotesFromDir(config.Master.TargetDir, config)
	if err != nil || len(items) == 0 {
		fmt.Fprint(w, `</body></html>`)
		return
	}
	origin := strings.TrimRight(config.Master.OriginURL, "/")
	fmt.Fprint(w, `<h2>Posts</h2><ul>`)
	for _, it := range items {
		slug := ""
		if it.Link != nil {
			slug = strings.TrimPrefix(it.Link.Href, origin+"/")
		}
		if slug == "" || strings.Contains(slug, "/") {
			continue
		}
		fmt.Fprintf(w, `<li><a href="/%s">%s</a> <time>%s</time></li>`,
			url.PathEscape(slug),
			html.EscapeString(it.Title),
			html.EscapeString(it.Created.Format("02.01.2006")))
	}
	fmt.Fprint(w, `</ul></body></html>`)
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
		if r.URL.Path == "/" {
			serveIndex(w, r, config)
			return
		}
		servePost(w, r, config, strings.TrimPrefix(r.URL.Path, "/"))
	})
}

package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

const DEFAULT_NOTE_DIR = "~/notes"

// the only supported date format is dd.mm.yyyy (e.g. 22.09.2026)
const noteDateLayout = "02.01.2006"

func expandHome(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}

func parseNoteDate(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	t, err := time.Parse(noteDateLayout, s)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func isEmptyLine(s string) bool {
	return strings.TrimSpace(s) == ""
}

// stripMarkdownInline removes a minimal set of inline markers for description preview
func stripMarkdownInline(s string) string {
	s = strings.TrimSpace(s)
	replacer := strings.NewReplacer(
		"**", "", "__", "", "~~", "",
		"`", "",
	)
	s = replacer.Replace(s)
	// drop leading heading/quote/list markers
	s = strings.TrimLeft(s, "#>*-+ \t")
	// unwrap [text](url) -> text, ![alt](url) -> alt
	for {
		open := strings.Index(s, "[")
		close := strings.Index(s, "]")
		parenOpen := strings.Index(s, "(")
		parenClose := strings.Index(s, ")")
		if open >= 0 && close > open {
			inner := s[open+1 : close]
			rest := ""
			if parenOpen == close+1 && parenClose > parenOpen {
				rest = s[parenClose+1:]
				s = s[:open] + inner + rest
			} else {
				s = s[:open] + inner + s[close+1:]
			}
			continue
		}
		break
	}
	return strings.TrimSpace(s)
}

func truncateRunes(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return strings.TrimSpace(string(runes[:max]))
}

// GetNote parses a single .md file:
// line 1: "# Title" (fallback: filename without .md)
// next non-empty line: date in dd.mm.yyyy (fallback: ModTime)
// next paragraph: description (first non-empty paragraph, ~300 runes)
// link: origin_url + "/" + slug (slug = filename without .md)
// content: full markdown minus title/date lines, rendered to HTML
func GetNote(path string, config *Config) (*Item, error) {
	expandedPath, err := expandHome(path)
	if err != nil {
		return nil, err
	}
	fileStat, err := os.Stat(expandedPath)
	if err != nil {
		return nil, err
	}
	modTime := fileStat.ModTime()

	raw, err := os.ReadFile(expandedPath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(raw), "\n")

	slug := strings.TrimSuffix(filepath.Base(expandedPath), filepath.Ext(expandedPath))
	if slug == "" {
		slug = strings.TrimSuffix(fileStat.Name(), ".md")
	}

	// 1. Title
	title := ""
	titleIdx := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if isEmptyLine(trimmed) {
			continue
		}
		if strings.HasPrefix(trimmed, "# ") {
			title = strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
			titleIdx = i
		} else if strings.HasPrefix(trimmed, "#") {
			title = strings.TrimSpace(strings.TrimLeft(trimmed, "# "))
			titleIdx = i
		}
		break
	}
	if title == "" {
		title = slug
	}

	// 2. Date: first non-empty line after title (or first line if no title)
	cursor := 0
	if titleIdx >= 0 {
		cursor = titleIdx + 1
	}
	dateIdx := -1
	dateVal := modTime
	for i := cursor; i < len(lines); i++ {
		if isEmptyLine(lines[i]) {
			continue
		}
		if t, ok := parseNoteDate(lines[i]); ok {
			dateVal = t
			dateIdx = i
			cursor = i + 1
		} else {
			cursor = i
		}
		break
	}

	// 3. Description: first non-empty paragraph from cursor
	var descLines []string
	for i := cursor; i < len(lines); i++ {
		if isEmptyLine(lines[i]) {
			if len(descLines) > 0 {
				break
			}
			continue
		}
		descLines = append(descLines, stripMarkdownInline(lines[i]))

		if i+1 < len(lines) && isEmptyLine(lines[i+1]) {
			break
		}
	}
	description := truncateRunes(strings.Join(descLines, " "), 300)

	// 4. Body for Content: everything except title line and date line
	var bodyLines []string
	for i, line := range lines {
		if i == titleIdx || i == dateIdx {
			continue
		}
		bodyLines = append(bodyLines, line)
	}
	body := strings.TrimSpace(strings.Join(bodyLines, "\n"))

	htmlContent := markdownToHTML(body)

	origin := strings.TrimRight(config.Master.OriginURL, "/")
	link := origin + "/" + slug

	updated := modTime
	if updated.Before(dateVal) {
		updated = dateVal
	}

	return &Item{
		Title:       title,
		Link:        &Link{Href: link},
		Description: description,
		Content:     htmlContent,
		Created:     dateVal,
		Updated:     updated,
	}, nil
}

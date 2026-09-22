package main

import (
	"fmt"
	"html"
	"strings"
)

// markdownToHTML converts a minimal subset of Markdown to HTML.
// supported: headings (# - ######), hr (---/***), blockquote (>),
// ul (-/*/+), ol (1.), fenced code (```), paragraphs, and inline
// `code`, **bold**/__bold__, ~~del~~, *italic*/_italic_, [text](url), ![alt](url)
func markdownToHTML(md string) string {
	md = strings.TrimSpace(md)
	if md == "" {
		return ""
	}
	lines := strings.Split(md, "\n")
	var buf strings.Builder
	inList := ""
	inCodeBlock := false

	flushList := func() {
		if inList != "" {
			buf.WriteString("</")
			buf.WriteString(inList)
			buf.WriteString(">\n")
			inList = ""
		}
	}
	var paraLines []string
	flushPara := func() {
		if len(paraLines) > 0 {
			text := strings.TrimSpace(strings.Join(paraLines, " "))
			buf.WriteString("<p>")
			buf.WriteString(renderInline(text))
			buf.WriteString("</p>\n")
			paraLines = nil
		}
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			flushPara()
			flushList()
			if !inCodeBlock {
				inCodeBlock = true
				buf.WriteString("<pre><code>")
			} else {
				inCodeBlock = false
				buf.WriteString("</code></pre>\n")
			}
			continue
		}
		if inCodeBlock {
			buf.WriteString(html.EscapeString(line))
			buf.WriteString("\n")
			continue
		}
		if trimmed == "" {
			flushPara()
			continue
		}
		// heading
		if strings.HasPrefix(trimmed, "#") {
			level := 0
			for level < len(trimmed) && trimmed[level] == '#' {
				level++
			}
			if level >= 1 && level <= 6 && len(trimmed) > level && trimmed[level] == ' ' {
				flushPara()
				flushList()
				content := strings.TrimSpace(trimmed[level:])
				buf.WriteString(fmt.Sprintf("<h%d>%s</h%d>\n", level, renderInline(content), level))
				continue
			}
		}
		// hr
		if trimmed == "---" || trimmed == "***" || trimmed == "___" {
			flushPara()
			flushList()
			buf.WriteString("<hr />\n")
			continue
		}
		// blockquote
		if strings.HasPrefix(trimmed, ">") {
			flushPara()
			flushList()
			content := strings.TrimSpace(strings.TrimPrefix(trimmed, ">"))
			var bqLines []string
			bqLines = append(bqLines, content)
			for i+1 < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i+1]), ">") {
				i++
				bqLines = append(bqLines, strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[i]), ">")))
			}
			buf.WriteString("<blockquote><p>")
			buf.WriteString(renderInline(strings.Join(bqLines, " ")))
			buf.WriteString("</p></blockquote>\n")
			continue
		}
		// unordered list
		isUL := false
		ulContent := ""
		for _, prefix := range []string{"- ", "* ", "+ "} {
			if strings.HasPrefix(trimmed, prefix) {
				isUL = true
				ulContent = strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
				break
			}
		}
		// ordered list
		isOL := false
		olContent := ""
		if !isUL {
			dot := strings.Index(trimmed, ". ")
			if dot > 0 && dot < 4 {
				numPart := trimmed[:dot]
				isNum := true
				for _, c := range numPart {
					if c < '0' || c > '9' {
						isNum = false
						break
					}
				}
				if isNum {
					isOL = true
					olContent = strings.TrimSpace(trimmed[dot+2:])
				}
			}
		}
		if isUL || isOL {
			flushPara()
			want := "ul"
			if isOL {
				want = "ol"
			}
			content := ulContent
			if isOL {
				content = olContent
			}
			if inList != want {
				flushList()
				buf.WriteString("<")
				buf.WriteString(want)
				buf.WriteString(">\n")
				inList = want
			}
			buf.WriteString("<li>")
			buf.WriteString(renderInline(content))
			buf.WriteString("</li>\n")
			continue
		}
		if inList != "" {
			flushList()
		}
		paraLines = append(paraLines, line)
	}
	flushPara()
	flushList()
	if inCodeBlock {
		buf.WriteString("</code></pre>\n")
	}
	return strings.TrimSpace(buf.String())
}

func renderInline(s string) string {
	var out strings.Builder
	i := 0
	for i < len(s) {
		// code span `code`
		if s[i] == '`' {
			end := strings.Index(s[i+1:], "`")
			if end >= 0 {
				content := s[i+1 : i+1+end]
				out.WriteString("<code>")
				out.WriteString(html.EscapeString(content))
				out.WriteString("</code>")
				i = i + 1 + end + 1
				continue
			}
		}
		// image ![alt](url)
		if strings.HasPrefix(s[i:], "![") {
			closeAlt := strings.Index(s[i+2:], "]")
			if closeAlt >= 0 {
				alt := s[i+2 : i+2+closeAlt]
				rest := s[i+2+closeAlt+1:]
				if strings.HasPrefix(rest, "(") {
					closeParen := strings.Index(rest, ")")
					if closeParen >= 0 {
						url := rest[1:closeParen]
						out.WriteString(`<img src="`)
						out.WriteString(html.EscapeString(url))
						out.WriteString(`" alt="`)
						out.WriteString(html.EscapeString(alt))
						out.WriteString(`" />`)
						i = i + 2 + closeAlt + 1 + 1 + closeParen + 1
						continue
					}
				}
			}
		}
		// link [text](url)
		if s[i] == '[' {
			closeText := strings.Index(s[i+1:], "]")
			if closeText >= 0 {
				text := s[i+1 : i+1+closeText]
				rest := s[i+1+closeText+1:]
				if strings.HasPrefix(rest, "(") {
					closeParen := strings.Index(rest, ")")
					if closeParen >= 0 {
						url := rest[1:closeParen]
						out.WriteString(`<a href="`)
						out.WriteString(html.EscapeString(url))
						out.WriteString(`">`)
						out.WriteString(renderInline(text))
						out.WriteString(`</a>`)
						i = i + 1 + closeText + 1 + 1 + closeParen + 1
						continue
					}
				}
			}
		}
		// bold **
		if strings.HasPrefix(s[i:], "**") {
			end := strings.Index(s[i+2:], "**")
			if end >= 0 {
				inner := s[i+2 : i+2+end]
				out.WriteString("<strong>")
				out.WriteString(renderInline(inner))
				out.WriteString("</strong>")
				i = i + 2 + end + 2
				continue
			}
		}
		if strings.HasPrefix(s[i:], "__") {
			end := strings.Index(s[i+2:], "__")
			if end >= 0 {
				inner := s[i+2 : i+2+end]
				out.WriteString("<strong>")
				out.WriteString(renderInline(inner))
				out.WriteString("</strong>")
				i = i + 2 + end + 2
				continue
			}
		}
		// strikethrough ~~
		if strings.HasPrefix(s[i:], "~~") {
			end := strings.Index(s[i+2:], "~~")
			if end >= 0 {
				inner := s[i+2 : i+2+end]
				out.WriteString("<del>")
				out.WriteString(renderInline(inner))
				out.WriteString("</del>")
				i = i + 2 + end + 2
				continue
			}
		}
		// italic * (single)
		if s[i] == '*' {
			// avoid ** already handled
			end := strings.Index(s[i+1:], "*")
			if end >= 0 {
				// ensure not empty and not containing newline
				inner := s[i+1 : i+1+end]
				if inner != "" {
					out.WriteString("<em>")
					out.WriteString(renderInline(inner))
					out.WriteString("</em>")
					i = i + 1 + end + 1
					continue
				}
			}
		}
		if s[i] == '_' {
			end := strings.Index(s[i+1:], "_")
			if end >= 0 && end > 0 {
				inner := s[i+1 : i+1+end]
				out.WriteString("<em>")
				out.WriteString(renderInline(inner))
				out.WriteString("</em>")
				i = i + 1 + end + 1
				continue
			}
		}
		c := s[i]
		if c == '&' {
			out.WriteString("&amp;")
		} else if c == '<' {
			out.WriteString("&lt;")
		} else if c == '>' {
			out.WriteString("&gt;")
		} else if c == '"' {
			out.WriteString("&quot;")
		} else {
			out.WriteByte(c)
		}
		i++
	}
	return out.String()
}

package core

import (
	"fmt"
	"regexp"
	"strings"
)

// WikiLink represents a parsed wiki-link from markdown content.
type WikiLink struct {
	Target      string // Full object ID: type/name-ulid
	DisplayText string // Optional display text from [[target|text]] syntax
}

// StoredWikiLink represents a wiki-link record stored in the database.
type StoredWikiLink struct {
	FromID      string
	ToID        string // Resolved full object ID (empty if broken)
	Target      string // Original DisplayID target
	DisplayText string
}

var wikiLinkPattern = regexp.MustCompile(`\[\[([^\]\|]+)(?:\|([^\]]+))?\]\]`)

// ParseWikiLinks extracts wiki-links from markdown body content.
// Duplicate targets are deduplicated, keeping the first occurrence.
func ParseWikiLinks(body string) []WikiLink {
	matches := wikiLinkPattern.FindAllStringSubmatch(body, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var links []WikiLink
	for _, m := range matches {
		target := m[1]
		if target == "" || seen[target] {
			continue
		}
		seen[target] = true
		displayText := ""
		if len(m) > 2 {
			displayText = m[2]
		}
		links = append(links, WikiLink{Target: target, DisplayText: displayText})
	}
	return links
}

// ListWikiLinks returns all wiki-links from the given object.
func (v *Vault) ListWikiLinks(objectID string) ([]StoredWikiLink, error) {
	if v.index == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	return v.index.ListWikiLinks(objectID)
}

// ListBacklinks returns all wiki-links pointing to the given object.
func (v *Vault) ListBacklinks(objectID string) ([]StoredWikiLink, error) {
	if v.index == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	return v.index.FindBacklinks(objectID)
}

// RenderWikiLinks replaces wiki-link syntax in body with plain display text.
// [[target|Display Text]] → Display Text
// [[target]] → DisplayID (target with ULID suffix stripped)
func RenderWikiLinks(body string) string {
	return RenderWikiLinksStyled(body, nil)
}

// RenderWikiLinksStyled replaces wiki-link syntax in body with styled display text.
// If style is non-nil, the display text is wrapped via style(text).
func RenderWikiLinksStyled(body string, style func(string) string) string {
	matches := wikiLinkPattern.FindAllStringSubmatchIndex(body, -1)
	if len(matches) == 0 {
		return body
	}

	var b strings.Builder
	b.Grow(len(body))
	last := 0
	for _, loc := range matches {
		b.WriteString(body[last:loc[0]])

		target := body[loc[2]:loc[3]]
		displayText := ""
		if loc[4] >= 0 {
			displayText = body[loc[4]:loc[5]]
		}
		if displayText == "" {
			displayText = StripULID(target)
		}
		if style != nil {
			displayText = style(displayText)
		}
		b.WriteString(displayText)
		last = loc[1]
	}
	b.WriteString(body[last:])
	return b.String()
}


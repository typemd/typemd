package core

import (
	"fmt"
	"regexp"
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
	if v.db == nil {
		return nil, fmt.Errorf("vault not opened")
	}

	rows, err := v.db.Query(
		"SELECT from_id, to_id, target, display_text FROM wikilinks WHERE from_id = ?",
		objectID,
	)
	if err != nil {
		return nil, fmt.Errorf("list wikilinks: %w", err)
	}
	defer rows.Close()

	var links []StoredWikiLink
	for rows.Next() {
		var l StoredWikiLink
		if err := rows.Scan(&l.FromID, &l.ToID, &l.Target, &l.DisplayText); err != nil {
			return nil, fmt.Errorf("scan wikilink: %w", err)
		}
		links = append(links, l)
	}
	return links, rows.Err()
}

// ListBacklinks returns all wiki-links pointing to the given object.
func (v *Vault) ListBacklinks(objectID string) ([]StoredWikiLink, error) {
	if v.db == nil {
		return nil, fmt.Errorf("vault not opened")
	}

	rows, err := v.db.Query(
		"SELECT from_id, to_id, target, display_text FROM wikilinks WHERE to_id = ?",
		objectID,
	)
	if err != nil {
		return nil, fmt.Errorf("list backlinks: %w", err)
	}
	defer rows.Close()

	var links []StoredWikiLink
	for rows.Next() {
		var l StoredWikiLink
		if err := rows.Scan(&l.FromID, &l.ToID, &l.Target, &l.DisplayText); err != nil {
			return nil, fmt.Errorf("scan backlink: %w", err)
		}
		links = append(links, l)
	}
	return links, rows.Err()
}

// resolveWikiLinkTarget checks if the target object ID exists in the database.
// Returns the target if it exists, empty string otherwise.
func (v *Vault) resolveWikiLinkTarget(target string) string {
	var id string
	err := v.db.QueryRow("SELECT id FROM objects WHERE id = ?", target).Scan(&id)
	if err != nil {
		return ""
	}
	return id
}

// syncWikiLinks extracts wiki-links from an object body and stores them in the DB.
func (v *Vault) syncWikiLinks(objectID, body string) error {
	// Delete existing wikilinks for this object
	if _, err := v.db.Exec("DELETE FROM wikilinks WHERE from_id = ?", objectID); err != nil {
		return fmt.Errorf("delete old wikilinks: %w", err)
	}

	links := ParseWikiLinks(body)
	for _, link := range links {
		toID := v.resolveWikiLinkTarget(link.Target)
		if _, err := v.db.Exec(
			"INSERT INTO wikilinks (from_id, to_id, target, display_text) VALUES (?, ?, ?, ?)",
			objectID, toID, link.Target, link.DisplayText,
		); err != nil {
			return fmt.Errorf("insert wikilink: %w", err)
		}
	}
	return nil
}

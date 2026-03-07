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

// queryWikiLinks is a shared helper for querying wikilinks by a specific column.
func (v *Vault) queryWikiLinks(column, objectID string) ([]StoredWikiLink, error) {
	if v.db == nil {
		return nil, fmt.Errorf("vault not opened")
	}

	query := fmt.Sprintf("SELECT from_id, to_id, target, display_text FROM wikilinks WHERE %s = ?", column)
	rows, err := v.db.Query(query, objectID)
	if err != nil {
		return nil, fmt.Errorf("query wikilinks by %s: %w", column, err)
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

// ListWikiLinks returns all wiki-links from the given object.
func (v *Vault) ListWikiLinks(objectID string) ([]StoredWikiLink, error) {
	return v.queryWikiLinks("from_id", objectID)
}

// ListBacklinks returns all wiki-links pointing to the given object.
func (v *Vault) ListBacklinks(objectID string) ([]StoredWikiLink, error) {
	return v.queryWikiLinks("to_id", objectID)
}

// syncWikiLinks extracts wiki-links from an object body and stores them in the DB.
// knownIDs is used for in-memory target resolution to avoid N+1 DB queries.
func (v *Vault) syncWikiLinks(objectID, body string, knownIDs map[string]bool) error {
	// Delete existing wikilinks for this object
	if _, err := v.db.Exec("DELETE FROM wikilinks WHERE from_id = ?", objectID); err != nil {
		return fmt.Errorf("delete old wikilinks: %w", err)
	}

	links := ParseWikiLinks(body)
	if len(links) == 0 {
		return nil
	}

	stmt, err := v.db.Prepare("INSERT INTO wikilinks (from_id, to_id, target, display_text) VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare wikilink insert: %w", err)
	}
	defer stmt.Close()

	for _, link := range links {
		toID := ""
		if knownIDs[link.Target] {
			toID = link.Target
		}
		if _, err := stmt.Exec(objectID, toID, link.Target, link.DisplayText); err != nil {
			return fmt.Errorf("insert wikilink: %w", err)
		}
	}
	return nil
}

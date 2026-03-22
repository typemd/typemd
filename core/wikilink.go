package core

import (
	"fmt"
	"regexp"
	"strings"
)

// WikiLink represents a parsed wiki-link from markdown content.
type WikiLink struct {
	Target      string // Link target: full ID (type/name-ulid), type-qualified (type/name), or bare name
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
	if v.Queries == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	return v.Queries.ListWikiLinks(objectID)
}

// ListBacklinks returns all wiki-links pointing to the given object.
func (v *Vault) ListBacklinks(objectID string) ([]StoredWikiLink, error) {
	if v.Queries == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	return v.Queries.ListBacklinks(objectID)
}

// FixWikiLinksResult holds the outcome of a FixWikiLinks operation.
type FixWikiLinksResult struct {
	Expanded           int
	UnresolvedWikiLinks []UnresolvedWikiLink
}

// FixWikiLinks walks all objects, resolves shorthand wiki-links to full IDs,
// and writes back the expanded targets to source files.
func (v *Vault) FixWikiLinks() (*FixWikiLinksResult, error) {
	objects, err := v.repo.Walk()
	if err != nil {
		return nil, fmt.Errorf("walk objects: %w", err)
	}

	ctx := buildSyncContext(objects)
	result := &FixWikiLinksResult{}
	schemaCache := make(map[string]*TypeSchema)

	for _, obj := range objects {
		links := ParseWikiLinks(obj.Body)
		if len(links) == 0 {
			continue
		}

		resolutions := make(map[string]string)
		for _, link := range links {
			res := resolveWikiLinkTarget(link.Target, obj.Type, ctx.diskIDs, ctx.nameIndex)
			if res.changed && res.resolvedID != "" {
				resolutions[link.Target] = res.resolvedID
			}
			if res.err != nil {
				result.UnresolvedWikiLinks = append(result.UnresolvedWikiLinks, UnresolvedWikiLink{
					ObjectID: obj.ID,
					Target:   link.Target,
					Reason:   resolveErrorReason(res.err),
					Matches:  resolveErrorMatches(res.err),
				})
			}
		}

		if len(resolutions) > 0 {
			newBody, count := expandWikiLinksInBody(obj.Body, resolutions)
			result.Expanded += count
			obj.Body = newBody
			schema, ok := schemaCache[obj.Type]
			if !ok {
				schema, _ = v.repo.GetSchema(obj.Type)
				schemaCache[obj.Type] = schema
			}
			keyOrder := OrderedPropKeys(obj.Properties, schema)
			if err := v.repo.Save(obj, keyOrder); err != nil {
				return nil, fmt.Errorf("write-back for %s: %w", obj.ID, err)
			}
		}
	}

	return result, nil
}

// wikiLinkResolution holds the result of resolving a single wiki-link target.
type wikiLinkResolution struct {
	resolvedID string // full object ID if resolved, empty if not
	changed    bool   // true if target was a shorthand that was resolved
	err        error  // non-nil if resolution failed (not found or ambiguous)
}

// resolveWikiLinkTarget classifies and resolves a wiki-link target string.
//
// Resolution order:
//  1. type/name-ulid (full ID) → exact match in diskIDs
//  2. type/name (type-qualified, no ULID) → resolveByName in nameIndex
//  3. name (no type prefix) → resolveByName in sourceType's nameIndex
func resolveWikiLinkTarget(target, sourceType string, diskIDs map[string]bool, nameIndex map[string]map[string][]string) wikiLinkResolution {
	parts := strings.SplitN(target, "/", 2)

	if len(parts) == 2 && parts[1] != "" {
		// Has type prefix: either full ID or type-qualified name
		if ulidSuffixPattern.MatchString(parts[1]) {
			// Full ID: exact match
			if diskIDs[target] {
				return wikiLinkResolution{resolvedID: target}
			}
			return wikiLinkResolution{} // broken link, not found
		}
		// Type-qualified name: resolve by name within the specified type
		fullID, err := resolveByName(nameIndex, parts[0], parts[1])
		if err != nil {
			return wikiLinkResolution{err: err}
		}
		return wikiLinkResolution{resolvedID: fullID, changed: true}
	}

	// No type prefix: same-type shorthand
	fullID, err := resolveByName(nameIndex, sourceType, target)
	if err != nil {
		return wikiLinkResolution{err: err}
	}
	return wikiLinkResolution{resolvedID: fullID, changed: true}
}

// expandWikiLinksInBody replaces shorthand wiki-link targets with their resolved full IDs.
// It returns the modified body and the number of replacements made.
func expandWikiLinksInBody(body string, resolutions map[string]string) (string, int) {
	if len(resolutions) == 0 {
		return body, 0
	}
	matches := wikiLinkPattern.FindAllStringSubmatchIndex(body, -1)
	if len(matches) == 0 {
		return body, 0
	}

	var b strings.Builder
	b.Grow(len(body))
	count := 0
	last := 0
	for _, loc := range matches {
		target := body[loc[2]:loc[3]]
		fullID, ok := resolutions[target]
		if !ok {
			continue
		}
		b.WriteString(body[last:loc[0]])
		b.WriteString("[[")
		b.WriteString(fullID)
		if loc[4] >= 0 {
			b.WriteString("|")
			b.WriteString(body[loc[4]:loc[5]])
		}
		b.WriteString("]]")
		count++
		last = loc[1]
	}
	if count == 0 {
		return body, 0
	}
	b.WriteString(body[last:])
	return b.String(), count
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


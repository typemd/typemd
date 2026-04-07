package core

import (
	"fmt"
	"log/slog"
	"strings"
)

// QueryService orchestrates read-side operations.
// It coordinates the repository (for source-of-truth reads by ID)
// and the index (for search and discovery).
type QueryService struct {
	repo  ObjectRepository
	index ObjectIndex
}

// NewQueryService creates a QueryService.
func NewQueryService(repo ObjectRepository, index ObjectIndex) *QueryService {
	return &QueryService{repo: repo, index: index}
}

// Get reads an object from its source file by known ID.
func (s *QueryService) Get(id string) (*Object, error) {
	return s.repo.Get(id)
}

// Resolve resolves a (possibly abbreviated) object ID to the full ID.
func (s *QueryService) Resolve(prefix string) (string, error) {
	parts := splitID(prefix)
	if parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("invalid object ID format: %q", prefix)
	}
	typeName, namePrefix := parts[0], parts[1]

	// 1. Exact match
	if _, err := s.repo.ModTime(prefix); err == nil {
		return prefix, nil
	}

	// 2. Glob for prefix matches
	matches, err := s.repo.GlobIDs(typeName, namePrefix)
	if err != nil {
		return "", fmt.Errorf("glob error: %w", err)
	}

	switch len(matches) {
	case 0:
		return "", fmt.Errorf("no object found matching %q", prefix)
	case 1:
		return matches[0], nil
	default:
		return "", &AmbiguousMatchError{Prefix: prefix, Matches: matches}
	}
}

// Query queries objects using structured filter rules with optional options.
// Archived objects are excluded by default unless QueryIncludeArchived() is passed.
// Falls back to filesystem scanning if the index is unavailable.
func (s *QueryService) Query(filter []FilterRule, opts ...QueryOption) ([]*Object, error) {
	cfg := buildQueryConfig(opts)

	// Exclude archived objects by default.
	// Defensive copy to avoid mutating the caller's backing array.
	if !cfg.includeArchived {
		filter = append(append([]FilterRule(nil), filter...), FilterRule{
			Property: ArchivedProperty,
			Operator: "is_not",
			Value:    "true",
		})
	}

	slog.Debug("query", "filters", len(filter))
	results, err := s.index.Query(filter, cfg.sort...)
	if err == nil {
		return objectResultsToObjects(results), nil
	}

	// Fallback: scan filesystem and filter in memory
	slog.Warn("index unavailable, using filesystem fallback", "op", "query", "error", err)
	objects, walkErr := s.repo.Walk()
	if walkErr != nil {
		return nil, fmt.Errorf("fallback walk: %w", walkErr)
	}

	var filtered []*Object
	for _, obj := range objects {
		if MatchFilters(obj, filter) {
			filtered = append(filtered, obj)
		}
	}

	SortObjects(filtered, cfg.sort)
	return filtered, nil
}

// Search performs full-text search.
// Falls back to case-insensitive substring matching if the index is unavailable.
func (s *QueryService) Search(keyword string) ([]*Object, error) {
	results, err := s.index.Search(keyword)
	if err == nil {
		slog.Debug("search", "keyword", keyword, "results", len(results))
		return objectResultsToObjects(results), nil
	}

	// Fallback: scan filesystem and substring match
	slog.Warn("index unavailable, using filesystem fallback", "op", "search", "error", err)
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, nil
	}

	objects, walkErr := s.repo.Walk()
	if walkErr != nil {
		return nil, fmt.Errorf("fallback walk: %w", walkErr)
	}

	var matched []*Object
	for _, obj := range objects {
		if matchesSearch(obj, keyword) {
			matched = append(matched, obj)
		}
	}
	return matched, nil
}

// matchesSearch checks if an object matches a search keyword via case-insensitive
// substring matching against name, description, and body.
// Note: narrower scope than FTS5 index search, which matches all property values.
func matchesSearch(obj *Object, keyword string) bool {
	return containsCI(obj.Properties["name"], keyword) ||
		containsCI(obj.Properties["description"], keyword) ||
		containsCI(obj.Body, keyword)
}

// ListRelations returns all relations for an object.
func (s *QueryService) ListRelations(objectID string) ([]Relation, error) {
	return s.index.FindRelations(objectID)
}

// ListWikiLinks returns all wiki-links from an object.
func (s *QueryService) ListWikiLinks(objectID string) ([]StoredWikiLink, error) {
	return s.index.ListWikiLinks(objectID)
}

// ListBacklinks returns all wiki-links pointing to an object.
func (s *QueryService) ListBacklinks(objectID string) ([]StoredWikiLink, error) {
	return s.index.FindBacklinks(objectID)
}

// BuildDisplayProperties assembles display-ready properties for an object.
func (s *QueryService) BuildDisplayProperties(obj *Object) ([]DisplayProperty, error) {
	if obj == nil {
		return nil, nil
	}

	schema, _ := s.repo.GetSchema(obj.Type)

	// Build merged properties map without mutating obj
	merged := make(map[string]any, len(obj.Properties))
	for k, v := range obj.Properties {
		merged[k] = v
	}

	// Single pass over schema: fill missing properties + build property lookup
	schemaProp := make(map[string]*Property)
	if schema != nil {
		for i, p := range schema.Properties {
			if _, ok := merged[p.Name]; !ok {
				merged[p.Name] = nil
			}
			schemaProp[p.Name] = &schema.Properties[i]
		}
	}

	// Get relations
	relations, err := s.index.FindRelations(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("list relations: %w", err)
	}

	// Build ordered properties, separating local (non-schema, non-system) properties
	propKeys := OrderedPropKeys(merged, schema)
	var result []DisplayProperty
	var localProps []DisplayProperty
	for _, k := range propKeys {
		dp := DisplayProperty{
			Key:   k,
			Value: merged[k],
		}
		if sp, ok := schemaProp[k]; ok {
			dp.Type = sp.Type
			dp.Emoji = sp.Emoji
			dp.Pin = sp.Pin
			dp.IsRelation = sp.Type == "relation"
			result = append(result, dp)
		} else if sysProp := lookupSystemProperty(k); sysProp != nil {
			dp.Type = sysProp.Type
			result = append(result, dp)
		} else {
			dp.IsLocal = true
			localProps = append(localProps, dp)
		}
	}

	// Append reverse relations
	for _, r := range relations {
		if r.ToID == obj.ID {
			result = append(result, DisplayProperty{
				Key:       r.Name,
				Value:     r.FromID,
				IsReverse: true,
				FromID:    r.FromID,
			})
		}
	}

	// Append outgoing links
	wikiLinks, err := s.index.ListWikiLinks(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("list wiki-links: %w", err)
	}
	for _, wl := range wikiLinks {
		if wl.ToID == "" {
			continue // skip broken links
		}
		result = append(result, DisplayProperty{
			Key:    LinksDisplayKey,
			IsLink: true,
			FromID: wl.ToID,
		})
	}

	// Append backlinks
	backlinks, err := s.index.FindBacklinks(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("list backlinks: %w", err)
	}
	for _, bl := range backlinks {
		result = append(result, DisplayProperty{
			Key:        BacklinksDisplayKey,
			Value:      bl.FromID,
			IsBacklink: true,
			FromID:     bl.FromID,
		})
	}

	// Append local properties last (after reverse relations and backlinks)
	result = append(result, localProps...)

	return result, nil
}

// splitID splits "type/name" into [type, name]. Returns ["",""] on invalid input.
func splitID(id string) [2]string {
	for i, c := range id {
		if c == '/' {
			return [2]string{id[:i], id[i+1:]}
		}
	}
	return [2]string{"", ""}
}

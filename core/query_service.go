package core

import "fmt"

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

// Query queries objects using key=value filter syntax.
func (s *QueryService) Query(filter string) ([]*Object, error) {
	results, err := s.index.Query(filter)
	if err != nil {
		return nil, err
	}
	return objectResultsToObjects(results), nil
}

// Search performs full-text search.
func (s *QueryService) Search(keyword string) ([]*Object, error) {
	results, err := s.index.Search(keyword)
	if err != nil {
		return nil, err
	}
	return objectResultsToObjects(results), nil
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

	// Build ordered properties
	propKeys := OrderedPropKeys(merged, schema)
	var result []DisplayProperty
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
		}
		result = append(result, dp)
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

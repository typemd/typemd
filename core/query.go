package core

import "fmt"

// objectResultToObject converts an ObjectResult to an Object.
func objectResultToObject(r *ObjectResult) *Object {
	return &Object{
		ID:         r.ID,
		Type:       r.Type,
		Filename:   r.Filename,
		Properties: r.Properties,
		Body:       r.Body,
	}
}

// objectResultsToObjects converts a slice of ObjectResult to Objects.
func objectResultsToObjects(results []*ObjectResult) []*Object {
	if results == nil {
		return nil
	}
	objects := make([]*Object, len(results))
	for i, r := range results {
		objects[i] = objectResultToObject(r)
	}
	return objects
}

// QueryObjects queries objects with optional sort. Delegates to QueryService.
func (v *Vault) QueryObjects(filter []FilterRule, sort ...SortRule) ([]*Object, error) {
	if v.Queries == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	return v.Queries.Query(filter, sort...)
}

// SearchObjects performs full-text search. Delegates to QueryService.
func (v *Vault) SearchObjects(keyword string) ([]*Object, error) {
	if v.Queries == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	return v.Queries.Search(keyword)
}

// VaultStats returns aggregate statistics for all types. Delegates to QueryService.
func (v *Vault) VaultStats() (*VaultStats, error) {
	if v.Queries == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	return v.Queries.VaultStats()
}

// TypeStats returns per-property statistics for a type. Delegates to QueryService.
func (v *Vault) TypeStats(typeName string) (*TypeStats, error) {
	if v.Queries == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	return v.Queries.TypeStats(typeName)
}

// RebuildIndex rebuilds the FTS5 index.
func (v *Vault) RebuildIndex() error {
	if v.index == nil {
		return fmt.Errorf("vault not opened")
	}
	return v.index.Rebuild()
}

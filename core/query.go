package core

import "fmt"

// QueryOption configures query behavior.
type QueryOption func(*queryConfig)

type queryConfig struct {
	sort            []SortRule
	includeArchived bool
}

func buildQueryConfig(opts []QueryOption) queryConfig {
	var cfg queryConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

// QuerySort adds sort rules to a query.
func QuerySort(rules ...SortRule) QueryOption {
	return func(cfg *queryConfig) {
		cfg.sort = append(cfg.sort, rules...)
	}
}

// QueryIncludeArchived includes archived objects in query results.
func QueryIncludeArchived() QueryOption {
	return func(cfg *queryConfig) {
		cfg.includeArchived = true
	}
}

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

// QueryObjects queries objects with optional sort and options. Delegates to QueryService.
func (v *Vault) QueryObjects(filter []FilterRule, opts ...QueryOption) ([]*Object, error) {
	if v.Queries == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	return v.Queries.Query(filter, opts...)
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

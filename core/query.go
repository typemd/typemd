package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

// filterCondition represents a single key=value query condition.
type filterCondition struct {
	Key   string
	Value string
}

// parseFilter parses a filter string like "type=book status=reading"
// into a slice of filterConditions. Empty string returns nil, nil.
// Returns an error if any token is not in key=value form.
func parseFilter(filter string) ([]filterCondition, error) {
	filter = strings.TrimSpace(filter)
	if filter == "" {
		return nil, nil
	}

	parts := strings.Fields(filter)
	conditions := make([]filterCondition, 0, len(parts))

	for _, part := range parts {
		idx := strings.Index(part, "=")
		if idx < 1 {
			return nil, fmt.Errorf("invalid filter condition: %q (expected key=value)", part)
		}
		conditions = append(conditions, filterCondition{
			Key:   part[:idx],
			Value: part[idx+1:],
		})
	}

	return conditions, nil
}

// scanObjects scans a sql.Rows result into a slice of Objects.
func scanObjects(rows *sql.Rows) ([]*Object, error) {
	var results []*Object
	for rows.Next() {
		var obj Object
		var propsJSON string
		if err := rows.Scan(&obj.ID, &obj.Type, &obj.Filename, &propsJSON, &obj.Body); err != nil {
			return nil, fmt.Errorf("scan object: %w", err)
		}
		if err := json.Unmarshal([]byte(propsJSON), &obj.Properties); err != nil {
			return nil, fmt.Errorf("unmarshal properties: %w", err)
		}
		results = append(results, &obj)
	}
	return results, rows.Err()
}

// objectResultToObject converts an ObjectResult to an Object.
// The Body field is empty since ObjectResult is a lightweight projection.
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

// QueryObjects queries objects using key=value filter syntax.
// Multiple conditions are combined with AND.
// "type" is a special key that filters on the objects.type column.
// Other keys filter on JSON properties using json_extract.
// An empty filter returns all objects.
func (v *Vault) QueryObjects(filter string) ([]*Object, error) {
	if v.index == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	results, err := v.index.Query(filter)
	if err != nil {
		return nil, err
	}
	return objectResultsToObjects(results), nil
}

// SearchObjects performs full-text search using FTS5.
// Searches across filename, properties, and body.
// Returns nil, nil for empty keyword.
func (v *Vault) SearchObjects(keyword string) ([]*Object, error) {
	if v.index == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	results, err := v.index.Search(keyword)
	if err != nil {
		return nil, err
	}
	return objectResultsToObjects(results), nil
}

// RebuildIndex rebuilds the FTS5 index from the objects table.
func (v *Vault) RebuildIndex() error {
	if v.index == nil {
		return fmt.Errorf("vault not opened")
	}
	return v.index.Rebuild()
}

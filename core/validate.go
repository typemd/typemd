package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ValidateAllObjects queries all objects and validates properties against their type schema.
// Returns a map of object ID to validation errors.
func ValidateAllObjects(v *Vault) map[string][]error {
	result := make(map[string][]error)
	objects, err := v.QueryObjects("")
	if err != nil {
		return result
	}
	for _, obj := range objects {
		schema, err := v.LoadType(obj.Type)
		if err != nil {
			result[obj.ID] = []error{fmt.Errorf("load type %q: %w", obj.Type, err)}
			continue
		}
		if errs := ValidateObject(obj.Properties, schema); len(errs) > 0 {
			result[obj.ID] = errs
		}
	}
	return result
}

// ValidateRelations checks that all relation endpoints reference existing objects.
func ValidateRelations(v *Vault) []error {
	var errs []error
	rows, err := v.db.Query("SELECT name, from_id, to_id FROM relations")
	if err != nil {
		return []error{fmt.Errorf("query relations: %w", err)}
	}
	defer rows.Close()

	for rows.Next() {
		var name, fromID, toID string
		if err := rows.Scan(&name, &fromID, &toID); err != nil {
			errs = append(errs, fmt.Errorf("scan relation: %w", err))
			continue
		}
		var count int
		if err := v.db.QueryRow("SELECT count(*) FROM objects WHERE id = ?", fromID).Scan(&count); err != nil {
			errs = append(errs, fmt.Errorf("check source %s: %w", fromID, err))
		} else if count == 0 {
			errs = append(errs, fmt.Errorf("%s -[%s]-> %s: source object not found", fromID, name, toID))
		}
		if err := v.db.QueryRow("SELECT count(*) FROM objects WHERE id = ?", toID).Scan(&count); err != nil {
			errs = append(errs, fmt.Errorf("check target %s: %w", toID, err))
		} else if count == 0 {
			errs = append(errs, fmt.Errorf("%s -[%s]-> %s: target object not found", fromID, name, toID))
		}
	}
	if err := rows.Err(); err != nil {
		errs = append(errs, fmt.Errorf("iterate relations: %w", err))
	}
	return errs
}

// ValidateWikiLinks checks for broken wiki-links (targets that don't resolve to existing objects).
func ValidateWikiLinks(v *Vault) []error {
	var errs []error
	rows, err := v.db.Query("SELECT from_id, target FROM wikilinks WHERE to_id = ''")
	if err != nil {
		return []error{fmt.Errorf("query broken wikilinks: %w", err)}
	}
	defer rows.Close()

	for rows.Next() {
		var fromID, target string
		if err := rows.Scan(&fromID, &target); err != nil {
			errs = append(errs, fmt.Errorf("scan wikilink: %w", err))
			continue
		}
		errs = append(errs, fmt.Errorf("%s: broken wiki-link [[%s]]", fromID, target))
	}
	if err := rows.Err(); err != nil {
		errs = append(errs, fmt.Errorf("iterate wikilinks: %w", err))
	}
	return errs
}

// ValidateNameUniqueness checks that no two objects of the same unique type share the same name.
// It scans all types with Unique: true and reports duplicate name values.
func ValidateNameUniqueness(v *Vault) []error {
	var errs []error

	// Collect all unique type names
	uniqueTypes := collectUniqueTypes(v)
	if len(uniqueTypes) == 0 {
		return nil
	}

	for _, typeName := range uniqueTypes {
		rows, err := v.db.Query(
			"SELECT id, json_extract(properties, '$.name') AS name FROM objects WHERE type = ?",
			typeName,
		)
		if err != nil {
			errs = append(errs, fmt.Errorf("query %s objects: %w", typeName, err))
			continue
		}

		seen := make(map[string]string) // name → first ID
		for rows.Next() {
			var id string
			var name *string
			if err := rows.Scan(&id, &name); err != nil {
				continue
			}
			if name == nil || *name == "" {
				continue
			}
			if firstID, ok := seen[*name]; ok {
				errs = append(errs, fmt.Errorf("duplicate %s name %q: %s and %s", typeName, *name, firstID, id))
			} else {
				seen[*name] = id
			}
		}
		if err := rows.Err(); err != nil {
			errs = append(errs, fmt.Errorf("iterate %s objects: %w", typeName, err))
		}
		rows.Close()
	}
	return errs
}

// collectUniqueTypes returns all type names that have Unique: true.
func collectUniqueTypes(v *Vault) []string {
	var uniqueTypes []string
	for _, name := range v.ListTypes() {
		schema, err := v.LoadType(name)
		if err != nil {
			continue
		}
		if schema.Unique {
			uniqueTypes = append(uniqueTypes, name)
		}
	}
	return uniqueTypes
}

// ValidateAllSchemas scans .typemd/types/*.yaml and validates each schema.
// Also validates shared properties if .typemd/properties.yaml exists.
// Returns a map of type name to validation errors.
func ValidateAllSchemas(v *Vault) map[string][]error {
	result := make(map[string][]error)

	// Load and validate shared properties
	sharedProps, err := v.LoadSharedProperties()
	if err != nil {
		result["_shared_properties"] = []error{err}
		return result
	}
	if sharedProps != nil {
		if errs := ValidateSharedProperties(sharedProps); len(errs) > 0 {
			result["_shared_properties"] = errs
		}
	}

	entries, err := os.ReadDir(v.TypesDir())
	if err != nil {
		return result
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		typeName := strings.TrimSuffix(entry.Name(), ".yaml")
		data, err := os.ReadFile(filepath.Join(v.TypesDir(), entry.Name()))
		if err != nil {
			result[typeName] = []error{fmt.Errorf("read file: %w", err)}
			continue
		}
		var schema TypeSchema
		if err := yaml.Unmarshal(data, &schema); err != nil {
			result[typeName] = []error{fmt.Errorf("parse YAML: %w", err)}
			continue
		}
		if errs := ValidateSchema(&schema, sharedProps); len(errs) > 0 {
			result[typeName] = errs
		}
	}
	return result
}

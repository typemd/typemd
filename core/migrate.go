package core

import "fmt"

// MigrateOptions configures a migration operation.
type MigrateOptions struct {
	DryRun  bool
	Renames map[string]string // old_name -> new_name
}

// MigrateChange describes the changes applied to a single object.
type MigrateChange struct {
	ObjectID string
	Added    []string
	Removed  []string
	Renamed  map[string]string // old -> new
}

// MigrateResult holds the outcome of a migration operation.
type MigrateResult struct {
	Changes []MigrateChange
}

// MigrateObjects updates all objects of the given type to match the current schema.
// It adds missing properties (with defaults), removes obsolete ones, and optionally
// renames properties specified in opts.Renames.
func (v *Vault) MigrateObjects(typeName string, opts MigrateOptions) (*MigrateResult, error) {
	if v.db == nil {
		return nil, fmt.Errorf("vault not opened")
	}

	schema, err := v.LoadType(typeName)
	if err != nil {
		return nil, fmt.Errorf("load type %q: %w", typeName, err)
	}

	// Build set of schema property names
	schemaProps := make(map[string]Property)
	for _, p := range schema.Properties {
		schemaProps[p.Name] = p
	}

	// Validate rename mappings
	for oldName, newName := range opts.Renames {
		if _, exists := schemaProps[newName]; !exists {
			return nil, fmt.Errorf("rename target %q not found in schema %q", newName, typeName)
		}
		if _, exists := schemaProps[oldName]; exists {
			return nil, fmt.Errorf("rename source %q still exists in schema %q (remove it from schema first)", oldName, typeName)
		}
	}

	// Query all objects of this type
	objects, err := v.QueryObjects("type=" + typeName)
	if err != nil {
		return nil, fmt.Errorf("query objects: %w", err)
	}

	result := &MigrateResult{}

	for _, obj := range objects {
		change := MigrateChange{ObjectID: obj.ID}

		// 1. Rename: move old property value to new name
		for oldName, newName := range opts.Renames {
			oldVal, hasOld := obj.Properties[oldName]
			_, hasNew := obj.Properties[newName]
			if hasOld && !hasNew {
				obj.Properties[newName] = oldVal
				delete(obj.Properties, oldName)
				if change.Renamed == nil {
					change.Renamed = make(map[string]string)
				}
				change.Renamed[oldName] = newName
			}
		}

		// 2. Add: schema properties missing from object
		for _, p := range schema.Properties {
			if _, exists := obj.Properties[p.Name]; !exists {
				obj.Properties[p.Name] = p.Default
				change.Added = append(change.Added, p.Name)
			}
		}

		// 3. Remove: object properties not in schema (and not a rename source)
		for key := range obj.Properties {
			if _, inSchema := schemaProps[key]; !inSchema {
				if _, isRenameSource := opts.Renames[key]; !isRenameSource {
					delete(obj.Properties, key)
					change.Removed = append(change.Removed, key)
				}
			}
		}

		// Skip if no changes
		if len(change.Added) == 0 && len(change.Removed) == 0 && len(change.Renamed) == 0 {
			continue
		}

		if !opts.DryRun {
			if err := v.writeObjectProperties(obj); err != nil {
				return nil, fmt.Errorf("write object %s: %w", obj.ID, err)
			}
		}

		result.Changes = append(result.Changes, change)
	}

	return result, nil
}

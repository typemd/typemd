package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// schemaPath returns the path to the type schema file, checking directory
// format first (types/<name>/schema.yaml), then single-file (types/<name>.yaml).
// Returns the path and whether it's in directory format.
func (r *LocalObjectRepository) schemaPath(name string) (string, bool) {
	dirPath := filepath.Join(r.typesDir(), name, "schema.yaml")
	if _, err := os.Stat(dirPath); err == nil {
		return dirPath, true
	}
	return filepath.Join(r.typesDir(), name+".yaml"), false
}

// migrateToDirectory migrates a single-file type schema to directory format.
func (r *LocalObjectRepository) migrateToDirectory(name string) error {
	singleFile := filepath.Join(r.typesDir(), name+".yaml")
	dirPath := filepath.Join(r.typesDir(), name)
	schemaFile := filepath.Join(dirPath, "schema.yaml")

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("create type directory %s: %w", name, err)
	}
	if err := os.Rename(singleFile, schemaFile); err != nil {
		return fmt.Errorf("migrate type schema %s: %w", name, err)
	}
	return nil
}

// GetSchema loads a type schema by name, resolving shared property references.
// Checks directory format first, then single file, auto-migrating on load.
func (r *LocalObjectRepository) GetSchema(name string) (*TypeSchema, error) {
	path, isDir := r.schemaPath(name)

	data, err := os.ReadFile(path)
	if err == nil {
		var schema TypeSchema
		if err := yaml.Unmarshal(data, &schema); err != nil {
			return nil, fmt.Errorf("parse type schema %s: %w", name, err)
		}

		// Extract name template from properties if present
		filtered := schema.Properties[:0]
		for _, prop := range schema.Properties {
			if prop.Name == NameProperty {
				schema.NameTemplate = prop.Template
				continue
			}
			filtered = append(filtered, prop)
		}
		schema.Properties = filtered

		// Normalize empty version to default
		if schema.Version == "" {
			schema.Version = DefaultSchemaVersion
		}

		// Resolve use entries if any exist
		if err := r.resolveSchemaUseEntries(&schema); err != nil {
			return nil, fmt.Errorf("resolve type schema %s: %w", name, err)
		}

		// Auto-migrate single file to directory format
		if !isDir {
			_ = r.migrateToDirectory(name)
		}

		return &schema, nil
	}

	if schema, ok := defaultTypes[name]; ok {
		if schema.Version == "" {
			schema.Version = DefaultSchemaVersion
		}
		return &schema, nil
	}

	return nil, fmt.Errorf("unknown type: %s", name)
}

// resolveSchemaUseEntries resolves use entries in a schema if any are present.
func (r *LocalObjectRepository) resolveSchemaUseEntries(schema *TypeSchema) error {
	hasUse := false
	for _, p := range schema.Properties {
		if p.Use != "" {
			hasUse = true
			break
		}
	}
	if !hasUse {
		return nil
	}

	_, err := r.GetSharedProperties()
	if err != nil {
		return err
	}
	return resolveUseEntries(schema, r.sharedPropsMap)
}

// WriteSchema writes raw schema data to a type schema file in directory format.
// If a single-file format exists, it is removed after writing.
func (r *LocalObjectRepository) WriteSchema(typeName string, data []byte) error {
	dirPath := filepath.Join(r.typesDir(), typeName)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("create type directory %s: %w", typeName, err)
	}
	schemaFile := filepath.Join(dirPath, "schema.yaml")
	if err := os.WriteFile(schemaFile, data, 0644); err != nil {
		return err
	}
	// Remove old single file if it exists
	singleFile := filepath.Join(r.typesDir(), typeName+".yaml")
	os.Remove(singleFile) // ignore error (file may not exist)
	return nil
}

// DeleteSchema removes a type schema, handling both directory and single-file formats.
func (r *LocalObjectRepository) DeleteSchema(typeName string) error {
	// Try directory format first
	dirPath := filepath.Join(r.typesDir(), typeName)
	if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
		if err := os.RemoveAll(dirPath); err != nil {
			return fmt.Errorf("remove type directory %q: %w", typeName, err)
		}
		return nil
	}

	// Fall back to single file
	path := filepath.Join(r.typesDir(), typeName+".yaml")
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("type schema %q does not exist", typeName)
		}
		return err
	}
	return nil
}

// ListSchemas returns the names of all available types (custom + built-in).
// Discovers both single-file (.yaml) and directory (dir/schema.yaml) formats.
func (r *LocalObjectRepository) ListSchemas() ([]string, error) {
	seen := make(map[string]bool)

	entries, err := os.ReadDir(r.typesDir())
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				// Directory format: check for schema.yaml inside
				schemaPath := filepath.Join(r.typesDir(), e.Name(), "schema.yaml")
				if _, err := os.Stat(schemaPath); err == nil {
					seen[e.Name()] = true
				}
			} else if strings.HasSuffix(e.Name(), ".yaml") {
				// Single-file format
				name := strings.TrimSuffix(e.Name(), ".yaml")
				seen[name] = true
			}
		}
	}

	// Built-in defaults
	for name := range defaultTypes {
		seen[name] = true
	}

	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// schemaPath returns the path to the type schema file (types/<name>/schema.yaml).
func (r *LocalObjectRepository) schemaPath(name string) string {
	return filepath.Join(r.typesDir(), name, "schema.yaml")
}

// GetSchema loads a type schema by name, resolving shared property references.
func (r *LocalObjectRepository) GetSchema(name string) (*TypeSchema, error) {
	path := r.schemaPath(name)

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

// WriteSchema writes raw schema data to a type schema file.
func (r *LocalObjectRepository) WriteSchema(typeName string, data []byte) error {
	dirPath := filepath.Join(r.typesDir(), typeName)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("create type directory %s: %w", typeName, err)
	}
	schemaFile := filepath.Join(dirPath, "schema.yaml")
	return os.WriteFile(schemaFile, data, 0644)
}

// DeleteSchema removes a type schema directory.
func (r *LocalObjectRepository) DeleteSchema(typeName string) error {
	dirPath := filepath.Join(r.typesDir(), typeName)
	if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
		if err := os.RemoveAll(dirPath); err != nil {
			return fmt.Errorf("remove type directory %q: %w", typeName, err)
		}
		return nil
	}
	return fmt.Errorf("type schema %q does not exist", typeName)
}

// ListSchemas returns the names of all available types (custom + built-in).
func (r *LocalObjectRepository) ListSchemas() ([]string, error) {
	seen := make(map[string]bool)

	entries, err := os.ReadDir(r.typesDir())
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				schemaPath := filepath.Join(r.typesDir(), e.Name(), "schema.yaml")
				if _, err := os.Stat(schemaPath); err == nil {
					seen[e.Name()] = true
				}
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

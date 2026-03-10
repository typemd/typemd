package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

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

		// 3. Remove: object properties not in schema (and not a rename source or system property)
		for key := range obj.Properties {
			if key == NameProperty {
				continue // system property, never remove
			}
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
			if err := v.saveObjectFile(obj); err != nil {
				return nil, fmt.Errorf("write object %s: %w", obj.ID, err)
			}
		}

		result.Changes = append(result.Changes, change)
	}

	return result, nil
}

// SchemaMigrateChange describes a schema-level migration for a single type.
type SchemaMigrateChange struct {
	TypeName   string
	Properties []string // property names that were converted
}

// SchemaMigrateResult holds the outcome of a schema migration.
type SchemaMigrateResult struct {
	Changes []SchemaMigrateChange
}

// MigrateSchemas scans .typemd/types/*.yaml and converts enum+values to select+options.
// When dryRun is true, it reports what would change without modifying files.
func (v *Vault) MigrateSchemas(dryRun bool) (*SchemaMigrateResult, error) {
	result := &SchemaMigrateResult{}

	entries, err := os.ReadDir(v.TypesDir())
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}
		return nil, fmt.Errorf("read types dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		path := filepath.Join(v.TypesDir(), entry.Name())
		typeName := strings.TrimSuffix(entry.Name(), ".yaml")

		change, err := migrateSchemaFile(path, typeName, dryRun)
		if err != nil {
			return nil, fmt.Errorf("migrate schema %s: %w", typeName, err)
		}
		if change != nil {
			result.Changes = append(result.Changes, *change)
		}
	}

	return result, nil
}

// migrateSchemaFile converts enum+values to select+options in a single schema file.
// Uses raw YAML node manipulation to preserve comments and formatting where possible.
func migrateSchemaFile(path, typeName string, dryRun bool) (*SchemaMigrateChange, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse YAML: %w", err)
	}

	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return nil, nil
	}

	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil, nil
	}

	propsNode := findYAMLKey(root, "properties")
	if propsNode == nil || propsNode.Kind != yaml.SequenceNode {
		return nil, nil
	}

	var converted []string

	for _, propNode := range propsNode.Content {
		if propNode.Kind != yaml.MappingNode {
			continue
		}

		propName := getYAMLValue(propNode, "name")
		typeVal := getYAMLValue(propNode, "type")

		if typeVal != "enum" {
			continue
		}

		// Get values array
		valuesNode := findYAMLKey(propNode, "values")
		if valuesNode == nil || valuesNode.Kind != yaml.SequenceNode {
			continue
		}

		// Build options from values
		optionsNode := &yaml.Node{
			Kind: yaml.SequenceNode,
			Tag:  "!!seq",
		}
		for _, vNode := range valuesNode.Content {
			optionMap := &yaml.Node{
				Kind: yaml.MappingNode,
				Tag:  "!!map",
				Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Tag: "!!str", Value: "value"},
					{Kind: yaml.ScalarNode, Tag: "!!str", Value: vNode.Value},
				},
			}
			optionsNode.Content = append(optionsNode.Content, optionMap)
		}

		// Replace type: enum -> type: select
		setYAMLValue(propNode, "type", "select")

		// Replace values key+node with options
		replaceYAMLKey(propNode, "values", "options", optionsNode)

		converted = append(converted, propName)
	}

	if len(converted) == 0 {
		return nil, nil
	}

	if !dryRun {
		out, err := yaml.Marshal(&doc)
		if err != nil {
			return nil, fmt.Errorf("marshal YAML: %w", err)
		}
		if err := os.WriteFile(path, out, 0644); err != nil {
			return nil, fmt.Errorf("write file: %w", err)
		}
	}

	return &SchemaMigrateChange{
		TypeName:   typeName,
		Properties: converted,
	}, nil
}

// findYAMLKey finds the value node for a given key in a mapping node.
func findYAMLKey(mapping *yaml.Node, key string) *yaml.Node {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}

// getYAMLValue gets a scalar string value for a key in a mapping node.
func getYAMLValue(mapping *yaml.Node, key string) string {
	node := findYAMLKey(mapping, key)
	if node == nil {
		return ""
	}
	return node.Value
}

// setYAMLValue sets a scalar string value for a key in a mapping node.
func setYAMLValue(mapping *yaml.Node, key, value string) {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == key {
			mapping.Content[i+1].Value = value
			return
		}
	}
}

// replaceYAMLKey replaces a key name and its value node in a mapping.
func replaceYAMLKey(mapping *yaml.Node, oldKey, newKey string, newValue *yaml.Node) {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == oldKey {
			mapping.Content[i].Value = newKey
			mapping.Content[i+1] = newValue
			return
		}
	}
}

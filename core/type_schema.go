package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// OrderedPropKeys returns property keys ordered by schema definition.
// Keys not in the schema are appended alphabetically at the end.
// If schema is nil, keys are sorted alphabetically.
func OrderedPropKeys(props map[string]any, schema *TypeSchema) []string {
	if schema == nil {
		keys := make([]string, 0, len(props))
		for k := range props {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys
	}

	seen := make(map[string]bool)
	var keys []string
	for _, p := range schema.Properties {
		if _, ok := props[p.Name]; ok {
			keys = append(keys, p.Name)
			seen[p.Name] = true
		}
	}
	var extra []string
	for k := range props {
		if !seen[k] {
			extra = append(extra, k)
		}
	}
	sort.Strings(extra)
	return append(keys, extra...)
}

// TypeSchema defines the schema for a type.
type TypeSchema struct {
	Name       string     `yaml:"name"`
	Properties []Property `yaml:"properties"`
}

// Property defines a single property in a type schema.
type Property struct {
	Name          string   `yaml:"name"`
	Type          string   `yaml:"type"`
	Values        []string `yaml:"values,omitempty"`
	Target        string   `yaml:"target,omitempty"`
	Default       any      `yaml:"default,omitempty"`
	Multiple      bool     `yaml:"multiple,omitempty"`
	Bidirectional bool     `yaml:"bidirectional,omitempty"`
	Inverse       string   `yaml:"inverse,omitempty"`
}


// defaultTypes contains built-in type schemas.
var defaultTypes = map[string]TypeSchema{
	"book": {
		Name: "book",
		Properties: []Property{
			{Name: "title", Type: "string"},
			{Name: "status", Type: "enum", Values: []string{"to-read", "reading", "done"}},
			{Name: "rating", Type: "number"},
		},
	},
	"person": {
		Name: "person",
		Properties: []Property{
			{Name: "name", Type: "string"},
			{Name: "role", Type: "string"},
		},
	},
	"note": {
		Name: "note",
		Properties: []Property{
			{Name: "title", Type: "string"},
			{Name: "tags", Type: "string"},
		},
	},
}

// LoadType loads a type schema by name.
// It first looks for a custom YAML file in .typemd/types/{name}.yaml,
// then falls back to built-in defaults.
func (v *Vault) LoadType(name string) (*TypeSchema, error) {
	path := filepath.Join(v.TypesDir(), name+".yaml")

	data, err := os.ReadFile(path)
	if err == nil {
		var schema TypeSchema
		if err := yaml.Unmarshal(data, &schema); err != nil {
			return nil, fmt.Errorf("parse type schema %s: %w", name, err)
		}
		return &schema, nil
	}

	if schema, ok := defaultTypes[name]; ok {
		return &schema, nil
	}

	return nil, fmt.Errorf("unknown type: %s", name)
}

// validPropertyTypes lists allowed property types.
var validPropertyTypes = map[string]bool{
	"string": true, "number": true, "enum": true, "relation": true,
}

// ValidateSchema validates a type schema itself for correctness.
func ValidateSchema(schema *TypeSchema) []error {
	var errs []error
	if schema.Name == "" {
		errs = append(errs, fmt.Errorf("schema missing required field: name"))
	}
	seen := make(map[string]bool)
	for i, prop := range schema.Properties {
		if prop.Name == "" {
			errs = append(errs, fmt.Errorf("property[%d]: missing required field: name", i))
			continue
		}
		if seen[prop.Name] {
			errs = append(errs, fmt.Errorf("property %q: duplicate property name", prop.Name))
		}
		seen[prop.Name] = true
		if prop.Type == "" {
			errs = append(errs, fmt.Errorf("property %q: missing required field: type", prop.Name))
			continue
		}
		if !validPropertyTypes[prop.Type] {
			errs = append(errs, fmt.Errorf("property %q: invalid type %q (valid: string, number, enum, relation)", prop.Name, prop.Type))
			continue
		}
		if prop.Type == "enum" && len(prop.Values) == 0 {
			errs = append(errs, fmt.Errorf("property %q: enum type requires non-empty values", prop.Name))
		}
		if prop.Type == "relation" && prop.Target == "" {
			errs = append(errs, fmt.Errorf("property %q: relation type requires target", prop.Name))
		}
	}
	return errs
}

// ValidateObject validates object properties against a type schema.
// Lenient mode: only validates properties defined in schema, ignores extra properties.
// Properties defined in schema but missing from props are also ignored.
func ValidateObject(props map[string]any, schema *TypeSchema) []error {
	var errs []error

	for _, prop := range schema.Properties {
		val, ok := props[prop.Name]
		if !ok {
			continue
		}

		switch prop.Type {
		case "string":
			if _, ok := val.(string); !ok {
				errs = append(errs, fmt.Errorf("property %q: expected string, got %T", prop.Name, val))
			}
		case "relation":
			if prop.Multiple {
				arr, ok := val.([]any)
				if !ok {
					errs = append(errs, fmt.Errorf("property %q: expected array, got %T", prop.Name, val))
					continue
				}
				for i, item := range arr {
					if _, ok := item.(string); !ok {
						errs = append(errs, fmt.Errorf("property %q[%d]: expected string, got %T", prop.Name, i, item))
					}
				}
			} else {
				if _, ok := val.(string); !ok {
					errs = append(errs, fmt.Errorf("property %q: expected string, got %T", prop.Name, val))
				}
			}
		case "number":
			switch val.(type) {
			case int, int64, float64:
				// valid
			default:
				errs = append(errs, fmt.Errorf("property %q: expected number, got %T", prop.Name, val))
			}
		case "enum":
			s, ok := val.(string)
			if !ok {
				errs = append(errs, fmt.Errorf("property %q: expected string for enum, got %T", prop.Name, val))
				continue
			}
			found := false
			for _, v := range prop.Values {
				if v == s {
					found = true
					break
				}
			}
			if !found {
				errs = append(errs, fmt.Errorf("property %q: value %q not in allowed values %v", prop.Name, s, prop.Values))
			}
		}
	}

	return errs
}

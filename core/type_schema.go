package core

import (
	"fmt"
	"regexp"
	"sort"
)

var dateRegexp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// OrderedPropKeys returns property keys ordered by schema definition.
// System properties appear first in registry order (name, created_at, updated_at),
// followed by schema-defined properties, then extras alphabetically.
// If schema is nil, keys are sorted alphabetically (with system properties first).
func OrderedPropKeys(props map[string]any, schema *TypeSchema) []string {
	// Collect stored system properties that are present, in registry order.
	// All system property names (stored + computed) go into sysSet so they
	// are excluded from the extra-keys section. Only stored properties
	// appear in prefix (frontmatter output).
	var prefix []string
	sysSet := make(map[string]bool, len(systemProperties))
	for i := range systemProperties {
		sp := &systemProperties[i]
		sysSet[sp.Name] = true
		if !sp.Computed {
			if _, ok := props[sp.Name]; ok {
				prefix = append(prefix, sp.Name)
			}
		}
	}

	if schema == nil {
		keys := make([]string, 0, len(props))
		for k := range props {
			if !sysSet[k] {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		return append(prefix, keys...)
	}

	seen := make(map[string]bool)
	for k := range sysSet {
		seen[k] = true
	}
	var keys []string
	for _, p := range schema.Properties {
		if _, ok := props[p.Name]; ok && !sysSet[p.Name] {
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
	keys = append(keys, extra...)
	return append(prefix, keys...)
}

// DefaultSchemaVersion is the default version for type schemas without an explicit version.
const DefaultSchemaVersion = "0.0"

// TypeSchema defines the schema for a type.
type TypeSchema struct {
	Name         string     `yaml:"name"`
	Plural       string     `yaml:"plural,omitempty"`
	Emoji        string     `yaml:"emoji,omitempty"`
	Color        string     `yaml:"color,omitempty"`
	Unique       bool       `yaml:"unique,omitempty"`
	Version      string     `yaml:"version,omitempty"`
	Description  string     `yaml:"description,omitempty"`
	Properties   []Property `yaml:"properties"`
	NameTemplate string     `yaml:"-"` // extracted from name property entry during load
}

// PluralName returns the plural form if set, otherwise falls back to Name.
func (s *TypeSchema) PluralName() string {
	if s.Plural != "" {
		return s.Plural
	}
	return s.Name
}

// Option defines a selectable value for select/multi_select properties.
type Option struct {
	Value string `yaml:"value"`
	Label string `yaml:"label,omitempty"`
}

// Property defines a single property in a type schema.
type Property struct {
	Name          string   `yaml:"name"`
	Use           string   `yaml:"use,omitempty"`
	Type          string   `yaml:"type"`
	Emoji         string   `yaml:"emoji,omitempty"`
	Description   string   `yaml:"description,omitempty"`
	Pin           int      `yaml:"pin,omitempty"`
	Options       []Option `yaml:"options,omitempty"`
	Target        string   `yaml:"target,omitempty"`
	Default       any      `yaml:"default,omitempty"`
	Multiple      bool     `yaml:"multiple,omitempty"`
	Bidirectional bool     `yaml:"bidirectional,omitempty"`
	Inverse       string   `yaml:"inverse,omitempty"`
	Template      string   `yaml:"template,omitempty"`
}

// LegacySharedPropertiesFile represents the old properties/properties.yaml format (pre-v0.9).
// Used only for migration from single-file to per-property-file format.
type LegacySharedPropertiesFile struct {
	Properties []Property `yaml:"properties"`
}

// FindProperty returns a pointer to the named property, or nil if not found.
func (s *TypeSchema) FindProperty(name string) *Property {
	for i, p := range s.Properties {
		if p.Name == name {
			return &s.Properties[i]
		}
	}
	return nil
}

// FindRelation returns a pointer to the named relation property, checking both
// schema-defined properties and system properties (e.g. tags).
func (s *TypeSchema) FindRelation(name string) *Property {
	for i, p := range s.Properties {
		if p.Name == name && p.Type == "relation" {
			return &s.Properties[i]
		}
	}
	return findSystemRelationProperty(name)
}

// PropertyNames returns the set of property names defined in the schema.
func (s *TypeSchema) PropertyNames() map[string]bool {
	names := make(map[string]bool, len(s.Properties))
	for _, p := range s.Properties {
		names[p.Name] = true
	}
	return names
}

// OptionValues returns a slice of option values (convenience for display/error messages).
func (p Property) OptionValues() []string {
	vals := make([]string, len(p.Options))
	for i, opt := range p.Options {
		vals[i] = opt.Value
	}
	return vals
}

// defaultTypes contains built-in type schemas.
// "tag" backs the "tags" system property; "page" is a general-purpose content container.
// All other types must be defined via types/<name>/schema.yaml files.
var defaultTypes = map[string]TypeSchema{
	TagTypeName: {
		Name:   TagTypeName,
		Plural: "tags",
		Emoji:  "🏷️",
		Unique: true,
		Properties: []Property{
			{Name: "color", Type: "string", Emoji: "🎨"},
			{Name: "icon", Type: "string", Emoji: "✨"},
		},
	},
	PageTypeName: {
		Name:   PageTypeName,
		Plural: "pages",
		Emoji:  "📄",
	},
}

// ValidPropertyTypeNames returns the list of allowed property type names in display order.
func ValidPropertyTypeNames() []string {
	return []string{
		"string", "number", "date", "datetime", "url",
		"checkbox", "select", "multi_select", "relation",
	}
}

// validPropertyTypes lists allowed property types.
var validPropertyTypes = map[string]bool{
	"string":       true,
	"number":       true,
	"date":         true,
	"datetime":     true,
	"url":          true,
	"checkbox":     true,
	"select":       true,
	"multi_select": true,
	"relation":     true,
}

// LoadType loads a type schema by name, using the in-memory cache when available.
func (v *Vault) LoadType(name string) (*TypeSchema, error) {
	if v.schemaCache != nil {
		if cached, ok := v.schemaCache[name]; ok {
			return cached, nil
		}
	}
	schema, err := v.repo.GetSchema(name)
	if err != nil {
		return nil, err
	}
	if v.schemaCache == nil {
		v.schemaCache = make(map[string]*TypeSchema)
	}
	v.schemaCache[name] = schema
	return schema, nil
}

// SaveType validates and persists a TypeSchema to types/<name>/schema.yaml.
func (v *Vault) SaveType(schema *TypeSchema) error {
	if errs := ValidateSchema(schema); len(errs) > 0 {
		return fmt.Errorf("invalid schema: %v", errs[0])
	}
	data, err := MarshalTypeSchema(schema)
	if err != nil {
		return fmt.Errorf("marshal type schema: %w", err)
	}
	if err := v.repo.WriteSchema(schema.Name, data); err != nil {
		return err
	}
	delete(v.schemaCache, schema.Name)
	return nil
}

// DeleteType removes a user-defined type schema. Built-in types cannot be deleted.
func (v *Vault) DeleteType(name string) error {
	if _, ok := defaultTypes[name]; ok {
		return fmt.Errorf("cannot delete built-in type %q", name)
	}
	if err := v.repo.DeleteSchema(name); err != nil {
		return err
	}
	delete(v.schemaCache, name)
	return nil
}

// InvalidateSchemaCache clears the entire schema cache,
// forcing subsequent LoadType calls to reload from disk.
func (v *Vault) InvalidateSchemaCache() {
	v.schemaCache = nil
}

// CountObjectsByType returns the number of objects of the given type.
func (v *Vault) CountObjectsByType(typeName string) (int, error) {
	results, err := v.index.Query(TypeFilter(typeName))
	if err != nil {
		return 0, err
	}
	return len(results), nil
}

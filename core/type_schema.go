package core

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"gopkg.in/yaml.v3"
)

var dateRegexp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// OrderedPropKeys returns property keys ordered by schema definition.
// System properties appear first in registry order (name, created_at, updated_at),
// followed by schema-defined properties, then extras alphabetically.
// If schema is nil, keys are sorted alphabetically (with system properties first).
func OrderedPropKeys(props map[string]any, schema *TypeSchema) []string {
	// Collect system properties that are present, in registry order
	sysNames := SystemPropertyNames()
	var prefix []string
	sysSet := make(map[string]bool)
	for _, name := range sysNames {
		sysSet[name] = true
		if _, ok := props[name]; ok {
			prefix = append(prefix, name)
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

// TypeSchema defines the schema for a type.
type TypeSchema struct {
	Name       string     `yaml:"name"`
	Emoji      string     `yaml:"emoji,omitempty"`
	Properties []Property `yaml:"properties"`
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
	Pin           int      `yaml:"pin,omitempty"`
	Options       []Option `yaml:"options,omitempty"`
	Target        string   `yaml:"target,omitempty"`
	Default       any      `yaml:"default,omitempty"`
	Multiple      bool     `yaml:"multiple,omitempty"`
	Bidirectional bool     `yaml:"bidirectional,omitempty"`
	Inverse       string   `yaml:"inverse,omitempty"`
}

// SharedPropertiesFile represents the .typemd/properties.yaml file.
type SharedPropertiesFile struct {
	Properties []Property `yaml:"properties"`
}


// PropertyNames returns the set of property names defined in the schema.
func (s *TypeSchema) PropertyNames() map[string]bool {
	names := make(map[string]bool, len(s.Properties))
	for _, p := range s.Properties {
		names[p.Name] = true
	}
	return names
}

// defaultTypes contains built-in type schemas.
var defaultTypes = map[string]TypeSchema{
	"book": {
		Name:  "book",
		Emoji: "📚",
		Properties: []Property{
			{Name: "title", Type: "string", Emoji: "📖"},
			{Name: "status", Type: "select", Emoji: "📋", Options: []Option{
				{Value: "to-read"},
				{Value: "reading"},
				{Value: "done"},
			}},
			{Name: "rating", Type: "number", Emoji: "⭐"},
		},
	},
	"person": {
		Name:  "person",
		Emoji: "👤",
		Properties: []Property{
			{Name: "role", Type: "string", Emoji: "💼"},
		},
	},
	"note": {
		Name:  "note",
		Emoji: "📝",
		Properties: []Property{
			{Name: "title", Type: "string", Emoji: "🏷️"},
			{Name: "tags", Type: "string", Emoji: "🔖"},
		},
	},
}

// LoadType loads a type schema by name.
// It first looks for a custom YAML file in .typemd/types/{name}.yaml,
// then falls back to built-in defaults.
// If the schema contains `use` entries, they are resolved from shared properties.
func (v *Vault) LoadType(name string) (*TypeSchema, error) {
	path := filepath.Join(v.TypesDir(), name+".yaml")

	data, err := os.ReadFile(path)
	if err == nil {
		var schema TypeSchema
		if err := yaml.Unmarshal(data, &schema); err != nil {
			return nil, fmt.Errorf("parse type schema %s: %w", name, err)
		}

		// Resolve use entries if any exist
		if err := v.resolveSchemaUseEntries(&schema); err != nil {
			return nil, fmt.Errorf("resolve type schema %s: %w", name, err)
		}

		return &schema, nil
	}

	if schema, ok := defaultTypes[name]; ok {
		return &schema, nil
	}

	return nil, fmt.Errorf("unknown type: %s", name)
}

// resolveSchemaUseEntries resolves use entries in a schema if any are present.
func (v *Vault) resolveSchemaUseEntries(schema *TypeSchema) error {
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

	_, err := v.LoadSharedProperties()
	if err != nil {
		return err
	}
	return resolveUseEntries(schema, v.sharedPropsMap)
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

// validatePropertyType validates the type field and related constraints of a property.
// Returns validation errors with the given property name prefix.
func validatePropertyType(prop Property, namePrefix string) []error {
	var errs []error

	if prop.Type == "" {
		errs = append(errs, fmt.Errorf("%s: missing required field: type", namePrefix))
		return errs
	}

	if prop.Type == "enum" {
		errs = append(errs, fmt.Errorf("%s: type \"enum\" is no longer supported, use \"select\" with \"options\" instead", namePrefix))
		return errs
	}

	if !validPropertyTypes[prop.Type] {
		errs = append(errs, fmt.Errorf("%s: invalid type %q (valid: string, number, date, datetime, url, checkbox, select, multi_select, relation)", namePrefix, prop.Type))
		return errs
	}

	if (prop.Type == "select" || prop.Type == "multi_select") && len(prop.Options) == 0 {
		errs = append(errs, fmt.Errorf("%s: %s type requires non-empty options", namePrefix, prop.Type))
	}

	if prop.Type == "relation" && prop.Target == "" {
		errs = append(errs, fmt.Errorf("%s: relation type requires target", namePrefix))
	}

	return errs
}

// ValidateSchema validates a type schema itself for correctness.
// If sharedProps is provided, it also validates `use` entries against shared properties.
func ValidateSchema(schema *TypeSchema, sharedProps ...[]Property) []error {
	var errs []error
	if schema.Name == "" {
		errs = append(errs, fmt.Errorf("schema missing required field: name"))
	}

	// Build shared properties map if provided
	var sharedMap map[string]Property
	if len(sharedProps) > 0 && sharedProps[0] != nil {
		sharedMap = SharedPropertiesMap(sharedProps[0])
	}

	seen := make(map[string]bool)
	seenEmoji := make(map[string]string) // emoji -> property name
	seenPin := make(map[int]string)      // pin -> property name
	for i, prop := range schema.Properties {
		// Handle `use` entries
		if prop.Use != "" {
			if prop.Name != "" {
				errs = append(errs, fmt.Errorf("property[%d]: \"use\" and \"name\" are mutually exclusive", i))
				continue
			}

			// Validate only pin and emoji overrides are present
			if err := validateUseOverrides(i, prop); err != nil {
				errs = append(errs, err)
				continue
			}

			// Validate reference exists in shared properties
			if sharedMap != nil {
				shared, ok := sharedMap[prop.Use]
				if !ok {
					errs = append(errs, fmt.Errorf("property[%d]: shared property %q not found", i, prop.Use))
					continue
				}

				// Use the shared property name for duplicate checking
				propName := shared.Name
				if seen[propName] {
					errs = append(errs, fmt.Errorf("property %q: duplicate property name", propName))
				}
				seen[propName] = true

				// Check emoji uniqueness (use override or shared)
				emoji := prop.Emoji
				if emoji == "" {
					emoji = shared.Emoji
				}
				if emoji != "" {
					if otherProp, ok := seenEmoji[emoji]; ok {
						errs = append(errs, fmt.Errorf("property %q: duplicate property emoji %q (already used by %q)", propName, emoji, otherProp))
					}
					seenEmoji[emoji] = propName
				}

				// Check pin uniqueness (use override or shared)
				pin := prop.Pin
				if pin == 0 {
					pin = shared.Pin
				}
				if pin < 0 {
					errs = append(errs, fmt.Errorf("property %q: pin value must be a positive integer, got %d", propName, pin))
				} else if pin > 0 {
					if otherProp, ok := seenPin[pin]; ok {
						errs = append(errs, fmt.Errorf("property %q: duplicate pin value %d (already used by %q)", propName, pin, otherProp))
					}
					seenPin[pin] = propName
				}
			}
			continue
		}

		if prop.Name == "" {
			errs = append(errs, fmt.Errorf("property[%d]: missing required field: name", i))
			continue
		}
		if IsSystemProperty(prop.Name) {
			errs = append(errs, fmt.Errorf("property %q: %q is a reserved system property and cannot be defined in type schemas", prop.Name, prop.Name))
			continue
		}
		// Check if local property name conflicts with a shared property name
		if sharedMap != nil {
			if _, ok := sharedMap[prop.Name]; ok {
				errs = append(errs, fmt.Errorf("property %q: conflicts with a shared property name", prop.Name))
			}
		}
		if seen[prop.Name] {
			errs = append(errs, fmt.Errorf("property %q: duplicate property name", prop.Name))
		}
		seen[prop.Name] = true
		if prop.Emoji != "" {
			if otherProp, ok := seenEmoji[prop.Emoji]; ok {
				errs = append(errs, fmt.Errorf("property %q: duplicate property emoji %q (already used by %q)", prop.Name, prop.Emoji, otherProp))
			}
			seenEmoji[prop.Emoji] = prop.Name
		}
		if prop.Pin < 0 {
			errs = append(errs, fmt.Errorf("property %q: pin value must be a positive integer, got %d", prop.Name, prop.Pin))
		} else if prop.Pin > 0 {
			if otherProp, ok := seenPin[prop.Pin]; ok {
				errs = append(errs, fmt.Errorf("property %q: duplicate pin value %d (already used by %q)", prop.Name, prop.Pin, otherProp))
			}
			seenPin[prop.Pin] = prop.Name
		}
		typeErrs := validatePropertyType(prop, fmt.Sprintf("property %q", prop.Name))
		errs = append(errs, typeErrs...)
		if len(typeErrs) > 0 && (prop.Type == "" || prop.Type == "enum" || !validPropertyTypes[prop.Type]) {
			continue
		}
	}
	return errs
}

// validateUseOverrides checks that a `use` property entry only has allowed override fields.
func validateUseOverrides(index int, prop Property) error {
	// Only pin and emoji overrides are allowed on use entries.
	disallowed := []struct {
		fieldName string
		isSet     bool
	}{
		{"type", prop.Type != ""},
		{"options", len(prop.Options) > 0},
		{"target", prop.Target != ""},
		{"default", prop.Default != nil},
		{"multiple", prop.Multiple},
		{"bidirectional", prop.Bidirectional},
		{"inverse", prop.Inverse != ""},
	}

	for _, f := range disallowed {
		if f.isSet {
			return fmt.Errorf("property[%d] use %q: only \"pin\" and \"emoji\" overrides are allowed on \"use\" entries, got %q", index, prop.Use, f.fieldName)
		}
	}
	return nil
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
		case "number":
			switch val.(type) {
			case int, int64, float64:
				// valid
			default:
				errs = append(errs, fmt.Errorf("property %q: expected number, got %T", prop.Name, val))
			}
		case "date":
			errs = append(errs, validateDate(prop.Name, val)...)
		case "datetime":
			errs = append(errs, validateDatetime(prop.Name, val)...)
		case "url":
			errs = append(errs, validateURL(prop.Name, val)...)
		case "checkbox":
			if _, ok := val.(bool); !ok {
				errs = append(errs, fmt.Errorf("property %q: expected boolean, got %T", prop.Name, val))
			}
		case "select":
			errs = append(errs, validateSelect(prop, val)...)
		case "multi_select":
			errs = append(errs, validateMultiSelect(prop, val)...)
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
		}
	}

	return errs
}

// validateDate validates a date property value (YYYY-MM-DD format).
// Handles time.Time values from YAML auto-parsing.
func validateDate(name string, val any) []error {
	switch v := val.(type) {
	case time.Time:
		return nil // YAML auto-parsed date
	case string:
		if !dateRegexp.MatchString(v) {
			return []error{fmt.Errorf("property %q: expected date in YYYY-MM-DD format, got %q", name, v)}
		}
		if _, err := time.Parse("2006-01-02", v); err != nil {
			return []error{fmt.Errorf("property %q: invalid date %q: %v", name, v, err)}
		}
		return nil
	default:
		return []error{fmt.Errorf("property %q: expected date string or time.Time, got %T", name, val)}
	}
}

// validateDatetime validates a datetime property value (ISO 8601 with time).
// Handles time.Time values from YAML auto-parsing.
func validateDatetime(name string, val any) []error {
	switch v := val.(type) {
	case time.Time:
		return nil // YAML auto-parsed datetime
	case string:
		formats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05",
			"2006-01-02T15:04:05Z",
		}
		for _, f := range formats {
			if _, err := time.Parse(f, v); err == nil {
				return nil
			}
		}
		return []error{fmt.Errorf("property %q: expected datetime in ISO 8601 format (e.g. 2006-01-02T15:04:05), got %q", name, v)}
	default:
		return []error{fmt.Errorf("property %q: expected datetime string or time.Time, got %T", name, val)}
	}
}

// validateURL validates a url property value (must have http:// or https:// scheme).
func validateURL(name string, val any) []error {
	s, ok := val.(string)
	if !ok {
		return []error{fmt.Errorf("property %q: expected string for url, got %T", name, val)}
	}
	u, err := url.Parse(s)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return []error{fmt.Errorf("property %q: url must start with http:// or https://, got %q", name, s)}
	}
	return nil
}

// validateSelect validates a select property value against options.
func validateSelect(prop Property, val any) []error {
	s, ok := val.(string)
	if !ok {
		return []error{fmt.Errorf("property %q: expected string for select, got %T", prop.Name, val)}
	}
	for _, opt := range prop.Options {
		if opt.Value == s {
			return nil
		}
	}
	return []error{fmt.Errorf("property %q: value %q not in allowed options %v", prop.Name, s, prop.OptionValues())}
}

// validateMultiSelect validates a multi_select property value.
// Accepts a list of values (each must be in options). Coerces a single string to a list.
func validateMultiSelect(prop Property, val any) []error {
	var items []string

	switch v := val.(type) {
	case string:
		items = []string{v}
	case []any:
		for i, item := range v {
			s, ok := item.(string)
			if !ok {
				return []error{fmt.Errorf("property %q[%d]: expected string, got %T", prop.Name, i, item)}
			}
			items = append(items, s)
		}
	case []string:
		items = v
	default:
		return []error{fmt.Errorf("property %q: expected string or array for multi_select, got %T", prop.Name, val)}
	}

	allowed := make(map[string]bool, len(prop.Options))
	for _, opt := range prop.Options {
		allowed[opt.Value] = true
	}

	var errs []error
	optionVals := prop.OptionValues()
	for _, item := range items {
		if !allowed[item] {
			errs = append(errs, fmt.Errorf("property %q: value %q not in allowed options %v", prop.Name, item, optionVals))
		}
	}
	return errs
}

// OptionValues returns a slice of option values (convenience for display/error messages).
func (p Property) OptionValues() []string {
	vals := make([]string, len(p.Options))
	for i, opt := range p.Options {
		vals[i] = opt.Value
	}
	return vals
}

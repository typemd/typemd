package core

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

// NameProperty is the reserved system property key for the object's display name.
const NameProperty = "name"

// AmbiguousMatchError is returned when a prefix matches multiple objects.
type AmbiguousMatchError struct {
	Prefix  string
	Matches []string
}

func (e *AmbiguousMatchError) Error() string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "ambiguous object ID %q matches %d objects:", e.Prefix, len(e.Matches))
	for _, m := range e.Matches {
		fmt.Fprintf(&buf, "\n  %s", m)
	}
	return buf.String()
}

// Object represents a typemd object.
type Object struct {
	ID         string
	Type       string
	Filename   string
	Properties map[string]any
	Body       string
}

// ParseID returns this object's ID as a parsed ObjectID Value Object.
func (o *Object) ParseID() ObjectID {
	return ObjectID{Type: o.Type, Filename: o.Filename}
}

// DisplayName returns the filename with ULID suffix stripped.
func (o *Object) DisplayName() string {
	return o.ParseID().DisplayName()
}

// GetName returns the object's display name from the NameProperty.
// If the property is missing or empty, it falls back to DisplayName().
func (o *Object) GetName() string {
	if name, ok := o.Properties[NameProperty].(string); ok && name != "" {
		return name
	}
	return o.DisplayName()
}

// DisplayID returns the object ID with ULID suffix stripped from the filename part.
func (o *Object) DisplayID() string {
	return o.ParseID().DisplayID()
}

// MarkUpdated sets the updated_at timestamp to now.
func (o *Object) MarkUpdated() {
	o.Properties[UpdatedAtProperty] = time.Now().Format(time.RFC3339)
}

// Validate checks this object's properties against the given type schema.
func (o *Object) Validate(schema *TypeSchema) []error {
	return ValidateObject(o.Properties, schema)
}

// SetProperty validates and sets a single property value, returning a domain event.
func (o *Object) SetProperty(key string, value any, schema *TypeSchema) (DomainEvent, error) {
	testProps := map[string]any{key: value}
	if errs := ValidateObject(testProps, schema); len(errs) > 0 {
		return nil, errs[0]
	}
	old := o.Properties[key]
	o.Properties[key] = value
	return PropertyChanged{ObjectID: o.ID, Key: key, Old: old, New: value}, nil
}

// LinkTo appends a relation target to this object's properties.
// For single-value relations, it overwrites. For multi-value, it appends.
// Returns an event and error; error if the target already exists in a multi-value relation.
func (o *Object) LinkTo(relName, targetID string, prop *Property) (DomainEvent, error) {
	if err := appendRelationValue(o.Properties, relName, targetID, prop.Multiple); err != nil {
		return nil, err
	}
	return ObjectLinked{FromID: o.ID, ToID: targetID, RelName: relName}, nil
}

// Unlink removes a relation target from this object's properties.
func (o *Object) Unlink(relName, targetID string, prop *Property) DomainEvent {
	removeRelationValue(o.Properties, relName, targetID, prop.Multiple)
	return ObjectUnlinked{FromID: o.ID, ToID: targetID, RelName: relName}
}

// ApplyTemplate applies template properties and body to this object.
// It filters template properties against the schema and skips immutable system properties.
func (o *Object) ApplyTemplate(tmpl *Template, schema *TypeSchema) {
	filtered := filterTemplateProperties(tmpl.Properties, schema)
	for key, val := range filtered {
		o.Properties[key] = val
	}
	o.Body = tmpl.Body
}

// writeFrontmatter writes properties and body as markdown with YAML frontmatter.
// keyOrder specifies the desired property output order. Properties not in keyOrder
// are appended at the end. If keyOrder is nil, properties are written in map iteration order.
func writeFrontmatter(props map[string]any, body string, keyOrder []string) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("---\n")
	if len(props) > 0 {
		if len(keyOrder) > 0 {
			written := make(map[string]bool)
			for _, key := range keyOrder {
				if val, ok := props[key]; ok {
					entry := map[string]any{key: val}
					yamlData, err := yaml.Marshal(entry)
					if err != nil {
						return nil, err
					}
					buf.Write(yamlData)
					written[key] = true
				}
			}
			// Write remaining keys not in keyOrder
			for key, val := range props {
				if !written[key] {
					entry := map[string]any{key: val}
					yamlData, err := yaml.Marshal(entry)
					if err != nil {
						return nil, err
					}
					buf.Write(yamlData)
				}
			}
		} else {
			yamlData, err := yaml.Marshal(props)
			if err != nil {
				return nil, err
			}
			buf.Write(yamlData)
		}
	}
	buf.WriteString("---\n")
	if body != "" {
		buf.WriteString("\n")
		buf.WriteString(body)
	}
	return buf.Bytes(), nil
}

// parseFrontmatter parses YAML frontmatter from markdown content.
func parseFrontmatter(data []byte) (map[string]any, string, error) {
	var props map[string]any
	rest, err := frontmatter.Parse(bytes.NewReader(data), &props)
	if err != nil {
		return nil, "", err
	}
	if props == nil {
		props = make(map[string]any)
	}
	body := strings.TrimLeft(string(rest), "\n")
	return props, body, nil
}

// NewObject creates a new object. Delegates to ObjectService.
func (v *Vault) NewObject(typeName, filename, templateName string) (*Object, error) {
	if v.Objects == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	return v.Objects.Create(typeName, filename, templateName)
}

// SetProperty updates a single property. Delegates to ObjectService.
func (v *Vault) SetProperty(id, key string, value any) error {
	if v.Objects == nil {
		return fmt.Errorf("vault not opened")
	}
	return v.Objects.SetProperty(id, key, value)
}

// GetObject reads an object by ID. Delegates to QueryService.
func (v *Vault) GetObject(id string) (*Object, error) {
	return v.Queries.Get(id)
}

// ResolveID resolves an abbreviated object ID. Delegates to QueryService.
func (v *Vault) ResolveID(prefix string) (string, error) {
	return v.Queries.Resolve(prefix)
}

// ResolveObject resolves a prefix and returns the object.
func (v *Vault) ResolveObject(prefix string) (*Object, error) {
	id, err := v.ResolveID(prefix)
	if err != nil {
		return nil, err
	}
	return v.GetObject(id)
}

// SaveObject persists an object. Delegates to ObjectService.
func (v *Vault) SaveObject(obj *Object) error {
	if v.Objects == nil {
		return fmt.Errorf("vault not opened")
	}
	return v.Objects.Save(obj)
}


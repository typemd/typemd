package core

import (
	"bytes"
	"encoding/json"
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

// DisplayName returns the filename with ULID suffix stripped.
func (o *Object) DisplayName() string {
	return StripULID(o.Filename)
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
	return o.Type + "/" + o.DisplayName()
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

// NewObject creates a new object with the given type and filename.
// If templateName is non-empty, the specified template is loaded and applied.
func (v *Vault) NewObject(typeName, filename, templateName string) (*Object, error) {
	if v.index == nil {
		return nil, fmt.Errorf("vault not opened")
	}

	schema, err := v.LoadType(typeName)
	if err != nil {
		return nil, fmt.Errorf("load type: %w", err)
	}

	// Load template if specified
	var tmpl *Template
	if templateName != "" {
		tmpl, err = v.LoadTemplate(typeName, templateName)
		if err != nil {
			return nil, fmt.Errorf("load template: %w", err)
		}
	}

	now := time.Now()

	// Handle empty name: check template, then name template, then error
	if filename == "" {
		if tmpl != nil {
			if nameVal, ok := tmpl.Properties[NameProperty]; ok {
				if s, ok := nameVal.(string); ok && s != "" {
					filename = s
				}
			}
		}
		if filename == "" {
			if schema.NameTemplate != "" {
				filename = EvaluateNameTemplate(schema.NameTemplate, now)
			} else {
				return nil, fmt.Errorf("name is required (type %q has no name template)", typeName)
			}
		}
	}

	// Enforce name uniqueness for types with unique constraint
	if schema.Unique {
		if err := v.checkNameUnique(typeName, filename); err != nil {
			return nil, err
		}
	}

	// Append ULID to filename for uniqueness
	slug := filename // preserve original slug for the name property
	ulidStr, err := GenerateULID()
	if err != nil {
		return nil, err
	}
	filename = slug + "-" + ulidStr
	id := typeName + "/" + filename

	// Create type directory
	if err := v.repo.EnsureDir(typeName); err != nil {
		return nil, fmt.Errorf("create directory: %w", err)
	}

	// Generate initial properties from schema defaults
	props := make(map[string]any)
	props[NameProperty] = slug
	nowStr := now.Format(time.RFC3339)
	props[CreatedAtProperty] = nowStr
	props[UpdatedAtProperty] = nowStr
	for _, p := range schema.Properties {
		if p.Default != nil {
			props[p.Name] = p.Default
		} else {
			props[p.Name] = nil
		}
	}

	// Apply template properties (overrides schema defaults)
	body := ""
	if tmpl != nil {
		filtered := filterTemplateProperties(tmpl.Properties, schema)
		for key, val := range filtered {
			props[key] = val
		}
		body = tmpl.Body
	}

	// Write .md file via repository (O_EXCL ensures atomic uniqueness check)
	newObj := &Object{
		ID:         id,
		Type:       typeName,
		Filename:   filename,
		Properties: props,
		Body:       body,
	}
	if err := v.repo.Create(newObj, OrderedPropKeys(props, schema)); err != nil {
		return nil, fmt.Errorf("create object file: %w", err)
	}

	// Insert into index
	propsJSON, err := json.Marshal(props)
	if err != nil {
		return nil, fmt.Errorf("marshal properties: %w", err)
	}
	if err := v.index.Upsert(id, typeName, filename, string(propsJSON), body); err != nil {
		return nil, fmt.Errorf("insert object: %w", err)
	}

	return newObj, nil
}

// saveObjectFile writes object properties to both .md file and index.
func (v *Vault) saveObjectFile(obj *Object) error {
	obj.Properties[UpdatedAtProperty] = time.Now().Format(time.RFC3339)
	// LoadType may fail for unknown types; nil schema is safe here because
	// OrderedPropKeys falls back to map iteration order when schema is nil.
	schema, _ := v.LoadType(obj.Type)
	keyOrder := OrderedPropKeys(obj.Properties, schema)

	if err := v.repo.Save(obj, keyOrder); err != nil {
		return fmt.Errorf("save object file: %w", err)
	}

	propsJSON, err := json.Marshal(obj.Properties)
	if err != nil {
		return fmt.Errorf("marshal properties: %w", err)
	}
	if err := v.index.Upsert(obj.ID, obj.Type, obj.Filename, string(propsJSON), obj.Body); err != nil {
		return fmt.Errorf("update index: %w", err)
	}

	return nil
}

// SetProperty updates a single property on an object.
func (v *Vault) SetProperty(id, key string, value any) error {
	if v.index == nil {
		return fmt.Errorf("vault not opened")
	}

	obj, err := v.GetObject(id)
	if err != nil {
		return fmt.Errorf("get object: %w", err)
	}

	schema, err := v.LoadType(obj.Type)
	if err != nil {
		return fmt.Errorf("load type: %w", err)
	}

	testProps := map[string]any{key: value}
	if errs := ValidateObject(testProps, schema); len(errs) > 0 {
		return errs[0]
	}

	obj.Properties[key] = value

	return v.saveObjectFile(obj)
}

// GetObject reads an object from its Markdown file.
func (v *Vault) GetObject(id string) (*Object, error) {
	return v.repo.Get(id)
}

// ResolveID resolves a (possibly abbreviated) object ID to the full ID.
// The input must be in "type/name" format. Exact matches take priority.
// If no exact match, a filesystem glob is used to find prefix matches.
func (v *Vault) ResolveID(prefix string) (string, error) {
	parts := strings.SplitN(prefix, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("invalid object ID format: %q", prefix)
	}
	typeName, namePrefix := parts[0], parts[1]

	// 1. Exact match — try to get the object directly
	if _, err := v.repo.ModTime(prefix); err == nil {
		return prefix, nil
	}

	// 2. Glob for prefix matches
	matches, err := v.repo.GlobIDs(typeName, namePrefix)
	if err != nil {
		return "", fmt.Errorf("glob error: %w", err)
	}

	switch len(matches) {
	case 0:
		return "", fmt.Errorf("no object found matching %q", prefix)
	case 1:
		return matches[0], nil
	default:
		return "", &AmbiguousMatchError{Prefix: prefix, Matches: matches}
	}
}

// ResolveObject resolves a prefix to a full ID and returns the object.
func (v *Vault) ResolveObject(prefix string) (*Object, error) {
	id, err := v.ResolveID(prefix)
	if err != nil {
		return nil, err
	}
	return v.GetObject(id)
}

// SaveObject persists an object's properties and body to the .md file and updates SQLite.
func (v *Vault) SaveObject(obj *Object) error {
	if v.index == nil {
		return fmt.Errorf("vault not opened")
	}
	return v.saveObjectFile(obj)
}


package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
func (v *Vault) NewObject(typeName, filename string) (*Object, error) {
	if v.db == nil {
		return nil, fmt.Errorf("vault not opened")
	}

	schema, err := v.LoadType(typeName)
	if err != nil {
		return nil, fmt.Errorf("load type: %w", err)
	}

	now := time.Now()

	// Handle empty name: use template or error
	if filename == "" {
		if schema.NameTemplate != "" {
			filename = EvaluateNameTemplate(schema.NameTemplate, now)
		} else {
			return nil, fmt.Errorf("name is required (type %q has no name template)", typeName)
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
	objPath := v.ObjectPath(typeName, filename)

	// Create type directory
	if err := os.MkdirAll(v.ObjectDir(typeName), 0755); err != nil {
		return nil, fmt.Errorf("create directory: %w", err)
	}

	// Generate initial properties from schema
	props := make(map[string]any)
	props[NameProperty] = slug // system property: display name from slug
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

	// Write .md file (O_EXCL ensures atomic uniqueness check)
	data, err := writeFrontmatter(props, "", OrderedPropKeys(props, schema))
	if err != nil {
		return nil, fmt.Errorf("write frontmatter: %w", err)
	}
	f, err := os.OpenFile(objPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf("object already exists: %s", id)
		}
		return nil, fmt.Errorf("create file: %w", err)
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return nil, fmt.Errorf("write file: %w", err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("close file: %w", err)
	}

	// Insert into SQLite
	propsJSON, err := json.Marshal(props)
	if err != nil {
		return nil, fmt.Errorf("marshal properties: %w", err)
	}
	_, err = v.db.Exec(
		"INSERT INTO objects (id, type, filename, properties, body) VALUES (?, ?, ?, ?, ?)",
		id, typeName, filename, string(propsJSON), "",
	)
	if err != nil {
		return nil, fmt.Errorf("insert object: %w", err)
	}

	return &Object{
		ID:         id,
		Type:       typeName,
		Filename:   filename,
		Properties: props,
		Body:       "",
	}, nil
}

// saveObjectFile writes object properties to both .md file and SQLite.
func (v *Vault) saveObjectFile(obj *Object) error {
	obj.Properties[UpdatedAtProperty] = time.Now().Format(time.RFC3339)
	schema, _ := v.LoadType(obj.Type)
	data, err := writeFrontmatter(obj.Properties, obj.Body, OrderedPropKeys(obj.Properties, schema))
	if err != nil {
		return fmt.Errorf("write frontmatter: %w", err)
	}
	objPath := v.ObjectPath(obj.Type, obj.Filename)
	if err := os.WriteFile(objPath, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	propsJSON, err := json.Marshal(obj.Properties)
	if err != nil {
		return fmt.Errorf("marshal properties: %w", err)
	}
	_, err = v.db.Exec(
		"UPDATE objects SET properties = ?, body = ? WHERE id = ?",
		string(propsJSON), obj.Body, obj.ID,
	)
	if err != nil {
		return fmt.Errorf("update object: %w", err)
	}

	return nil
}

// SetProperty updates a single property on an object.
func (v *Vault) SetProperty(id, key string, value any) error {
	if v.db == nil {
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
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid object ID: %s", id)
	}
	typeName, filename := parts[0], parts[1]

	objPath := v.ObjectPath(typeName, filename)
	data, err := os.ReadFile(objPath)
	if err != nil {
		return nil, fmt.Errorf("read object %s: %w", id, err)
	}

	props, body, err := parseFrontmatter(data)
	if err != nil {
		return nil, fmt.Errorf("parse object %s: %w", id, err)
	}

	return &Object{
		ID:         id,
		Type:       typeName,
		Filename:   filename,
		Properties: props,
		Body:       body,
	}, nil
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

	// 1. Exact match
	exactPath := v.ObjectPath(typeName, namePrefix)
	if _, err := os.Stat(exactPath); err == nil {
		return prefix, nil
	}

	// 2. Glob for prefix matches
	pattern := filepath.Join(v.ObjectDir(typeName), namePrefix+"*.md")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("glob error: %w", err)
	}

	switch len(matches) {
	case 0:
		return "", fmt.Errorf("no object found matching %q", prefix)
	case 1:
		// Extract ID from path: strip dir prefix and .md suffix
		base := filepath.Base(matches[0])
		filename := strings.TrimSuffix(base, ".md")
		return typeName + "/" + filename, nil
	default:
		candidates := make([]string, len(matches))
		for i, m := range matches {
			base := filepath.Base(m)
			filename := strings.TrimSuffix(base, ".md")
			candidates[i] = typeName + "/" + filename
		}
		return "", &AmbiguousMatchError{Prefix: prefix, Matches: candidates}
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
	if v.db == nil {
		return fmt.Errorf("vault not opened")
	}
	return v.saveObjectFile(obj)
}


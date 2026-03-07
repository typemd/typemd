package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

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

// schemaKeyOrder returns property names in schema-defined order.
func schemaKeyOrder(schema *TypeSchema) []string {
	keys := make([]string, len(schema.Properties))
	for i, p := range schema.Properties {
		keys[i] = p.Name
	}
	return keys
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

	// Append ULID to filename for uniqueness
	ulidStr, err := GenerateULID()
	if err != nil {
		return nil, err
	}
	filename = filename + "-" + ulidStr
	id := typeName + "/" + filename
	objPath := v.ObjectPath(typeName, filename)

	// Create type directory
	if err := os.MkdirAll(v.ObjectDir(typeName), 0755); err != nil {
		return nil, fmt.Errorf("create directory: %w", err)
	}

	// Generate initial properties from schema
	props := make(map[string]any)
	for _, p := range schema.Properties {
		if p.Default != nil {
			props[p.Name] = p.Default
		} else {
			props[p.Name] = nil
		}
	}

	// Write .md file (O_EXCL ensures atomic uniqueness check)
	data, err := writeFrontmatter(props, "", schemaKeyOrder(schema))
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

// writeObjectProperties writes object properties to both .md file and SQLite.
func (v *Vault) writeObjectProperties(obj *Object) error {
	var keyOrder []string
	if schema, err := v.LoadType(obj.Type); err == nil {
		keyOrder = schemaKeyOrder(schema)
	}
	data, err := writeFrontmatter(obj.Properties, obj.Body, keyOrder)
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

	return v.writeObjectProperties(obj)
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

// SaveObject persists an object's properties and body to the .md file and updates SQLite.
func (v *Vault) SaveObject(obj *Object) error {
	if v.db == nil {
		return fmt.Errorf("vault not opened")
	}
	return v.writeObjectProperties(obj)
}


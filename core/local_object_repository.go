package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// LocalObjectRepository implements ObjectRepository using the local filesystem.
// It encapsulates all path conventions and file I/O for vault entities.
type LocalObjectRepository struct {
	root string

	// Shared properties cache
	sharedProperties  []Property
	sharedPropsMap    map[string]Property
	sharedPropsLoaded bool
}

// NewLocalObjectRepository creates a new LocalObjectRepository rooted at the given directory.
func NewLocalObjectRepository(root string) *LocalObjectRepository {
	return &LocalObjectRepository{root: root}
}

// --- Path conventions (private) ---

func (r *LocalObjectRepository) vaultDir() string {
	return filepath.Join(r.root, ".typemd")
}

func (r *LocalObjectRepository) typesDir() string {
	return filepath.Join(r.vaultDir(), "types")
}

func (r *LocalObjectRepository) sharedPropertiesPath() string {
	return filepath.Join(r.vaultDir(), "properties.yaml")
}

func (r *LocalObjectRepository) objectsDir() string {
	return filepath.Join(r.root, "objects")
}

func (r *LocalObjectRepository) objectDir(typeName string) string {
	return filepath.Join(r.objectsDir(), typeName)
}

func (r *LocalObjectRepository) objectPath(typeName, filename string) string {
	return filepath.Join(r.objectDir(typeName), filename+".md")
}

func (r *LocalObjectRepository) templatesDir() string {
	return filepath.Join(r.root, "templates")
}

func (r *LocalObjectRepository) typeTemplatesDir(typeName string) string {
	return filepath.Join(r.templatesDir(), typeName)
}

func (r *LocalObjectRepository) templatePath(typeName, name string) string {
	return filepath.Join(r.typeTemplatesDir(typeName), name+".md")
}

// --- Object entity operations ---

// Get reads and parses an object file by ID, returning a fully populated domain entity.
func (r *LocalObjectRepository) Get(id string) (*Object, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid object ID: %s", id)
	}
	typeName, filename := parts[0], parts[1]

	objPath := r.objectPath(typeName, filename)
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

// Save serializes an object entity to its .md file.
func (r *LocalObjectRepository) Save(obj *Object, keyOrder []string) error {
	data, err := writeFrontmatter(obj.Properties, obj.Body, keyOrder)
	if err != nil {
		return fmt.Errorf("write frontmatter: %w", err)
	}
	objPath := r.objectPath(obj.Type, obj.Filename)
	if err := os.WriteFile(objPath, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

// Create writes a new object file with exclusive creation semantics (fails if file exists).
func (r *LocalObjectRepository) Create(obj *Object, keyOrder []string) error {
	data, err := writeFrontmatter(obj.Properties, obj.Body, keyOrder)
	if err != nil {
		return fmt.Errorf("write frontmatter: %w", err)
	}
	objPath := r.objectPath(obj.Type, obj.Filename)
	f, err := os.OpenFile(objPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("object already exists: %s", obj.ID)
		}
		return fmt.Errorf("create file: %w", err)
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return fmt.Errorf("write file: %w", err)
	}
	return f.Close()
}

// Walk traverses all object files and returns parsed domain entities.
// Unparseable files are skipped silently.
func (r *LocalObjectRepository) Walk() ([]*Object, error) {
	objsDir := r.objectsDir()
	if _, err := os.Stat(objsDir); os.IsNotExist(err) {
		return nil, nil
	}

	objects := make([]*Object, 0) // non-nil: directory exists but may be empty
	err := filepath.Walk(objsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		rel, err := filepath.Rel(objsDir, path)
		if err != nil {
			return nil
		}

		parts := strings.SplitN(rel, string(os.PathSeparator), 2)
		if len(parts) != 2 {
			return nil
		}
		typeName := parts[0]
		filename := strings.TrimSuffix(parts[1], ".md")
		id := typeName + "/" + filename

		data, err := os.ReadFile(path)
		if err != nil {
			return nil // skip unreadable files
		}

		props, body, err := parseFrontmatter(data)
		if err != nil {
			return nil // skip unparseable files
		}

		objects = append(objects, &Object{
			ID:         id,
			Type:       typeName,
			Filename:   filename,
			Properties: props,
			Body:       body,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk objects: %w", err)
	}

	return objects, nil
}

// GlobIDs finds object IDs matching a prefix pattern within a type directory.
func (r *LocalObjectRepository) GlobIDs(typeName, prefix string) ([]string, error) {
	pattern := filepath.Join(r.objectDir(typeName), prefix+"*.md")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob error: %w", err)
	}

	ids := make([]string, len(matches))
	for i, m := range matches {
		base := filepath.Base(m)
		filename := strings.TrimSuffix(base, ".md")
		ids[i] = typeName + "/" + filename
	}
	return ids, nil
}

// ModTime returns the last modification time of an object file.
func (r *LocalObjectRepository) ModTime(id string) (time.Time, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid object ID: %s", id)
	}
	objPath := r.objectPath(parts[0], parts[1])
	info, err := os.Stat(objPath)
	if err != nil {
		return time.Time{}, fmt.Errorf("stat object %s: %w", id, err)
	}
	return info.ModTime(), nil
}

// EnsureDir creates the type's object directory if it doesn't exist.
func (r *LocalObjectRepository) EnsureDir(typeName string) error {
	return os.MkdirAll(r.objectDir(typeName), 0755)
}

// --- Type schema operations ---

// GetSchema loads a type schema by name, resolving shared property references.
func (r *LocalObjectRepository) GetSchema(name string) (*TypeSchema, error) {
	path := filepath.Join(r.typesDir(), name+".yaml")

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

		// Resolve use entries if any exist
		if err := r.resolveSchemaUseEntries(&schema); err != nil {
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
	path := filepath.Join(r.typesDir(), typeName+".yaml")
	return os.WriteFile(path, data, 0644)
}

// DeleteSchema removes a type schema YAML file.
func (r *LocalObjectRepository) DeleteSchema(typeName string) error {
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
func (r *LocalObjectRepository) ListSchemas() ([]string, error) {
	seen := make(map[string]bool)

	// Custom types from YAML files
	entries, err := os.ReadDir(r.typesDir())
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".yaml") {
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

// --- Template operations ---

// GetTemplate reads and parses a template file.
func (r *LocalObjectRepository) GetTemplate(typeName, name string) (*Template, error) {
	path := r.templatePath(typeName, name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("template %q not found for type %q", name, typeName)
		}
		return nil, fmt.Errorf("read template: %w", err)
	}

	// Handle template files with no frontmatter delimiter
	content := string(data)
	if len(data) > 0 && !strings.HasPrefix(content, "---") {
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		return &Template{
			Name:       name,
			Properties: make(map[string]any),
			Body:       content,
		}, nil
	}

	props, body, err := parseFrontmatter(data)
	if err != nil {
		return nil, fmt.Errorf("parse template frontmatter: %w", err)
	}

	return &Template{
		Name:       name,
		Properties: props,
		Body:       body,
	}, nil
}

// ListTemplates returns the names of all templates available for the given type.
func (r *LocalObjectRepository) ListTemplates(typeName string) ([]string, error) {
	dir := r.typeTemplatesDir(typeName)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read templates directory: %w", err)
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if name, ok := strings.CutSuffix(e.Name(), ".md"); ok {
			names = append(names, name)
		}
	}
	return names, nil
}

// --- Shared property operations ---

// GetSharedProperties loads shared property definitions, with caching.
func (r *LocalObjectRepository) GetSharedProperties() ([]Property, error) {
	if r.sharedPropsLoaded {
		return r.sharedProperties, nil
	}

	data, err := os.ReadFile(r.sharedPropertiesPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			r.sharedProperties = nil
			r.sharedPropsLoaded = true
			return nil, nil
		}
		return nil, fmt.Errorf("read shared properties: %w", err)
	}

	var file SharedPropertiesFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse shared properties: %w", err)
	}

	r.sharedProperties = file.Properties
	r.sharedPropsMap = SharedPropertiesMap(file.Properties)
	r.sharedPropsLoaded = true
	return r.sharedProperties, nil
}

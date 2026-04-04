package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

func (r *LocalObjectRepository) typesDir() string {
	return filepath.Join(r.root, "types")
}

func (r *LocalObjectRepository) sharedPropertiesPath() string {
	return filepath.Join(r.root, "properties", "properties.yaml")
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
		if errors.Is(err, os.ErrExist) {
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

// CorruptedFile represents a file that could not be parsed.
type CorruptedFile struct {
	Path  string // relative path like "book/broken-01abc.md"
	Error error
}

// walkObjects is the shared implementation for Walk and WalkAll.
// When reportCorrupted is true, unparseable files are collected; otherwise they are skipped.
func (r *LocalObjectRepository) walkObjects(reportCorrupted bool) ([]*Object, []CorruptedFile, error) {
	objsDir := r.objectsDir()
	if _, err := os.Stat(objsDir); errors.Is(err, os.ErrNotExist) {
		return nil, nil, nil
	}

	objects := make([]*Object, 0)
	var corrupted []CorruptedFile

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

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			if reportCorrupted {
				corrupted = append(corrupted, CorruptedFile{Path: rel, Error: readErr})
			}
			return nil
		}

		// Detect files without frontmatter delimiters before parsing,
		// since the frontmatter library silently returns empty props for these.
		if reportCorrupted && !strings.HasPrefix(string(data), "---") {
			corrupted = append(corrupted, CorruptedFile{
				Path:  rel,
				Error: fmt.Errorf("missing frontmatter delimiters"),
			})
			return nil
		}

		props, body, parseErr := parseFrontmatter(data)
		if parseErr != nil {
			if reportCorrupted {
				corrupted = append(corrupted, CorruptedFile{Path: rel, Error: parseErr})
			}
			return nil
		}

		objects = append(objects, &Object{
			ID:         typeName + "/" + filename,
			Type:       typeName,
			Filename:   filename,
			Properties: props,
			Body:       body,
		})
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("walk objects: %w", err)
	}

	return objects, corrupted, nil
}

// Walk traverses all object files and returns parsed domain entities.
// Unparseable files are skipped silently.
func (r *LocalObjectRepository) Walk() ([]*Object, error) {
	objects, _, err := r.walkObjects(false)
	return objects, err
}

// WalkAll traverses all object files and returns both parsed entities and corrupted files.
// Unlike Walk(), unparseable files are reported as CorruptedFile entries instead of silently skipped.
func (r *LocalObjectRepository) WalkAll() ([]*Object, []CorruptedFile, error) {
	return r.walkObjects(true)
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

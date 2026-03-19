package core

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

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

// SaveTemplate writes a template file for the given type.
// If properties are non-empty, they are written as YAML frontmatter.
// The type template directory is created if it doesn't exist.
func (r *LocalObjectRepository) SaveTemplate(typeName, name string, tmpl *Template) error {
	dir := r.typeTemplatesDir(typeName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create template directory: %w", err)
	}

	var buf strings.Builder
	if len(tmpl.Properties) > 0 {
		buf.WriteString("---\n")
		yamlData, err := yaml.Marshal(tmpl.Properties)
		if err != nil {
			return fmt.Errorf("marshal template properties: %w", err)
		}
		buf.Write(yamlData)
		buf.WriteString("---\n")
		if tmpl.Body != "" {
			buf.WriteString("\n")
			buf.WriteString(tmpl.Body)
		}
	} else {
		buf.WriteString(tmpl.Body)
	}

	path := r.templatePath(typeName, name)
	if err := os.WriteFile(path, []byte(buf.String()), 0644); err != nil {
		return fmt.Errorf("write template file: %w", err)
	}
	return nil
}

// DeleteTemplate removes a template file for the given type.
// If the type template directory becomes empty after deletion, it is also removed.
func (r *LocalObjectRepository) DeleteTemplate(typeName, name string) error {
	path := r.templatePath(typeName, name)
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("template %q not found for type %q", name, typeName)
		}
		return fmt.Errorf("delete template: %w", err)
	}

	// Clean up empty type template directory
	dir := r.typeTemplatesDir(typeName)
	entries, err := os.ReadDir(dir)
	if err == nil && len(entries) == 0 {
		os.Remove(dir)
	}

	return nil
}

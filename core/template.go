package core

import (
	"fmt"
	"os"
	"strings"
)

// Template represents a parsed object template with frontmatter properties and body.
type Template struct {
	Name       string
	Properties map[string]any
	Body       string
}

// ListTemplates returns the names of all templates available for the given type.
// Template names are derived from filenames (without .md extension) in templates/<type>/.
// Returns an empty slice if no templates directory exists or it contains no .md files.
func (v *Vault) ListTemplates(typeName string) ([]string, error) {
	dir := v.TypeTemplatesDir(typeName)
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

// LoadTemplate reads and parses a template file for the given type and template name.
// Returns the parsed template with frontmatter properties and body content.
func (v *Vault) LoadTemplate(typeName, templateName string) (*Template, error) {
	path := v.TemplatePath(typeName, templateName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("template %q not found for type %q", templateName, typeName)
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
			Name:       templateName,
			Properties: make(map[string]any),
			Body:       content,
		}, nil
	}

	props, body, err := parseFrontmatter(data)
	if err != nil {
		return nil, fmt.Errorf("parse template frontmatter: %w", err)
	}

	return &Template{
		Name:       templateName,
		Properties: props,
		Body:       body,
	}, nil
}

// filterTemplateProperties returns only the template properties that should be applied.
// It excludes immutable system properties and properties not defined in the schema.
func filterTemplateProperties(tmplProps map[string]any, schema *TypeSchema) map[string]any {
	if len(tmplProps) == 0 {
		return nil
	}

	filtered := make(map[string]any)
	validProps := schema.PropertyNames()

	for key, val := range tmplProps {
		// Skip immutable system properties
		if IsImmutableSystemProperty(key) {
			continue
		}
		// Allow mutable system properties
		if IsSystemProperty(key) {
			filtered[key] = val
			continue
		}
		// Allow schema-defined properties
		if validProps[key] {
			filtered[key] = val
		}
	}

	return filtered
}

package core

// Template represents a parsed object template with frontmatter properties and body.
type Template struct {
	Name       string
	Properties map[string]any
	Body       string
}

// ListTemplates returns the names of all templates available for the given type.
func (v *Vault) ListTemplates(typeName string) ([]string, error) {
	return v.repo.ListTemplates(typeName)
}

// LoadTemplate reads and parses a template file for the given type and template name.
func (v *Vault) LoadTemplate(typeName, templateName string) (*Template, error) {
	return v.repo.GetTemplate(typeName, templateName)
}

// SaveTemplate writes a template file for the given type and template name.
func (v *Vault) SaveTemplate(typeName, templateName string, tmpl *Template) error {
	return v.repo.SaveTemplate(typeName, templateName, tmpl)
}

// DeleteTemplate removes a template file for the given type and template name.
func (v *Vault) DeleteTemplate(typeName, templateName string) error {
	return v.repo.DeleteTemplate(typeName, templateName)
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

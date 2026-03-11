package core

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadSharedProperties loads shared property definitions from .typemd/properties.yaml.
// Returns an empty slice if the file does not exist.
// Results are cached on the Vault for reuse across multiple LoadType() calls.
func (v *Vault) LoadSharedProperties() ([]Property, error) {
	if v.sharedPropsLoaded {
		return v.sharedProperties, nil
	}

	data, err := os.ReadFile(v.SharedPropertiesPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			v.sharedProperties = nil
			v.sharedPropsLoaded = true
			return nil, nil
		}
		return nil, fmt.Errorf("read shared properties: %w", err)
	}

	var file SharedPropertiesFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse shared properties: %w", err)
	}

	v.sharedProperties = file.Properties
	v.sharedPropsMap = SharedPropertiesMap(file.Properties)
	v.sharedPropsLoaded = true
	return v.sharedProperties, nil
}

// SharedPropertiesMap returns shared properties as a map keyed by name.
func SharedPropertiesMap(props []Property) map[string]Property {
	m := make(map[string]Property, len(props))
	for _, p := range props {
		m[p.Name] = p
	}
	return m
}

// ValidateSharedProperties validates the shared properties file for correctness.
func ValidateSharedProperties(props []Property) []error {
	var errs []error
	seen := make(map[string]bool)

	for i, prop := range props {
		if prop.Name == "" {
			errs = append(errs, fmt.Errorf("shared property[%d]: missing required field: name", i))
			continue
		}
		if IsSystemProperty(prop.Name) {
			errs = append(errs, fmt.Errorf("shared property %q: %q is a reserved system property", prop.Name, prop.Name))
			continue
		}
		if seen[prop.Name] {
			errs = append(errs, fmt.Errorf("shared property %q: duplicate shared property name", prop.Name))
		}
		seen[prop.Name] = true

		// Validate using the same rules as type schema properties
		typeErrs := validatePropertyType(prop, fmt.Sprintf("shared property %q", prop.Name))
		errs = append(errs, typeErrs...)
		if len(typeErrs) > 0 && (prop.Type == "" || prop.Type == "enum" || !validPropertyTypes[prop.Type]) {
			continue
		}
	}

	return errs
}

// resolveUseEntries resolves `use` references in a type schema's properties
// by replacing them with fully resolved Property objects from shared properties.
func resolveUseEntries(schema *TypeSchema, sharedMap map[string]Property) error {
	for i, prop := range schema.Properties {
		if prop.Use == "" {
			continue
		}
		shared, ok := sharedMap[prop.Use]
		if !ok {
			return fmt.Errorf("property use %q: shared property not found", prop.Use)
		}

		// Copy shared property and apply overrides
		resolved := shared
		resolved.Use = "" // Clear the Use field after resolution
		if prop.Pin != 0 {
			resolved.Pin = prop.Pin
		}
		if prop.Emoji != "" {
			resolved.Emoji = prop.Emoji
		}
		schema.Properties[i] = resolved
	}
	return nil
}

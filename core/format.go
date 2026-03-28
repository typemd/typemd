package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// FormatResult contains the results of a format operation.
type FormatResult struct {
	// Changed lists the paths of files that were (or would be) reformatted.
	Changed []string
}

// FormatAll formats both objects and schemas, returning combined results.
// If typeName is non-empty, only objects and schemas of that type are formatted.
// If dryRun is true, files are not modified; Changed lists what would change.
func (v *Vault) FormatAll(typeName string, dryRun bool) (*FormatResult, error) {
	// Validate type name if provided
	if typeName != "" {
		if _, err := v.LoadType(typeName); err != nil {
			return nil, fmt.Errorf("unknown type %q", typeName)
		}
	}

	objResult, err := v.FormatObjects(typeName, dryRun)
	if err != nil {
		return nil, err
	}

	schemaResult, err := v.FormatSchemas(typeName, dryRun)
	if err != nil {
		return nil, err
	}

	combined := &FormatResult{
		Changed: append(objResult.Changed, schemaResult.Changed...),
	}
	return combined, nil
}

// FormatObjects formats object files by re-serializing frontmatter with
// canonical property ordering and YAML style.
// If typeName is non-empty, only objects of that type are formatted.
// If dryRun is true, files are not modified; Changed lists what would change.
func (v *Vault) FormatObjects(typeName string, dryRun bool) (*FormatResult, error) {
	objects, err := v.repo.Walk()
	if err != nil {
		return nil, fmt.Errorf("walk objects: %w", err)
	}

	result := &FormatResult{}

	for _, obj := range objects {
		if typeName != "" && obj.Type != typeName {
			continue
		}

		schema, _ := v.LoadType(obj.Type)
		keyOrder := OrderedPropKeys(obj.Properties, schema)

		formatted, err := writeFrontmatter(obj.Properties, obj.Body, keyOrder)
		if err != nil {
			return nil, fmt.Errorf("format object %s: %w", obj.ID, err)
		}

		// Read original file for comparison
		objPath := v.ObjectPath(obj.Type, obj.Filename)
		original, err := os.ReadFile(objPath)
		if err != nil {
			return nil, fmt.Errorf("read object %s: %w", obj.ID, err)
		}

		if bytes.Equal(original, formatted) {
			continue
		}

		result.Changed = append(result.Changed, obj.ID)

		if !dryRun {
			if err := os.WriteFile(objPath, formatted, 0644); err != nil {
				return nil, fmt.Errorf("write object %s: %w", obj.ID, err)
			}
		}
	}

	return result, nil
}

// FormatSchemas formats type schema files by round-tripping through
// MarshalTypeSchema to produce canonical YAML output.
// If typeName is non-empty, only that schema is formatted.
// If dryRun is true, files are not modified; Changed lists what would change.
func (v *Vault) FormatSchemas(typeName string, dryRun bool) (*FormatResult, error) {
	names, err := v.repo.ListSchemas()
	if err != nil {
		return nil, fmt.Errorf("list schemas: %w", err)
	}

	result := &FormatResult{}

	for _, name := range names {
		if typeName != "" && name != typeName {
			continue
		}

		// Find the schema file path; skip built-in types without files
		schemaPath := v.schemaFilePath(name)
		if schemaPath == "" {
			continue
		}

		original, err := os.ReadFile(schemaPath)
		if err != nil {
			return nil, fmt.Errorf("read schema %s: %w", name, err)
		}

		schema, err := v.LoadType(name)
		if err != nil {
			return nil, fmt.Errorf("load schema %s: %w", name, err)
		}

		formatted, err := MarshalTypeSchema(schema)
		if err != nil {
			return nil, fmt.Errorf("marshal schema %s: %w", name, err)
		}

		if bytes.Equal(original, formatted) {
			continue
		}

		result.Changed = append(result.Changed, name)

		if !dryRun {
			if err := os.WriteFile(schemaPath, formatted, 0644); err != nil {
				return nil, fmt.Errorf("write schema %s: %w", name, err)
			}
		}
	}

	return result, nil
}

// schemaFilePath returns the path to the schema file for a type, or empty
// string if the type is built-in without a custom file.
func (v *Vault) schemaFilePath(name string) string {
	dirPath := filepath.Join(v.TypesDir(), name, "schema.yaml")
	if _, err := os.Stat(dirPath); err == nil {
		return dirPath
	}
	return ""
}


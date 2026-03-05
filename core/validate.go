package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ValidateAllObjects queries all objects and validates properties against their type schema.
// Returns a map of object ID to validation errors.
func ValidateAllObjects(v *Vault) map[string][]error {
	result := make(map[string][]error)
	objects, err := v.QueryObjects("")
	if err != nil {
		return result
	}
	for _, obj := range objects {
		schema, err := v.LoadType(obj.Type)
		if err != nil {
			result[obj.ID] = []error{fmt.Errorf("load type %q: %w", obj.Type, err)}
			continue
		}
		if errs := ValidateObject(obj.Properties, schema); len(errs) > 0 {
			result[obj.ID] = errs
		}
	}
	return result
}

// ValidateAllSchemas scans .typemd/types/*.yaml and validates each schema.
// Returns a map of type name to validation errors.
func ValidateAllSchemas(v *Vault) map[string][]error {
	result := make(map[string][]error)
	entries, err := os.ReadDir(v.TypesDir())
	if err != nil {
		return result
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		typeName := strings.TrimSuffix(entry.Name(), ".yaml")
		data, err := os.ReadFile(filepath.Join(v.TypesDir(), entry.Name()))
		if err != nil {
			result[typeName] = []error{fmt.Errorf("read file: %w", err)}
			continue
		}
		var schema TypeSchema
		if err := yaml.Unmarshal(data, &schema); err != nil {
			result[typeName] = []error{fmt.Errorf("parse YAML: %w", err)}
			continue
		}
		if errs := ValidateSchema(&schema); len(errs) > 0 {
			result[typeName] = errs
		}
	}
	return result
}

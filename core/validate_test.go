package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateAllSchemas_Valid(t *testing.T) {
	v := setupTestVault(t)
	schema := []byte("name: book\nproperties:\n  - name: title\n    type: string\n")
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), schema, 0644)

	result := ValidateAllSchemas(v)
	if errs, ok := result["book"]; ok && len(errs) > 0 {
		t.Errorf("expected no errors for book, got %v", errs)
	}
}

func TestValidateAllSchemas_Invalid(t *testing.T) {
	v := setupTestVault(t)
	schema := []byte("name: bad\nproperties:\n  - name: status\n    type: enum\n")
	os.WriteFile(filepath.Join(v.TypesDir(), "bad.yaml"), schema, 0644)

	result := ValidateAllSchemas(v)
	if errs, ok := result["bad"]; !ok || len(errs) == 0 {
		t.Error("expected validation errors for bad schema")
	}
}

func TestValidateAllSchemas_MalformedYAML(t *testing.T) {
	v := setupTestVault(t)
	os.WriteFile(filepath.Join(v.TypesDir(), "broken.yaml"), []byte(":\ninvalid yaml["), 0644)

	result := ValidateAllSchemas(v)
	if errs, ok := result["broken"]; !ok || len(errs) == 0 {
		t.Error("expected parse error for malformed YAML")
	}
}

func TestValidateAllSchemas_NoSchemas(t *testing.T) {
	v := setupTestVault(t)
	result := ValidateAllSchemas(v)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

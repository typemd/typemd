package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListTemplates_IgnoresNonMdFiles(t *testing.T) {
	v := setupTestVault(t)

	dir := v.TypeTemplatesDir("book")
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "review.md"), []byte("## Review\n"), 0644)
	os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("not a template"), 0644)
	os.WriteFile(filepath.Join(dir, ".hidden"), []byte("hidden file"), 0644)

	names, err := v.ListTemplates("book")
	if err != nil {
		t.Fatalf("ListTemplates() error = %v", err)
	}
	if len(names) != 1 {
		t.Fatalf("expected 1 template, got %d: %v", len(names), names)
	}
	if names[0] != "review" {
		t.Errorf("expected 'review', got %q", names[0])
	}
}

func TestListTemplates_IgnoresNestedDirectories(t *testing.T) {
	v := setupTestVault(t)

	dir := v.TypeTemplatesDir("book")
	os.MkdirAll(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "review.md"), []byte("## Review\n"), 0644)
	os.WriteFile(filepath.Join(dir, "subdir", "nested.md"), []byte("## Nested\n"), 0644)

	names, err := v.ListTemplates("book")
	if err != nil {
		t.Fatalf("ListTemplates() error = %v", err)
	}
	if len(names) != 1 {
		t.Fatalf("expected 1 template, got %d: %v", len(names), names)
	}
	if names[0] != "review" {
		t.Errorf("expected 'review', got %q", names[0])
	}
}

func TestFilterTemplateProperties_SkipsImmutable(t *testing.T) {
	schema := &TypeSchema{
		Properties: []Property{{Name: "status", Type: "string"}},
	}
	tmplProps := map[string]any{
		"status":     "draft",
		"created_at": "2020-01-01T00:00:00Z",
		"updated_at": "2020-01-01T00:00:00Z",
	}

	filtered := filterTemplateProperties(tmplProps, schema)
	if _, ok := filtered["created_at"]; ok {
		t.Error("created_at should be filtered out")
	}
	if _, ok := filtered["updated_at"]; ok {
		t.Error("updated_at should be filtered out")
	}
	if filtered["status"] != "draft" {
		t.Errorf("status = %v, want 'draft'", filtered["status"])
	}
}

func TestFilterTemplateProperties_AllowsMutableSystemProps(t *testing.T) {
	schema := &TypeSchema{}
	tmplProps := map[string]any{
		"name":        "custom-name",
		"description": "A description",
	}

	filtered := filterTemplateProperties(tmplProps, schema)
	if filtered["name"] != "custom-name" {
		t.Errorf("name = %v, want 'custom-name'", filtered["name"])
	}
	if filtered["description"] != "A description" {
		t.Errorf("description = %v, want 'A description'", filtered["description"])
	}
}

func TestFilterTemplateProperties_SkipsUnknownProps(t *testing.T) {
	schema := &TypeSchema{
		Properties: []Property{{Name: "status", Type: "string"}},
	}
	tmplProps := map[string]any{
		"status":  "draft",
		"unknown": "value",
	}

	filtered := filterTemplateProperties(tmplProps, schema)
	if _, ok := filtered["unknown"]; ok {
		t.Error("unknown property should be filtered out")
	}
	if filtered["status"] != "draft" {
		t.Errorf("status = %v, want 'draft'", filtered["status"])
	}
}

package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupFormatVault(t *testing.T) *Vault {
	t.Helper()
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { v.Close() })
	return v
}

func writeRawObject(t *testing.T, v *Vault, typeName, slug, content string) string {
	t.Helper()
	objDir := filepath.Join(v.ObjectsDir(), typeName)
	os.MkdirAll(objDir, 0755)
	ulid, _ := GenerateULID()
	filename := slug + "-" + ulid
	objPath := filepath.Join(objDir, filename+".md")
	if err := os.WriteFile(objPath, []byte(content), 0644); err != nil {
		t.Fatalf("write object: %v", err)
	}
	return typeName + "/" + filename
}

func TestFormatObjects_EmptyFrontmatter(t *testing.T) {
	v := setupFormatVault(t)

	writeRawObject(t, v, "page", "empty", "---\n---\n")

	result, err := v.FormatObjects("", false)
	if err != nil {
		t.Fatalf("FormatObjects: %v", err)
	}

	// Empty frontmatter should still be formatted (may add/remove newlines)
	// The key thing is it doesn't error
	_ = result
}

func TestFormatObjects_WithoutSchema(t *testing.T) {
	v := setupFormatVault(t)

	// Create an object for a type with no custom schema (page is built-in)
	content := "---\nname: Test\ncreated_at: \"2025-01-01T00:00:00Z\"\nupdated_at: \"2025-01-01T00:00:00Z\"\n---\n"
	writeRawObject(t, v, "page", "test", content)

	result, err := v.FormatObjects("", false)
	if err != nil {
		t.Fatalf("FormatObjects: %v", err)
	}

	// Should succeed with built-in type schema
	_ = result
}

func TestFormatObjects_NilPropertyValue(t *testing.T) {
	v := setupFormatVault(t)

	v.SaveType(&TypeSchema{
		Name:       "book",
		Properties: []Property{{Name: "author", Type: "string"}},
	})

	content := "---\nname: Test\nauthor: null\ncreated_at: \"2025-01-01T00:00:00Z\"\nupdated_at: \"2025-01-01T00:00:00Z\"\n---\n"
	writeRawObject(t, v, "book", "test", content)

	result, err := v.FormatObjects("", false)
	if err != nil {
		t.Fatalf("FormatObjects: %v", err)
	}

	_ = result
}

func TestFormatObjects_DryRunDoesNotModify(t *testing.T) {
	v := setupFormatVault(t)

	v.SaveType(&TypeSchema{
		Name:       "book",
		Properties: []Property{{Name: "author", Type: "string"}},
	})

	// Write with wrong order
	content := "---\nauthor: Alice\nname: Test\ncreated_at: \"2025-01-01T00:00:00Z\"\nupdated_at: \"2025-01-01T00:00:00Z\"\n---\n"
	id := writeRawObject(t, v, "book", "test", content)

	result, err := v.FormatObjects("", true)
	if err != nil {
		t.Fatalf("FormatObjects: %v", err)
	}
	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(result.Changed))
	}

	// Verify file was NOT modified
	parts := strings.SplitN(id, "/", 2)
	data, _ := os.ReadFile(v.ObjectPath(parts[0], parts[1]))
	if !strings.HasPrefix(string(data), "---\nauthor:") {
		t.Fatalf("file was modified during dry-run")
	}
}

func TestFormatSchemas_WithNameTemplate(t *testing.T) {
	v := setupFormatVault(t)

	v.SaveType(&TypeSchema{
		Name:         "note",
		NameTemplate: "{{ date:2006-01-02 }}",
		Properties:   []Property{{Name: "topic", Type: "string"}},
	})

	// Schema was just saved via SaveType, so it should already be canonical.
	result, err := v.FormatSchemas("", false)
	if err != nil {
		t.Fatalf("FormatSchemas: %v", err)
	}
	if len(result.Changed) != 0 {
		t.Fatalf("expected 0 changed (just saved), got %d: %v", len(result.Changed), result.Changed)
	}
}

func TestFormatSchemas_BuiltinTypeSkipped(t *testing.T) {
	v := setupFormatVault(t)

	// Format schemas for "tag" type — it's built-in, no file
	result, err := v.FormatSchemas("tag", false)
	if err != nil {
		t.Fatalf("FormatSchemas: %v", err)
	}
	if len(result.Changed) != 0 {
		t.Fatalf("expected 0 changed for built-in type, got %d", len(result.Changed))
	}
}

func TestFormatAll_InvalidType(t *testing.T) {
	v := setupFormatVault(t)

	_, err := v.FormatAll("nonexistent", false)
	if err == nil {
		t.Fatal("expected error for invalid type")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Fatalf("expected error to mention type name, got: %v", err)
	}
}

func TestFormatAll_CombinesResults(t *testing.T) {
	v := setupFormatVault(t)

	v.SaveType(&TypeSchema{
		Name:       "book",
		Properties: []Property{{Name: "author", Type: "string"}},
	})

	// Write object with wrong order
	content := "---\nauthor: Alice\nname: Test\ncreated_at: \"2025-01-01T00:00:00Z\"\nupdated_at: \"2025-01-01T00:00:00Z\"\n---\n"
	writeRawObject(t, v, "book", "test", content)

	result, err := v.FormatAll("", false)
	if err != nil {
		t.Fatalf("FormatAll: %v", err)
	}

	// Should have at least 1 changed (the mis-ordered object)
	if len(result.Changed) < 1 {
		t.Fatalf("expected at least 1 changed, got %d", len(result.Changed))
	}
}

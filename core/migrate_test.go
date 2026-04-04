package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupMigrateTestVault(t *testing.T) *Vault {
	t.Helper()
	v := setupTestVault(t)

	// Initial schema with title and status
	schema := []byte(`name: book
properties:
  - name: title
    type: string
  - name: status
    type: string
`)
	mustWriteTypeSchema(v, "book", schema)
	return v
}

func TestVault_MigrateObjects_AddProperty(t *testing.T) {
	v := setupMigrateTestVault(t)

	// Create objects with original schema
	objA, _ := v.NewObject("book", "book-a", "")
	objB, _ := v.NewObject("book", "book-b", "")

	// Update schema: add isbn with default
	newSchema := []byte(`name: book
properties:
  - name: title
    type: string
  - name: status
    type: string
  - name: isbn
    type: string
    default: "unknown"
`)
	mustWriteTypeSchema(v, "book", newSchema)

	result, err := v.MigrateObjects("book", MigrateOptions{})
	if err != nil {
		t.Fatalf("MigrateObjects() error = %v", err)
	}

	if len(result.Changes) != 2 {
		t.Fatalf("len(Changes) = %d, want 2", len(result.Changes))
	}

	// Verify both objects now have isbn
	for _, id := range []string{objA.ID, objB.ID} {
		obj, err := v.GetObject(id)
		if err != nil {
			t.Fatalf("GetObject(%s) error = %v", id, err)
		}
		if obj.Properties["isbn"] != "unknown" {
			t.Errorf("%s isbn = %v, want %q", id, obj.Properties["isbn"], "unknown")
		}
	}
}

func TestVault_MigrateObjects_RemoveProperty(t *testing.T) {
	v := setupMigrateTestVault(t)

	created, _ := v.NewObject("book", "test", "")

	// Update schema: remove status
	newSchema := []byte(`name: book
properties:
  - name: title
    type: string
`)
	mustWriteTypeSchema(v, "book", newSchema)

	result, err := v.MigrateObjects("book", MigrateOptions{})
	if err != nil {
		t.Fatalf("MigrateObjects() error = %v", err)
	}

	if len(result.Changes) != 1 {
		t.Fatalf("len(Changes) = %d, want 1", len(result.Changes))
	}
	if len(result.Changes[0].Removed) != 1 || result.Changes[0].Removed[0] != "status" {
		t.Errorf("Removed = %v, want [status]", result.Changes[0].Removed)
	}

	obj, _ := v.GetObject(created.ID)
	if _, exists := obj.Properties["status"]; exists {
		t.Error("status property should have been removed")
	}
}

func TestVault_MigrateObjects_RenameProperty(t *testing.T) {
	v := setupMigrateTestVault(t)

	created, _ := v.NewObject("book", "test", "")
	v.SetProperty(created.ID, "status", "reading")

	// Update schema: rename status -> reading_status
	newSchema := []byte(`name: book
properties:
  - name: title
    type: string
  - name: reading_status
    type: string
`)
	mustWriteTypeSchema(v, "book", newSchema)

	result, err := v.MigrateObjects("book", MigrateOptions{
		Renames: map[string]string{"status": "reading_status"},
	})
	if err != nil {
		t.Fatalf("MigrateObjects() error = %v", err)
	}

	if len(result.Changes) != 1 {
		t.Fatalf("len(Changes) = %d, want 1", len(result.Changes))
	}
	if result.Changes[0].Renamed["status"] != "reading_status" {
		t.Errorf("Renamed = %v, want status->reading_status", result.Changes[0].Renamed)
	}

	updated, _ := v.GetObject(created.ID)
	if updated.Properties["reading_status"] != "reading" {
		t.Errorf("reading_status = %v, want %q", updated.Properties["reading_status"], "reading")
	}
	if _, exists := updated.Properties["status"]; exists {
		t.Error("old property 'status' should have been removed")
	}
}

func TestVault_MigrateObjects_DryRun(t *testing.T) {
	v := setupMigrateTestVault(t)

	created, _ := v.NewObject("book", "test", "")

	// Update schema: add isbn
	newSchema := []byte(`name: book
properties:
  - name: title
    type: string
  - name: status
    type: string
  - name: isbn
    type: string
    default: "unknown"
`)
	mustWriteTypeSchema(v, "book", newSchema)

	result, err := v.MigrateObjects("book", MigrateOptions{DryRun: true})
	if err != nil {
		t.Fatalf("MigrateObjects() error = %v", err)
	}

	if len(result.Changes) != 1 {
		t.Fatalf("len(Changes) = %d, want 1", len(result.Changes))
	}

	// Verify file was NOT modified
	obj, _ := v.GetObject(created.ID)
	if _, exists := obj.Properties["isbn"]; exists {
		t.Error("dry-run should not modify files")
	}
}

func TestVault_MigrateObjects_NoChanges(t *testing.T) {
	v := setupMigrateTestVault(t)

	v.NewObject("book", "test", "")

	// Schema unchanged — no migration needed
	result, err := v.MigrateObjects("book", MigrateOptions{})
	if err != nil {
		t.Fatalf("MigrateObjects() error = %v", err)
	}

	if len(result.Changes) != 0 {
		t.Errorf("len(Changes) = %d, want 0 (no changes needed)", len(result.Changes))
	}
}

func TestVault_MigrateObjects_TypeNotFound(t *testing.T) {
	v := setupMigrateTestVault(t)

	_, err := v.MigrateObjects("nonexistent", MigrateOptions{})
	if err == nil {
		t.Fatal("expected error for nonexistent type, got nil")
	}
}

func TestVault_MigrateObjects_RenameTargetNotInSchema(t *testing.T) {
	v := setupMigrateTestVault(t)

	_, err := v.MigrateObjects("book", MigrateOptions{
		Renames: map[string]string{"status": "nonexistent"},
	})
	if err == nil {
		t.Fatal("expected error when rename target not in schema, got nil")
	}
}

func TestVault_MigrateObjects_RenameSourceStillInSchema(t *testing.T) {
	v := setupMigrateTestVault(t)

	// status still exists in schema — should error
	_, err := v.MigrateObjects("book", MigrateOptions{
		Renames: map[string]string{"status": "title"},
	})
	if err == nil {
		t.Fatal("expected error when rename source still in schema, got nil")
	}
}

func TestVault_MigrateObjects_DBNotOpen(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()

	_, err := v.MigrateObjects("book", MigrateOptions{})
	if err == nil {
		t.Fatal("expected error when DB not opened, got nil")
	}
}

// ── Directory layout migration tests ────────────────────────────────────────

func TestMigrateDirectoryLayout_EmptyOldTypesDir(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()

	// Create empty old types dir, remove new types dir
	os.MkdirAll(filepath.Join(dir, ".typemd", "types"), 0755)
	os.Remove(filepath.Join(dir, "types"))

	err := v.migrateDirectoryLayout()
	if err != nil {
		t.Fatalf("migrateDirectoryLayout() error = %v", err)
	}

	// New types dir should exist (even if empty)
	if _, err := os.Stat(filepath.Join(dir, "types")); os.IsNotExist(err) {
		t.Error("expected types/ to exist at root")
	}
	// Old types dir should be gone
	if _, err := os.Stat(filepath.Join(dir, ".typemd", "types")); !os.IsNotExist(err) {
		t.Error("expected .typemd/types/ to be removed")
	}
}

func TestMigrateDirectoryLayout_CreatesPropertiesDir(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()

	// Create old properties file, remove new properties dir
	os.WriteFile(filepath.Join(dir, ".typemd", "properties.yaml"), []byte("- name: rating\n  type: number\n"), 0644)
	os.Remove(filepath.Join(dir, "properties"))

	err := v.migrateDirectoryLayout()
	if err != nil {
		t.Fatalf("migrateDirectoryLayout() error = %v", err)
	}

	// Properties dir should be created and file moved
	data, err := os.ReadFile(filepath.Join(dir, "properties", "properties.yaml"))
	if err != nil {
		t.Fatalf("expected properties/properties.yaml to exist: %v", err)
	}
	if !strings.Contains(string(data), "rating") {
		t.Error("expected properties content to be preserved")
	}
	// Old file should be gone
	if _, err := os.Stat(filepath.Join(dir, ".typemd", "properties.yaml")); !os.IsNotExist(err) {
		t.Error("expected .typemd/properties.yaml to be removed")
	}
}

func TestMigrateDirectoryLayout_NothingToMigrate(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()

	err := v.migrateDirectoryLayout()
	if err != nil {
		t.Fatalf("migrateDirectoryLayout() error = %v", err)
	}
}

func TestMigrateDirectoryLayout_TypesConflict(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()

	// Both old and new types dirs exist with content
	os.MkdirAll(filepath.Join(dir, ".typemd", "types", "book"), 0755)
	os.WriteFile(filepath.Join(dir, ".typemd", "types", "book", "schema.yaml"), []byte("name: book\n"), 0644)
	os.MkdirAll(filepath.Join(dir, "types", "note"), 0755)
	os.WriteFile(filepath.Join(dir, "types", "note", "schema.yaml"), []byte("name: note\n"), 0644)

	err := v.migrateDirectoryLayout()
	if err == nil {
		t.Fatal("expected conflict error")
	}
	if !strings.Contains(err.Error(), "conflict") {
		t.Errorf("expected error to mention 'conflict', got: %v", err)
	}
}

func TestMigrateDirectoryLayout_PropertiesConflict(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()

	// Both old and new properties exist
	os.WriteFile(filepath.Join(dir, ".typemd", "properties.yaml"), []byte("old\n"), 0644)
	os.MkdirAll(filepath.Join(dir, "properties"), 0755)
	os.WriteFile(filepath.Join(dir, "properties", "properties.yaml"), []byte("new\n"), 0644)

	err := v.migrateDirectoryLayout()
	if err == nil {
		t.Fatal("expected conflict error")
	}
	if !strings.Contains(err.Error(), "conflict") {
		t.Errorf("expected error to mention 'conflict', got: %v", err)
	}
}

func TestInit_CreatesRootLevelDirs(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// types/ should be at root, not .typemd/types/
	if _, err := os.Stat(filepath.Join(dir, "types")); os.IsNotExist(err) {
		t.Error("expected types/ at vault root")
	}
	if _, err := os.Stat(filepath.Join(dir, ".typemd", "types")); !os.IsNotExist(err) {
		t.Error("expected .typemd/types/ to NOT exist")
	}

	// properties/ should be at root
	if _, err := os.Stat(filepath.Join(dir, "properties")); os.IsNotExist(err) {
		t.Error("expected properties/ at vault root")
	}
}

// ── Schema migration tests (enum → select) ─────────────────────────────────

func TestVault_MigrateSchemas_EnumToSelect(t *testing.T) {
	v := setupMigrateTestVault(t)

	enumSchema := []byte(`name: book
properties:
  - name: title
    type: string
  - name: status
    type: enum
    values:
      - to-read
      - reading
      - done
`)
	mustWriteTypeSchema(v, "book", enumSchema)

	result, err := v.MigrateSchemas(false)
	if err != nil {
		t.Fatalf("MigrateSchemas() error = %v", err)
	}

	if len(result.Changes) != 1 {
		t.Fatalf("len(Changes) = %d, want 1", len(result.Changes))
	}
	if result.Changes[0].TypeName != "book" {
		t.Errorf("TypeName = %q, want %q", result.Changes[0].TypeName, "book")
	}
	if len(result.Changes[0].Properties) != 1 || result.Changes[0].Properties[0] != "status" {
		t.Errorf("Properties = %v, want [status]", result.Changes[0].Properties)
	}

	// Verify the file was rewritten correctly
	schema, err := v.LoadType("book")
	if err != nil {
		t.Fatalf("LoadType() error = %v", err)
	}
	statusProp := schema.Properties[1]
	if statusProp.Type != "select" {
		t.Errorf("status.Type = %q, want %q", statusProp.Type, "select")
	}
	if len(statusProp.Options) != 3 {
		t.Fatalf("len(status.Options) = %d, want 3", len(statusProp.Options))
	}
	if statusProp.Options[0].Value != "to-read" {
		t.Errorf("Options[0].Value = %q, want %q", statusProp.Options[0].Value, "to-read")
	}
}

func TestVault_MigrateSchemas_DryRun(t *testing.T) {
	v := setupMigrateTestVault(t)

	enumSchema := []byte(`name: book
properties:
  - name: status
    type: enum
    values: [to-read, reading, done]
`)
	mustWriteTypeSchema(v, "book", enumSchema)

	result, err := v.MigrateSchemas(true)
	if err != nil {
		t.Fatalf("MigrateSchemas() error = %v", err)
	}

	if len(result.Changes) != 1 {
		t.Fatalf("len(Changes) = %d, want 1", len(result.Changes))
	}

	// Verify file was NOT modified
	data, _ := os.ReadFile(filepath.Join(v.TypesDir(), "book", "schema.yaml"))
	if !strings.Contains(string(data), "type: enum") {
		t.Error("dry-run should not modify the schema file")
	}
}

func TestVault_MigrateSchemas_NoEnums(t *testing.T) {
	v := setupMigrateTestVault(t)

	schema := []byte(`name: book
properties:
  - name: title
    type: string
  - name: status
    type: select
    options:
      - value: to-read
`)
	mustWriteTypeSchema(v, "book", schema)

	result, err := v.MigrateSchemas(false)
	if err != nil {
		t.Fatalf("MigrateSchemas() error = %v", err)
	}

	if len(result.Changes) != 0 {
		t.Errorf("len(Changes) = %d, want 0", len(result.Changes))
	}
}

func TestVault_MigrateSchemas_MultipleEnums(t *testing.T) {
	v := setupMigrateTestVault(t)

	schema := []byte(`name: book
properties:
  - name: status
    type: enum
    values: [to-read, reading, done]
  - name: category
    type: enum
    values: [fiction, non-fiction]
`)
	mustWriteTypeSchema(v, "book", schema)

	result, err := v.MigrateSchemas(false)
	if err != nil {
		t.Fatalf("MigrateSchemas() error = %v", err)
	}

	if len(result.Changes) != 1 {
		t.Fatalf("len(Changes) = %d, want 1", len(result.Changes))
	}
	if len(result.Changes[0].Properties) != 2 {
		t.Errorf("Properties = %v, want 2 properties", result.Changes[0].Properties)
	}
}

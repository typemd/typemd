package core

import (
	"os"
	"path/filepath"
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
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), schema, 0644)
	return v
}

func TestVault_MigrateObjects_AddProperty(t *testing.T) {
	v := setupMigrateTestVault(t)

	// Create objects with original schema
	objA, _ := v.NewObject("book", "book-a")
	objB, _ := v.NewObject("book", "book-b")

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
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), newSchema, 0644)

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

	created, _ := v.NewObject("book", "test")

	// Update schema: remove status
	newSchema := []byte(`name: book
properties:
  - name: title
    type: string
`)
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), newSchema, 0644)

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

	created, _ := v.NewObject("book", "test")
	v.SetProperty(created.ID, "status", "reading")

	// Update schema: rename status -> reading_status
	newSchema := []byte(`name: book
properties:
  - name: title
    type: string
  - name: reading_status
    type: string
`)
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), newSchema, 0644)

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

	created, _ := v.NewObject("book", "test")

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
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), newSchema, 0644)

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

	v.NewObject("book", "test")

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

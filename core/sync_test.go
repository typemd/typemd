package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVault_SyncIndex_NewFile(t *testing.T) {
	v := setupTestVault(t)

	// Manually create an object file (bypassing NewObject API)
	typeDir := filepath.Join(v.ObjectsDir(), "book")
	os.MkdirAll(typeDir, 0755)
	os.WriteFile(filepath.Join(typeDir, "test-book.md"), []byte("---\ntitle: Test Book\n---\n\nHello world.\n"), 0644)

	// Also need type schema for the type directory to be valid
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), []byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)

	if _, err := v.SyncIndex(); err != nil {
		t.Fatalf("SyncIndex() error = %v", err)
	}

	// Verify object is now in DB
	objs, err := v.QueryObjects("type=book")
	if err != nil {
		t.Fatalf("QueryObjects() error = %v", err)
	}
	if len(objs) != 1 {
		t.Fatalf("len(objs) = %d, want 1", len(objs))
	}
	if objs[0].ID != "book/test-book" {
		t.Errorf("ID = %q, want %q", objs[0].ID, "book/test-book")
	}
}

func TestVault_SyncIndex_UpdatedFile(t *testing.T) {
	v := setupTestVault(t)

	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), []byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)

	// Create via API (body is empty in DB)
	v.NewObject("book", "test-book")

	// Manually edit the file to add body
	objPath := v.ObjectPath("book", "test-book")
	os.WriteFile(objPath, []byte("---\ntitle: Updated\n---\n\nNew body content.\n"), 0644)

	if _, err := v.SyncIndex(); err != nil {
		t.Fatalf("SyncIndex() error = %v", err)
	}

	objs, err := v.QueryObjects("type=book")
	if err != nil {
		t.Fatalf("QueryObjects() error = %v", err)
	}
	if len(objs) != 1 {
		t.Fatalf("len(objs) = %d, want 1", len(objs))
	}
	if strings.TrimSpace(objs[0].Body) != "New body content." {
		t.Errorf("Body = %q, want %q", objs[0].Body, "New body content.")
	}
}

func TestVault_SyncIndex_DeletedFile(t *testing.T) {
	v := setupTestVault(t)

	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), []byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)

	// Create via API
	v.NewObject("book", "test-book")

	// Delete the file
	os.Remove(v.ObjectPath("book", "test-book"))

	result, err := v.SyncIndex()
	if err != nil {
		t.Fatalf("SyncIndex() error = %v", err)
	}
	if result.Deleted != 1 {
		t.Errorf("Deleted = %d, want 1", result.Deleted)
	}

	objs, err := v.QueryObjects("type=book")
	if err != nil {
		t.Fatalf("QueryObjects() error = %v", err)
	}
	if len(objs) != 0 {
		t.Errorf("len(objs) = %d, want 0 (deleted file should be removed from DB)", len(objs))
	}
}

func TestVault_SyncIndex_DBNotOpen(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()

	_, err := v.SyncIndex()
	if err == nil {
		t.Fatal("expected error when DB not opened, got nil")
	}
}

func TestVault_SyncIndex_OrphanedRelations(t *testing.T) {
	v := setupRelationTestVault(t)

	// Create two objects and link them
	v.NewObject("book", "golang-in-action")
	v.NewObject("person", "alan-donovan")
	v.LinkObjects("book/golang-in-action", "author", "person/alan-donovan")

	// Verify relations exist
	var count int
	v.db.QueryRow("SELECT COUNT(*) FROM relations").Scan(&count)
	if count == 0 {
		t.Fatal("expected relations to exist before deletion")
	}

	// Delete the person file from disk (simulating user deletion)
	os.Remove(v.ObjectPath("person", "alan-donovan"))

	// SyncIndex should detect and clean orphaned relations
	result, err := v.SyncIndex()
	if err != nil {
		t.Fatalf("SyncIndex() error = %v", err)
	}

	if len(result.Orphaned) == 0 {
		t.Fatal("expected orphaned relations, got none")
	}

	// Verify orphaned relations contain the right data
	found := false
	for _, o := range result.Orphaned {
		if o.ToID == "person/alan-donovan" || o.FromID == "person/alan-donovan" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected orphan referencing person/alan-donovan, got %+v", result.Orphaned)
	}

	// Verify relations table is now clean
	v.db.QueryRow("SELECT COUNT(*) FROM relations").Scan(&count)
	if count != 0 {
		t.Errorf("relations count after cleanup = %d, want 0", count)
	}
}

func TestVault_SyncIndex_NoOrphansWhenAllExist(t *testing.T) {
	v := setupRelationTestVault(t)

	v.NewObject("book", "test-book")
	v.NewObject("person", "alan")
	v.LinkObjects("book/test-book", "author", "person/alan")

	result, err := v.SyncIndex()
	if err != nil {
		t.Fatalf("SyncIndex() error = %v", err)
	}

	if len(result.Orphaned) != 0 {
		t.Errorf("expected no orphans, got %+v", result.Orphaned)
	}

	// Relations should still exist
	var count int
	v.db.QueryRow("SELECT COUNT(*) FROM relations").Scan(&count)
	if count == 0 {
		t.Error("expected relations to still exist")
	}
}

func TestVault_SyncIndex_OrphanFromSourceDeletion(t *testing.T) {
	v := setupRelationTestVault(t)

	v.NewObject("book", "test-book")
	v.NewObject("person", "alan")
	v.LinkObjects("book/test-book", "author", "person/alan")

	// Delete the source (book) instead of the target
	os.Remove(v.ObjectPath("book", "test-book"))

	result, err := v.SyncIndex()
	if err != nil {
		t.Fatalf("SyncIndex() error = %v", err)
	}

	if len(result.Orphaned) == 0 {
		t.Fatal("expected orphaned relations when source is deleted")
	}

	// Verify relations table is clean
	var count int
	v.db.QueryRow("SELECT COUNT(*) FROM relations").Scan(&count)
	if count != 0 {
		t.Errorf("relations count after cleanup = %d, want 0", count)
	}
}

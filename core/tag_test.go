package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveTagReference_ByID(t *testing.T) {
	v := setupTestVault(t)
	tag, _ := v.NewObject("tag", "go")

	diskTags := map[string]*Object{tag.ID: tag}
	nameIndex := map[string]string{"go": tag.ID}
	resolved, ok := v.resolveTagReference(tag.ID, diskTags, nameIndex)
	if !ok {
		t.Fatal("expected tag reference to resolve by ID")
	}
	if resolved != tag.ID {
		t.Errorf("resolved = %q, want %q", resolved, tag.ID)
	}
}

func TestResolveTagReference_ByName(t *testing.T) {
	v := setupTestVault(t)
	tag, _ := v.NewObject("tag", "go")

	diskTags := map[string]*Object{tag.ID: tag}
	nameIndex := map[string]string{"go": tag.ID}
	resolved, ok := v.resolveTagReference("tag/go", diskTags, nameIndex)
	if !ok {
		t.Fatal("expected tag reference to resolve by name")
	}
	if resolved != tag.ID {
		t.Errorf("resolved = %q, want %q", resolved, tag.ID)
	}
}

func TestResolveTagReference_MissingName(t *testing.T) {
	v := setupTestVault(t)
	diskTags := map[string]*Object{}
	nameIndex := map[string]string{}
	_, ok := v.resolveTagReference("tag/nonexistent", diskTags, nameIndex)
	if ok {
		t.Error("expected tag reference not to resolve")
	}
}

func TestResolveTagReference_BrokenFullID(t *testing.T) {
	v := setupTestVault(t)
	diskTags := map[string]*Object{}
	nameIndex := map[string]string{}
	// A full ID that ends with ULID but doesn't exist
	_, ok := v.resolveTagReference("tag/go-01jqr3k5mpbvn8e0f2g7h9txyz", diskTags, nameIndex)
	if ok {
		t.Error("expected broken full-ID reference not to resolve")
	}
}

func TestResolveTagReference_NoPrefix(t *testing.T) {
	v := setupTestVault(t)
	diskTags := map[string]*Object{}
	nameIndex := map[string]string{}
	_, ok := v.resolveTagReference("go", diskTags, nameIndex)
	if ok {
		t.Error("expected reference without tag/ prefix not to resolve")
	}
}

func TestSyncIndex_TagRelationsWritten(t *testing.T) {
	v := setupTestVault(t)

	// Create a tag object
	tag, _ := v.NewObject("tag", "go")

	// Create a book with tags referencing the tag by name
	book, _ := v.NewObject("book", "golang-book")
	book.Properties[TagsProperty] = []any{"tag/go"}
	v.SaveObject(book)

	// Sync
	_, err := v.SyncIndex()
	if err != nil {
		t.Fatalf("SyncIndex error = %v", err)
	}

	// Check relations table
	rels, err := v.ListRelations(book.ID)
	if err != nil {
		t.Fatalf("ListRelations error = %v", err)
	}

	found := false
	for _, r := range rels {
		if r.Name == TagsProperty && r.FromID == book.ID && r.ToID == tag.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected tag relation from %s to %s, got %v", book.ID, tag.ID, rels)
	}
}

func TestSyncIndex_AutoCreateTag(t *testing.T) {
	v := setupTestVault(t)

	// Create a book referencing a tag that doesn't exist
	book, _ := v.NewObject("book", "golang-book")
	book.Properties[TagsProperty] = []any{"tag/auto-created"}
	v.SaveObject(book)

	_, err := v.SyncIndex()
	if err != nil {
		t.Fatalf("SyncIndex error = %v", err)
	}

	// Check that the auto-created tag exists
	entries, _ := filepath.Glob(filepath.Join(v.ObjectDir("tag"), "auto-created-*.md"))
	if len(entries) == 0 {
		t.Error("expected auto-created tag object on disk")
	}

	// Check relations table
	rels, err := v.ListRelations(book.ID)
	if err != nil {
		t.Fatalf("ListRelations error = %v", err)
	}
	found := false
	for _, r := range rels {
		if r.Name == TagsProperty {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected tag relation to be created for auto-created tag")
	}
}

func TestSyncIndex_TagPreservedInFiltering(t *testing.T) {
	v := setupTestVault(t)

	tag, _ := v.NewObject("tag", "go")

	// Write a raw book file with tags in frontmatter
	ulid, _ := GenerateULID()
	filename := "test-book-" + ulid
	objPath := v.ObjectPath("book", filename)
	os.MkdirAll(filepath.Dir(objPath), 0755)
	content := "---\nname: test-book\ntags:\n  - " + tag.ID + "\ntitle: Test\n---\n"
	os.WriteFile(objPath, []byte(content), 0644)

	_, err := v.SyncIndex()
	if err != nil {
		t.Fatalf("SyncIndex error = %v", err)
	}

	// Read back from DB
	var propsJSON string
	err = v.db.QueryRow("SELECT properties FROM objects WHERE id = ?", "book/"+filename).Scan(&propsJSON)
	if err != nil {
		t.Fatalf("query error = %v", err)
	}

	// tags should be preserved
	if !strings.Contains(propsJSON, TagsProperty) {
		t.Errorf("indexed properties should contain %q, got: %s", TagsProperty, propsJSON)
	}
}

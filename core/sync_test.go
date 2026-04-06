package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mustReconcileAndProject runs a full Reconcile + Project cycle.
func mustReconcileAndProject(t *testing.T, v *Vault) *ReconcileResult {
	t.Helper()
	events, result, err := v.Reconcile()
	if err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}
	if err := v.Project(events); err != nil {
		t.Fatalf("Project() error = %v", err)
	}
	return result
}

// mustQueryRelationCount returns the number of rows in the relations table.
func mustQueryRelationCount(t *testing.T, v *Vault) int {
	t.Helper()
	var count int
	if err := v.db.QueryRow("SELECT COUNT(*) FROM relations").Scan(&count); err != nil {
		t.Fatalf("query relation count: %v", err)
	}
	return count
}

func TestVault_SyncIndex_NewFile(t *testing.T) {
	v := setupTestVault(t)

	// setupTestVault already writes a "book" schema with title property.
	// Manually create an object file (bypassing NewObject API).
	typeDir := filepath.Join(v.ObjectsDir(), "book")
	if err := os.MkdirAll(typeDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(typeDir, "test-book.md"), []byte("---\ntitle: Test Book\n---\n\nHello world.\n"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	mustReconcileAndProject(t, v)

	objs, err := v.QueryObjects(TypeFilter("book"))
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

	obj, err := v.NewObject("book", "test-book", "")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	// Manually edit the file to add body
	objPath := v.ObjectPath(obj.Type, obj.Filename)
	if err := os.WriteFile(objPath, []byte("---\ntitle: Updated\n---\n\nNew body content.\n"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	mustReconcileAndProject(t, v)

	objs, err := v.QueryObjects(TypeFilter("book"))
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

	obj, err := v.NewObject("book", "test-book", "")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	if err := os.Remove(v.ObjectPath(obj.Type, obj.Filename)); err != nil {
		t.Fatalf("os.Remove() error = %v", err)
	}

	result := mustReconcileAndProject(t, v)
	if result.Deleted != 1 {
		t.Errorf("Deleted = %d, want 1", result.Deleted)
	}

	objs, err := v.QueryObjects(TypeFilter("book"))
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
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	_, _, err := v.Reconcile()
	if err == nil {
		t.Fatal("expected error when DB not opened, got nil")
	}
}

func TestVault_SyncIndex_OrphanedRelations(t *testing.T) {
	tests := []struct {
		name         string
		deleteSource bool // true: delete source (book); false: delete target (person)
	}{
		{name: "delete target cleans up relations", deleteSource: false},
		{name: "delete source cleans up relations", deleteSource: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := setupRelationTestVault(t)

			book, err := v.NewObject("book", "golang-in-action", "")
			if err != nil {
				t.Fatalf("NewObject(book) error = %v", err)
			}
			person, err := v.NewObject("person", "alan-donovan", "")
			if err != nil {
				t.Fatalf("NewObject(person) error = %v", err)
			}
			if err := v.LinkObjects(book.ID, "author", person.ID); err != nil {
				t.Fatalf("LinkObjects() error = %v", err)
			}

			if count := mustQueryRelationCount(t, v); count == 0 {
				t.Fatal("expected relations to exist before deletion")
			}

			// Delete either source or target
			toDelete := person
			if tc.deleteSource {
				toDelete = book
			}
			if err := os.Remove(v.ObjectPath(toDelete.Type, toDelete.Filename)); err != nil {
				t.Fatalf("os.Remove() error = %v", err)
			}

			result := mustReconcileAndProject(t, v)
			if result.Deleted != 1 {
				t.Errorf("Deleted = %d, want 1", result.Deleted)
			}
			if count := mustQueryRelationCount(t, v); count != 0 {
				t.Errorf("relations count after cleanup = %d, want 0", count)
			}
		})
	}
}

func TestVault_SyncIndex_NoOrphansWhenAllExist(t *testing.T) {
	v := setupRelationTestVault(t)

	book, err := v.NewObject("book", "test-book", "")
	if err != nil {
		t.Fatalf("NewObject(book) error = %v", err)
	}
	alan, err := v.NewObject("person", "alan", "")
	if err != nil {
		t.Fatalf("NewObject(person) error = %v", err)
	}
	if err := v.LinkObjects(book.ID, "author", alan.ID); err != nil {
		t.Fatalf("LinkObjects() error = %v", err)
	}

	mustReconcileAndProject(t, v)

	if count := mustQueryRelationCount(t, v); count == 0 {
		t.Error("expected relations to still exist")
	}
}

// ── Incremental sync (ReconcileFiles) ─────────────────────────────────────

func TestReconcileFiles_UpsertsNewObject(t *testing.T) {
	v := setupTestVault(t)

	obj, err := v.NewObject("book", "incrementalbook", "")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	objPath := v.ObjectPath(obj.Type, obj.Filename)
	events, _, err := v.ReconcileFiles([]string{objPath})
	if err != nil {
		t.Fatalf("ReconcileFiles() error = %v", err)
	}
	if err := v.Project(events); err != nil {
		t.Fatalf("Project() error = %v", err)
	}

	results, err := v.SearchObjects("incrementalbook")
	if err != nil {
		t.Fatalf("SearchObjects() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].Type != "book" {
		t.Errorf("Type = %q, want %q", results[0].Type, "book")
	}
}

func TestReconcileFiles_RemovesDeletedObject(t *testing.T) {
	v := setupTestVault(t)

	obj, err := v.NewObject("book", "to-be-deleted", "")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	mustReconcileAndProject(t, v)

	objs, err := v.QueryObjects(TypeFilter("book"))
	if err != nil {
		t.Fatalf("QueryObjects() error = %v", err)
	}
	if len(objs) != 1 {
		t.Fatalf("expected 1 object before deletion, got %d", len(objs))
	}

	objPath := v.ObjectPath(obj.Type, obj.Filename)
	if err := os.Remove(objPath); err != nil {
		t.Fatalf("os.Remove() error = %v", err)
	}

	events, result, err := v.ReconcileFiles([]string{objPath})
	if err != nil {
		t.Fatalf("ReconcileFiles() error = %v", err)
	}
	if err := v.Project(events); err != nil {
		t.Fatalf("Project() error = %v", err)
	}
	if result.Deleted != 1 {
		t.Errorf("Deleted = %d, want 1", result.Deleted)
	}

	objs, err = v.QueryObjects(TypeFilter("book"))
	if err != nil {
		t.Fatalf("QueryObjects() error = %v", err)
	}
	if len(objs) != 0 {
		t.Errorf("len(objs) = %d, want 0 (deleted file should be removed from index)", len(objs))
	}
}

func TestReconcileFiles_NilPathsFallsBackToFullSync(t *testing.T) {
	v := setupTestVault(t)

	_, err := v.NewObject("book", "nilpathsbook", "")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	events, _, err := v.ReconcileFiles(nil)
	if err != nil {
		t.Fatalf("ReconcileFiles(nil) error = %v", err)
	}
	if err := v.Project(events); err != nil {
		t.Fatalf("Project() error = %v", err)
	}

	results, err := v.SearchObjects("nilpathsbook")
	if err != nil {
		t.Fatalf("SearchObjects() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].Type != "book" {
		t.Errorf("Type = %q, want %q", results[0].Type, "book")
	}
}

// ── Property filtering during sync ────────────────────────────────────────

// getIndexedProps bypasses the query layer to verify raw indexed data in SQLite.
func getIndexedProps(t *testing.T, v *Vault, objectID string) map[string]any {
	t.Helper()
	var propsJSON string
	if err := v.db.QueryRow("SELECT properties FROM objects WHERE id = ?", objectID).Scan(&propsJSON); err != nil {
		t.Fatalf("query indexed properties for %q: %v", objectID, err)
	}
	var props map[string]any
	if err := json.Unmarshal([]byte(propsJSON), &props); err != nil {
		t.Fatalf("unmarshal indexed properties for %q: %v", objectID, err)
	}
	return props
}

func TestSync_FiltersUndefinedProperties(t *testing.T) {
	v := setupTestVault(t)

	// setupTestVault already writes a "book" schema with properties: title, status, rating.
	// Write a raw object file that includes an extra undefined property "mood".
	typeDir := filepath.Join(v.ObjectsDir(), "book")
	if err := os.MkdirAll(typeDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(typeDir, "clean-code.md"), []byte(
		"---\ntitle: Clean Code\nstatus: reading\nmood: happy\n---\nSome body.\n",
	), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	mustReconcileAndProject(t, v)

	props := getIndexedProps(t, v, "book/clean-code")

	if _, ok := props["title"]; !ok {
		t.Error("indexed properties should contain \"title\"")
	}
	if _, ok := props["status"]; !ok {
		t.Error("indexed properties should contain \"status\"")
	}
	if _, ok := props["mood"]; ok {
		t.Error("indexed properties should NOT contain \"mood\" (undefined in schema)")
	}
}

func TestSync_NoSchemaRetainsAllProperties(t *testing.T) {
	v := setupTestVault(t)

	// "recipe" has no type schema — all frontmatter properties should be indexed.
	typeDir := filepath.Join(v.ObjectsDir(), "recipe")
	if err := os.MkdirAll(typeDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(typeDir, "pasta.md"), []byte(
		"---\nname: Pasta\ntime: 30min\n---\nBoil water.\n",
	), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	mustReconcileAndProject(t, v)

	props := getIndexedProps(t, v, "recipe/pasta")

	if _, ok := props["name"]; !ok {
		t.Error("indexed properties should contain \"name\"")
	}
	if _, ok := props["time"]; !ok {
		t.Error("indexed properties should contain \"time\"")
	}
}

// setupRelationSyncVault creates a vault with book→person schemas that include
// both author and editor relation properties, suitable for testing multi-property
// expansion and SyncResult tracking.
func setupRelationSyncVault(t *testing.T) *Vault {
	t.Helper()
	v := setupTestVault(t)

	bookSchema := []byte(`name: book
properties:
  - name: title
    type: string
  - name: author
    type: relation
    target: person
    bidirectional: true
    inverse: books
  - name: editor
    type: relation
    target: person
`)
	mustWriteTypeSchema(v, "book", bookSchema)

	personSchema := []byte(`name: person
properties:
  - name: books
    type: relation
    target: book
    multiple: true
    bidirectional: true
    inverse: author
`)
	mustWriteTypeSchema(v, "person", personSchema)

	return v
}

func TestReconcile_MultiPropertyExpansion(t *testing.T) {
	v := setupRelationSyncVault(t)

	person1, err := v.NewObject("person", "john-doe", "")
	if err != nil {
		t.Fatalf("NewObject(person1) error = %v", err)
	}
	person2, err := v.NewObject("person", "jane-smith", "")
	if err != nil {
		t.Fatalf("NewObject(person2) error = %v", err)
	}

	// Create book with name references (not full IDs)
	book, err := v.NewObject("book", "clean-code", "")
	if err != nil {
		t.Fatalf("NewObject(book) error = %v", err)
	}
	book.Properties["author"] = "person/john-doe"
	book.Properties["editor"] = "person/jane-smith"
	if err := v.SaveObject(book); err != nil {
		t.Fatalf("SaveObject() error = %v", err)
	}

	result := mustReconcileAndProject(t, v)

	if result.Expanded != 2 {
		t.Errorf("Expanded = %d, want 2", result.Expanded)
	}

	// Verify both properties were written back with full IDs
	freshBook, err := v.GetObject(book.ID)
	if err != nil {
		t.Fatalf("GetObject() error = %v", err)
	}
	authorVal, _ := freshBook.Properties["author"].(string)
	if authorVal != person1.ID {
		t.Errorf("author = %q, want %q", authorVal, person1.ID)
	}
	editorVal, _ := freshBook.Properties["editor"].(string)
	if editorVal != person2.ID {
		t.Errorf("editor = %q, want %q", editorVal, person2.ID)
	}
}

func TestReconcile_SyncResultExpansionCount(t *testing.T) {
	v := setupRelationSyncVault(t)

	if _, err := v.NewObject("person", "john-doe", ""); err != nil {
		t.Fatalf("NewObject(person) error = %v", err)
	}

	// Create book with a name reference
	book, err := v.NewObject("book", "clean-code", "")
	if err != nil {
		t.Fatalf("NewObject(book) error = %v", err)
	}
	book.Properties["author"] = "person/john-doe"
	if err := v.SaveObject(book); err != nil {
		t.Fatalf("SaveObject() error = %v", err)
	}

	result := mustReconcileAndProject(t, v)

	if result.Expanded != 1 {
		t.Errorf("Expanded = %d, want 1", result.Expanded)
	}

	// Second sync should have 0 expansions (already full IDs)
	result2 := mustReconcileAndProject(t, v)
	if result2.Expanded != 0 {
		t.Errorf("second sync Expanded = %d, want 0", result2.Expanded)
	}
}

func TestReconcile_SyncResultUnresolvedCount(t *testing.T) {
	v := setupRelationSyncVault(t)

	// Create book referencing a non-existent person by name
	book, err := v.NewObject("book", "clean-code", "")
	if err != nil {
		t.Fatalf("NewObject(book) error = %v", err)
	}
	book.Properties["author"] = "person/nobody"
	if err := v.SaveObject(book); err != nil {
		t.Fatalf("SaveObject() error = %v", err)
	}

	result := mustReconcileAndProject(t, v)

	if len(result.Unresolved) != 1 {
		t.Errorf("Unresolved = %d, want 1", len(result.Unresolved))
	}
}

func TestSync_FrontmatterNotModifiedDuringSync(t *testing.T) {
	v := setupTestVault(t)

	// Write a raw object file with an extra undefined property "mood".
	typeDir := filepath.Join(v.ObjectsDir(), "book")
	if err := os.MkdirAll(typeDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	objPath := filepath.Join(typeDir, "tidy-code.md")
	if err := os.WriteFile(objPath, []byte(
		"---\ntitle: Tidy Code\nmood: excited\n---\nGreat read.\n",
	), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	mustReconcileAndProject(t, v)

	// The file on disk must still contain the extra property — sync should not strip it.
	data, err := os.ReadFile(objPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "mood: excited") {
		t.Errorf("file should still contain \"mood: excited\" in frontmatter, got:\n%s", content)
	}
	if !strings.Contains(content, "title: Tidy Code") {
		t.Errorf("file should still contain \"title: Tidy Code\" in frontmatter, got:\n%s", content)
	}
}

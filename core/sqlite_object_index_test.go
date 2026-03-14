package core

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestIndex(t *testing.T) *SQLiteObjectIndex {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	idx := NewSQLiteObjectIndex(db)
	if err := idx.EnsureSchema(); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}
	return idx
}

func TestSQLiteObjectIndex_EnsureSchema(t *testing.T) {
	idx := setupTestIndex(t)
	// Calling EnsureSchema again should be idempotent
	if err := idx.EnsureSchema(); err != nil {
		t.Errorf("second EnsureSchema failed: %v", err)
	}
}

func TestSQLiteObjectIndex_UpsertAndQuery(t *testing.T) {
	idx := setupTestIndex(t)

	if err := idx.Upsert("book/test-01abc", "book", "test-01abc", `{"name":"Test"}`, "body"); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	results, err := idx.Query("")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].ID != "book/test-01abc" {
		t.Errorf("ID = %q, want %q", results[0].ID, "book/test-01abc")
	}
	if results[0].Type != "book" {
		t.Errorf("Type = %q, want %q", results[0].Type, "book")
	}
}

func TestSQLiteObjectIndex_UpsertUpdatesExisting(t *testing.T) {
	idx := setupTestIndex(t)

	idx.Upsert("book/test-01abc", "book", "test-01abc", `{"name":"Old"}`, "")
	idx.Upsert("book/test-01abc", "book", "test-01abc", `{"name":"New"}`, "")

	results, _ := idx.Query("")
	if len(results) != 1 {
		t.Fatalf("expected 1 result after upsert, got %d", len(results))
	}
	if results[0].Properties["name"] != "New" {
		t.Errorf("name = %v, want %q", results[0].Properties["name"], "New")
	}
}

func TestSQLiteObjectIndex_QueryByType(t *testing.T) {
	idx := setupTestIndex(t)

	idx.Upsert("book/a-01abc", "book", "a-01abc", `{"name":"A"}`, "")
	idx.Upsert("note/b-01abc", "note", "b-01abc", `{"name":"B"}`, "")

	results, err := idx.Query("type=book")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 book, got %d", len(results))
	}
	if results[0].Type != "book" {
		t.Errorf("Type = %q, want %q", results[0].Type, "book")
	}
}

func TestSQLiteObjectIndex_QueryByProperty(t *testing.T) {
	idx := setupTestIndex(t)

	idx.Upsert("book/a-01abc", "book", "a-01abc", `{"status":"reading"}`, "")
	idx.Upsert("book/b-01abc", "book", "b-01abc", `{"status":"done"}`, "")

	results, err := idx.Query("status=reading")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].ID != "book/a-01abc" {
		t.Errorf("ID = %q, want %q", results[0].ID, "book/a-01abc")
	}
}

func TestSQLiteObjectIndex_Search(t *testing.T) {
	idx := setupTestIndex(t)

	idx.Upsert("book/golang-01abc", "book", "golang-01abc", `{"name":"Go Programming"}`, "Learn golang")
	idx.Upsert("book/python-01abc", "book", "python-01abc", `{"name":"Python Basics"}`, "Learn python")

	// Rebuild FTS for inserted data
	idx.Rebuild()

	results, err := idx.Search("golang")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].ID != "book/golang-01abc" {
		t.Errorf("ID = %q, want %q", results[0].ID, "book/golang-01abc")
	}
}

func TestSQLiteObjectIndex_SearchEmpty(t *testing.T) {
	idx := setupTestIndex(t)

	results, err := idx.Search("")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if results != nil {
		t.Errorf("expected nil for empty search, got %v", results)
	}
}

func TestSQLiteObjectIndex_RemoveAndListIDs(t *testing.T) {
	idx := setupTestIndex(t)

	idx.Upsert("book/a-01abc", "book", "a-01abc", `{}`, "")
	idx.Upsert("book/b-01abc", "book", "b-01abc", `{}`, "")

	ids, _ := idx.ListIDs()
	if len(ids) != 2 {
		t.Fatalf("expected 2 IDs, got %d", len(ids))
	}

	if err := idx.Remove("book/a-01abc"); err != nil {
		t.Fatalf("remove: %v", err)
	}

	ids, _ = idx.ListIDs()
	if len(ids) != 1 {
		t.Fatalf("expected 1 ID after remove, got %d", len(ids))
	}
	if ids[0] != "book/b-01abc" {
		t.Errorf("remaining ID = %q, want %q", ids[0], "book/b-01abc")
	}
}

func TestSQLiteObjectIndex_NeedsSync(t *testing.T) {
	idx := setupTestIndex(t)

	needs, err := idx.NeedsSync()
	if err != nil {
		t.Fatalf("needsSync: %v", err)
	}
	if !needs {
		t.Error("expected NeedsSync=true on empty index")
	}

	idx.Upsert("book/a-01abc", "book", "a-01abc", `{}`, "")

	needs, _ = idx.NeedsSync()
	if needs {
		t.Error("expected NeedsSync=false after upsert")
	}
}

func TestSQLiteObjectIndex_Relations(t *testing.T) {
	idx := setupTestIndex(t)

	idx.Upsert("book/a-01abc", "book", "a-01abc", `{}`, "")
	idx.Upsert("person/b-01abc", "person", "b-01abc", `{}`, "")

	if err := idx.InsertRelation("author", "book/a-01abc", "person/b-01abc"); err != nil {
		t.Fatalf("insert relation: %v", err)
	}

	rels, err := idx.FindRelations("book/a-01abc")
	if err != nil {
		t.Fatalf("find relations: %v", err)
	}
	if len(rels) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(rels))
	}
	if rels[0].Name != "author" || rels[0].FromID != "book/a-01abc" || rels[0].ToID != "person/b-01abc" {
		t.Errorf("unexpected relation: %+v", rels[0])
	}

	// Also findable from target side
	rels, _ = idx.FindRelations("person/b-01abc")
	if len(rels) != 1 {
		t.Fatalf("expected 1 relation from target side, got %d", len(rels))
	}

	// Delete
	if err := idx.DeleteRelation("author", "book/a-01abc", "person/b-01abc"); err != nil {
		t.Fatalf("delete relation: %v", err)
	}
	rels, _ = idx.FindRelations("book/a-01abc")
	if len(rels) != 0 {
		t.Errorf("expected 0 relations after delete, got %d", len(rels))
	}
}

func TestSQLiteObjectIndex_DeleteRelationsByName(t *testing.T) {
	idx := setupTestIndex(t)

	idx.Upsert("book/a-01abc", "book", "a-01abc", `{}`, "")
	idx.Upsert("person/b-01abc", "person", "b-01abc", `{}`, "")
	idx.InsertRelation("tags", "book/a-01abc", "person/b-01abc")
	idx.InsertRelation("author", "book/a-01abc", "person/b-01abc")

	if err := idx.DeleteRelationsByName("tags"); err != nil {
		t.Fatalf("delete by name: %v", err)
	}

	rels, _ := idx.FindRelations("book/a-01abc")
	if len(rels) != 1 {
		t.Fatalf("expected 1 relation remaining, got %d", len(rels))
	}
	if rels[0].Name != "author" {
		t.Errorf("remaining relation = %q, want %q", rels[0].Name, "author")
	}
}

func TestSQLiteObjectIndex_CleanOrphanedRelations(t *testing.T) {
	idx := setupTestIndex(t)

	idx.Upsert("book/a-01abc", "book", "a-01abc", `{}`, "")
	idx.InsertRelation("author", "book/a-01abc", "person/gone-01abc")

	orphaned, err := idx.CleanOrphanedRelations()
	if err != nil {
		t.Fatalf("clean orphaned: %v", err)
	}
	if len(orphaned) != 1 {
		t.Fatalf("expected 1 orphaned, got %d", len(orphaned))
	}
	if orphaned[0].ToID != "person/gone-01abc" {
		t.Errorf("orphaned ToID = %q, want %q", orphaned[0].ToID, "person/gone-01abc")
	}

	// Should be cleaned up
	rels, _ := idx.FindRelations("book/a-01abc")
	if len(rels) != 0 {
		t.Errorf("expected 0 relations after cleanup, got %d", len(rels))
	}
}

func TestSQLiteObjectIndex_WikiLinks(t *testing.T) {
	idx := setupTestIndex(t)

	idx.Upsert("note/a-01abc", "note", "a-01abc", `{}`, "")
	idx.Upsert("person/b-01abc", "person", "b-01abc", `{}`, "")

	links := []WikiLinkEntry{
		{ToID: "person/b-01abc", Target: "person/b-01abc", DisplayText: "Bob"},
	}
	if err := idx.SyncWikiLinks("note/a-01abc", links); err != nil {
		t.Fatalf("sync wikilinks: %v", err)
	}

	// Forward links
	stored, err := idx.ListWikiLinks("note/a-01abc")
	if err != nil {
		t.Fatalf("list wikilinks: %v", err)
	}
	if len(stored) != 1 {
		t.Fatalf("expected 1 wikilink, got %d", len(stored))
	}
	if stored[0].DisplayText != "Bob" {
		t.Errorf("display text = %q, want %q", stored[0].DisplayText, "Bob")
	}

	// Backlinks
	backlinks, err := idx.FindBacklinks("person/b-01abc")
	if err != nil {
		t.Fatalf("find backlinks: %v", err)
	}
	if len(backlinks) != 1 {
		t.Fatalf("expected 1 backlink, got %d", len(backlinks))
	}
	if backlinks[0].FromID != "note/a-01abc" {
		t.Errorf("backlink FromID = %q, want %q", backlinks[0].FromID, "note/a-01abc")
	}

	// Delete
	if err := idx.DeleteWikiLinks("note/a-01abc"); err != nil {
		t.Fatalf("delete wikilinks: %v", err)
	}
	stored, _ = idx.ListWikiLinks("note/a-01abc")
	if len(stored) != 0 {
		t.Errorf("expected 0 wikilinks after delete, got %d", len(stored))
	}
}

func TestSQLiteObjectIndex_SyncWikiLinksReplaces(t *testing.T) {
	idx := setupTestIndex(t)

	idx.Upsert("note/a-01abc", "note", "a-01abc", `{}`, "")

	// First sync
	idx.SyncWikiLinks("note/a-01abc", []WikiLinkEntry{
		{ToID: "x", Target: "x", DisplayText: ""},
		{ToID: "y", Target: "y", DisplayText: ""},
	})

	// Second sync with different links — should replace
	idx.SyncWikiLinks("note/a-01abc", []WikiLinkEntry{
		{ToID: "z", Target: "z", DisplayText: "only this"},
	})

	stored, _ := idx.ListWikiLinks("note/a-01abc")
	if len(stored) != 1 {
		t.Fatalf("expected 1 wikilink after re-sync, got %d", len(stored))
	}
	if stored[0].Target != "z" {
		t.Errorf("target = %q, want %q", stored[0].Target, "z")
	}
}

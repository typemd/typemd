package core

import (
	"testing"
)

// ── parseFilter tests ────────────────────────────────────────────────────────

func TestParseFilter_Single(t *testing.T) {
	got, err := parseFilter("type=book")
	if err != nil {
		t.Fatalf("parseFilter() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Key != "type" || got[0].Value != "book" {
		t.Errorf("got[0] = %+v, want {type book}", got[0])
	}
}

func TestParseFilter_Multiple(t *testing.T) {
	got, err := parseFilter("type=book status=reading")
	if err != nil {
		t.Fatalf("parseFilter() error = %v", err)
	}
	want := []filterCondition{
		{Key: "type", Value: "book"},
		{Key: "status", Value: "reading"},
	}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestParseFilter_Empty(t *testing.T) {
	got, err := parseFilter("")
	if err != nil {
		t.Fatalf("parseFilter(\"\") error = %v", err)
	}
	if got != nil {
		t.Errorf("got = %v, want nil", got)
	}
}

func TestParseFilter_Whitespace(t *testing.T) {
	got, err := parseFilter("   ")
	if err != nil {
		t.Fatalf("parseFilter(whitespace) error = %v", err)
	}
	if got != nil {
		t.Errorf("got = %v, want nil", got)
	}
}

func TestParseFilter_Invalid_NoEquals(t *testing.T) {
	_, err := parseFilter("invalid-no-equals")
	if err == nil {
		t.Fatal("expected error for missing '=', got nil")
	}
}

func TestParseFilter_Invalid_LeadingEquals(t *testing.T) {
	_, err := parseFilter("=value")
	if err == nil {
		t.Fatal("expected error for leading '=', got nil")
	}
}

func TestParseFilter_ValueWithEquals(t *testing.T) {
	// value contains '='; only split on the first one
	got, err := parseFilter("key=a=b")
	if err != nil {
		t.Fatalf("parseFilter() error = %v", err)
	}
	if got[0].Key != "key" || got[0].Value != "a=b" {
		t.Errorf("got[0] = %+v, want {key a=b}", got[0])
	}
}

// ── QueryObjects tests ───────────────────────────────────────────────────────

func TestVault_QueryObjects_ByType(t *testing.T) {
	v := setupTestVault(t)

	v.NewObject("book", "book1", "")   //nolint:errcheck
	v.NewObject("book", "book2", "")   //nolint:errcheck
	v.NewObject("person", "alice", "") //nolint:errcheck

	results, err := v.QueryObjects("type=book")
	if err != nil {
		t.Fatalf("QueryObjects() error = %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("QueryObjects(type=book) len = %d, want 2", len(results))
	}
	for _, obj := range results {
		if obj.Type != "book" {
			t.Errorf("obj.Type = %q, want %q", obj.Type, "book")
		}
	}
}

func TestVault_QueryObjects_ByProperty(t *testing.T) {
	v := setupTestVault(t)

	b1, _ := v.NewObject("book", "book1", "")
	v.SetProperty(b1.ID, "status", "reading") //nolint:errcheck

	b2, _ := v.NewObject("book", "book2", "")
	v.SetProperty(b2.ID, "status", "done") //nolint:errcheck

	results, err := v.QueryObjects("type=book status=reading")
	if err != nil {
		t.Fatalf("QueryObjects() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("QueryObjects(status=reading) len = %d, want 1", len(results))
	}
	if results[0].ID != b1.ID {
		t.Errorf("results[0].ID = %q, want %q", results[0].ID, b1.ID)
	}
}

func TestVault_QueryObjects_NoResults(t *testing.T) {
	v := setupTestVault(t)

	v.NewObject("book", "book1", "") //nolint:errcheck

	results, err := v.QueryObjects("type=person")
	if err != nil {
		t.Fatalf("QueryObjects() error = %v", err)
	}
	if len(results) != 0 {
		t.Errorf("QueryObjects(type=person) len = %d, want 0", len(results))
	}
}

func TestVault_QueryObjects_EmptyFilter(t *testing.T) {
	v := setupTestVault(t)

	v.NewObject("book", "book1", "")   //nolint:errcheck
	v.NewObject("person", "alice", "") //nolint:errcheck

	results, err := v.QueryObjects("")
	if err != nil {
		t.Fatalf("QueryObjects(\"\") error = %v", err)
	}
	if len(results) != 2 {
		t.Errorf("QueryObjects(\"\") len = %d, want 2", len(results))
	}
}

func TestVault_QueryObjects_InvalidFilter(t *testing.T) {
	v := setupTestVault(t)

	_, err := v.QueryObjects("invalid-filter")
	if err == nil {
		t.Fatal("expected error for invalid filter, got nil")
	}
}

func TestVault_QueryObjects_DBNotOpen(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	_, err := v.QueryObjects("type=book")
	if err == nil {
		t.Fatal("expected error when DB not opened, got nil")
	}
}

func TestVault_QueryObjects_MultipleConditions(t *testing.T) {
	v := setupTestVault(t)

	b1, _ := v.NewObject("book", "b1", "")
	v.SetProperty(b1.ID, "status", "reading") //nolint:errcheck
	v.SetProperty(b1.ID, "title", "Go")       //nolint:errcheck

	b2, _ := v.NewObject("book", "b2", "")
	v.SetProperty(b2.ID, "status", "reading") //nolint:errcheck
	v.SetProperty(b2.ID, "title", "Rust")     //nolint:errcheck

	b3, _ := v.NewObject("book", "b3", "")
	v.SetProperty(b3.ID, "status", "done") //nolint:errcheck
	v.SetProperty(b3.ID, "title", "Go")    //nolint:errcheck

	// type=book status=reading title=Go — should match only b1
	results, err := v.QueryObjects("type=book status=reading title=Go")
	if err != nil {
		t.Fatalf("QueryObjects() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("QueryObjects(3 conditions) len = %d, want 1", len(results))
	}
	if results[0].ID != b1.ID {
		t.Errorf("results[0].ID = %q, want %q", results[0].ID, b1.ID)
	}
}

// ── SearchObjects tests ──────────────────────────────────────────────────────

func TestVault_SearchObjects_ByFilename(t *testing.T) {
	v := setupTestVault(t)

	concurrency, _ := v.NewObject("book", "concurrency-in-go", "")
	v.NewObject("book", "clean-code", "") //nolint:errcheck

	results, err := v.SearchObjects("concurrency")
	if err != nil {
		t.Fatalf("SearchObjects() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("SearchObjects(concurrency) len = %d, want 1", len(results))
	}
	if results[0].ID != concurrency.ID {
		t.Errorf("results[0].ID = %q, want %q", results[0].ID, concurrency.ID)
	}
}

func TestVault_SearchObjects_ByBody(t *testing.T) {
	v := setupTestVault(t)

	obj, err := v.NewObject("book", "mybook", "")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}
	obj.Body = "This book covers goroutines and channels."
	if err := v.SaveObject(obj); err != nil {
		t.Fatalf("saveObjectFile() error = %v", err)
	}

	v.NewObject("book", "other", "") //nolint:errcheck

	results, err := v.SearchObjects("goroutines")
	if err != nil {
		t.Fatalf("SearchObjects() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("SearchObjects(goroutines) len = %d, want 1", len(results))
	}
	if results[0].ID != obj.ID {
		t.Errorf("results[0].ID = %q, want %q", results[0].ID, obj.ID)
	}
}

func TestVault_SearchObjects_ByProperty(t *testing.T) {
	v := setupTestVault(t)

	b1, _ := v.NewObject("book", "book1", "")
	v.SetProperty(b1.ID, "title", "Mastering Go") //nolint:errcheck

	b2, _ := v.NewObject("book", "book2", "")
	v.SetProperty(b2.ID, "title", "Python Basics") //nolint:errcheck

	results, err := v.SearchObjects("Mastering")
	if err != nil {
		t.Fatalf("SearchObjects() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("SearchObjects(Mastering) len = %d, want 1", len(results))
	}
	if results[0].ID != b1.ID {
		t.Errorf("results[0].ID = %q, want %q", results[0].ID, b1.ID)
	}
}

func TestVault_SearchObjects_NoResults(t *testing.T) {
	v := setupTestVault(t)

	v.NewObject("book", "testbook", "") //nolint:errcheck

	results, err := v.SearchObjects("xyzzy")
	if err != nil {
		t.Fatalf("SearchObjects() error = %v", err)
	}
	if len(results) != 0 {
		t.Errorf("SearchObjects(xyzzy) len = %d, want 0", len(results))
	}
}

func TestVault_SearchObjects_EmptyKeyword(t *testing.T) {
	v := setupTestVault(t)

	v.NewObject("book", "book1", "") //nolint:errcheck

	results, err := v.SearchObjects("")
	if err != nil {
		t.Fatalf("SearchObjects(\"\") error = %v", err)
	}
	if results != nil {
		t.Errorf("SearchObjects(\"\") = %v, want nil", results)
	}
}

func TestVault_SearchObjects_DBNotOpen(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	_, err := v.SearchObjects("test")
	if err == nil {
		t.Fatal("expected error when DB not opened, got nil")
	}
}

func TestVault_SearchObjects_UpdateSync(t *testing.T) {
	v := setupTestVault(t)

	obj, err := v.NewObject("book", "evolving", "")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	// initial search should find nothing
	results, _ := v.SearchObjects("blockchain")
	if len(results) != 0 {
		t.Errorf("before update: expected 0 results, got %d", len(results))
	}

	// UPDATE should sync FTS via trigger
	obj.Body = "This chapter covers blockchain fundamentals."
	if err := v.SaveObject(obj); err != nil {
		t.Fatalf("saveObjectFile() error = %v", err)
	}

	results, err = v.SearchObjects("blockchain")
	if err != nil {
		t.Fatalf("SearchObjects() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("after update: expected 1 result, got %d", len(results))
	}
}

// ── RebuildIndex tests ───────────────────────────────────────────────────────

func TestVault_RebuildIndex(t *testing.T) {
	v := setupTestVault(t)

	obj, err := v.NewObject("book", "golang-in-action", "")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}
	obj.Body = "concurrency patterns in Go"
	if err := v.SaveObject(obj); err != nil {
		t.Fatalf("saveObjectFile() error = %v", err)
	}

	if err := v.RebuildIndex(); err != nil {
		t.Fatalf("RebuildIndex() error = %v", err)
	}

	results, err := v.SearchObjects("concurrency")
	if err != nil {
		t.Fatalf("SearchObjects() after rebuild error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("SearchObjects() after rebuild len = %d, want 1", len(results))
	}
}

func TestVault_RebuildIndex_DBNotOpen(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	err := v.RebuildIndex()
	if err == nil {
		t.Fatal("expected error when DB not opened, got nil")
	}
}

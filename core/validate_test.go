package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func TestValidateAllObjects_Valid(t *testing.T) {
	v := setupTestVault(t)
	schema := []byte("name: book\nproperties:\n  - name: title\n    type: string\n")
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), schema, 0644)

	obj, err := v.NewObject("book", "test-book", "")
	if err != nil {
		t.Fatalf("NewObject error = %v", err)
	}
	v.SetProperty(obj.ID, "title", "Go in Action")

	result := ValidateAllObjects(v)
	if errs, ok := result[obj.ID]; ok && len(errs) > 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateAllObjects_Invalid(t *testing.T) {
	v := setupTestVault(t)
	schema := []byte("name: book\nproperties:\n  - name: rating\n    type: number\n")
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), schema, 0644)

	obj, err := v.NewObject("book", "bad-book", "")
	if err != nil {
		t.Fatalf("NewObject error = %v", err)
	}
	obj.Properties["rating"] = "not-a-number"
	v.SaveObject(obj)

	result := ValidateAllObjects(v)
	if errs, ok := result[obj.ID]; !ok || len(errs) == 0 {
		t.Error("expected validation errors for invalid object")
	}
}

func TestValidateAllObjects_NoObjects(t *testing.T) {
	v := setupTestVault(t)
	result := ValidateAllObjects(v)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestValidateRelations_Valid(t *testing.T) {
	v := setupTestVault(t)
	schema := []byte("name: book\nproperties:\n  - name: author\n    type: relation\n    target: person\n")
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), schema, 0644)
	personSchema := []byte("name: person\nproperties:\n  - name: name\n    type: string\n")
	os.WriteFile(filepath.Join(v.TypesDir(), "person.yaml"), personSchema, 0644)

	book, _ := v.NewObject("book", "test-book", "")
	alice, _ := v.NewObject("person", "alice", "")
	v.LinkObjects(book.ID, "author", alice.ID)

	errs := ValidateRelations(v)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateRelations_OrphanedTarget(t *testing.T) {
	v := setupTestVault(t)
	v.db.Exec("INSERT INTO relations (name, from_id, to_id) VALUES (?, ?, ?)",
		"author", "book/test-book", "person/ghost")
	v.db.Exec("INSERT INTO objects (id, type, filename, properties, body) VALUES (?, ?, ?, ?, ?)",
		"book/test-book", "book", "test-book", "{}", "")

	errs := ValidateRelations(v)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateRelations_OrphanedSource(t *testing.T) {
	v := setupTestVault(t)
	v.db.Exec("INSERT INTO relations (name, from_id, to_id) VALUES (?, ?, ?)",
		"author", "book/ghost", "person/alice")
	v.db.Exec("INSERT INTO objects (id, type, filename, properties, body) VALUES (?, ?, ?, ?, ?)",
		"person/alice", "person", "alice", "{}", "")

	errs := ValidateRelations(v)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateAllObjects_UnknownType(t *testing.T) {
	v := setupTestVault(t)
	v.db.Exec("INSERT INTO objects (id, type, filename, properties, body) VALUES (?, ?, ?, ?, ?)",
		"unknown/test", "unknown", "test", "{}", "")

	result := ValidateAllObjects(v)
	if errs, ok := result["unknown/test"]; !ok || len(errs) == 0 {
		t.Error("expected error for object with unknown type")
	}
}

func TestValidateRelations_NoRelations(t *testing.T) {
	v := setupTestVault(t)
	errs := ValidateRelations(v)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateWikiLinks_NoBrokenLinks(t *testing.T) {
	v := setupTestVault(t)

	os.WriteFile(filepath.Join(v.TypesDir(), "note.yaml"),
		[]byte("name: note\nproperties:\n  - name: title\n    type: string\n"), 0644)

	noteA, _ := v.NewObject("note", "alpha", "")
	noteB, _ := v.NewObject("note", "beta", "")

	body := fmt.Sprintf("---\ntitle: Alpha\n---\n\nSee [[%s]].\n", noteB.ID)
	os.WriteFile(v.ObjectPath(noteA.Type, noteA.Filename), []byte(body), 0644)
	v.SyncIndex()

	errs := ValidateWikiLinks(v)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateWikiLinks_BrokenLink(t *testing.T) {
	v := setupTestVault(t)

	os.WriteFile(filepath.Join(v.TypesDir(), "note.yaml"),
		[]byte("name: note\nproperties:\n  - name: title\n    type: string\n"), 0644)

	note, _ := v.NewObject("note", "alpha", "")

	body := "---\ntitle: Alpha\n---\n\nSee [[note/nonexistent-01jjjjjjjjjjjjjjjjjjjjjjjj]].\n"
	os.WriteFile(v.ObjectPath(note.Type, note.Filename), []byte(body), 0644)
	v.SyncIndex()

	errs := ValidateWikiLinks(v)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Error(), "broken wiki-link") {
		t.Errorf("error = %q, want it to contain 'broken wiki-link'", errs[0])
	}
	if !strings.Contains(errs[0].Error(), "note/nonexistent") {
		t.Errorf("error = %q, want it to contain target", errs[0])
	}
}

func TestValidateWikiLinks_NoWikiLinks(t *testing.T) {
	v := setupTestVault(t)
	errs := ValidateWikiLinks(v)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

// ── Name uniqueness tests ────────────────────────────────────────────────────

func TestValidateNameUniqueness_NoDuplicates(t *testing.T) {
	v := setupTestVault(t)
	v.NewObject("tag", "go", "")
	v.NewObject("tag", "rust", "")

	errs := ValidateNameUniqueness(v)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateNameUniqueness_CaseSensitive(t *testing.T) {
	v := setupTestVault(t)
	v.NewObject("tag", "Go", "")
	v.NewObject("tag", "go", "")

	errs := ValidateNameUniqueness(v)
	if len(errs) != 0 {
		t.Errorf("expected no errors (case-sensitive), got %v", errs)
	}
}

func TestNewObject_TagNameDuplicate(t *testing.T) {
	v := setupTestVault(t)
	_, err := v.NewObject("tag", "go", "")
	if err != nil {
		t.Fatalf("first NewObject error = %v", err)
	}
	_, err = v.NewObject("tag", "go", "")
	if err == nil {
		t.Fatal("expected error for duplicate tag name, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error = %q, want it to contain 'already exists'", err)
	}
}

func TestNewObject_NonTagDuplicateNameAllowed(t *testing.T) {
	v := setupTestVault(t)
	_, err := v.NewObject("book", "test", "")
	if err != nil {
		t.Fatalf("first book error = %v", err)
	}
	_, err = v.NewObject("book", "test", "")
	if err != nil {
		t.Errorf("second book with same name should be allowed, got %v", err)
	}
}


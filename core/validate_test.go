package core

import (
	"os"
	"path/filepath"
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

	obj, err := v.NewObject("book", "test-book")
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

	obj, err := v.NewObject("book", "bad-book")
	if err != nil {
		t.Fatalf("NewObject error = %v", err)
	}
	obj.Properties["rating"] = "not-a-number"
	v.writeObjectProperties(obj)

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

	v.NewObject("book", "test-book")
	v.NewObject("person", "alice")
	v.LinkObjects("book/test-book", "author", "person/alice")

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

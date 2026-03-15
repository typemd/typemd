package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestRepo(t *testing.T) *LocalObjectRepository {
	t.Helper()
	dir := t.TempDir()

	// Create vault structure
	os.MkdirAll(filepath.Join(dir, ".typemd", "types"), 0755)
	os.MkdirAll(filepath.Join(dir, "objects"), 0755)

	return NewLocalObjectRepository(dir)
}

func TestLocalObjectRepository_GetAndSave(t *testing.T) {
	repo := setupTestRepo(t)

	// Create a test object file
	os.MkdirAll(filepath.Join(repo.root, "objects", "book"), 0755)
	os.WriteFile(
		filepath.Join(repo.root, "objects", "book", "test-01abc.md"),
		[]byte("---\nname: Test\nstatus: reading\n---\n\nHello world.\n"),
		0644,
	)

	obj, err := repo.Get("book/test-01abc")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if obj.ID != "book/test-01abc" {
		t.Errorf("ID = %q, want %q", obj.ID, "book/test-01abc")
	}
	if obj.Type != "book" {
		t.Errorf("Type = %q, want %q", obj.Type, "book")
	}
	if obj.Properties["name"] != "Test" {
		t.Errorf("name = %v, want %q", obj.Properties["name"], "Test")
	}
	if obj.Body != "Hello world.\n" {
		t.Errorf("Body = %q, want %q", obj.Body, "Hello world.\n")
	}

	// Save with modified body
	obj.Body = "Updated body.\n"
	if err := repo.Save(obj, []string{"name", "status"}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Re-read and verify
	obj2, _ := repo.Get("book/test-01abc")
	if obj2.Body != "Updated body.\n" {
		t.Errorf("Body after save = %q, want %q", obj2.Body, "Updated body.\n")
	}
}

func TestLocalObjectRepository_GetNonExistent(t *testing.T) {
	repo := setupTestRepo(t)

	_, err := repo.Get("book/nonexistent-01abc")
	if err == nil {
		t.Error("expected error for non-existent object")
	}
}

func TestLocalObjectRepository_Create(t *testing.T) {
	repo := setupTestRepo(t)
	os.MkdirAll(filepath.Join(repo.root, "objects", "book"), 0755)

	obj := &Object{
		ID:         "book/new-01abc",
		Type:       "book",
		Filename:   "new-01abc",
		Properties: map[string]any{"name": "New Book"},
		Body:       "Content here.\n",
	}

	if err := repo.Create(obj, []string{"name"}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	// Verify file exists
	got, err := repo.Get("book/new-01abc")
	if err != nil {
		t.Fatalf("Get after Create: %v", err)
	}
	if got.Properties["name"] != "New Book" {
		t.Errorf("name = %v, want %q", got.Properties["name"], "New Book")
	}
}

func TestLocalObjectRepository_CreateFailsIfExists(t *testing.T) {
	repo := setupTestRepo(t)
	os.MkdirAll(filepath.Join(repo.root, "objects", "book"), 0755)
	os.WriteFile(
		filepath.Join(repo.root, "objects", "book", "existing-01abc.md"),
		[]byte("---\nname: Existing\n---\n"),
		0644,
	)

	obj := &Object{
		ID:       "book/existing-01abc",
		Type:     "book",
		Filename: "existing-01abc",
		Properties: map[string]any{"name": "Dupe"},
	}

	err := repo.Create(obj, nil)
	if err == nil {
		t.Error("expected error when creating existing file")
	}
}

func TestLocalObjectRepository_Walk(t *testing.T) {
	repo := setupTestRepo(t)

	// Create a few object files
	os.MkdirAll(filepath.Join(repo.root, "objects", "book"), 0755)
	os.MkdirAll(filepath.Join(repo.root, "objects", "note"), 0755)
	os.WriteFile(filepath.Join(repo.root, "objects", "book", "a-01abc.md"), []byte("---\nname: A\n---\n"), 0644)
	os.WriteFile(filepath.Join(repo.root, "objects", "note", "b-01abc.md"), []byte("---\nname: B\n---\n"), 0644)
	objects, err := repo.Walk()
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(objects) != 2 {
		t.Fatalf("Walk returned %d objects, want 2", len(objects))
	}
}

func TestLocalObjectRepository_WalkEmptyVault(t *testing.T) {
	repo := setupTestRepo(t)
	// Don't create objects dir
	os.RemoveAll(filepath.Join(repo.root, "objects"))

	objects, err := repo.Walk()
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if objects != nil {
		t.Errorf("Walk on missing dir should return nil, got %d objects", len(objects))
	}
}

func TestLocalObjectRepository_GlobIDs(t *testing.T) {
	repo := setupTestRepo(t)

	os.MkdirAll(filepath.Join(repo.root, "objects", "book"), 0755)
	os.WriteFile(filepath.Join(repo.root, "objects", "book", "clean-code-01abc.md"), []byte("---\n---\n"), 0644)
	os.WriteFile(filepath.Join(repo.root, "objects", "book", "clean-arch-02def.md"), []byte("---\n---\n"), 0644)
	os.WriteFile(filepath.Join(repo.root, "objects", "book", "golang-03ghi.md"), []byte("---\n---\n"), 0644)

	ids, err := repo.GlobIDs("book", "clean")
	if err != nil {
		t.Fatalf("GlobIDs: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("GlobIDs returned %d, want 2", len(ids))
	}
}

func TestLocalObjectRepository_ModTime(t *testing.T) {
	repo := setupTestRepo(t)

	os.MkdirAll(filepath.Join(repo.root, "objects", "book"), 0755)
	os.WriteFile(filepath.Join(repo.root, "objects", "book", "test-01abc.md"), []byte("---\n---\n"), 0644)

	mt, err := repo.ModTime("book/test-01abc")
	if err != nil {
		t.Fatalf("ModTime: %v", err)
	}
	if mt.IsZero() {
		t.Error("ModTime returned zero time")
	}
	if time.Since(mt) > 5*time.Second {
		t.Error("ModTime seems too old for a just-created file")
	}
}

func TestLocalObjectRepository_EnsureDir(t *testing.T) {
	repo := setupTestRepo(t)

	if err := repo.EnsureDir("newtype"); err != nil {
		t.Fatalf("EnsureDir: %v", err)
	}

	info, err := os.Stat(filepath.Join(repo.root, "objects", "newtype"))
	if err != nil {
		t.Fatalf("directory should exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected a directory")
	}
}

func TestLocalObjectRepository_GetSchema(t *testing.T) {
	repo := setupTestRepo(t)

	os.WriteFile(
		filepath.Join(repo.root, ".typemd", "types", "book.yaml"),
		[]byte("name: book\nproperties:\n  - name: title\n    type: string\n  - name: rating\n    type: number\n"),
		0644,
	)

	schema, err := repo.GetSchema("book")
	if err != nil {
		t.Fatalf("GetSchema: %v", err)
	}
	if schema.Name != "book" {
		t.Errorf("Name = %q, want %q", schema.Name, "book")
	}
	if len(schema.Properties) != 2 {
		t.Errorf("Properties count = %d, want 2", len(schema.Properties))
	}
}

func TestLocalObjectRepository_GetSchemaBuiltIn(t *testing.T) {
	repo := setupTestRepo(t)

	schema, err := repo.GetSchema("tag")
	if err != nil {
		t.Fatalf("GetSchema(tag): %v", err)
	}
	if schema.Name != "tag" {
		t.Errorf("Name = %q, want %q", schema.Name, "tag")
	}
}

func TestLocalObjectRepository_GetSchemaUnknown(t *testing.T) {
	repo := setupTestRepo(t)

	_, err := repo.GetSchema("nonexistent")
	if err == nil {
		t.Error("expected error for unknown type")
	}
}

func TestLocalObjectRepository_ListSchemas(t *testing.T) {
	repo := setupTestRepo(t)

	os.WriteFile(filepath.Join(repo.root, ".typemd", "types", "book.yaml"), []byte("name: book\n"), 0644)
	os.WriteFile(filepath.Join(repo.root, ".typemd", "types", "note.yaml"), []byte("name: note\n"), 0644)

	names, err := repo.ListSchemas()
	if err != nil {
		t.Fatalf("ListSchemas: %v", err)
	}
	// Should include custom types + built-in tag
	if len(names) < 3 {
		t.Errorf("ListSchemas returned %d, want at least 3 (book, note, tag)", len(names))
	}
}

func TestLocalObjectRepository_GetTemplate(t *testing.T) {
	repo := setupTestRepo(t)

	os.MkdirAll(filepath.Join(repo.root, "templates", "book"), 0755)
	os.WriteFile(
		filepath.Join(repo.root, "templates", "book", "fiction.md"),
		[]byte("---\ngenre: fiction\n---\n\nOnce upon a time.\n"),
		0644,
	)

	tmpl, err := repo.GetTemplate("book", "fiction")
	if err != nil {
		t.Fatalf("GetTemplate: %v", err)
	}
	if tmpl.Name != "fiction" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "fiction")
	}
	if tmpl.Properties["genre"] != "fiction" {
		t.Errorf("genre = %v, want %q", tmpl.Properties["genre"], "fiction")
	}
	if tmpl.Body != "Once upon a time.\n" {
		t.Errorf("Body = %q, want %q", tmpl.Body, "Once upon a time.\n")
	}
}

func TestLocalObjectRepository_ListTemplates(t *testing.T) {
	repo := setupTestRepo(t)

	os.MkdirAll(filepath.Join(repo.root, "templates", "book"), 0755)
	os.WriteFile(filepath.Join(repo.root, "templates", "book", "fiction.md"), []byte("---\n---\n"), 0644)
	os.WriteFile(filepath.Join(repo.root, "templates", "book", "nonfic.md"), []byte("---\n---\n"), 0644)

	names, err := repo.ListTemplates("book")
	if err != nil {
		t.Fatalf("ListTemplates: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("ListTemplates returned %d, want 2", len(names))
	}
}

func TestLocalObjectRepository_ListTemplatesNoDir(t *testing.T) {
	repo := setupTestRepo(t)

	names, err := repo.ListTemplates("book")
	if err != nil {
		t.Fatalf("ListTemplates: %v", err)
	}
	if names != nil {
		t.Errorf("expected nil for missing templates dir, got %v", names)
	}
}

func TestLocalObjectRepository_GetSharedProperties(t *testing.T) {
	repo := setupTestRepo(t)

	os.WriteFile(
		filepath.Join(repo.root, ".typemd", "properties.yaml"),
		[]byte("properties:\n  - name: status\n    type: select\n    options:\n      - value: active\n"),
		0644,
	)

	props, err := repo.GetSharedProperties()
	if err != nil {
		t.Fatalf("GetSharedProperties: %v", err)
	}
	if len(props) != 1 {
		t.Fatalf("expected 1 shared property, got %d", len(props))
	}
	if props[0].Name != "status" {
		t.Errorf("name = %q, want %q", props[0].Name, "status")
	}

	// Second call should use cache
	props2, _ := repo.GetSharedProperties()
	if len(props2) != 1 {
		t.Error("cached result should match")
	}
}

func TestLocalObjectRepository_GetSharedPropertiesNoFile(t *testing.T) {
	repo := setupTestRepo(t)

	props, err := repo.GetSharedProperties()
	if err != nil {
		t.Fatalf("GetSharedProperties: %v", err)
	}
	if props != nil {
		t.Errorf("expected nil for missing file, got %v", props)
	}
}

// --- WalkAll tests ---

func TestWalkAll_EmptyDir(t *testing.T) {
	repo := setupTestRepo(t)
	// Remove objects dir so it doesn't exist
	os.RemoveAll(filepath.Join(repo.root, "objects"))

	objects, corrupted, err := repo.WalkAll()
	if err != nil {
		t.Fatalf("WalkAll: %v", err)
	}
	if objects != nil {
		t.Errorf("expected nil objects, got %d", len(objects))
	}
	if corrupted != nil {
		t.Errorf("expected nil corrupted, got %d", len(corrupted))
	}
}

func TestWalkAll_ValidFiles(t *testing.T) {
	repo := setupTestRepo(t)

	os.MkdirAll(filepath.Join(repo.root, "objects", "book"), 0755)
	os.MkdirAll(filepath.Join(repo.root, "objects", "note"), 0755)
	os.WriteFile(filepath.Join(repo.root, "objects", "book", "a-01abc.md"), []byte("---\nname: A\n---\n"), 0644)
	os.WriteFile(filepath.Join(repo.root, "objects", "note", "b-01abc.md"), []byte("---\nname: B\n---\n"), 0644)

	objects, corrupted, err := repo.WalkAll()
	if err != nil {
		t.Fatalf("WalkAll: %v", err)
	}
	if len(objects) != 2 {
		t.Errorf("expected 2 objects, got %d", len(objects))
	}
	if len(corrupted) != 0 {
		t.Errorf("expected 0 corrupted, got %d", len(corrupted))
	}
}

func TestWalkAll_MixedFiles(t *testing.T) {
	repo := setupTestRepo(t)

	os.MkdirAll(filepath.Join(repo.root, "objects", "book"), 0755)
	// Valid file
	os.WriteFile(filepath.Join(repo.root, "objects", "book", "good-01abc.md"), []byte("---\nname: Good\n---\n"), 0644)
	// Corrupted file: invalid YAML
	os.WriteFile(filepath.Join(repo.root, "objects", "book", "bad-02def.md"), []byte("---\n: :\n  bad yaml[[\n---\n"), 0644)

	objects, corrupted, err := repo.WalkAll()
	if err != nil {
		t.Fatalf("WalkAll: %v", err)
	}
	if len(objects) != 1 {
		t.Errorf("expected 1 valid object, got %d", len(objects))
	}
	if len(corrupted) != 1 {
		t.Errorf("expected 1 corrupted file, got %d", len(corrupted))
	}
	if len(corrupted) > 0 {
		expected := filepath.Join("book", "bad-02def.md")
		if corrupted[0].Path != expected {
			t.Errorf("corrupted path = %q, want %q", corrupted[0].Path, expected)
		}
		if corrupted[0].Error == nil {
			t.Error("expected non-nil error on corrupted file")
		}
	}
}

func TestWalkAll_MalformedYAML(t *testing.T) {
	repo := setupTestRepo(t)

	os.MkdirAll(filepath.Join(repo.root, "objects", "note"), 0755)
	// Frontmatter delimiters present, but YAML content is malformed
	os.WriteFile(filepath.Join(repo.root, "objects", "note", "broken-01abc.md"), []byte("---\n[invalid: yaml: content\n---\n"), 0644)

	objects, corrupted, err := repo.WalkAll()
	if err != nil {
		t.Fatalf("WalkAll: %v", err)
	}
	if len(objects) != 0 {
		t.Errorf("expected 0 objects, got %d", len(objects))
	}
	if len(corrupted) != 1 {
		t.Fatalf("expected 1 corrupted, got %d", len(corrupted))
	}
	if corrupted[0].Error == nil {
		t.Error("expected non-nil error for malformed YAML")
	}
}

func TestWalkAll_NoFrontmatterDelimiters(t *testing.T) {
	repo := setupTestRepo(t)

	os.MkdirAll(filepath.Join(repo.root, "objects", "book"), 0755)
	// File with no frontmatter delimiters at all
	os.WriteFile(filepath.Join(repo.root, "objects", "book", "nofm-01abc.md"), []byte("Just some plain text with no frontmatter.\n"), 0644)

	objects, corrupted, err := repo.WalkAll()
	if err != nil {
		t.Fatalf("WalkAll: %v", err)
	}
	if len(objects) != 0 {
		t.Errorf("expected 0 objects, got %d", len(objects))
	}
	if len(corrupted) != 1 {
		t.Fatalf("expected 1 corrupted, got %d", len(corrupted))
	}
	expected := filepath.Join("book", "nofm-01abc.md")
	if corrupted[0].Path != expected {
		t.Errorf("corrupted path = %q, want %q", corrupted[0].Path, expected)
	}
	if corrupted[0].Error == nil {
		t.Error("expected non-nil error for file without frontmatter delimiters")
	}
}

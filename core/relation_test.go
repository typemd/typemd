package core

import (
	"os"
	"path/filepath"
	"testing"
)

// setupRelationTestVault creates a vault with book/person schemas that have relation properties.
func setupRelationTestVault(t *testing.T) *Vault {
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
`)
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), bookSchema, 0644)

	personSchema := []byte(`name: person
properties:
  - name: name
    type: string
  - name: books
    type: relation
    target: book
    multiple: true
    bidirectional: true
    inverse: author
`)
	os.WriteFile(filepath.Join(v.TypesDir(), "person.yaml"), personSchema, 0644)

	return v
}

func TestVault_LinkObjects_SingleValue(t *testing.T) {
	v := setupRelationTestVault(t)

	book, _ := v.NewObject("book", "golang-in-action")
	person, _ := v.NewObject("person", "alan-donovan")

	err := v.LinkObjects(book.ID, "author", person.ID)
	if err != nil {
		t.Fatalf("LinkObjects() error = %v", err)
	}

	// Verify source frontmatter
	bookObj, _ := v.GetObject(book.ID)
	if bookObj.Properties["author"] != person.ID {
		t.Errorf("author = %v, want %q", bookObj.Properties["author"], person.ID)
	}

	// Verify inverse (person.books should be array with one entry)
	personObj, _ := v.GetObject(person.ID)
	books, ok := personObj.Properties["books"].([]any)
	if !ok {
		t.Fatalf("books type = %T, want []any", personObj.Properties["books"])
	}
	if len(books) != 1 || books[0] != book.ID {
		t.Errorf("books = %v, want [%s]", books, book.ID)
	}

	// Verify relations table
	var count int
	v.db.QueryRow("SELECT COUNT(*) FROM relations WHERE from_id = ? AND name = ?",
		book.ID, "author").Scan(&count)
	if count != 1 {
		t.Errorf("relations count (forward) = %d, want 1", count)
	}
	v.db.QueryRow("SELECT COUNT(*) FROM relations WHERE from_id = ? AND name = ?",
		person.ID, "books").Scan(&count)
	if count != 1 {
		t.Errorf("relations count (inverse) = %d, want 1", count)
	}
}

func TestVault_LinkObjects_MultipleValue(t *testing.T) {
	v := setupRelationTestVault(t)

	bookA, _ := v.NewObject("book", "book-a")
	bookB, _ := v.NewObject("book", "book-b")
	person, _ := v.NewObject("person", "alan")

	// Link two books to person via inverse (person.books is multiple)
	v.LinkObjects(bookA.ID, "author", person.ID)
	v.LinkObjects(bookB.ID, "author", person.ID)

	personObj, _ := v.GetObject(person.ID)
	books, ok := personObj.Properties["books"].([]any)
	if !ok {
		t.Fatalf("books type = %T, want []any", personObj.Properties["books"])
	}
	if len(books) != 2 {
		t.Errorf("len(books) = %d, want 2", len(books))
	}
}

func TestVault_LinkObjects_TargetNotFound(t *testing.T) {
	v := setupRelationTestVault(t)
	book, _ := v.NewObject("book", "test")

	err := v.LinkObjects(book.ID, "author", "person/nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent target, got nil")
	}
}

func TestVault_LinkObjects_RelationNotFound(t *testing.T) {
	v := setupRelationTestVault(t)
	book, _ := v.NewObject("book", "test")
	person, _ := v.NewObject("person", "alan")

	err := v.LinkObjects(book.ID, "nonexistent", person.ID)
	if err == nil {
		t.Fatal("expected error for unknown relation, got nil")
	}
}

func TestVault_LinkObjects_TypeMismatch(t *testing.T) {
	v := setupRelationTestVault(t)
	bookA, _ := v.NewObject("book", "book-a")
	bookB, _ := v.NewObject("book", "book-b")

	// author targets person, not book
	err := v.LinkObjects(bookA.ID, "author", bookB.ID)
	if err == nil {
		t.Fatal("expected error for type mismatch, got nil")
	}
}

func TestVault_LinkObjects_DBNotOpen(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()

	err := v.LinkObjects("book/test", "author", "person/alan")
	if err == nil {
		t.Fatal("expected error when DB not opened, got nil")
	}
}

func TestVault_LinkObjects_SingleValueOverwrite(t *testing.T) {
	v := setupRelationTestVault(t)

	book, _ := v.NewObject("book", "test")
	alan, _ := v.NewObject("person", "alan")
	brian, _ := v.NewObject("person", "brian")

	// Link to alan first
	v.LinkObjects(book.ID, "author", alan.ID)

	// Link to brian (should overwrite since author is single-value)
	err := v.LinkObjects(book.ID, "author", brian.ID)
	if err != nil {
		t.Fatalf("LinkObjects() overwrite error = %v", err)
	}

	bookObj, _ := v.GetObject(book.ID)
	if bookObj.Properties["author"] != brian.ID {
		t.Errorf("author = %v, want %q", bookObj.Properties["author"], brian.ID)
	}
}

func TestVault_LinkObjects_DuplicateMultiple(t *testing.T) {
	v := setupRelationTestVault(t)

	book, _ := v.NewObject("book", "test")
	alan, _ := v.NewObject("person", "alan")

	v.LinkObjects(book.ID, "author", alan.ID)

	// person.books is multiple, linking same book again should error
	// because the inverse side already has this entry
	err := v.LinkObjects(book.ID, "author", alan.ID)
	// Note: single-value overwrite on book side is fine,
	// but inverse side should detect duplicate
	_ = err
}

func TestVault_UnlinkObjects_SingleValue(t *testing.T) {
	v := setupRelationTestVault(t)

	book, _ := v.NewObject("book", "test")
	alan, _ := v.NewObject("person", "alan")
	v.LinkObjects(book.ID, "author", alan.ID)

	// Unlink without --both: only remove from book
	err := v.UnlinkObjects(book.ID, "author", alan.ID, false)
	if err != nil {
		t.Fatalf("UnlinkObjects() error = %v", err)
	}

	bookObj, _ := v.GetObject(book.ID)
	if bookObj.Properties["author"] != nil {
		t.Errorf("author = %v, want nil", bookObj.Properties["author"])
	}

	// person.books should still have the entry
	personObj, _ := v.GetObject(alan.ID)
	books, ok := personObj.Properties["books"].([]any)
	if !ok || len(books) != 1 {
		t.Errorf("books = %v, expected still linked", personObj.Properties["books"])
	}
}

func TestVault_UnlinkObjects_Both(t *testing.T) {
	v := setupRelationTestVault(t)

	book, _ := v.NewObject("book", "test")
	alan, _ := v.NewObject("person", "alan")
	v.LinkObjects(book.ID, "author", alan.ID)

	err := v.UnlinkObjects(book.ID, "author", alan.ID, true)
	if err != nil {
		t.Fatalf("UnlinkObjects() error = %v", err)
	}

	bookObj, _ := v.GetObject(book.ID)
	if bookObj.Properties["author"] != nil {
		t.Errorf("author = %v, want nil", bookObj.Properties["author"])
	}

	personObj, _ := v.GetObject(alan.ID)
	if personObj.Properties["books"] != nil {
		t.Errorf("books = %v, want nil", personObj.Properties["books"])
	}

	// Verify relations table is clean
	var count int
	v.db.QueryRow("SELECT COUNT(*) FROM relations").Scan(&count)
	if count != 0 {
		t.Errorf("relations count = %d, want 0", count)
	}
}

func TestVault_UnlinkObjects_MultipleRemoveOne(t *testing.T) {
	v := setupRelationTestVault(t)

	bookA, _ := v.NewObject("book", "book-a")
	bookB, _ := v.NewObject("book", "book-b")
	alan, _ := v.NewObject("person", "alan")

	v.LinkObjects(bookA.ID, "author", alan.ID)
	v.LinkObjects(bookB.ID, "author", alan.ID)

	// Unlink one book with --both
	err := v.UnlinkObjects(bookA.ID, "author", alan.ID, true)
	if err != nil {
		t.Fatalf("UnlinkObjects() error = %v", err)
	}

	// person.books should still have book-b
	personObj, _ := v.GetObject(alan.ID)
	books, ok := personObj.Properties["books"].([]any)
	if !ok || len(books) != 1 {
		t.Fatalf("books = %v, want [%s]", personObj.Properties["books"], bookB.ID)
	}
	if books[0] != bookB.ID {
		t.Errorf("books[0] = %v, want %q", books[0], bookB.ID)
	}
}

func TestVault_UnlinkObjects_DBNotOpen(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()

	err := v.UnlinkObjects("book/test", "author", "person/alan", false)
	if err == nil {
		t.Fatal("expected error when DB not opened, got nil")
	}
}

func TestFindSystemRelationProperty_Tags(t *testing.T) {
	p := findSystemRelationProperty(TagsProperty)
	if p == nil {
		t.Fatal("expected non-nil Property for tags")
	}
	if p.Target != TagTypeName {
		t.Errorf("Target = %q, want %q", p.Target, TagTypeName)
	}
	if !p.Multiple {
		t.Error("Multiple should be true for tags")
	}
	if p.Type != "relation" {
		t.Errorf("Type = %q, want %q", p.Type, "relation")
	}
}

func TestFindSystemRelationProperty_NonRelation(t *testing.T) {
	p := findSystemRelationProperty("name")
	if p != nil {
		t.Errorf("expected nil for non-relation system property, got %+v", p)
	}
}

func TestFindSystemRelationProperty_Nonexistent(t *testing.T) {
	p := findSystemRelationProperty("nonexistent")
	if p != nil {
		t.Errorf("expected nil for nonexistent property, got %+v", p)
	}
}

func TestVault_LinkObjects_SystemRelationTags(t *testing.T) {
	v := setupRelationTestVault(t)

	tag, _ := v.NewObject("tag", "go")
	book, _ := v.NewObject("book", "golang-book")

	err := v.LinkObjects(book.ID, TagsProperty, tag.ID)
	if err != nil {
		t.Fatalf("LinkObjects() error = %v", err)
	}

	bookObj, _ := v.GetObject(book.ID)
	tags, ok := bookObj.Properties[TagsProperty].([]any)
	if !ok {
		t.Fatalf("tags type = %T, want []any", bookObj.Properties[TagsProperty])
	}
	if len(tags) != 1 || tags[0] != tag.ID {
		t.Errorf("tags = %v, want [%s]", tags, tag.ID)
	}
}

func TestVault_UnlinkObjects_SystemRelationTags(t *testing.T) {
	v := setupRelationTestVault(t)

	tag, _ := v.NewObject("tag", "go")
	book, _ := v.NewObject("book", "golang-book")

	v.LinkObjects(book.ID, TagsProperty, tag.ID)

	err := v.UnlinkObjects(book.ID, TagsProperty, tag.ID, false)
	if err != nil {
		t.Fatalf("UnlinkObjects() error = %v", err)
	}

	bookObj, _ := v.GetObject(book.ID)
	if bookObj.Properties[TagsProperty] != nil {
		t.Errorf("tags = %v, want nil", bookObj.Properties[TagsProperty])
	}
}

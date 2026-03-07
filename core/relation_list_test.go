package core

import "testing"

func TestVault_ListRelations_Empty(t *testing.T) {
	v := setupRelationTestVault(t)
	book, _ := v.NewObject("book", "test")

	rels, err := v.ListRelations(book.ID)
	if err != nil {
		t.Fatalf("ListRelations() error = %v", err)
	}
	if len(rels) != 0 {
		t.Errorf("len(rels) = %d, want 0", len(rels))
	}
}

func TestVault_ListRelations_Forward(t *testing.T) {
	v := setupRelationTestVault(t)
	book, _ := v.NewObject("book", "test")
	alan, _ := v.NewObject("person", "alan")
	v.LinkObjects(book.ID, "author", alan.ID)

	// Bidirectional link creates 2 DB rows: forward (author) + inverse (books)
	// ListRelations returns both since book is from_id on forward and to_id on inverse
	rels, err := v.ListRelations(book.ID)
	if err != nil {
		t.Fatalf("ListRelations() error = %v", err)
	}
	if len(rels) != 2 {
		t.Fatalf("len(rels) = %d, want 2", len(rels))
	}

	// Verify forward relation exists
	found := false
	for _, r := range rels {
		if r.Name == "author" && r.FromID == book.ID && r.ToID == alan.ID {
			found = true
		}
	}
	if !found {
		t.Errorf("expected forward relation {author, %s, %s} in %+v", book.ID, alan.ID, rels)
	}
}

func TestVault_ListRelations_BothDirections(t *testing.T) {
	v := setupRelationTestVault(t)
	book, _ := v.NewObject("book", "test")
	alan, _ := v.NewObject("person", "alan")
	v.LinkObjects(book.ID, "author", alan.ID)

	// person/alan should have both: inverse (books, from_id=person/alan) + forward (author, to_id=person/alan)
	rels, err := v.ListRelations(alan.ID)
	if err != nil {
		t.Fatalf("ListRelations() error = %v", err)
	}
	if len(rels) != 2 {
		t.Fatalf("len(rels) = %d, want 2", len(rels))
	}

	// Verify inverse relation exists
	found := false
	for _, r := range rels {
		if r.Name == "books" && r.FromID == alan.ID && r.ToID == book.ID {
			found = true
		}
	}
	if !found {
		t.Errorf("expected inverse relation {books, %s, %s} in %+v", alan.ID, book.ID, rels)
	}
}

func TestVault_ListRelations_DBNotOpen(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()

	_, err := v.ListRelations("book/test")
	if err == nil {
		t.Fatal("expected error when DB not opened, got nil")
	}
}

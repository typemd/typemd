package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestVault_SyncIndex_WikiLinks(t *testing.T) {
	v := setupTestVault(t)

	// Create type schemas
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"),
		[]byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)
	os.WriteFile(filepath.Join(v.TypesDir(), "person.yaml"),
		[]byte("name: person\nproperties:\n  - name: name\n    type: string\n"), 0644)

	// Create objects
	book, _ := v.NewObject("book", "clean-code", "")
	person, _ := v.NewObject("person", "robert-martin", "")

	// Edit book's body to include a wiki-link using full object ID
	bookPath := v.ObjectPath(book.Type, book.Filename)
	body := fmt.Sprintf("---\ntitle: Clean Code\n---\n\nWritten by [[%s]].\n", person.ID)
	os.WriteFile(bookPath, []byte(body), 0644)

	// Sync
	if _, err := v.SyncIndex(); err != nil {
		t.Fatalf("SyncIndex() error = %v", err)
	}

	// Verify wikilinks table has the link
	links, err := v.ListWikiLinks(book.ID)
	if err != nil {
		t.Fatalf("ListWikiLinks() error = %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("got %d wikilinks, want 1", len(links))
	}
	if links[0].Target != person.ID {
		t.Errorf("Target = %q, want %q", links[0].Target, person.ID)
	}
	if links[0].ToID != person.ID {
		t.Errorf("ToID = %q, want %q", links[0].ToID, person.ID)
	}

	// Verify backlinks on the person
	backlinks, err := v.ListBacklinks(person.ID)
	if err != nil {
		t.Fatalf("ListBacklinks() error = %v", err)
	}
	if len(backlinks) != 1 {
		t.Fatalf("got %d backlinks, want 1", len(backlinks))
	}
	if backlinks[0].FromID != book.ID {
		t.Errorf("FromID = %q, want %q", backlinks[0].FromID, book.ID)
	}
}

func TestVault_SyncIndex_WikiLinks_BrokenLink(t *testing.T) {
	v := setupTestVault(t)

	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"),
		[]byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)

	book, _ := v.NewObject("book", "clean-code", "")

	// Link to a non-existent object (full ID that doesn't exist)
	bookPath := v.ObjectPath(book.Type, book.Filename)
	os.WriteFile(bookPath, []byte("---\ntitle: Clean Code\n---\n\nSee [[person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj]].\n"), 0644)

	if _, err := v.SyncIndex(); err != nil {
		t.Fatalf("SyncIndex() error = %v", err)
	}

	links, err := v.ListWikiLinks(book.ID)
	if err != nil {
		t.Fatalf("ListWikiLinks() error = %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("got %d wikilinks, want 1", len(links))
	}
	if links[0].ToID != "" {
		t.Errorf("ToID = %q, want empty (broken link)", links[0].ToID)
	}
	if links[0].Target != "person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj" {
		t.Errorf("Target = %q, want %q", links[0].Target, "person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj")
	}
}

func TestVault_SyncIndex_WikiLinks_Updated(t *testing.T) {
	v := setupTestVault(t)

	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"),
		[]byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)
	os.WriteFile(filepath.Join(v.TypesDir(), "person.yaml"),
		[]byte("name: person\nproperties:\n  - name: name\n    type: string\n"), 0644)

	book, _ := v.NewObject("book", "clean-code", "")
	personA, _ := v.NewObject("person", "alice", "")
	personB, _ := v.NewObject("person", "bob", "")

	// First sync: link to alice
	bookPath := v.ObjectPath(book.Type, book.Filename)
	os.WriteFile(bookPath, []byte(fmt.Sprintf("---\ntitle: Clean Code\n---\n\nBy [[%s]].\n", personA.ID)), 0644)
	v.SyncIndex()

	// Second sync: change link to bob
	os.WriteFile(bookPath, []byte(fmt.Sprintf("---\ntitle: Clean Code\n---\n\nBy [[%s]].\n", personB.ID)), 0644)
	v.SyncIndex()

	links, err := v.ListWikiLinks(book.ID)
	if err != nil {
		t.Fatalf("ListWikiLinks() error = %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("got %d wikilinks, want 1", len(links))
	}
	if links[0].ToID != personB.ID {
		t.Errorf("ToID = %q, want %q", links[0].ToID, personB.ID)
	}

	// Alice should have no backlinks
	backlinks, err := v.ListBacklinks(personA.ID)
	if err != nil {
		t.Fatalf("ListBacklinks() error = %v", err)
	}
	if len(backlinks) != 0 {
		t.Errorf("alice backlinks = %d, want 0", len(backlinks))
	}
}

func TestVault_SyncIndex_WikiLinks_DisplayText(t *testing.T) {
	v := setupTestVault(t)

	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"),
		[]byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)
	os.WriteFile(filepath.Join(v.TypesDir(), "person.yaml"),
		[]byte("name: person\nproperties:\n  - name: name\n    type: string\n"), 0644)

	book, _ := v.NewObject("book", "clean-code", "")
	person, _ := v.NewObject("person", "robert-martin", "")

	bookPath := v.ObjectPath(book.Type, book.Filename)
	body := fmt.Sprintf("---\ntitle: Clean Code\n---\n\nBy [[%s|Uncle Bob]].\n", person.ID)
	os.WriteFile(bookPath, []byte(body), 0644)

	v.SyncIndex()

	links, err := v.ListWikiLinks(book.ID)
	if err != nil {
		t.Fatalf("ListWikiLinks() error = %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("got %d wikilinks, want 1", len(links))
	}
	if links[0].DisplayText != "Uncle Bob" {
		t.Errorf("DisplayText = %q, want %q", links[0].DisplayText, "Uncle Bob")
	}
}

func TestVault_SyncIndex_WikiLinks_CleanedOnObjectDeletion(t *testing.T) {
	v := setupTestVault(t)

	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"),
		[]byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)
	os.WriteFile(filepath.Join(v.TypesDir(), "person.yaml"),
		[]byte("name: person\nproperties:\n  - name: name\n    type: string\n"), 0644)

	book, _ := v.NewObject("book", "clean-code", "")
	person, _ := v.NewObject("person", "robert-martin", "")

	bookPath := v.ObjectPath(book.Type, book.Filename)
	body := fmt.Sprintf("---\ntitle: Clean Code\n---\n\nBy [[%s]].\n", person.ID)
	os.WriteFile(bookPath, []byte(body), 0644)
	v.SyncIndex()

	// Delete the book
	os.Remove(bookPath)
	v.SyncIndex()

	// Person should have no backlinks
	backlinks, err := v.ListBacklinks(person.ID)
	if err != nil {
		t.Fatalf("ListBacklinks() error = %v", err)
	}
	if len(backlinks) != 0 {
		t.Errorf("backlinks = %d, want 0 after source deleted", len(backlinks))
	}

	// Wikilinks table should have no entries for the deleted book
	var count int
	v.db.QueryRow("SELECT COUNT(*) FROM wikilinks WHERE from_id = ?", book.ID).Scan(&count)
	if count != 0 {
		t.Errorf("wikilinks count = %d, want 0", count)
	}
}

func TestVault_ListWikiLinks_DBNotOpen(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()

	_, err := v.ListWikiLinks("book/test")
	if err == nil {
		t.Fatal("expected error when DB not opened")
	}
}

func TestVault_ListBacklinks_DBNotOpen(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()

	_, err := v.ListBacklinks("book/test")
	if err == nil {
		t.Fatal("expected error when DB not opened")
	}
}

func TestVault_ListBacklinks_MultipleSourcesSorted(t *testing.T) {
	v := setupTestVault(t)

	os.WriteFile(filepath.Join(v.TypesDir(), "note.yaml"),
		[]byte("name: note\nproperties:\n  - name: title\n    type: string\n"), 0644)

	target, _ := v.NewObject("note", "target", "")
	noteA, _ := v.NewObject("note", "alpha", "")
	noteB, _ := v.NewObject("note", "beta", "")

	// Both notes link to target using full ID
	os.WriteFile(v.ObjectPath(noteA.Type, noteA.Filename),
		[]byte(fmt.Sprintf("---\ntitle: Alpha\n---\n\nSee [[%s]].\n", target.ID)), 0644)
	os.WriteFile(v.ObjectPath(noteB.Type, noteB.Filename),
		[]byte(fmt.Sprintf("---\ntitle: Beta\n---\n\nSee [[%s]].\n", target.ID)), 0644)

	v.SyncIndex()

	backlinks, err := v.ListBacklinks(target.ID)
	if err != nil {
		t.Fatalf("ListBacklinks() error = %v", err)
	}
	if len(backlinks) != 2 {
		t.Fatalf("got %d backlinks, want 2", len(backlinks))
	}

	// Verify both sources are present
	fromIDs := []string{backlinks[0].FromID, backlinks[1].FromID}
	sort.Strings(fromIDs)
	expectedIDs := []string{noteA.ID, noteB.ID}
	sort.Strings(expectedIDs)
	for i, id := range fromIDs {
		if id != expectedIDs[i] {
			t.Errorf("fromIDs[%d] = %q, want %q", i, id, expectedIDs[i])
		}
	}
}

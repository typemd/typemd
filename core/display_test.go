package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildDisplayProperties(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatal(err)
	}
	writeCommonTestTypeSchemas(v)
	if err := v.Open(); err != nil {
		t.Fatal(err)
	}
	defer v.Close()

	// Create two objects and link them
	book, err := v.NewObject("book", "test-book")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("basic properties from schema", func(t *testing.T) {
		props, err := v.BuildDisplayProperties(book)
		if err != nil {
			t.Fatal(err)
		}
		// Should have all schema-defined properties (title, status, rating)
		if len(props) < 3 {
			t.Errorf("expected at least 3 properties, got %d", len(props))
		}
		// All should be forward (not reverse relations)
		for _, p := range props {
			if p.IsReverse {
				t.Errorf("unexpected reverse relation: %s", p.Key)
			}
		}
	})

	t.Run("nil object returns empty", func(t *testing.T) {
		props, err := v.BuildDisplayProperties(nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(props) != 0 {
			t.Errorf("expected 0 properties for nil object, got %d", len(props))
		}
	})
}

func TestBuildDisplayPropertiesWithRelations(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatal(err)
	}

	// Write custom schemas
	schemaYAML := `name: article
properties:
  - name: title
    type: string
  - name: author
    type: relation
    target: person
`
	if err := os.WriteFile(filepath.Join(v.TypesDir(), "article.yaml"), []byte(schemaYAML), 0644); err != nil {
		t.Fatal(err)
	}
	writeCommonTestTypeSchemas(v)

	if err := v.Open(); err != nil {
		t.Fatal(err)
	}
	defer v.Close()

	article, err := v.NewObject("article", "test-article")
	if err != nil {
		t.Fatal(err)
	}
	person, err := v.NewObject("person", "test-person")
	if err != nil {
		t.Fatal(err)
	}
	if err := v.LinkObjects(article.ID, "author", person.ID); err != nil {
		t.Fatal(err)
	}

	t.Run("forward relation marked correctly", func(t *testing.T) {
		article, _ = v.GetObject(article.ID) // re-read to get updated properties
		props, err := v.BuildDisplayProperties(article)
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, p := range props {
			if p.Key == "author" {
				found = true
				if !p.IsRelation {
					t.Error("author should be marked as relation")
				}
				if p.IsReverse {
					t.Error("author should not be marked as reverse")
				}
			}
		}
		if !found {
			t.Error("author property not found")
		}
	})

	t.Run("reverse relation on target", func(t *testing.T) {
		props, err := v.BuildDisplayProperties(person)
		if err != nil {
			t.Fatal(err)
		}
		foundReverse := false
		for _, p := range props {
			if p.IsReverse && p.Key == "author" {
				foundReverse = true
			}
		}
		if !foundReverse {
			t.Error("expected reverse relation 'author' on person")
		}
	})
}

func TestBuildDisplayPropertiesWithBacklinks(t *testing.T) {
	v := setupTestVault(t)

	os.WriteFile(filepath.Join(v.TypesDir(), "note.yaml"),
		[]byte("name: note\nproperties:\n  - name: title\n    type: string\n"), 0644)

	noteA, _ := v.NewObject("note", "alpha")
	noteB, _ := v.NewObject("note", "beta")

	// noteA links to noteB via wiki-link
	bodyA := fmt.Sprintf("---\ntitle: Alpha\n---\n\nSee [[%s]].\n", noteB.ID)
	os.WriteFile(v.ObjectPath(noteA.Type, noteA.Filename), []byte(bodyA), 0644)
	v.SyncIndex()

	// noteB should have a backlink from noteA
	noteB, _ = v.GetObject(noteB.ID)
	props, err := v.BuildDisplayProperties(noteB)
	if err != nil {
		t.Fatal(err)
	}

	foundBacklink := false
	for _, p := range props {
		if p.IsBacklink {
			foundBacklink = true
			if p.Key != BacklinksDisplayKey {
				t.Errorf("backlink Key = %q, want %q", p.Key, BacklinksDisplayKey)
			}
			if p.FromID != noteA.ID {
				t.Errorf("backlink FromID = %q, want %q", p.FromID, noteA.ID)
			}
		}
	}
	if !foundBacklink {
		t.Error("expected backlink property on noteB")
	}
}

func TestBuildDisplayPropertiesNoBacklinks(t *testing.T) {
	v := setupTestVault(t)

	os.WriteFile(filepath.Join(v.TypesDir(), "note.yaml"),
		[]byte("name: note\nproperties:\n  - name: title\n    type: string\n"), 0644)

	note, _ := v.NewObject("note", "lonely")

	props, err := v.BuildDisplayProperties(note)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range props {
		if p.IsBacklink {
			t.Error("expected no backlinks for object with no incoming wiki-links")
		}
	}
}

func TestBacklinkFormat(t *testing.T) {
	p := DisplayProperty{
		Key:        BacklinksDisplayKey,
		Value:      "note/alpha-01abc",
		IsBacklink: true,
		FromID:     "note/alpha-01abc",
	}
	got := p.Format()
	expected := "backlinks: ⟵ note/alpha-01abc"
	if got != expected {
		t.Errorf("Format() = %q, want %q", got, expected)
	}
}

func TestFormat_RelationDisplayID(t *testing.T) {
	p := DisplayProperty{
		Key:        "author",
		Value:      "person/robert-martin-01kk39c30y47xb1dvbs8ywqv50",
		IsRelation: true,
	}
	got := p.Format()
	expected := "author: → person/robert-martin"
	if got != expected {
		t.Errorf("Format() = %q, want %q", got, expected)
	}
}

func TestFormat_ReverseRelationDisplayID(t *testing.T) {
	p := DisplayProperty{
		Key:       "books",
		Value:     "book/clean-code-01kk39c30y47xb1dvbs8ywqv50",
		IsReverse: true,
		FromID:    "book/clean-code-01kk39c30y47xb1dvbs8ywqv50",
	}
	got := p.Format()
	expected := "books: ← book/clean-code"
	if got != expected {
		t.Errorf("Format() = %q, want %q", got, expected)
	}
}

func TestFormat_BacklinkDisplayID(t *testing.T) {
	p := DisplayProperty{
		Key:        BacklinksDisplayKey,
		Value:      "note/my-note-01kk39c30y47xb1dvbs8ywqv50",
		IsBacklink: true,
		FromID:     "note/my-note-01kk39c30y47xb1dvbs8ywqv50",
	}
	got := p.Format()
	expected := "backlinks: ⟵ note/my-note"
	if got != expected {
		t.Errorf("Format() = %q, want %q", got, expected)
	}
}

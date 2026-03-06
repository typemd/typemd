package core

import (
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

	// Write a custom schema with a relation property
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

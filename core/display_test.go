package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
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
	book, err := v.NewObject("book", "test-book", "")
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
	mustWriteTypeSchema(v, "article", []byte(schemaYAML))
	writeCommonTestTypeSchemas(v)

	if err := v.Open(); err != nil {
		t.Fatal(err)
	}
	defer v.Close()

	article, err := v.NewObject("article", "test-article", "")
	if err != nil {
		t.Fatal(err)
	}
	person, err := v.NewObject("person", "test-person", "")
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

	mustWriteTypeSchema(v, "note",
		[]byte("name: note\nproperties:\n  - name: title\n    type: string\n"))

	noteA, _ := v.NewObject("note", "alpha", "")
	noteB, _ := v.NewObject("note", "beta", "")

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

	mustWriteTypeSchema(v, "note",
		[]byte("name: note\nproperties:\n  - name: title\n    type: string\n"))

	note, _ := v.NewObject("note", "lonely", "")

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

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name     string
		prop     DisplayProperty
		expected string
	}{
		{
			name:     "string value",
			prop:     DisplayProperty{Key: "title", Value: "Hello World", Type: "string"},
			expected: "Hello World",
		},
		{
			name:     "checkbox true",
			prop:     DisplayProperty{Key: "active", Value: true, Type: "checkbox"},
			expected: "☑",
		},
		{
			name:     "checkbox false",
			prop:     DisplayProperty{Key: "active", Value: false, Type: "checkbox"},
			expected: "☐",
		},
		{
			name:     "checkbox nil",
			prop:     DisplayProperty{Key: "active", Value: nil, Type: "checkbox"},
			expected: "☐",
		},
		{
			name:     "date",
			prop:     DisplayProperty{Key: "published", Value: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), Type: "date"},
			expected: "2024-01-15",
		},
		{
			name:     "datetime",
			prop:     DisplayProperty{Key: "created", Value: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), Type: "datetime"},
			expected: "2024-01-15 10:30:00",
		},
		{
			name:     "multi_select",
			prop:     DisplayProperty{Key: "tags", Value: []any{"go", "cli"}, Type: "multi_select"},
			expected: "[go, cli]",
		},
		{
			name:     "relation",
			prop:     DisplayProperty{Key: "author", Value: "person/robert-martin-01kk39c30y47xb1dvbs8ywqv50", IsRelation: true},
			expected: "→ person/robert-martin",
		},
		{
			name:     "backlink",
			prop:     DisplayProperty{Key: BacklinksDisplayKey, IsBacklink: true, FromID: "note/my-note-01kk39c30y47xb1dvbs8ywqv50"},
			expected: "⟵ note/my-note",
		},
		{
			name:     "reverse relation",
			prop:     DisplayProperty{Key: "books", IsReverse: true, FromID: "book/clean-code-01kk39c30y47xb1dvbs8ywqv50"},
			expected: "← book/clean-code",
		},
		{
			name:     "nil value",
			prop:     DisplayProperty{Key: "empty", Value: nil},
			expected: "",
		},
		{
			name:     "integer value",
			prop:     DisplayProperty{Key: "count", Value: 42, Type: "number"},
			expected: "42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.prop.FormatValue()
			if got != tt.expected {
				t.Errorf("FormatValue() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestFormatValueWithCustomDateFormat(t *testing.T) {
	tests := []struct {
		name     string
		prop     DisplayProperty
		expected string
	}{
		{
			name:     "custom date format DD/MM/YYYY",
			prop:     DisplayProperty{Key: "d", Value: time.Date(2026, 3, 28, 0, 0, 0, 0, time.UTC), Type: "date", DateFormat: "DD/MM/YYYY"},
			expected: "28/03/2026",
		},
		{
			name:     "custom datetime format with space",
			prop:     DisplayProperty{Key: "dt", Value: time.Date(2026, 3, 28, 14, 30, 0, 0, time.UTC), Type: "datetime", DatetimeFormat: "DD/MM/YYYY HH:mm:ss"},
			expected: "28/03/2026 14:30:00",
		},
		{
			name:     "empty date format uses default",
			prop:     DisplayProperty{Key: "d", Value: time.Date(2026, 3, 28, 0, 0, 0, 0, time.UTC), Type: "date", DateFormat: ""},
			expected: "2026-03-28",
		},
		{
			name:     "empty datetime format uses default with space separator",
			prop:     DisplayProperty{Key: "dt", Value: time.Date(2026, 3, 28, 14, 30, 0, 0, time.UTC), Type: "datetime", DatetimeFormat: ""},
			expected: "2026-03-28 14:30:00",
		},
		{
			name:     "unrecognized tokens pass through",
			prop:     DisplayProperty{Key: "d", Value: time.Date(2026, 3, 28, 0, 0, 0, 0, time.UTC), Type: "date", DateFormat: "YYYY年MM月DD日"},
			expected: "2026年03月28日",
		},
		{
			name:     "date value as string",
			prop:     DisplayProperty{Key: "d", Value: "2026-03-28", Type: "date", DateFormat: "MM/DD/YYYY"},
			expected: "03/28/2026",
		},
		{
			name:     "datetime value as string",
			prop:     DisplayProperty{Key: "dt", Value: "2026-03-28T14:30:00", Type: "datetime", DatetimeFormat: "YYYY/MM/DD HH:mm"},
			expected: "2026/03/28 14:30",
		},
		{
			name:     "datetime value as RFC3339 string",
			prop:     DisplayProperty{Key: "dt", Value: "2026-03-28T14:30:00+08:00", Type: "datetime", DatetimeFormat: "YYYY/MM/DD HH:mm"},
			expected: "2026/03/28 14:30",
		},
		{
			name:     "nil date value returns empty",
			prop:     DisplayProperty{Key: "d", Value: nil, Type: "date", DateFormat: "MM/DD/YYYY"},
			expected: "",
		},
		{
			name:     "custom datetime format omitting seconds",
			prop:     DisplayProperty{Key: "dt", Value: time.Date(2026, 3, 28, 14, 30, 45, 0, time.UTC), Type: "datetime", DatetimeFormat: "YYYY-MM-DD HH:mm"},
			expected: "2026-03-28 14:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.prop.FormatValue()
			if got != tt.expected {
				t.Errorf("FormatValue() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestFormatDelegatesToFormatValue(t *testing.T) {
	tests := []struct {
		name     string
		prop     DisplayProperty
		expected string
	}{
		{
			name:     "string",
			prop:     DisplayProperty{Key: "title", Value: "Hello", Type: "string"},
			expected: "title: Hello",
		},
		{
			name:     "checkbox true",
			prop:     DisplayProperty{Key: "active", Value: true, Type: "checkbox"},
			expected: "active: ☑",
		},
		{
			name:     "checkbox false",
			prop:     DisplayProperty{Key: "active", Value: false, Type: "checkbox"},
			expected: "active: ☐",
		},
		{
			name:     "date",
			prop:     DisplayProperty{Key: "published", Value: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), Type: "date"},
			expected: "published: 2024-01-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.prop.Format()
			if got != tt.expected {
				t.Errorf("Format() = %q, want %q", got, tt.expected)
			}
		})
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

func TestBuildDisplayProperties_IsLocal(t *testing.T) {
	v := setupTestVault(t)

	t.Run("extra property marked as local", func(t *testing.T) {
		book, err := v.NewObject("book", "local-test", "")
		if err != nil {
			t.Fatal(err)
		}
		// Set an extra property not in the book schema
		if err := v.SetProperty(book.ID, "custom_field", "hello"); err != nil {
			t.Fatal(err)
		}
		book, _ = v.GetObject(book.ID)
		props, err := v.BuildDisplayProperties(book)
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range props {
			if p.Key == "custom_field" {
				if !p.IsLocal {
					t.Error("custom_field should be marked as IsLocal")
				}
				return
			}
		}
		t.Error("custom_field not found in display properties")
	})

	t.Run("schema property not marked as local", func(t *testing.T) {
		book, err := v.NewObject("book", "schema-test", "")
		if err != nil {
			t.Fatal(err)
		}
		if err := v.SetProperty(book.ID, "title", "Go Book"); err != nil {
			t.Fatal(err)
		}
		book, _ = v.GetObject(book.ID)
		props, err := v.BuildDisplayProperties(book)
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range props {
			if p.Key == "title" {
				if p.IsLocal {
					t.Error("title should NOT be marked as IsLocal")
				}
				return
			}
		}
		t.Error("title not found in display properties")
	})

	t.Run("system properties not marked as local", func(t *testing.T) {
		book, err := v.NewObject("book", "sys-test", "")
		if err != nil {
			t.Fatal(err)
		}
		props, err := v.BuildDisplayProperties(book)
		if err != nil {
			t.Fatal(err)
		}
		sysKeys := []string{"name", "created_at", "updated_at"}
		for _, sk := range sysKeys {
			for _, p := range props {
				if p.Key == sk && p.IsLocal {
					t.Errorf("system property %q should NOT be marked as IsLocal", sk)
				}
			}
		}
	})

	t.Run("object with only local properties", func(t *testing.T) {
		// Create a book, add extra props, then remove the schema
		book, err := v.NewObject("book", "only-local", "")
		if err != nil {
			t.Fatal(err)
		}
		if err := v.SetProperty(book.ID, "extra1", "val1"); err != nil {
			t.Fatal(err)
		}
		if err := v.SetProperty(book.ID, "extra2", "val2"); err != nil {
			t.Fatal(err)
		}
		// Remove the book schema so all non-system properties become local
		os.RemoveAll(filepath.Join(v.TypesDir(), "book"))
		v.InvalidateSchemaCache()

		book, _ = v.GetObject(book.ID)
		props, err := v.BuildDisplayProperties(book)
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range props {
			if IsSystemProperty(p.Key) {
				if p.IsLocal {
					t.Errorf("system property %q should NOT be local even without schema", p.Key)
				}
			} else {
				if !p.IsLocal {
					t.Errorf("non-system property %q should be local when schema is missing", p.Key)
				}
			}
		}
	})
}


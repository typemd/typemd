package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// ── MarshalTypeSchema unit tests ────────────────────────────────────────────

func TestMarshalTypeSchema_EmptyProperties(t *testing.T) {
	schema := &TypeSchema{Name: "empty", Properties: nil}
	data, err := MarshalTypeSchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(data), "name: empty") {
		t.Errorf("expected name field, got:\n%s", data)
	}
	if !strings.Contains(string(data), "properties: []") {
		t.Errorf("expected empty properties, got:\n%s", data)
	}
}

func TestMarshalTypeSchema_RelationWithAllFields(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{
				Name:          "author",
				Type:          "relation",
				Target:        "person",
				Multiple:      true,
				Bidirectional: true,
				Inverse:       "books",
				Emoji:         "👤",
				Pin:           1,
			},
		},
	}
	data, err := MarshalTypeSchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(data)
	for _, want := range []string{"target: person", "multiple: true", "bidirectional: true", "inverse: books", "emoji:", "pin: 1"} {
		if !strings.Contains(s, want) {
			t.Errorf("expected YAML to contain %q, got:\n%s", want, s)
		}
	}
}

func TestMarshalTypeSchema_SelectWithOptions(t *testing.T) {
	schema := &TypeSchema{
		Name: "task",
		Properties: []Property{
			{
				Name: "status",
				Type: "select",
				Options: []Option{
					{Value: "todo", Label: "To Do"},
					{Value: "done"},
				},
			},
		},
	}
	data, err := MarshalTypeSchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(data)
	if !strings.Contains(s, "value: todo") {
		t.Errorf("expected option value, got:\n%s", s)
	}
	if !strings.Contains(s, "label: To Do") {
		t.Errorf("expected option label, got:\n%s", s)
	}
}

func TestMarshalTypeSchema_NameTemplateRoundTrip(t *testing.T) {
	original := &TypeSchema{
		Name:         "journal",
		NameTemplate: "{{ date:YYYY-MM-DD }}",
		Properties: []Property{
			{Name: "mood", Type: "string"},
		},
	}
	data, err := MarshalTypeSchema(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	// Simulate the load path: unmarshal + extract NameTemplate
	var raw TypeSchema
	if err := yaml.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	filtered := raw.Properties[:0]
	for _, p := range raw.Properties {
		if p.Name == NameProperty {
			raw.NameTemplate = p.Template
			continue
		}
		filtered = append(filtered, p)
	}
	raw.Properties = filtered

	if raw.NameTemplate != original.NameTemplate {
		t.Errorf("NameTemplate = %q, want %q", raw.NameTemplate, original.NameTemplate)
	}
	if len(raw.Properties) != 1 || raw.Properties[0].Name != "mood" {
		t.Errorf("expected 1 property 'mood', got %v", raw.Properties)
	}
}

// ── DeleteSchema unit tests ─────────────────────────────────────────────────

func TestDeleteSchema_RemovesFile(t *testing.T) {
	dir := t.TempDir()
	repo := NewLocalObjectRepository(dir)
	typesDir := filepath.Join(dir, "types")
	os.MkdirAll(typesDir, 0755)
	os.MkdirAll(filepath.Join(typesDir, "test"), 0755)
	os.WriteFile(filepath.Join(typesDir, "test", "schema.yaml"), []byte("name: test\n"), 0644)

	if err := repo.DeleteSchema("test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(typesDir, "test")); !os.IsNotExist(err) {
		t.Error("expected directory to be deleted")
	}
}

func TestDeleteSchema_NonExistent(t *testing.T) {
	dir := t.TempDir()
	repo := NewLocalObjectRepository(dir)
	os.MkdirAll(filepath.Join(dir, "types"), 0755)

	err := repo.DeleteSchema("nope")
	if err == nil {
		t.Error("expected error for non-existent schema")
	}
}

// ── Vault type CRUD unit tests ──────────────────────────────────────────────

func TestSaveType_EmptyName(t *testing.T) {
	v := setupTestVault(t)

	err := v.SaveType(&TypeSchema{Name: ""})
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestDeleteType_BuiltInTag(t *testing.T) {
	v := setupTestVault(t)

	err := v.DeleteType("tag")
	if err == nil {
		t.Error("expected error for built-in type")
	}
	if !strings.Contains(err.Error(), "cannot delete built-in type") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCountObjectsByType_EmptyIndex(t *testing.T) {
	v := setupTestVault(t)

	count, err := v.CountObjectsByType("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

// ── MarshalTypeSchema: zero-value omission ─────────────────────────────────

func TestMarshalTypeSchema_OmitsZeroValueFields(t *testing.T) {
	schema := &TypeSchema{Name: "note"}
	data, err := MarshalTypeSchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(data)
	for _, field := range []string{"plural:", "unique:", "emoji:"} {
		if strings.Contains(s, field) {
			t.Errorf("expected YAML to omit %q, got:\n%s", field, s)
		}
	}
}

// ── Domain event tests ──────────────────────────────────────────────────────

// captureEvents subscribes to vault events and returns a pointer to the captured slice.
func captureEvents(v *Vault) *[]DomainEvent {
	var captured []DomainEvent
	v.Events.Subscribe(func(e DomainEvent) {
		captured = append(captured, e)
	})
	return &captured
}

func TestSaveType_EmitsTypeSavedEvent(t *testing.T) {
	v := setupTestVault(t)
	captured := captureEvents(v)

	if err := v.SaveType(&TypeSchema{Name: "project"}); err != nil {
		t.Fatalf("SaveType() error = %v", err)
	}

	found := false
	for _, e := range *captured {
		if ts, ok := e.(TypeSaved); ok {
			if ts.Schema.Name != "project" {
				t.Errorf("TypeSaved.Schema.Name = %q, want %q", ts.Schema.Name, "project")
			}
			found = true
		}
	}
	if !found {
		t.Error("expected TypeSaved event to be emitted")
	}
}

func TestSaveType_NoEventOnValidationFailure(t *testing.T) {
	v := setupTestVault(t)
	captured := captureEvents(v)

	err := v.SaveType(&TypeSchema{Name: ""})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if len(*captured) != 0 {
		t.Errorf("expected no events, got %d", len(*captured))
	}
}

func TestDeleteType_EmitsTypeDeletedEvent(t *testing.T) {
	v := setupTestVault(t)
	mustWriteTypeSchema(v, "scratch", []byte("name: scratch\nproperties: []\n"))
	captured := captureEvents(v)

	if err := v.DeleteType("scratch"); err != nil {
		t.Fatalf("DeleteType() error = %v", err)
	}

	found := false
	for _, e := range *captured {
		if td, ok := e.(TypeDeleted); ok {
			if td.Name != "scratch" {
				t.Errorf("TypeDeleted.Name = %q, want %q", td.Name, "scratch")
			}
			found = true
		}
	}
	if !found {
		t.Error("expected TypeDeleted event to be emitted")
	}
}

func TestDeleteType_NoEventForBuiltInType(t *testing.T) {
	v := setupTestVault(t)
	captured := captureEvents(v)

	err := v.DeleteType("tag")
	if err == nil {
		t.Fatal("expected error for built-in type")
	}
	if len(*captured) != 0 {
		t.Errorf("expected no events, got %d", len(*captured))
	}
}

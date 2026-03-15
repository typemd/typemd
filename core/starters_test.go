package core

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestStarterTypes_Count(t *testing.T) {
	starters := StarterTypes()
	if len(starters) != 3 {
		t.Fatalf("expected 3 starter types, got %d", len(starters))
	}
}

func TestStarterTypes_Names(t *testing.T) {
	starters := StarterTypes()
	names := make(map[string]bool)
	for _, st := range starters {
		names[st.Name] = true
	}
	for _, want := range []string{"idea", "note", "book"} {
		if !names[want] {
			t.Errorf("expected starter type %q, not found", want)
		}
	}
}

func TestStarterTypes_EmbedIntegrity(t *testing.T) {
	starters := StarterTypes()
	for _, st := range starters {
		if len(st.YAML) == 0 {
			t.Errorf("starter type %q has empty YAML", st.Name)
		}
		var schema TypeSchema
		if err := yaml.Unmarshal(st.YAML, &schema); err != nil {
			t.Errorf("starter type %q YAML parse error: %v", st.Name, err)
		}
		if schema.Name != st.Name {
			t.Errorf("starter type %q: YAML name = %q, metadata name = %q", st.Name, schema.Name, st.Name)
		}
	}
}

func TestStarterTypes_SchemaValidation(t *testing.T) {
	starters := StarterTypes()
	for _, st := range starters {
		var schema TypeSchema
		if err := yaml.Unmarshal(st.YAML, &schema); err != nil {
			t.Fatalf("starter type %q YAML parse error: %v", st.Name, err)
		}
		errs := ValidateSchema(&schema)
		if len(errs) > 0 {
			t.Errorf("starter type %q validation errors: %v", st.Name, errs)
		}
	}
}

func TestStarterTypes_EmojiPresent(t *testing.T) {
	starters := StarterTypes()
	for _, st := range starters {
		if st.Emoji == "" {
			t.Errorf("starter type %q has empty emoji", st.Name)
		}
	}
}

func TestWriteStarterTypes_EmptyNames(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := v.WriteStarterTypes(nil); err != nil {
		t.Fatalf("WriteStarterTypes(nil): %v", err)
	}
	entries, _ := os.ReadDir(v.TypesDir())
	if len(entries) != 0 {
		t.Errorf("expected 0 files, got %d", len(entries))
	}
}

func TestWriteStarterTypes_UnknownName(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := v.WriteStarterTypes([]string{"nonexistent"}); err != nil {
		t.Fatalf("WriteStarterTypes: %v", err)
	}
	entries, _ := os.ReadDir(v.TypesDir())
	if len(entries) != 0 {
		t.Errorf("expected 0 files for unknown name, got %d", len(entries))
	}
}

func TestWriteStarterTypes_DuplicateNames(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := v.WriteStarterTypes([]string{"idea", "idea"}); err != nil {
		t.Fatalf("WriteStarterTypes: %v", err)
	}
	path := filepath.Join(v.TypesDir(), "idea.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected idea.yaml to exist")
	}
}

func TestWriteStarterTypes_FilesAreLoadable(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := v.WriteStarterTypes([]string{"book"}); err != nil {
		t.Fatalf("WriteStarterTypes: %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("open: %v", err)
	}
	defer v.Close()
	schema, err := v.LoadType("book")
	if err != nil {
		t.Fatalf("LoadType: %v", err)
	}
	if schema.Name != "book" {
		t.Errorf("expected name 'book', got %q", schema.Name)
	}
	if schema.Emoji != "📚" {
		t.Errorf("expected emoji '📚', got %q", schema.Emoji)
	}
}

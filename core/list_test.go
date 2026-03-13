package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVault_ListTypes_DefaultsOnly(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	types := v.ListTypes()
	if len(types) != len(defaultTypes) {
		t.Errorf("expected %d types, got %d", len(defaultTypes), len(types))
	}

	// Should be sorted
	for i := 1; i < len(types); i++ {
		if types[i] < types[i-1] {
			t.Errorf("types not sorted: %v", types)
			break
		}
	}
}

func TestVault_ListTypes_WithCustomType(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Create a custom type
	customType := `name: project
properties:
  - name: title
    type: string
`
	if err := os.WriteFile(filepath.Join(v.TypesDir(), "project.yaml"), []byte(customType), 0644); err != nil {
		t.Fatalf("write custom type: %v", err)
	}

	types := v.ListTypes()
	found := false
	for _, name := range types {
		if name == "project" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("custom type 'project' not found in %v", types)
	}

	// Should include both defaults and custom
	if len(types) != len(defaultTypes)+1 {
		t.Errorf("expected %d types, got %d", len(defaultTypes)+1, len(types))
	}
}

func TestVault_ListTypes_CustomOverridesDefault(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Create a custom tag type that overrides the built-in tag
	customTag := `name: tag
properties:
  - name: color
    type: string
  - name: icon
    type: string
  - name: priority
    type: number
`
	if err := os.WriteFile(filepath.Join(v.TypesDir(), "tag.yaml"), []byte(customTag), 0644); err != nil {
		t.Fatalf("write custom type: %v", err)
	}

	types := v.ListTypes()
	// Should not have duplicates — custom tag overrides built-in tag
	if len(types) != len(defaultTypes) {
		t.Errorf("expected %d types (no duplicates), got %d: %v", len(defaultTypes), len(types), types)
	}
}

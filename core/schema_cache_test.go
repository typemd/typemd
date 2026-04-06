package core

import (
	"testing"
)

func TestSchemaCache_LoadTypeReturnsConsistentResults(t *testing.T) {
	v := setupTestVault(t)

	typeName := "project"
	mustWriteTypeSchema(v, typeName, []byte("name: project\nemoji: \"📁\"\nproperties: []\n"))

	schema, err := v.LoadType(typeName)
	if err != nil {
		t.Fatalf("LoadType(%q) error = %v", typeName, err)
	}
	if schema.Emoji != "📁" {
		t.Errorf("LoadType(%q).Emoji = %q, want %q", typeName, schema.Emoji, "📁")
	}
}

func TestSchemaCache_SaveTypeUpdatesCache(t *testing.T) {
	v := setupTestVault(t)

	typeName := "project"
	mustWriteTypeSchema(v, typeName, []byte("name: project\nemoji: \"📁\"\nproperties: []\n"))

	// Load to populate cache
	schema, err := v.LoadType(typeName)
	if err != nil {
		t.Fatalf("LoadType(%q) error = %v", typeName, err)
	}
	if schema.Emoji != "📁" {
		t.Fatalf("initial Emoji = %q, want %q", schema.Emoji, "📁")
	}

	// SaveType with a different emoji
	if err := v.SaveType(&TypeSchema{Name: typeName, Emoji: "🚀"}); err != nil {
		t.Fatalf("SaveType() error = %v", err)
	}

	// Load again — should reflect the updated emoji
	schema2, err := v.LoadType(typeName)
	if err != nil {
		t.Fatalf("LoadType(%q) after SaveType error = %v", typeName, err)
	}
	if schema2.Emoji != "🚀" {
		t.Errorf("LoadType(%q).Emoji after SaveType = %q, want %q", typeName, schema2.Emoji, "🚀")
	}
}

func TestSchemaCache_DeleteTypeMakesUnavailable(t *testing.T) {
	v := setupTestVault(t)

	typeName := "project"
	mustWriteTypeSchema(v, typeName, []byte("name: project\nemoji: \"📁\"\nproperties: []\n"))

	// Load to populate cache
	if _, err := v.LoadType(typeName); err != nil {
		t.Fatalf("LoadType(%q) error = %v", typeName, err)
	}

	// Delete the type
	if err := v.DeleteType(typeName); err != nil {
		t.Fatalf("DeleteType(%q) error = %v", typeName, err)
	}

	// Load again — should return an error
	_, err := v.LoadType(typeName)
	if err == nil {
		t.Errorf("LoadType(%q) after DeleteType: expected error, got nil", typeName)
	}
}

func TestSchemaCache_InvalidateForcesReloadFromDisk(t *testing.T) {
	v := setupTestVault(t)

	typeName := "project"
	mustWriteTypeSchema(v, typeName, []byte("name: project\nemoji: \"📁\"\nproperties: []\n"))

	// Load to populate cache
	schema, err := v.LoadType(typeName)
	if err != nil {
		t.Fatalf("LoadType(%q) error = %v", typeName, err)
	}
	if schema.Emoji != "📁" {
		t.Fatalf("initial Emoji = %q, want %q", schema.Emoji, "📁")
	}

	// Overwrite the file on disk directly (bypass cache)
	mustWriteTypeSchema(v, typeName, []byte("name: project\nemoji: \"🚀\"\nproperties: []\n"))

	// Invalidate the cache
	v.InvalidateSchemaCache()

	// Load again — should reflect the on-disk change
	schema2, err := v.LoadType(typeName)
	if err != nil {
		t.Fatalf("LoadType(%q) after invalidate error = %v", typeName, err)
	}
	if schema2.Emoji != "🚀" {
		t.Errorf("LoadType(%q).Emoji after invalidate = %q, want %q", typeName, schema2.Emoji, "🚀")
	}
}

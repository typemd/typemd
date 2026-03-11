package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// ── Registry tests (1.4) ───────────────────────────────────────────────────

func TestIsSystemProperty_EmptyString(t *testing.T) {
	if IsSystemProperty("") {
		t.Error("empty string should not be a system property")
	}
}

func TestIsSystemProperty_CaseSensitive(t *testing.T) {
	if IsSystemProperty("Name") {
		t.Error("system property check should be case-sensitive")
	}
	if IsSystemProperty("CREATED_AT") {
		t.Error("system property check should be case-sensitive")
	}
}

func TestSystemPropertyNames_Order(t *testing.T) {
	names := SystemPropertyNames()
	expected := []string{"name", "created_at", "updated_at"}
	if len(names) != len(expected) {
		t.Fatalf("SystemPropertyNames() returned %d names, want %d", len(names), len(expected))
	}
	for i, name := range expected {
		if names[i] != name {
			t.Errorf("SystemPropertyNames()[%d] = %q, want %q", i, names[i], name)
		}
	}
}

func TestSystemPropertyNames_ReturnsNewSlice(t *testing.T) {
	a := SystemPropertyNames()
	b := SystemPropertyNames()
	a[0] = "modified"
	if b[0] == "modified" {
		t.Error("SystemPropertyNames should return a new slice each time")
	}
}

// ── Validation tests (2.5) ─────────────────────────────────────────────────

func TestValidateSchema_RejectsAllSystemProperties(t *testing.T) {
	for _, name := range SystemPropertyNames() {
		schema := &TypeSchema{
			Name: "test",
			Properties: []Property{
				{Name: name, Type: "string"},
			},
		}
		errs := ValidateSchema(schema)
		if len(errs) == 0 {
			t.Errorf("expected error for system property %q in schema", name)
		}
		found := false
		for _, err := range errs {
			if strings.Contains(err.Error(), "reserved system property") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("error for %q should mention 'reserved system property', got %v", name, errs)
		}
	}
}

func TestValidateSharedProperties_RejectsAllSystemProperties(t *testing.T) {
	for _, name := range SystemPropertyNames() {
		props := []Property{
			{Name: name, Type: "string"},
		}
		errs := ValidateSharedProperties(props)
		if len(errs) == 0 {
			t.Errorf("expected error for system property %q in shared properties", name)
		}
	}
}

// ── Timestamp creation tests (3.4) ─────────────────────────────────────────

func TestNewObject_TimestampFormat(t *testing.T) {
	v := setupTestVault(t)

	obj, err := v.NewObject("book", "test")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	for _, prop := range []string{"created_at", "updated_at"} {
		val, ok := obj.Properties[prop].(string)
		if !ok {
			t.Errorf("%s should be a string, got %T", prop, obj.Properties[prop])
			continue
		}
		parsed, err := time.Parse(time.RFC3339, val)
		if err != nil {
			t.Errorf("%s = %q is not valid RFC 3339: %v", prop, val, err)
			continue
		}
		// Should have timezone offset (not UTC "Z" unless running in UTC)
		if !strings.Contains(val, "+") && !strings.Contains(val, "Z") {
			t.Errorf("%s = %q should contain timezone offset", prop, val)
		}
		// Should be recent (within last 5 seconds)
		if time.Since(parsed) > 5*time.Second {
			t.Errorf("%s = %q is not recent", prop, val)
		}
	}
}

func TestNewObject_CreatedAtAndUpdatedAtMatch(t *testing.T) {
	v := setupTestVault(t)

	obj, err := v.NewObject("book", "test")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	created := obj.Properties["created_at"].(string)
	updated := obj.Properties["updated_at"].(string)
	if created != updated {
		t.Errorf("on creation, created_at (%s) should equal updated_at (%s)", created, updated)
	}
}

// ── Save timestamp tests (4.4) ─────────────────────────────────────────────

func TestSaveObject_UpdatesUpdatedAt(t *testing.T) {
	v := setupTestVault(t)

	obj, err := v.NewObject("book", "test")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	originalCreated := obj.Properties["created_at"].(string)
	originalUpdated := obj.Properties["updated_at"].(string)

	// Modify and save
	obj.Properties["title"] = "New Title"
	if err := v.SaveObject(obj); err != nil {
		t.Fatalf("SaveObject() error = %v", err)
	}

	// Re-read from disk
	got, err := v.GetObject(obj.ID)
	if err != nil {
		t.Fatalf("GetObject() error = %v", err)
	}

	// created_at should not change
	if got.Properties["created_at"] != originalCreated {
		t.Errorf("created_at changed: was %v, now %v", originalCreated, got.Properties["created_at"])
	}

	// updated_at should be >= original
	newUpdated, ok := got.Properties["updated_at"].(string)
	if !ok {
		t.Fatalf("updated_at should be string, got %T", got.Properties["updated_at"])
	}
	newTime, _ := time.Parse(time.RFC3339, newUpdated)
	origTime, _ := time.Parse(time.RFC3339, originalUpdated)
	if newTime.Before(origTime) {
		t.Errorf("updated_at went backwards: was %s, now %s", originalUpdated, newUpdated)
	}
}

func TestSetProperty_UpdatesUpdatedAt(t *testing.T) {
	v := setupTestVault(t)

	obj, err := v.NewObject("book", "test")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	originalCreated := obj.Properties["created_at"].(string)

	if err := v.SetProperty(obj.ID, "title", "Updated"); err != nil {
		t.Fatalf("SetProperty() error = %v", err)
	}

	got, err := v.GetObject(obj.ID)
	if err != nil {
		t.Fatalf("GetObject() error = %v", err)
	}

	if got.Properties["created_at"] != originalCreated {
		t.Errorf("created_at changed after SetProperty")
	}

	if _, ok := got.Properties["updated_at"].(string); !ok {
		t.Error("updated_at should be a string after SetProperty")
	}
}

// ── Ordering tests (5.5) ───────────────────────────────────────────────────

func TestOrderedPropKeys_SystemPropertiesFirst(t *testing.T) {
	props := map[string]any{
		"name":       "test",
		"created_at": "2026-03-11T10:00:00+08:00",
		"updated_at": "2026-03-11T10:00:00+08:00",
		"title":      "Test",
		"rating":     5,
	}
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "title", Type: "string"},
			{Name: "rating", Type: "number"},
		},
	}

	keys := OrderedPropKeys(props, schema)
	expected := []string{"name", "created_at", "updated_at", "title", "rating"}
	if len(keys) != len(expected) {
		t.Fatalf("OrderedPropKeys returned %d keys, want %d: %v", len(keys), len(expected), keys)
	}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("keys[%d] = %q, want %q", i, keys[i], k)
		}
	}
}

func TestOrderedPropKeys_MissingTimestamps(t *testing.T) {
	props := map[string]any{
		"name":  "test",
		"title": "Test",
	}
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "title", Type: "string"},
		},
	}

	keys := OrderedPropKeys(props, schema)
	expected := []string{"name", "title"}
	if len(keys) != len(expected) {
		t.Fatalf("OrderedPropKeys returned %d keys, want %d: %v", len(keys), len(expected), keys)
	}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("keys[%d] = %q, want %q", i, keys[i], k)
		}
	}
}

func TestOrderedPropKeys_NoSchema(t *testing.T) {
	props := map[string]any{
		"name":       "test",
		"created_at": "2026-03-11T10:00:00+08:00",
		"updated_at": "2026-03-11T10:00:00+08:00",
		"zebra":      "z",
		"apple":      "a",
	}

	keys := OrderedPropKeys(props, nil)
	// System properties first in registry order, then extras alphabetically
	expected := []string{"name", "created_at", "updated_at", "apple", "zebra"}
	if len(keys) != len(expected) {
		t.Fatalf("OrderedPropKeys returned %d keys, want %d: %v", len(keys), len(expected), keys)
	}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("keys[%d] = %q, want %q", i, keys[i], k)
		}
	}
}

// ── Graceful absence tests (6.4) ───────────────────────────────────────────

func TestGetObject_WithoutTimestamps(t *testing.T) {
	v := setupTestVault(t)

	// Create a raw object file without timestamps
	typeName := "book"
	ulid, _ := GenerateULID()
	filename := "old-book-" + ulid
	objPath := v.ObjectPath(typeName, filename)
	os.MkdirAll(filepath.Dir(objPath), 0755)
	content := "---\nname: old-book\ntitle: Old Book\n---\n"
	os.WriteFile(objPath, []byte(content), 0644)

	obj, err := v.GetObject(typeName + "/" + filename)
	if err != nil {
		t.Fatalf("GetObject() error = %v", err)
	}

	if _, ok := obj.Properties["created_at"]; ok {
		t.Error("old object should not have created_at")
	}
	if _, ok := obj.Properties["updated_at"]; ok {
		t.Error("old object should not have updated_at")
	}
	if obj.Properties["name"] != "old-book" {
		t.Errorf("name = %v, want %q", obj.Properties["name"], "old-book")
	}
}

func TestSyncIndex_DoesNotAddTimestampsToExistingObjects(t *testing.T) {
	v := setupTestVault(t)

	// Create raw object file without timestamps
	typeName := "book"
	ulid, _ := GenerateULID()
	filename := "legacy-" + ulid
	objPath := v.ObjectPath(typeName, filename)
	os.MkdirAll(filepath.Dir(objPath), 0755)
	content := "---\nname: legacy\ntitle: Legacy Book\n---\n"
	os.WriteFile(objPath, []byte(content), 0644)

	// Sync
	_, err := v.SyncIndex()
	if err != nil {
		t.Fatalf("SyncIndex() error = %v", err)
	}

	// Read file back — should not have timestamps
	data, err := os.ReadFile(objPath)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}
	if strings.Contains(string(data), "created_at") || strings.Contains(string(data), "updated_at") {
		t.Errorf("SyncIndex added timestamps to existing object:\n%s", string(data))
	}
}

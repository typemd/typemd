package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestTypeSchema_UnmarshalYAML(t *testing.T) {
	data := []byte(`
name: book
properties:
  - name: title
    type: string
  - name: status
    type: select
    options:
      - value: to-read
      - value: reading
      - value: done
  - name: rating
    type: number
  - name: author
    type: relation
    target: person
`)

	var schema TypeSchema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if schema.Name != "book" {
		t.Errorf("Name = %q, want %q", schema.Name, "book")
	}
	if len(schema.Properties) != 4 {
		t.Fatalf("len(Properties) = %d, want 4", len(schema.Properties))
	}

	// select property
	status := schema.Properties[1]
	if status.Type != "select" {
		t.Errorf("status.Type = %q, want %q", status.Type, "select")
	}
	if len(status.Options) != 3 {
		t.Errorf("len(status.Options) = %d, want 3", len(status.Options))
	}
	if status.Options[0].Value != "to-read" {
		t.Errorf("status.Options[0].Value = %q, want %q", status.Options[0].Value, "to-read")
	}

	// relation property
	author := schema.Properties[3]
	if author.Target != "person" {
		t.Errorf("author.Target = %q, want %q", author.Target, "person")
	}
}

func TestTypeSchema_OptionWithLabel(t *testing.T) {
	data := []byte(`
name: task
properties:
  - name: priority
    type: select
    options:
      - value: low
        label: Low Priority
      - value: high
        label: High Priority
`)

	var schema TypeSchema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	opts := schema.Properties[0].Options
	if opts[0].Label != "Low Priority" {
		t.Errorf("opts[0].Label = %q, want %q", opts[0].Label, "Low Priority")
	}
	if opts[1].Label != "High Priority" {
		t.Errorf("opts[1].Label = %q, want %q", opts[1].Label, "High Priority")
	}
}

func TestDefaultTypes(t *testing.T) {
	for _, name := range []string{"book", "person", "note"} {
		schema, ok := defaultTypes[name]
		if !ok {
			t.Errorf("defaultTypes missing %q", name)
			continue
		}
		if schema.Name != name {
			t.Errorf("defaultTypes[%q].Name = %q", name, schema.Name)
		}
		if len(schema.Properties) == 0 {
			t.Errorf("defaultTypes[%q] has no properties", name)
		}
	}
}

func TestDefaultTypes_BookUsesSelect(t *testing.T) {
	book := defaultTypes["book"]
	status := book.Properties[1]
	if status.Type != "select" {
		t.Errorf("book.status.Type = %q, want %q", status.Type, "select")
	}
	if len(status.Options) != 3 {
		t.Errorf("len(book.status.Options) = %d, want 3", len(status.Options))
	}
}

func TestVault_LoadType_CustomYAML(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	yamlContent := []byte(`name: book
properties:
  - name: title
    type: string
  - name: isbn
    type: string
`)
	if err := os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), yamlContent, 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	schema, err := v.LoadType("book")
	if err != nil {
		t.Fatalf("LoadType() error = %v", err)
	}

	if len(schema.Properties) != 2 {
		t.Errorf("len(Properties) = %d, want 2", len(schema.Properties))
	}
	if schema.Properties[1].Name != "isbn" {
		t.Errorf("Properties[1].Name = %q, want %q", schema.Properties[1].Name, "isbn")
	}
}

func TestVault_LoadType_BuiltinFallback(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	schema, err := v.LoadType("book")
	if err != nil {
		t.Fatalf("LoadType() error = %v", err)
	}

	if schema.Name != "book" {
		t.Errorf("Name = %q, want %q", schema.Name, "book")
	}
	if len(schema.Properties) != 3 {
		t.Errorf("len(Properties) = %d, want 3", len(schema.Properties))
	}
}

func TestVault_LoadType_NotFound(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	_, err := v.LoadType("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown type, got nil")
	}
}

func TestValidateObject_ValidProps(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "title", Type: "string"},
			{Name: "status", Type: "select", Options: []Option{
				{Value: "to-read"}, {Value: "reading"}, {Value: "done"},
			}},
			{Name: "rating", Type: "number"},
			{Name: "author", Type: "relation", Target: "person"},
		},
	}

	props := map[string]any{
		"title":  "Go in Action",
		"status": "reading",
		"rating": 4.5,
		"author": "person/rob-pike",
	}

	errs := ValidateObject(props, schema)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateObject_TypeErrors(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "title", Type: "string"},
			{Name: "rating", Type: "number"},
		},
	}

	props := map[string]any{
		"title":  123,
		"rating": "abc",
	}

	errs := ValidateObject(props, schema)
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d: %v", len(errs), errs)
	}
}

func TestValidateObject_InvalidSelect(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "status", Type: "select", Options: []Option{
				{Value: "to-read"}, {Value: "reading"}, {Value: "done"},
			}},
		},
	}

	props := map[string]any{
		"status": "invalid-value",
	}

	errs := ValidateObject(props, schema)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateObject_ExtraPropsIgnored(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "title", Type: "string"},
		},
	}

	props := map[string]any{
		"title":        "Go in Action",
		"custom_field": "anything",
	}

	errs := ValidateObject(props, schema)
	if len(errs) != 0 {
		t.Errorf("expected no errors for extra props, got %v", errs)
	}
}

func TestValidateObject_MissingPropsIgnored(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "title", Type: "string"},
			{Name: "rating", Type: "number"},
		},
	}

	props := map[string]any{
		"title": "Go in Action",
	}

	errs := ValidateObject(props, schema)
	if len(errs) != 0 {
		t.Errorf("expected no errors for missing props, got %v", errs)
	}
}

func TestProperty_RelationFields(t *testing.T) {
	data := []byte(`
name: book
properties:
  - name: author
    type: relation
    target: person
    multiple: true
    bidirectional: true
    inverse: books
`)

	var schema TypeSchema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	author := schema.Properties[0]
	if !author.Multiple {
		t.Error("expected Multiple = true")
	}
	if !author.Bidirectional {
		t.Error("expected Bidirectional = true")
	}
	if author.Inverse != "books" {
		t.Errorf("Inverse = %q, want %q", author.Inverse, "books")
	}
}

func TestValidateObject_RelationMultiple(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "author", Type: "relation", Target: "person", Multiple: true},
		},
	}

	props := map[string]any{
		"author": []any{"person/alan", "person/brian"},
	}
	errs := ValidateObject(props, schema)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}

	props2 := map[string]any{
		"author": []any{"person/alan", 123},
	}
	errs2 := ValidateObject(props2, schema)
	if len(errs2) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(errs2), errs2)
	}
}

func TestProperty_Default(t *testing.T) {
	data := []byte(`
name: book
properties:
  - name: status
    type: select
    options:
      - value: to-read
      - value: reading
      - value: done
    default: to-read
  - name: title
    type: string
  - name: rating
    type: number
    default: 0
`)

	var schema TypeSchema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	status := schema.Properties[0]
	if status.Default != "to-read" {
		t.Errorf("status.Default = %v, want %q", status.Default, "to-read")
	}

	title := schema.Properties[1]
	if title.Default != nil {
		t.Errorf("title.Default = %v, want nil", title.Default)
	}

	rating := schema.Properties[2]
	if rating.Default != 0 {
		t.Errorf("rating.Default = %v, want 0", rating.Default)
	}
}

func TestValidateSchema_Valid(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "title", Type: "string"},
			{Name: "status", Type: "select", Options: []Option{
				{Value: "to-read"}, {Value: "reading"}, {Value: "done"},
			}},
			{Name: "author", Type: "relation", Target: "person"},
			{Name: "rating", Type: "number"},
		},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateSchema_MissingName(t *testing.T) {
	schema := &TypeSchema{Properties: []Property{{Name: "title", Type: "string"}}}
	errs := ValidateSchema(schema)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateSchema_InvalidPropertyType(t *testing.T) {
	schema := &TypeSchema{
		Name:       "test",
		Properties: []Property{{Name: "field", Type: "boolean"}},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateSchema_EnumRejectedWithGuidance(t *testing.T) {
	schema := &TypeSchema{
		Name:       "test",
		Properties: []Property{{Name: "status", Type: "enum"}},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	errMsg := errs[0].Error()
	if !strings.Contains(errMsg, "no longer supported") || !strings.Contains(errMsg, "select") {
		t.Errorf("expected guidance message about select, got %q", errMsg)
	}
}

func TestValidateSchema_SelectWithoutOptions(t *testing.T) {
	schema := &TypeSchema{
		Name:       "test",
		Properties: []Property{{Name: "status", Type: "select"}},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateSchema_MultiSelectWithoutOptions(t *testing.T) {
	schema := &TypeSchema{
		Name:       "test",
		Properties: []Property{{Name: "tags", Type: "multi_select"}},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateSchema_RelationWithoutTarget(t *testing.T) {
	schema := &TypeSchema{
		Name:       "test",
		Properties: []Property{{Name: "author", Type: "relation"}},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateSchema_DuplicatePropertyName(t *testing.T) {
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Name: "title", Type: "string"},
			{Name: "title", Type: "number"},
		},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateSchema_MissingPropertyName(t *testing.T) {
	schema := &TypeSchema{
		Name:       "test",
		Properties: []Property{{Type: "string"}},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateSchema_MissingPropertyType(t *testing.T) {
	schema := &TypeSchema{
		Name:       "test",
		Properties: []Property{{Name: "title"}},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateSchema_AllNewTypes(t *testing.T) {
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Name: "title", Type: "string"},
			{Name: "count", Type: "number"},
			{Name: "published", Type: "date"},
			{Name: "created_at", Type: "datetime"},
			{Name: "homepage", Type: "url"},
			{Name: "active", Type: "checkbox"},
			{Name: "status", Type: "select", Options: []Option{{Value: "a"}}},
			{Name: "tags", Type: "multi_select", Options: []Option{{Value: "b"}}},
			{Name: "author", Type: "relation", Target: "person"},
		},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 0 {
		t.Errorf("expected no errors for all 9 types, got %v", errs)
	}
}

func TestTypeSchema_EmojiField(t *testing.T) {
	data := []byte(`
name: book
emoji: 📚
properties:
  - name: title
    type: string
`)

	var schema TypeSchema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if schema.Emoji != "📚" {
		t.Errorf("Emoji = %q, want %q", schema.Emoji, "📚")
	}
}

func TestTypeSchema_EmojiFieldOmitted(t *testing.T) {
	data := []byte(`
name: book
properties:
  - name: title
    type: string
`)

	var schema TypeSchema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if schema.Emoji != "" {
		t.Errorf("Emoji = %q, want empty string", schema.Emoji)
	}
}

func TestDefaultTypes_HaveEmoji(t *testing.T) {
	expected := map[string]string{
		"book":   "📚",
		"person": "👤",
		"note":   "📝",
	}

	for name, wantEmoji := range expected {
		schema, ok := defaultTypes[name]
		if !ok {
			t.Errorf("defaultTypes missing %q", name)
			continue
		}
		if schema.Emoji != wantEmoji {
			t.Errorf("defaultTypes[%q].Emoji = %q, want %q", name, schema.Emoji, wantEmoji)
		}
	}
}

func TestVault_LoadType_CustomEmojiOverride(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	yamlContent := []byte(`name: book
emoji: 📖
properties:
  - name: title
    type: string
`)
	if err := os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), yamlContent, 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	schema, err := v.LoadType("book")
	if err != nil {
		t.Fatalf("LoadType() error = %v", err)
	}

	if schema.Emoji != "📖" {
		t.Errorf("Emoji = %q, want %q", schema.Emoji, "📖")
	}
}

// ── New property type validation tests ──────────────────────────────────────

func TestValidateObject_Date(t *testing.T) {
	schema := &TypeSchema{
		Name:       "event",
		Properties: []Property{{Name: "date", Type: "date"}},
	}

	// Valid string date
	errs := ValidateObject(map[string]any{"date": "2026-01-15"}, schema)
	if len(errs) != 0 {
		t.Errorf("valid date: expected no errors, got %v", errs)
	}

	// Valid time.Time (YAML auto-parsed)
	errs = ValidateObject(map[string]any{"date": time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)}, schema)
	if len(errs) != 0 {
		t.Errorf("time.Time date: expected no errors, got %v", errs)
	}

	// Invalid format
	errs = ValidateObject(map[string]any{"date": "01/15/2026"}, schema)
	if len(errs) != 1 {
		t.Errorf("invalid date format: expected 1 error, got %d: %v", len(errs), errs)
	}

	// Invalid type
	errs = ValidateObject(map[string]any{"date": 123}, schema)
	if len(errs) != 1 {
		t.Errorf("invalid date type: expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateObject_Datetime(t *testing.T) {
	schema := &TypeSchema{
		Name:       "event",
		Properties: []Property{{Name: "at", Type: "datetime"}},
	}

	// RFC3339
	errs := ValidateObject(map[string]any{"at": "2026-01-15T10:30:00Z"}, schema)
	if len(errs) != 0 {
		t.Errorf("RFC3339: expected no errors, got %v", errs)
	}

	// RFC3339 with offset
	errs = ValidateObject(map[string]any{"at": "2026-01-15T10:30:00+08:00"}, schema)
	if len(errs) != 0 {
		t.Errorf("RFC3339 offset: expected no errors, got %v", errs)
	}

	// Without timezone
	errs = ValidateObject(map[string]any{"at": "2026-01-15T10:30:00"}, schema)
	if len(errs) != 0 {
		t.Errorf("no tz: expected no errors, got %v", errs)
	}

	// time.Time
	errs = ValidateObject(map[string]any{"at": time.Now()}, schema)
	if len(errs) != 0 {
		t.Errorf("time.Time: expected no errors, got %v", errs)
	}

	// Invalid
	errs = ValidateObject(map[string]any{"at": "not-a-datetime"}, schema)
	if len(errs) != 1 {
		t.Errorf("invalid datetime: expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateObject_URL(t *testing.T) {
	schema := &TypeSchema{
		Name:       "bookmark",
		Properties: []Property{{Name: "link", Type: "url"}},
	}

	errs := ValidateObject(map[string]any{"link": "https://example.com"}, schema)
	if len(errs) != 0 {
		t.Errorf("valid https: expected no errors, got %v", errs)
	}

	errs = ValidateObject(map[string]any{"link": "http://example.com/path"}, schema)
	if len(errs) != 0 {
		t.Errorf("valid http: expected no errors, got %v", errs)
	}

	errs = ValidateObject(map[string]any{"link": "ftp://example.com"}, schema)
	if len(errs) != 1 {
		t.Errorf("ftp: expected 1 error, got %d: %v", len(errs), errs)
	}

	errs = ValidateObject(map[string]any{"link": "not-a-url"}, schema)
	if len(errs) != 1 {
		t.Errorf("no scheme: expected 1 error, got %d: %v", len(errs), errs)
	}

	errs = ValidateObject(map[string]any{"link": 123}, schema)
	if len(errs) != 1 {
		t.Errorf("non-string: expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateObject_Checkbox(t *testing.T) {
	schema := &TypeSchema{
		Name:       "task",
		Properties: []Property{{Name: "done", Type: "checkbox"}},
	}

	errs := ValidateObject(map[string]any{"done": true}, schema)
	if len(errs) != 0 {
		t.Errorf("true: expected no errors, got %v", errs)
	}

	errs = ValidateObject(map[string]any{"done": false}, schema)
	if len(errs) != 0 {
		t.Errorf("false: expected no errors, got %v", errs)
	}

	// String "true" should be rejected
	errs = ValidateObject(map[string]any{"done": "true"}, schema)
	if len(errs) != 1 {
		t.Errorf("string true: expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateObject_MultiSelect(t *testing.T) {
	schema := &TypeSchema{
		Name: "item",
		Properties: []Property{
			{Name: "tags", Type: "multi_select", Options: []Option{
				{Value: "go"}, {Value: "rust"}, {Value: "python"},
			}},
		},
	}

	// Valid array
	errs := ValidateObject(map[string]any{"tags": []any{"go", "rust"}}, schema)
	if len(errs) != 0 {
		t.Errorf("valid array: expected no errors, got %v", errs)
	}

	// Single string coerced to list
	errs = ValidateObject(map[string]any{"tags": "go"}, schema)
	if len(errs) != 0 {
		t.Errorf("single string: expected no errors, got %v", errs)
	}

	// Invalid option
	errs = ValidateObject(map[string]any{"tags": []any{"go", "java"}}, schema)
	if len(errs) != 1 {
		t.Errorf("invalid option: expected 1 error, got %d: %v", len(errs), errs)
	}

	// Non-string item in array
	errs = ValidateObject(map[string]any{"tags": []any{"go", 123}}, schema)
	if len(errs) != 1 {
		t.Errorf("non-string item: expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidateObject_DateInvalidDate(t *testing.T) {
	schema := &TypeSchema{
		Name:       "event",
		Properties: []Property{{Name: "date", Type: "date"}},
	}

	// Valid format but invalid date (Feb 30)
	errs := ValidateObject(map[string]any{"date": "2026-02-30"}, schema)
	if len(errs) != 1 {
		t.Errorf("Feb 30: expected 1 error, got %d: %v", len(errs), errs)
	}
}

// ── Property emoji tests ────────────────────────────────────────────────────

func TestProperty_EmojiField(t *testing.T) {
	data := []byte(`
name: book
properties:
  - name: title
    type: string
    emoji: 📖
  - name: rating
    type: number
    emoji: ⭐
`)

	var schema TypeSchema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if schema.Properties[0].Emoji != "📖" {
		t.Errorf("title.Emoji = %q, want %q", schema.Properties[0].Emoji, "📖")
	}
	if schema.Properties[1].Emoji != "⭐" {
		t.Errorf("rating.Emoji = %q, want %q", schema.Properties[1].Emoji, "⭐")
	}
}

func TestProperty_EmojiFieldOmitted(t *testing.T) {
	data := []byte(`
name: book
properties:
  - name: title
    type: string
`)

	var schema TypeSchema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if schema.Properties[0].Emoji != "" {
		t.Errorf("title.Emoji = %q, want empty string", schema.Properties[0].Emoji)
	}
}

func TestValidateSchema_UniquePropertyEmojis(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "title", Type: "string", Emoji: "📖"},
			{Name: "rating", Type: "number", Emoji: "⭐"},
		},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 0 {
		t.Errorf("expected no errors for unique emojis, got %v", errs)
	}
}

func TestValidateSchema_DuplicatePropertyEmoji(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "title", Type: "string", Emoji: "👤"},
			{Name: "author", Type: "string", Emoji: "👤"},
		},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error for duplicate emoji, got %d: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Error(), "duplicate property emoji") {
		t.Errorf("expected duplicate emoji error, got %q", errs[0].Error())
	}
}

func TestValidateSchema_EmptyEmojisDoNotConflict(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "title", Type: "string"},
			{Name: "author", Type: "string"},
			{Name: "rating", Type: "number", Emoji: "⭐"},
		},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 0 {
		t.Errorf("expected no errors for empty emojis, got %v", errs)
	}
}

func TestValidateSchema_PositivePinAccepted(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "status", Type: "string", Pin: 1},
			{Name: "rating", Type: "number", Pin: 2},
		},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid pin values, got %v", errs)
	}
}

func TestValidateSchema_NegativePinRejected(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "status", Type: "string", Pin: -1},
		},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error for negative pin, got %d: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Error(), "pin value must be a positive integer") {
		t.Errorf("expected pin validation error, got %q", errs[0].Error())
	}
}

func TestValidateSchema_DuplicatePinRejected(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "status", Type: "string", Pin: 1},
			{Name: "rating", Type: "number", Pin: 1},
		},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error for duplicate pin, got %d: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Error(), "duplicate pin value") {
		t.Errorf("expected duplicate pin error, got %q", errs[0].Error())
	}
}

func TestValidateSchema_UnpinnedDoNotConflict(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "title", Type: "string"},
			{Name: "author", Type: "string"},
			{Name: "status", Type: "string", Pin: 1},
		},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 0 {
		t.Errorf("expected no errors for unpinned properties, got %v", errs)
	}
}

func TestTypeSchema_PropertyNames(t *testing.T) {
	t.Run("empty schema", func(t *testing.T) {
		schema := &TypeSchema{Name: "empty"}
		names := schema.PropertyNames()
		if len(names) != 0 {
			t.Errorf("expected empty map, got %v", names)
		}
	})

	t.Run("multi-property schema", func(t *testing.T) {
		schema := &TypeSchema{
			Name: "book",
			Properties: []Property{
				{Name: "title", Type: "string"},
				{Name: "status", Type: "select"},
				{Name: "rating", Type: "number"},
			},
		}
		names := schema.PropertyNames()
		if len(names) != 3 {
			t.Fatalf("expected 3 names, got %d", len(names))
		}
		for _, name := range []string{"title", "status", "rating"} {
			if !names[name] {
				t.Errorf("expected %q in PropertyNames", name)
			}
		}
	})
}

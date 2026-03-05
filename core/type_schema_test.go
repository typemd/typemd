package core

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTypeSchema_UnmarshalYAML(t *testing.T) {
	data := []byte(`
name: book
properties:
  - name: title
    type: string
  - name: status
    type: enum
    values: [to-read, reading, done]
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

	// enum property
	status := schema.Properties[1]
	if status.Type != "enum" {
		t.Errorf("status.Type = %q, want %q", status.Type, "enum")
	}
	if len(status.Values) != 3 {
		t.Errorf("len(status.Values) = %d, want 3", len(status.Values))
	}

	// relation property
	author := schema.Properties[3]
	if author.Target != "person" {
		t.Errorf("author.Target = %q, want %q", author.Target, "person")
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

func TestVault_LoadType_CustomYAML(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// 寫入自訂 schema
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

	// 應該載入自訂版本（有 isbn），不是內建預設
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

	// 沒有自訂 book.yaml，應該回傳內建預設
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
			{Name: "status", Type: "enum", Values: []string{"to-read", "reading", "done"}},
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
		"title":  123,   // 應該是 string
		"rating": "abc", // 應該是 number
	}

	errs := ValidateObject(props, schema)
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d: %v", len(errs), errs)
	}
}

func TestValidateObject_InvalidEnum(t *testing.T) {
	schema := &TypeSchema{
		Name: "book",
		Properties: []Property{
			{Name: "status", Type: "enum", Values: []string{"to-read", "reading", "done"}},
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
		// rating 不存在，寬鬆模式不報錯
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

	// valid: array of strings
	props := map[string]any{
		"author": []any{"person/alan", "person/brian"},
	}
	errs := ValidateObject(props, schema)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}

	// invalid: array with non-string
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
    type: enum
    values: [to-read, reading, done]
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

	// status has default
	status := schema.Properties[0]
	if status.Default != "to-read" {
		t.Errorf("status.Default = %v, want %q", status.Default, "to-read")
	}

	// title has no default
	title := schema.Properties[1]
	if title.Default != nil {
		t.Errorf("title.Default = %v, want nil", title.Default)
	}

	// rating has default 0
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
			{Name: "status", Type: "enum", Values: []string{"to-read", "reading", "done"}},
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

func TestValidateSchema_EnumWithoutValues(t *testing.T) {
	schema := &TypeSchema{
		Name:       "test",
		Properties: []Property{{Name: "status", Type: "enum"}},
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


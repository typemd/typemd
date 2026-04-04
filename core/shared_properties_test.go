package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupSharedPropsTestVault(t *testing.T) *Vault {
	t.Helper()
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("vault init: %v", err)
	}
	return v
}

func writeSharedProperty(t *testing.T, v *Vault, name, content string) {
	t.Helper()
	dir := v.SharedPropertiesDir()
	os.MkdirAll(dir, 0755)
	if err := os.WriteFile(filepath.Join(dir, name+".yaml"), []byte(content), 0644); err != nil {
		t.Fatalf("write shared property %q: %v", name, err)
	}
}

func TestLoadSharedProperties_FileExists(t *testing.T) {
	v := setupSharedPropsTestVault(t)
	writeSharedProperty(t, v, "due_date", "type: date\nemoji: 📅\n")
	writeSharedProperty(t, v, "priority", "type: select\noptions:\n  - value: high\n  - value: low\n")

	props, err := v.LoadSharedProperties()
	if err != nil {
		t.Fatalf("LoadSharedProperties error = %v", err)
	}
	if len(props) != 2 {
		t.Fatalf("len(props) = %d, want 2", len(props))
	}
	// Properties are sorted by filename
	byName := make(map[string]Property)
	for _, p := range props {
		byName[p.Name] = p
	}
	if p, ok := byName["due_date"]; !ok || p.Type != "date" {
		t.Errorf("due_date = %+v", byName["due_date"])
	}
	if p, ok := byName["priority"]; !ok || p.Type != "select" {
		t.Errorf("priority = %+v", byName["priority"])
	}
}

func TestLoadSharedProperties_FileNotExist(t *testing.T) {
	v := setupSharedPropsTestVault(t)
	// Remove the properties directory that Init creates
	os.RemoveAll(v.SharedPropertiesDir())

	props, err := v.LoadSharedProperties()
	if err != nil {
		t.Fatalf("LoadSharedProperties error = %v", err)
	}
	if props != nil {
		t.Errorf("expected nil, got %v", props)
	}
}

func TestLoadSharedProperties_EmptyDir(t *testing.T) {
	v := setupSharedPropsTestVault(t)
	// Dir exists (from Init) but no yaml files

	props, err := v.LoadSharedProperties()
	if err != nil {
		t.Fatalf("LoadSharedProperties error = %v", err)
	}
	if props != nil {
		t.Errorf("expected nil for empty dir, got %v", props)
	}
}

func TestLoadSharedProperties_MalformedYAML(t *testing.T) {
	v := setupSharedPropsTestVault(t)
	writeSharedProperty(t, v, "bad", "not: [valid: yaml: {")

	_, err := v.LoadSharedProperties()
	if err == nil {
		t.Fatal("expected error for malformed YAML")
	}
}

func TestLoadSharedProperties_Caching(t *testing.T) {
	v := setupSharedPropsTestVault(t)
	writeSharedProperty(t, v, "due_date", "type: date\n")

	props1, _ := v.LoadSharedProperties()
	// Modify file after first load
	writeSharedProperty(t, v, "due_date", "type: string\n")

	props2, _ := v.LoadSharedProperties()
	// Should return cached result
	if len(props1) != len(props2) {
		t.Errorf("caching failed: first=%d, second=%d", len(props1), len(props2))
	}
	if props2[0].Name != "due_date" {
		t.Errorf("expected cached due_date, got %q", props2[0].Name)
	}
	if props2[0].Type != "date" {
		t.Errorf("expected cached type date, got %q", props2[0].Type)
	}
}

func TestLoadSharedProperties_NameDerivedFromFilename(t *testing.T) {
	v := setupSharedPropsTestVault(t)
	// File contains a name field that should be ignored
	writeSharedProperty(t, v, "due_date", "name: something_else\ntype: date\n")

	props, err := v.LoadSharedProperties()
	if err != nil {
		t.Fatalf("LoadSharedProperties error = %v", err)
	}
	if len(props) != 1 {
		t.Fatalf("len(props) = %d, want 1", len(props))
	}
	if props[0].Name != "due_date" {
		t.Errorf("name = %q, want %q (should come from filename)", props[0].Name, "due_date")
	}
}

func TestLoadSharedProperties_NonYAMLFilesIgnored(t *testing.T) {
	v := setupSharedPropsTestVault(t)
	writeSharedProperty(t, v, "rating", "type: number\n")
	// Write non-YAML files
	dir := v.SharedPropertiesDir()
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("not a property"), 0644)
	os.WriteFile(filepath.Join(dir, ".DS_Store"), []byte(""), 0644)

	props, err := v.LoadSharedProperties()
	if err != nil {
		t.Fatalf("LoadSharedProperties error = %v", err)
	}
	if len(props) != 1 {
		t.Fatalf("len(props) = %d, want 1", len(props))
	}
	if props[0].Name != "rating" {
		t.Errorf("name = %q, want %q", props[0].Name, "rating")
	}
}

func TestValidateSharedProperties_DuplicateNames(t *testing.T) {
	props := []Property{
		{Name: "due_date", Type: "date"},
		{Name: "due_date", Type: "string"},
	}
	errs := ValidateSharedProperties(props)
	if len(errs) == 0 {
		t.Fatal("expected error for duplicate names")
	}
	if !strings.Contains(errs[0].Error(), "duplicate") {
		t.Errorf("error = %v, want mention of duplicate", errs[0])
	}
}

func TestValidateSharedProperties_ReservedName(t *testing.T) {
	props := []Property{
		{Name: "name", Type: "string"},
	}
	errs := ValidateSharedProperties(props)
	if len(errs) == 0 {
		t.Fatal("expected error for reserved name")
	}
}

func TestValidateSharedProperties_InvalidType(t *testing.T) {
	props := []Property{
		{Name: "bad", Type: "invalid"},
	}
	errs := ValidateSharedProperties(props)
	if len(errs) == 0 {
		t.Fatal("expected error for invalid type")
	}
}

func TestValidateSharedProperties_SelectWithoutOptions(t *testing.T) {
	props := []Property{
		{Name: "status", Type: "select"},
	}
	errs := ValidateSharedProperties(props)
	if len(errs) == 0 {
		t.Fatal("expected error for select without options")
	}
}

func TestValidateSharedProperties_RelationWithoutTarget(t *testing.T) {
	props := []Property{
		{Name: "author", Type: "relation"},
	}
	errs := ValidateSharedProperties(props)
	if len(errs) == 0 {
		t.Fatal("expected error for relation without target")
	}
}

func TestValidateSharedProperties_MissingName(t *testing.T) {
	props := []Property{
		{Type: "string"},
	}
	errs := ValidateSharedProperties(props)
	if len(errs) == 0 {
		t.Fatal("expected error for missing name")
	}
}

func TestValidateSchema_UseWithTypeField(t *testing.T) {
	shared := []Property{{Name: "due_date", Type: "date"}}
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Use: "due_date", Type: "string"},
		},
	}
	errs := ValidateSchema(schema, shared)
	if len(errs) == 0 {
		t.Fatal("expected error for use with type field")
	}
	if !strings.Contains(errs[0].Error(), "only \"pin\", \"emoji\", and \"description\"") {
		t.Errorf("error = %v", errs[0])
	}
}

func TestValidateSchema_UseWithOptionsField(t *testing.T) {
	shared := []Property{{Name: "priority", Type: "select", Options: []Option{{Value: "high"}}}}
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Use: "priority", Options: []Option{{Value: "urgent"}}},
		},
	}
	errs := ValidateSchema(schema, shared)
	if len(errs) == 0 {
		t.Fatal("expected error for use with options field")
	}
}

func TestValidateSchema_UseWithDefaultField(t *testing.T) {
	shared := []Property{{Name: "due_date", Type: "date"}}
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Use: "due_date", Default: "2026-01-01"},
		},
	}
	errs := ValidateSchema(schema, shared)
	if len(errs) == 0 {
		t.Fatal("expected error for use with default field")
	}
}

func TestValidateSchema_UseBothUseAndName(t *testing.T) {
	shared := []Property{{Name: "due_date", Type: "date"}}
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Use: "due_date", Name: "my_date"},
		},
	}
	errs := ValidateSchema(schema, shared)
	if len(errs) == 0 {
		t.Fatal("expected error for both use and name")
	}
	if !strings.Contains(errs[0].Error(), "mutually exclusive") {
		t.Errorf("error = %v", errs[0])
	}
}

func TestValidateSchema_UseNonexistent(t *testing.T) {
	shared := []Property{{Name: "due_date", Type: "date"}}
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Use: "nonexistent"},
		},
	}
	errs := ValidateSchema(schema, shared)
	if len(errs) == 0 {
		t.Fatal("expected error for nonexistent use")
	}
	if !strings.Contains(errs[0].Error(), "not found") {
		t.Errorf("error = %v", errs[0])
	}
}

func TestValidateSchema_LocalPropertyConflictsWithShared(t *testing.T) {
	shared := []Property{{Name: "due_date", Type: "date"}}
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Name: "due_date", Type: "string"},
		},
	}
	errs := ValidateSchema(schema, shared)
	if len(errs) == 0 {
		t.Fatal("expected error for name conflict")
	}
	if !strings.Contains(errs[0].Error(), "conflicts with a shared property") {
		t.Errorf("error = %v", errs[0])
	}
}

func TestValidateSchema_UseWithNoSharedProps(t *testing.T) {
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Use: "due_date"},
		},
	}
	errs := ValidateSchema(schema)
	if len(errs) != 0 {
		t.Errorf("expected no errors without shared props context, got %v", errs)
	}
}

func TestResolveUseEntries_NoOverride(t *testing.T) {
	shared := map[string]Property{
		"due_date": {Name: "due_date", Type: "date", Emoji: "📅"},
	}
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Use: "due_date"},
		},
	}
	if err := resolveUseEntries(schema, shared); err != nil {
		t.Fatalf("resolveUseEntries error = %v", err)
	}
	p := schema.Properties[0]
	if p.Name != "due_date" || p.Type != "date" || p.Emoji != "📅" || p.Use != "" {
		t.Errorf("resolved = %+v", p)
	}
}

func TestResolveUseEntries_WithPinOverride(t *testing.T) {
	shared := map[string]Property{
		"due_date": {Name: "due_date", Type: "date", Emoji: "📅"},
	}
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Use: "due_date", Pin: 2},
		},
	}
	if err := resolveUseEntries(schema, shared); err != nil {
		t.Fatalf("resolveUseEntries error = %v", err)
	}
	if schema.Properties[0].Pin != 2 {
		t.Errorf("pin = %d, want 2", schema.Properties[0].Pin)
	}
}

func TestResolveUseEntries_WithEmojiOverride(t *testing.T) {
	shared := map[string]Property{
		"due_date": {Name: "due_date", Type: "date", Emoji: "📅"},
	}
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Use: "due_date", Emoji: "🗓️"},
		},
	}
	if err := resolveUseEntries(schema, shared); err != nil {
		t.Fatalf("resolveUseEntries error = %v", err)
	}
	if schema.Properties[0].Emoji != "🗓️" {
		t.Errorf("emoji = %q, want 🗓️", schema.Properties[0].Emoji)
	}
}

func TestResolveUseEntries_NotFound(t *testing.T) {
	shared := map[string]Property{}
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Use: "nonexistent"},
		},
	}
	if err := resolveUseEntries(schema, shared); err == nil {
		t.Fatal("expected error for nonexistent use")
	}
}

func TestResolveUseEntries_UseClearedAfterResolution(t *testing.T) {
	shared := map[string]Property{
		"due_date": {Name: "due_date", Type: "date"},
	}
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Use: "due_date"},
		},
	}
	resolveUseEntries(schema, shared)
	if schema.Properties[0].Use != "" {
		t.Errorf("Use field = %q, want empty after resolution", schema.Properties[0].Use)
	}
}

func TestResolveUseEntries_PreservesOrder(t *testing.T) {
	shared := map[string]Property{
		"due_date": {Name: "due_date", Type: "date"},
	}
	schema := &TypeSchema{
		Name: "test",
		Properties: []Property{
			{Name: "title", Type: "string"},
			{Use: "due_date"},
			{Name: "budget", Type: "number"},
		},
	}
	if err := resolveUseEntries(schema, shared); err != nil {
		t.Fatalf("resolveUseEntries error = %v", err)
	}
	names := []string{
		schema.Properties[0].Name,
		schema.Properties[1].Name,
		schema.Properties[2].Name,
	}
	if names[0] != "title" || names[1] != "due_date" || names[2] != "budget" {
		t.Errorf("order = %v, want [title due_date budget]", names)
	}
}

func TestLoadType_WithUseResolution(t *testing.T) {
	v := setupSharedPropsTestVault(t)
	if err := v.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer v.Close()

	writeSharedProperty(t, v, "due_date", "type: date\nemoji: 📅\n")

	mustWriteTypeSchema(v, "project", []byte(`name: project
properties:
  - name: title
    type: string
  - use: due_date
    pin: 1
`))

	schema, err := v.LoadType("project")
	if err != nil {
		t.Fatalf("LoadType error = %v", err)
	}
	if len(schema.Properties) != 2 {
		t.Fatalf("properties = %d, want 2", len(schema.Properties))
	}
	dd := schema.Properties[1]
	if dd.Name != "due_date" || dd.Type != "date" || dd.Pin != 1 || dd.Emoji != "📅" || dd.Use != "" {
		t.Errorf("resolved due_date = %+v", dd)
	}
}

func TestValidateAllSchemas_WithSharedProperties(t *testing.T) {
	v := setupSharedPropsTestVault(t)

	writeSharedProperty(t, v, "due_date", "type: date\n")

	mustWriteTypeSchema(v, "project", []byte(`name: project
properties:
  - use: due_date
`))

	result := ValidateAllSchemas(v)
	for k, errs := range result {
		t.Errorf("unexpected errors for %q: %v", k, errs)
	}
}

func TestValidateAllSchemas_SharedPropertiesErrors(t *testing.T) {
	v := setupSharedPropsTestVault(t)

	// Create property with reserved name
	writeSharedProperty(t, v, "name", "type: string\n")

	result := ValidateAllSchemas(v)
	if _, ok := result["_shared_properties"]; !ok {
		t.Fatal("expected _shared_properties errors")
	}
}

func TestLoadType_UseDescriptionOverride(t *testing.T) {
	v := setupSharedPropsTestVault(t)
	if err := v.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer v.Close()

	writeSharedProperty(t, v, "due_date", "type: date\nemoji: 📅\ndescription: \"A date something is due\"\n")

	mustWriteTypeSchema(v, "project", []byte(`name: project
properties:
  - use: due_date
    description: "Project deadline"
`))

	schema, err := v.LoadType("project")
	if err != nil {
		t.Fatalf("LoadType error = %v", err)
	}
	dd := schema.Properties[0]
	if dd.Description != "Project deadline" {
		t.Errorf("description = %q, want %q", dd.Description, "Project deadline")
	}
}

func TestLoadType_UseDescriptionPreservedWhenNoOverride(t *testing.T) {
	v := setupSharedPropsTestVault(t)
	if err := v.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer v.Close()

	writeSharedProperty(t, v, "due_date", "type: date\ndescription: \"A date something is due\"\n")

	mustWriteTypeSchema(v, "project", []byte(`name: project
properties:
  - use: due_date
`))

	schema, err := v.LoadType("project")
	if err != nil {
		t.Fatalf("LoadType error = %v", err)
	}
	dd := schema.Properties[0]
	if dd.Description != "A date something is due" {
		t.Errorf("description = %q, want %q", dd.Description, "A date something is due")
	}
}

func TestValidateUseOverrides_DescriptionAllowed(t *testing.T) {
	prop := Property{Use: "due_date", Description: "Project deadline"}
	if err := validateUseOverrides(0, prop); err != nil {
		t.Errorf("validateUseOverrides() = %v, want nil", err)
	}
}

func TestValidateUseOverrides_AllAllowedOverrides(t *testing.T) {
	prop := Property{Use: "due_date", Pin: 1, Emoji: "🗓️", Description: "Project deadline"}
	if err := validateUseOverrides(0, prop); err != nil {
		t.Errorf("validateUseOverrides() = %v, want nil", err)
	}
}

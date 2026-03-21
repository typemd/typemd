package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

// formatContext holds state for format BDD scenarios.
type formatContext struct {
	objectFormatResult *FormatResult
	schemaFormatResult *FormatResult
}

func newFormatContext() *formatContext {
	return &formatContext{}
}

// formatSetupTypeWithProperties creates a type schema. Opens the vault if
// not already opened (Background "a vault is initialized" only calls Init).
func (dc *domainContext) formatSetupTypeWithProperties(typeName string, table *godog.Table) {
	if dc.vault.Objects == nil {
		if err := dc.vault.Open(); err != nil {
			panic(fmt.Sprintf("vault open failed: %v", err))
		}
	}
	schema := &TypeSchema{
		Name: typeName,
	}
	for i, row := range table.Rows {
		if i == 0 {
			continue // skip header
		}
		prop := Property{
			Name: row.Cells[0].Value,
			Type: row.Cells[1].Value,
		}
		schema.Properties = append(schema.Properties, prop)
	}

	if err := dc.vault.SaveType(schema); err != nil {
		panic(fmt.Sprintf("save type failed: %v", err))
	}
}

func (dc *domainContext) formatWriteRawObject(objectRef string, content *godog.DocString) {
	parts := strings.SplitN(objectRef, "/", 2)
	if len(parts) != 2 {
		panic(fmt.Sprintf("invalid object ref: %s", objectRef))
	}
	typeName, slug := parts[0], parts[1]

	objDir := filepath.Join(dc.vault.ObjectsDir(), typeName)
	os.MkdirAll(objDir, 0755)

	ulid := mustULID()
	filename := slug + "-" + ulid
	objPath := filepath.Join(objDir, filename+".md")

	if err := os.WriteFile(objPath, []byte(content.Content), 0644); err != nil {
		panic(fmt.Sprintf("write object failed: %v", err))
	}

	// Track the object for later reference
	dc.objects[objectRef] = &Object{
		ID:       typeName + "/" + filename,
		Type:     typeName,
		Filename: filename,
	}
}

func (dc *domainContext) formatCreateObjectViaVault(objectRef, name string) {
	parts := strings.SplitN(objectRef, "/", 2)
	if len(parts) != 2 {
		panic(fmt.Sprintf("invalid object ref: %s", objectRef))
	}
	typeName := parts[0]

	obj, err := dc.vault.NewObject(typeName, name, "")
	if err != nil {
		panic(fmt.Sprintf("create object failed: %v", err))
	}
	dc.objects[objectRef] = obj
}

func (dc *domainContext) formatAllObjects(fc *formatContext) error {
	result, err := dc.vault.FormatObjects("", false)
	fc.objectFormatResult = result
	dc.lastErr = err
	return nil
}

func (dc *domainContext) formatObjectsOfType(fc *formatContext, typeName string) error {
	if dc.vault.Objects == nil {
		if err := dc.vault.Open(); err != nil {
			dc.lastErr = err
			return nil
		}
	}
	result, err := dc.vault.FormatAll(typeName, false)
	fc.objectFormatResult = result
	dc.lastErr = err
	return nil
}

func (dc *domainContext) formatAllObjectsDryRun(fc *formatContext) error {
	result, err := dc.vault.FormatObjects("", true)
	fc.objectFormatResult = result
	dc.lastErr = err
	return nil
}

func (fc *formatContext) nObjectsShouldBeFormatted(n int) error {
	result := fc.objectFormatResult
	if result == nil {
		if n == 0 {
			return nil
		}
		return fmt.Errorf("expected %d formatted objects, got nil result", n)
	}
	if len(result.Changed) != n {
		return fmt.Errorf("expected %d formatted objects, got %d: %v", n, len(result.Changed), result.Changed)
	}
	return nil
}

func (dc *domainContext) formatCheckPropertyOrder(objectRef, first, second string) error {
	obj := dc.objects[objectRef]
	if obj == nil {
		return fmt.Errorf("object %s not found", objectRef)
	}

	data, err := os.ReadFile(dc.vault.ObjectPath(obj.Type, obj.Filename))
	if err != nil {
		return fmt.Errorf("read object: %w", err)
	}

	content := string(data)
	firstIdx := strings.Index(content, first+":")
	secondIdx := strings.Index(content, second+":")

	if firstIdx < 0 {
		return fmt.Errorf("property %q not found in file", first)
	}
	if secondIdx < 0 {
		return fmt.Errorf("property %q not found in file", second)
	}
	if firstIdx >= secondIdx {
		return fmt.Errorf("expected %q (at %d) before %q (at %d)", first, firstIdx, second, secondIdx)
	}
	return nil
}

func (dc *domainContext) formatCheckBody(objectRef, expectedBody string) error {
	obj := dc.objects[objectRef]
	if obj == nil {
		return fmt.Errorf("object %s not found", objectRef)
	}

	data, err := os.ReadFile(dc.vault.ObjectPath(obj.Type, obj.Filename))
	if err != nil {
		return fmt.Errorf("read object: %w", err)
	}

	_, body, err := parseFrontmatter(data)
	if err != nil {
		return fmt.Errorf("parse frontmatter: %w", err)
	}

	trimmed := strings.TrimSpace(body)
	if trimmed != expectedBody {
		return fmt.Errorf("expected body %q, got %q", expectedBody, trimmed)
	}
	return nil
}

func (dc *domainContext) formatCheckProperty(objectRef, propName, expected string) error {
	obj := dc.objects[objectRef]
	if obj == nil {
		return fmt.Errorf("object %s not found", objectRef)
	}

	data, err := os.ReadFile(dc.vault.ObjectPath(obj.Type, obj.Filename))
	if err != nil {
		return fmt.Errorf("read object: %w", err)
	}

	props, _, err := parseFrontmatter(data)
	if err != nil {
		return fmt.Errorf("parse frontmatter: %w", err)
	}

	val := fmt.Sprintf("%v", props[propName])
	if val != expected {
		return fmt.Errorf("expected property %q = %q, got %q", propName, expected, val)
	}
	return nil
}

// formatSetupTypeWithRawSchema writes a schema file with raw YAML content.
func (dc *domainContext) formatSetupTypeWithRawSchema(typeName string, content *godog.DocString) {
	if dc.vault.Objects == nil {
		if err := dc.vault.Open(); err != nil {
			panic(fmt.Sprintf("vault open failed: %v", err))
		}
	}
	schemaDir := filepath.Join(dc.vault.TypesDir(), typeName)
	os.MkdirAll(schemaDir, 0755)
	schemaPath := filepath.Join(schemaDir, "schema.yaml")
	if err := os.WriteFile(schemaPath, []byte(content.Content), 0644); err != nil {
		panic(fmt.Sprintf("write schema failed: %v", err))
	}
}

func (dc *domainContext) formatAllSchemas(fc *formatContext) error {
	result, err := dc.vault.FormatSchemas("", false)
	fc.schemaFormatResult = result
	dc.lastErr = err
	return nil
}

func (fc *formatContext) nSchemasShouldBeFormatted(n int) error {
	result := fc.schemaFormatResult
	if result == nil {
		if n == 0 {
			return nil
		}
		return fmt.Errorf("expected %d formatted schemas, got nil result", n)
	}
	if len(result.Changed) != n {
		return fmt.Errorf("expected %d formatted schemas, got %d: %v", n, len(result.Changed), result.Changed)
	}
	return nil
}

func initFormatSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	fc := newFormatContext()

	ctx.Step(`^a type "([^"]*)" with properties:$`, dc.formatSetupTypeWithProperties)
	ctx.Step(`^an object "([^"]*)" with raw frontmatter:$`, dc.formatWriteRawObject)
	ctx.Step(`^an object "([^"]*)" created through the vault with name "([^"]*)"$`, dc.formatCreateObjectViaVault)
	ctx.Step(`^I format all objects$`, func() error { return dc.formatAllObjects(fc) })
	ctx.Step(`^I format objects of type "([^"]*)"$`, func(t string) error { return dc.formatObjectsOfType(fc, t) })
	ctx.Step(`^I format all objects in dry-run mode$`, func() error { return dc.formatAllObjectsDryRun(fc) })
	ctx.Step(`^(\d+) objects? should be formatted$`, fc.nObjectsShouldBeFormatted)
	ctx.Step(`^the object "([^"]*)" frontmatter should have "([^"]*)" before "([^"]*)"$`, dc.formatCheckPropertyOrder)
	ctx.Step(`^the object "([^"]*)" file should still have "([^"]*)" before "([^"]*)"$`, dc.formatCheckPropertyOrder)
	ctx.Step(`^the object "([^"]*)" body should be "([^"]*)"$`, dc.formatCheckBody)
	ctx.Step(`^the object "([^"]*)" property "([^"]*)" should be "([^"]*)"$`, dc.formatCheckProperty)
	ctx.Step(`^a type "([^"]*)" with raw schema:$`, dc.formatSetupTypeWithRawSchema)
	ctx.Step(`^I format all schemas$`, func() error { return dc.formatAllSchemas(fc) })
	ctx.Step(`^(\d+) schemas? should be formatted$`, fc.nSchemasShouldBeFormatted)
}

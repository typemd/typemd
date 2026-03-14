package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

// ── Property type steps ─────────────────────────────────────────────────────

func (dc *domainContext) aTypeSchemaWithAll9PropertyTypes() {
	schema := `name: complete
properties:
  - name: title
    type: string
  - name: count
    type: number
  - name: published
    type: date
  - name: due_at
    type: datetime
  - name: homepage
    type: url
  - name: active
    type: checkbox
  - name: status
    type: select
    options:
      - value: draft
      - value: published
  - name: labels
    type: multi_select
    options:
      - value: go
      - value: rust
  - name: author
    type: relation
    target: person
`
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "complete.yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithAnEnumProperty(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: status
    type: enum
    values:
      - to-read
      - reading
      - done
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithADateProperty(typeName string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: date\n    type: date\n", typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithAURLProperty(typeName string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: link\n    type: url\n", typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithASelectStatusProperty(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: status
    type: select
    options:
      - value: to-read
      - value: reading
      - value: done
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

// aObjectNamedExistsWithRawProperty creates an object and writes a property directly
// to the file, bypassing SetProperty validation. This is needed for negative test cases.
func (dc *domainContext) aObjectNamedExistsWithRawProperty(typeName, name, prop, value string) {
	dc.aObjectNamedExists(typeName, name)
	dc.currentObject.Properties[prop] = value
	if err := dc.vault.SaveObject(dc.currentObject); err != nil {
		panic(fmt.Sprintf("saveObjectFile failed: %v", err))
	}
	// Re-sync to update DB
	dc.vault.SyncIndex()
}

func (dc *domainContext) iValidateTheObjectAgainstItsSchema() {
	if dc.currentObject == nil {
		dc.lastErr = fmt.Errorf("no current object")
		return
	}
	schema, err := dc.vault.LoadType(dc.currentObject.Type)
	if err != nil {
		dc.lastErr = err
		return
	}
	obj, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		dc.lastErr = err
		return
	}
	dc.objectValidationErrors = ValidateObject(obj.Properties, schema)
}

func (dc *domainContext) theObjectShouldHaveNoValidationErrors() error {
	if len(dc.objectValidationErrors) != 0 {
		return fmt.Errorf("expected no validation errors, got %v", dc.objectValidationErrors)
	}
	return nil
}

func (dc *domainContext) theObjectShouldHaveValidationErrors() error {
	if len(dc.objectValidationErrors) == 0 {
		return fmt.Errorf("expected validation errors, got none")
	}
	return nil
}

func (dc *domainContext) iMigrateSchemas() {
	result, err := dc.vault.MigrateSchemas(false)
	dc.lastErr = err
	dc.schemaMigrateResult = result
}

func (dc *domainContext) theSchemaShouldUseSelectInsteadOfEnum(typeName string) error {
	schema, err := dc.vault.LoadType(typeName)
	if err != nil {
		return fmt.Errorf("LoadType(%q) error: %v", typeName, err)
	}
	for _, p := range schema.Properties {
		if p.Type == "enum" {
			return fmt.Errorf("property %q still uses type \"enum\"", p.Name)
		}
	}
	// Verify at least one select property exists (was converted)
	for _, p := range schema.Properties {
		if p.Type == "select" {
			return nil
		}
	}
	return fmt.Errorf("no select property found in schema %q", typeName)
}

func initPropertyTypeSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^an? "([^"]*)" object named "([^"]*)" exists with raw property "([^"]*)" set to "([^"]*)"$`, dc.aObjectNamedExistsWithRawProperty)
	ctx.Step(`^a type schema with all 9 property types$`, dc.aTypeSchemaWithAll9PropertyTypes)
	ctx.Step(`^a type schema "([^"]*)" with an enum property$`, dc.aTypeSchemaWithAnEnumProperty)
	ctx.Step(`^a type schema "([^"]*)" with a date property$`, dc.aTypeSchemaWithADateProperty)
	ctx.Step(`^a type schema "([^"]*)" with a url property$`, dc.aTypeSchemaWithAURLProperty)
	ctx.Step(`^a type schema "([^"]*)" with a select status property$`, dc.aTypeSchemaWithASelectStatusProperty)
	ctx.Step(`^I validate the object against its schema$`, dc.iValidateTheObjectAgainstItsSchema)
	ctx.Step(`^the object should have no validation errors$`, dc.theObjectShouldHaveNoValidationErrors)
	ctx.Step(`^the object should have validation errors$`, dc.theObjectShouldHaveValidationErrors)
	ctx.Step(`^I migrate schemas$`, dc.iMigrateSchemas)
	ctx.Step(`^the "([^"]*)" schema should use select instead of enum$`, dc.theSchemaShouldUseSelectInsteadOfEnum)
}

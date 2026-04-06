package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

// ── Type CRUD step state ────────────────────────────────────────────────────

type typeCrudContext struct {
	dc          *domainContext
	schema      *TypeSchema
	objectCount int
}

func newTypeCrudContext(dc *domainContext) *typeCrudContext {
	return &typeCrudContext{dc: dc}
}

// ── Given steps ─────────────────────────────────────────────────────────────

func (tc *typeCrudContext) aTypeSchemaWithNoExtraFields(name string) {
	tc.schema = &TypeSchema{
		Name: name,
	}
}

func (tc *typeCrudContext) theSchemaHasAStringProperty(propName string) {
	tc.schema.Properties = append(tc.schema.Properties, Property{
		Name: propName,
		Type: "string",
	})
}

func (tc *typeCrudContext) theSchemaHasASelectPropertyWithOptions(propName, optionsCSV string) {
	opts := []Option{}
	for _, v := range strings.Split(optionsCSV, ",") {
		opts = append(opts, Option{Value: strings.TrimSpace(v)})
	}
	tc.schema.Properties = append(tc.schema.Properties, Property{
		Name:    propName,
		Type:    "select",
		Options: opts,
	})
}

func (tc *typeCrudContext) aTypeSchemaFileExistsOnDisk(name string) {
	data := fmt.Sprintf("name: %s\nproperties: []\n", name)
	dirPath := filepath.Join(tc.dc.vault.TypesDir(), name)
	os.MkdirAll(dirPath, 0755)
	os.WriteFile(filepath.Join(dirPath, "schema.yaml"), []byte(data), 0644)
}

func (tc *typeCrudContext) iAddANumberPropertyToTheSchema(propName string) {
	tc.schema.Properties = append(tc.schema.Properties, Property{
		Name: propName,
		Type: "number",
	})
}

func (tc *typeCrudContext) aCustomTypeSchemaWithEmoji(typeName, emoji string) {
	data := fmt.Sprintf("name: %s\nemoji: %s\nproperties: []\n", typeName, emoji)
	mustWriteTypeSchema(tc.dc.vault, typeName, []byte(data))
}

func (tc *typeCrudContext) aCustomTagTypeSchemaWithoutUniqueField() {
	data := "name: tag\nemoji: \"🏷️\"\nproperties:\n  - name: color\n    type: string\n  - name: icon\n    type: string\n"
	mustWriteTypeSchema(tc.dc.vault, "tag", []byte(data))
}

// ── When steps ──────────────────────────────────────────────────────────────

func (tc *typeCrudContext) iDeleteSchema(name string) {
	tc.dc.lastErr = tc.dc.vault.repo.DeleteSchema(name)
}

func (tc *typeCrudContext) iSaveTheTypeSchema() {
	tc.dc.lastErr = tc.dc.vault.SaveType(tc.schema)
}

func (tc *typeCrudContext) iDeleteType(name string) {
	tc.dc.lastErr = tc.dc.vault.DeleteType(name)
}

func (tc *typeCrudContext) iCountObjectsOfType(typeName string) {
	count, err := tc.dc.vault.CountObjectsByType(typeName)
	tc.objectCount = count
	tc.dc.lastErr = err
}

// ── Then steps ──────────────────────────────────────────────────────────────

func (tc *typeCrudContext) theTypeSchemaFileShouldNotExistOnDisk(name string) error {
	dirPath := filepath.Join(tc.dc.vault.TypesDir(), name, "schema.yaml")
	if _, err := os.Stat(dirPath); err == nil {
		return fmt.Errorf("expected type schema %s/schema.yaml to not exist", name)
	}
	return nil
}

func (tc *typeCrudContext) theTypeSchemaFileShouldExistOnDisk(name string) error {
	dirPath := filepath.Join(tc.dc.vault.TypesDir(), name, "schema.yaml")
	if _, err := os.Stat(dirPath); err == nil {
		return nil
	}
	return fmt.Errorf("expected type schema %s/schema.yaml to exist", name)
}

func (tc *typeCrudContext) loadingTypeShouldReturnASchemaWithNProperties(name string, n int) error {
	schema, err := tc.dc.vault.LoadType(name)
	if err != nil {
		return fmt.Errorf("LoadType(%q) failed: %v", name, err)
	}
	if len(schema.Properties) != n {
		return fmt.Errorf("expected %d properties, got %d", n, len(schema.Properties))
	}
	return nil
}

func (tc *typeCrudContext) theErrorMessageShouldContain(substr string) error {
	if tc.dc.lastErr == nil {
		return fmt.Errorf("expected an error, got nil")
	}
	if !strings.Contains(tc.dc.lastErr.Error(), substr) {
		return fmt.Errorf("expected error to contain %q, got %q", substr, tc.dc.lastErr.Error())
	}
	return nil
}

func (tc *typeCrudContext) theCountShouldBe(expected int) error {
	if tc.objectCount != expected {
		return fmt.Errorf("expected count %d, got %d", expected, tc.objectCount)
	}
	return nil
}

// ── Init ────────────────────────────────────────────────────────────────────

func initTypeCrudSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	tc := newTypeCrudContext(dc)

	// Given
	ctx.Step(`^a type schema "([^"]*)" with no extra fields$`, tc.aTypeSchemaWithNoExtraFields)
	ctx.Step(`^the schema has a "([^"]*)" string property$`, tc.theSchemaHasAStringProperty)
	ctx.Step(`^the schema has a "([^"]*)" select property with options "([^"]*)"$`, tc.theSchemaHasASelectPropertyWithOptions)
	ctx.Step(`^a type schema file "([^"]*)" exists on disk$`, tc.aTypeSchemaFileExistsOnDisk)
	ctx.Step(`^I add a "([^"]*)" number property to the schema$`, tc.iAddANumberPropertyToTheSchema)
	ctx.Step(`^a custom "([^"]*)" type schema with emoji "([^"]*)"$`, tc.aCustomTypeSchemaWithEmoji)
	ctx.Step(`^a custom tag type schema without unique field$`, tc.aCustomTagTypeSchemaWithoutUniqueField)

	// When
	ctx.Step(`^I delete schema "([^"]*)"$`, tc.iDeleteSchema)
	ctx.Step(`^I save the type schema$`, tc.iSaveTheTypeSchema)
	ctx.Step(`^I delete type "([^"]*)"$`, tc.iDeleteType)
	ctx.Step(`^I count objects of type "([^"]*)"$`, tc.iCountObjectsOfType)

	// Then
	ctx.Step(`^the type schema file "([^"]*)" should not exist on disk$`, tc.theTypeSchemaFileShouldNotExistOnDisk)
	ctx.Step(`^the type schema file "([^"]*)" should exist on disk$`, tc.theTypeSchemaFileShouldExistOnDisk)
	ctx.Step(`^loading type "([^"]*)" should return a schema with (\d+) propert(?:y|ies)$`, tc.loadingTypeShouldReturnASchemaWithNProperties)
	ctx.Step(`^the error message should contain "([^"]*)"$`, tc.theErrorMessageShouldContain)
	ctx.Step(`^the count should be (\d+)$`, tc.theCountShouldBe)
}

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

func (dc *domainContext) aTypeSchemaWithNameTemplate(typeName, template string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: name\n    template: %q\n  - name: content\n    type: string\n", typeName, template)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) iCreateAObjectWithNoName(typeName string) {
	obj, err := dc.vault.NewObject(typeName, "")
	dc.lastErr = err
	if err == nil {
		dc.currentObject = obj
	}
}

func (dc *domainContext) theObjectPropertyShouldContain(key, substr string) error {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	val := fmt.Sprintf("%v", got.Properties[key])
	if !strings.Contains(val, substr) {
		return fmt.Errorf("property %q = %q, does not contain %q", key, val, substr)
	}
	return nil
}

func (dc *domainContext) aTypeSchemaWithNamePropertyAndType(typeName, propType string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: name\n    type: %s\n", typeName, propType)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func initNameTemplateSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a type schema "([^"]*)" with name template "([^"]*)"$`, dc.aTypeSchemaWithNameTemplate)
	ctx.Step(`^I create a "([^"]*)" object with no name$`, dc.iCreateAObjectWithNoName)
	ctx.Step(`^the object property "([^"]*)" should contain "([^"]*)"$`, dc.theObjectPropertyShouldContain)
	ctx.Step(`^a type schema "([^"]*)" with name property and type "([^"]*)"$`, dc.aTypeSchemaWithNamePropertyAndType)
}

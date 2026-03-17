package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

func (dc *domainContext) aTypeSchemaWithDescription(typeName, description string) {
	schema := fmt.Sprintf("name: %s\ndescription: %q\nproperties:\n  - name: title\n    type: string\n", typeName, description)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithPropertyHavingDescription(typeName, propName, description string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: %s\n    type: string\n    description: %q\n", typeName, propName, description)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) theLoadedSchemaDescriptionShouldBe(expected string) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("no schema loaded")
	}
	if dc.loadedSchema.Description != expected {
		return fmt.Errorf("expected schema description %q, got %q", expected, dc.loadedSchema.Description)
	}
	return nil
}

func (dc *domainContext) theLoadedPropertyDescriptionShouldBe(propName, expected string) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("no schema loaded")
	}
	for _, p := range dc.loadedSchema.Properties {
		if p.Name == propName {
			if p.Description != expected {
				return fmt.Errorf("expected property %q description %q, got %q", propName, expected, p.Description)
			}
			return nil
		}
	}
	return fmt.Errorf("property %q not found in loaded schema", propName)
}

func initDescriptionSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a type schema "([^"]*)" with description "([^"]*)"$`, dc.aTypeSchemaWithDescription)
	ctx.Step(`^a type schema "([^"]*)" with property "([^"]*)" having description "([^"]*)"$`, dc.aTypeSchemaWithPropertyHavingDescription)
	ctx.Step(`^the loaded schema description should be "([^"]*)"$`, dc.theLoadedSchemaDescriptionShouldBe)
	ctx.Step(`^the loaded property "([^"]*)" description should be "([^"]*)"$`, dc.theLoadedPropertyDescriptionShouldBe)
}

package core

import (
	"fmt"

	"github.com/cucumber/godog"
)

// ── Plural display name steps ───────────────────────────────────────────────

func (dc *domainContext) aTypeSchemaWithPlural(typeName, plural string) {
	schema := fmt.Sprintf("name: %s\nplural: %s\nproperties:\n  - name: title\n    type: string\n", typeName, plural)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func (dc *domainContext) aTypeSchemaWithoutPlural(typeName string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: title\n    type: string\n", typeName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func (dc *domainContext) theLoadedSchemaPluralShouldBe(expected string) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("no schema loaded")
	}
	if dc.loadedSchema.Plural != expected {
		return fmt.Errorf("expected Plural %q, got %q", expected, dc.loadedSchema.Plural)
	}
	return nil
}

func (dc *domainContext) theLoadedSchemaPluralNameShouldBe(expected string) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("no schema loaded")
	}
	got := dc.loadedSchema.PluralName()
	if got != expected {
		return fmt.Errorf("expected PluralName() %q, got %q", expected, got)
	}
	return nil
}

func initPluralSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a type schema "([^"]*)" with plural "([^"]*)"$`, dc.aTypeSchemaWithPlural)
	ctx.Step(`^a type schema "([^"]*)" without plural$`, dc.aTypeSchemaWithoutPlural)
	ctx.Step(`^the loaded schema plural should be "([^"]*)"$`, dc.theLoadedSchemaPluralShouldBe)
	ctx.Step(`^the loaded schema PluralName should be "([^"]*)"$`, dc.theLoadedSchemaPluralNameShouldBe)
}

package core

import (
	"fmt"

	"github.com/cucumber/godog"
)

// ── Pinned property steps ───────────────────────────────────────────────────

func (dc *domainContext) aTypeSchemaWithPropertyHavingPin(typeName, propName string, pin int) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: %s\n    type: string\n    pin: %d\n", typeName, propName, pin)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func (dc *domainContext) aTypeSchemaWithPropertiesHavingUniquePins(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: status
    type: string
    pin: 1
  - name: rating
    type: number
    pin: 2
`, typeName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func (dc *domainContext) aTypeSchemaWithPropertiesHavingDuplicatePins(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: status
    type: string
    pin: 1
  - name: rating
    type: number
    pin: 1
`, typeName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func (dc *domainContext) aTypeSchemaWithSomePropertiesUnpinned(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: title
    type: string
  - name: author
    type: string
  - name: status
    type: string
    pin: 1
`, typeName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func initPinnedSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a type schema "([^"]*)" with property "([^"]*)" having pin (-?\d+)$`, dc.aTypeSchemaWithPropertyHavingPin)
	ctx.Step(`^a type schema "([^"]*)" with properties having unique pins$`, dc.aTypeSchemaWithPropertiesHavingUniquePins)
	ctx.Step(`^a type schema "([^"]*)" with properties having duplicate pins$`, dc.aTypeSchemaWithPropertiesHavingDuplicatePins)
	ctx.Step(`^a type schema "([^"]*)" with some properties unpinned$`, dc.aTypeSchemaWithSomePropertiesUnpinned)
}

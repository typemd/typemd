package core

import (
	"fmt"

	"github.com/cucumber/godog"
)

func (dc *domainContext) aTypeSchemaWithADateProperty(typeName string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: date\n    type: date\n", typeName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func (dc *domainContext) aTypeSchemaWithADatetimeProperty(typeName string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: due_at\n    type: datetime\n", typeName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func (dc *domainContext) theDisplayPropertyShouldHaveFormattedValue(key, expected string) error {
	if dc.displayProps == nil {
		return fmt.Errorf("no display properties built")
	}
	for _, dp := range dc.displayProps {
		if dp.Key == key {
			got := dp.FormatValue()
			if got != expected {
				return fmt.Errorf("display property %q FormatValue() = %q, want %q", key, got, expected)
			}
			return nil
		}
	}
	return fmt.Errorf("display property %q not found", key)
}

func initDateDisplayFormatSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a type schema "([^"]*)" with a date property$`, dc.aTypeSchemaWithADateProperty)
	ctx.Step(`^a type schema "([^"]*)" with a datetime property$`, dc.aTypeSchemaWithADatetimeProperty)
	ctx.Step(`^the display property "([^"]*)" should have formatted value "([^"]*)"$`, dc.theDisplayPropertyShouldHaveFormattedValue)
}

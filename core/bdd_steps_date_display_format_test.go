package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

func (dc *domainContext) aTypeSchemaWithADatetimeProperty(typeName string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: due_at\n    type: datetime\n", typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
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
	ctx.Step(`^a type schema "([^"]*)" with a datetime property$`, dc.aTypeSchemaWithADatetimeProperty)
	ctx.Step(`^the display property "([^"]*)" should have formatted value "([^"]*)"$`, dc.theDisplayPropertyShouldHaveFormattedValue)
}

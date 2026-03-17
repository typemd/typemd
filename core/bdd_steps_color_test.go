package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

func (dc *domainContext) aTypeSchemaWithColor(typeName, color string) {
	schema := fmt.Sprintf("name: %s\ncolor: \"%s\"\nproperties:\n  - name: title\n    type: string\n", typeName, color)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func initColorSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a type schema "([^"]*)" with color "([^"]*)"$`, dc.aTypeSchemaWithColor)
}

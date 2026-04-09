package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/cucumber/godog"
)

// ── Aliases steps ──────────────────────────────────────────────────────────

func (dc *domainContext) iSetTheAliasesOfToAnd(name, alias1, alias2 string) {
	obj := dc.objects[name]
	if obj == nil {
		panic(fmt.Sprintf("object %q not found", name))
	}
	obj.Properties[AliasesProperty] = []string{alias1, alias2}
}

func (dc *domainContext) iSetTheAliasOf(name, alias string) {
	obj := dc.objects[name]
	if obj == nil {
		panic(fmt.Sprintf("object %q not found", name))
	}
	obj.Properties[AliasesProperty] = []string{alias}
}

func (dc *domainContext) iSaveTheObjectNamed(name string) {
	obj := dc.objects[name]
	if obj == nil {
		dc.lastErr = fmt.Errorf("object %q not found", name)
		return
	}
	dc.lastErr = dc.vault.SaveObject(obj)
}

func (dc *domainContext) iSetTheNameOf(name, newName string) {
	obj := dc.objects[name]
	if obj == nil {
		panic(fmt.Sprintf("object %q not found", name))
	}
	obj.Properties[NameProperty] = newName
}

func (dc *domainContext) theFileOfShouldContainAfter(name, needle, after string) error {
	obj := dc.objects[name]
	if obj == nil {
		return fmt.Errorf("object %q not found", name)
	}
	data, err := os.ReadFile(dc.vault.ObjectPath(obj.Type, obj.Filename))
	if err != nil {
		return fmt.Errorf("ReadFile error: %v", err)
	}
	content := string(data)
	afterIdx := strings.Index(content, after)
	if afterIdx == -1 {
		return fmt.Errorf("file does not contain %q", after)
	}
	needleIdx := strings.Index(content[afterIdx:], needle)
	if needleIdx == -1 {
		return fmt.Errorf("file does not contain %q after %q", needle, after)
	}
	return nil
}

func (dc *domainContext) theFileOfShouldNotContain(name, needle string) error {
	obj := dc.objects[name]
	if obj == nil {
		return fmt.Errorf("object %q not found", name)
	}
	data, err := os.ReadFile(dc.vault.ObjectPath(obj.Type, obj.Filename))
	if err != nil {
		return fmt.Errorf("ReadFile error: %v", err)
	}
	if strings.Contains(string(data), needle) {
		return fmt.Errorf("file should not contain %q", needle)
	}
	return nil
}

func (dc *domainContext) theObjectShouldHaveAlias(name, alias string) error {
	obj, err := dc.vault.GetObject(dc.objects[name].ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	// Check []any form (from YAML parse)
	if aliases, ok := obj.Properties[AliasesProperty].([]any); ok {
		for _, a := range aliases {
			if fmt.Sprintf("%v", a) == alias {
				return nil
			}
		}
	}
	// Check []string form (in-memory)
	if aliasesStr, ok := obj.Properties[AliasesProperty].([]string); ok {
		for _, a := range aliasesStr {
			if a == alias {
				return nil
			}
		}
	}
	return fmt.Errorf("object %q does not have alias %q (aliases: %v)", name, alias, obj.Properties[AliasesProperty])
}

func (dc *domainContext) iDefineATypeSchemaWithAPropertyNamed(propName string) {
	schema := &TypeSchema{
		Name: "testtype",
		Properties: []Property{
			{Name: propName, Type: "string"},
		},
	}
	dc.lastErr = dc.vault.SaveType(schema)
}

func (dc *domainContext) schemaValidationShouldReportAReservedSystemPropertyErrorFor(propName string) error {
	if dc.lastErr == nil {
		return fmt.Errorf("expected an error, got nil")
	}
	if !strings.Contains(dc.lastErr.Error(), "reserved system property") {
		return fmt.Errorf("expected error about reserved system property, got %q", dc.lastErr.Error())
	}
	if !strings.Contains(dc.lastErr.Error(), propName) {
		return fmt.Errorf("expected error to mention %q, got %q", propName, dc.lastErr.Error())
	}
	return nil
}

func (dc *domainContext) theFileOfShouldContain(name, needle string) error {
	obj := dc.objects[name]
	if obj == nil {
		return fmt.Errorf("object %q not found", name)
	}
	data, err := os.ReadFile(dc.vault.ObjectPath(obj.Type, obj.Filename))
	if err != nil {
		return fmt.Errorf("ReadFile error: %v", err)
	}
	if !strings.Contains(string(data), needle) {
		return fmt.Errorf("file does not contain %q", needle)
	}
	return nil
}

func initAliasesSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^I set the aliases of "([^"]*)" to "([^"]*)" and "([^"]*)"$`, dc.iSetTheAliasesOfToAnd)
	ctx.Step(`^I set the aliases of "([^"]*)" to "([^"]*)"$`, dc.iSetTheAliasOf)
	ctx.Step(`^I save the object "([^"]*)"$`, dc.iSaveTheObjectNamed)
	ctx.Step(`^I set the name of "([^"]*)" to "([^"]*)"$`, dc.iSetTheNameOf)
	ctx.Step(`^the file of "([^"]*)" should contain "([^"]*)" after "([^"]*)"$`, dc.theFileOfShouldContainAfter)
	ctx.Step(`^the file of "([^"]*)" should contain "([^"]*)"$`, dc.theFileOfShouldContain)
	ctx.Step(`^the file of "([^"]*)" should not contain "([^"]*)"$`, dc.theFileOfShouldNotContain)
	ctx.Step(`^the object "([^"]*)" should have alias "([^"]*)"$`, dc.theObjectShouldHaveAlias)
	ctx.Step(`^I define a type schema with a property named "([^"]*)"$`, dc.iDefineATypeSchemaWithAPropertyNamed)
	ctx.Step(`^schema validation should report a reserved system property error for "([^"]*)"$`, dc.schemaValidationShouldReportAReservedSystemPropertyErrorFor)
}


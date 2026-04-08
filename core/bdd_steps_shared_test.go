package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

// ── Shared properties steps (per-property files) ────────────────────────

func (dc *domainContext) aSharedPropertyFileWithTypeAndEmoji(name, propType, emoji string) {
	dir := dc.vault.SharedPropertiesDir()
	os.MkdirAll(dir, 0755)
	content := fmt.Sprintf("type: %s\nemoji: %s\n", propType, emoji)
	os.WriteFile(filepath.Join(dir, name+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aSharedPropertyFileWithTypeAndOptions(name, propType, optionsCSV string) {
	dir := dc.vault.SharedPropertiesDir()
	os.MkdirAll(dir, 0755)
	var optLines string
	for _, opt := range strings.Split(optionsCSV, ",") {
		optLines += fmt.Sprintf("  - value: %s\n", strings.TrimSpace(opt))
	}
	content := fmt.Sprintf("type: %s\noptions:\n%s", propType, optLines)
	os.WriteFile(filepath.Join(dir, name+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aSharedPropertyFileWithType(name, propType string) {
	dir := dc.vault.SharedPropertiesDir()
	os.MkdirAll(dir, 0755)
	content := fmt.Sprintf("type: %s\n", propType)
	os.WriteFile(filepath.Join(dir, name+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) anEmptySharedPropertiesDirectory() {
	os.MkdirAll(dc.vault.SharedPropertiesDir(), 0755)
}

func (dc *domainContext) aNonYAMLFileInPropertiesDirectory(filename string) {
	dir := dc.vault.SharedPropertiesDir()
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, filename), []byte("not a property"), 0644)
}

func (dc *domainContext) aSharedPropertyFileWithNameOverride(name, propType, nameOverride string) {
	dir := dc.vault.SharedPropertiesDir()
	os.MkdirAll(dir, 0755)
	content := fmt.Sprintf("name: %s\ntype: %s\n", nameOverride, propType)
	os.WriteFile(filepath.Join(dir, name+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aSharedPropertiesFileWithDescribedProperties() {
	dir := dc.vault.SharedPropertiesDir()
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "due_date.yaml"), []byte("type: date\nemoji: 📅\ndescription: \"A date something is due\"\n"), 0644)
	os.WriteFile(filepath.Join(dir, "priority.yaml"), []byte("type: select\ndescription: \"How important this is\"\noptions:\n  - value: high\n  - value: medium\n  - value: low\n"), 0644)
}

func (dc *domainContext) iLoadSharedProperties() {
	dc.sharedProperties, dc.lastErr = dc.vault.LoadSharedProperties()
}

func (dc *domainContext) sharedPropertiesShouldContainNEntries(expected int) error {
	got := len(dc.sharedProperties)
	if got != expected {
		return fmt.Errorf("shared properties count = %d, want %d", got, expected)
	}
	return nil
}

func (dc *domainContext) sharedPropertyShouldHaveType(name, expectedType string) error {
	for _, p := range dc.sharedProperties {
		if p.Name == name {
			if p.Type != expectedType {
				return fmt.Errorf("shared property %q type = %q, want %q", name, p.Type, expectedType)
			}
			return nil
		}
	}
	return fmt.Errorf("shared property %q not found", name)
}

func (dc *domainContext) sharedPropertiesShouldHaveNoErrors() error {
	if errs, ok := dc.schemaErrors["_shared_properties"]; ok && len(errs) > 0 {
		return fmt.Errorf("expected no shared properties errors, got %v", errs)
	}
	return nil
}

func (dc *domainContext) sharedPropertiesShouldHaveErrors() error {
	errs, ok := dc.schemaErrors["_shared_properties"]
	if !ok || len(errs) == 0 {
		return fmt.Errorf("expected shared properties errors, got none")
	}
	return nil
}

func (dc *domainContext) aTypeSchemaWithUse(typeName, useName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
`, typeName, useName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(content))
}

func (dc *domainContext) aTypeSchemaWithUseAndPin(typeName, useName string, pin int) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
    pin: %d
`, typeName, useName, pin)
	mustWriteTypeSchema(dc.vault, typeName, []byte(content))
}

func (dc *domainContext) aTypeSchemaWithUseAndEmoji(typeName, useName, emoji string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
    emoji: %s
`, typeName, useName, emoji)
	mustWriteTypeSchema(dc.vault, typeName, []byte(content))
}

func (dc *domainContext) aTypeSchemaWithUseAndDisallowedTypeOverride(typeName, useName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
    type: string
`, typeName, useName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(content))
}

func (dc *domainContext) aTypeSchemaWithLocalProperty(typeName, propName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - name: %s
    type: string
`, typeName, propName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(content))
}

func (dc *domainContext) aTypeSchemaWithDuplicateUse(typeName, useName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
  - use: %s
`, typeName, useName, useName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(content))
}

func (dc *domainContext) aTypeSchemaWithBothUseAndNameOnSameEntry(typeName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: due_date
    name: my_date
`, typeName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(content))
}

func (dc *domainContext) iLoadType(typeName string) {
	dc.loadedSchema, dc.lastErr = dc.vault.LoadType(typeName)
}

func (dc *domainContext) theLoadedTypeShouldHaveNProperties(expected int) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("loaded schema is nil")
	}
	got := len(dc.loadedSchema.Properties)
	if got != expected {
		return fmt.Errorf("loaded type properties = %d, want %d", got, expected)
	}
	return nil
}

func (dc *domainContext) theLoadedPropertyShouldHaveType(propName, expectedType string) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("loaded schema is nil")
	}
	for _, p := range dc.loadedSchema.Properties {
		if p.Name == propName {
			if p.Type != expectedType {
				return fmt.Errorf("loaded property %q type = %q, want %q", propName, p.Type, expectedType)
			}
			return nil
		}
	}
	return fmt.Errorf("loaded property %q not found", propName)
}

func (dc *domainContext) theLoadedPropertyShouldHaveEmoji(propName, expectedEmoji string) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("loaded schema is nil")
	}
	for _, p := range dc.loadedSchema.Properties {
		if p.Name == propName {
			if p.Emoji != expectedEmoji {
				return fmt.Errorf("loaded property %q emoji = %q, want %q", propName, p.Emoji, expectedEmoji)
			}
			return nil
		}
	}
	return fmt.Errorf("loaded property %q not found", propName)
}

func (dc *domainContext) theLoadedPropertyShouldHavePin(propName string, expectedPin int) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("loaded schema is nil")
	}
	for _, p := range dc.loadedSchema.Properties {
		if p.Name == propName {
			if p.Pin != expectedPin {
				return fmt.Errorf("loaded property %q pin = %d, want %d", propName, p.Pin, expectedPin)
			}
			return nil
		}
	}
	return fmt.Errorf("loaded property %q not found", propName)
}

func (dc *domainContext) aTypeSchemaWithUseAndDescription(typeName, useName, description string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
    description: %q
`, typeName, useName, description)
	mustWriteTypeSchema(dc.vault, typeName, []byte(content))
}

func (dc *domainContext) aTypeSchemaWithMixedUseAndNameProperties(typeName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - name: title
    type: string
  - use: due_date
  - name: budget
    type: number
`, typeName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(content))
}

func (dc *domainContext) theLoadedSchemaPluralShouldBe(expected string) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("no schema loaded")
	}
	if dc.loadedSchema.Plural != expected {
		return fmt.Errorf("expected plural %q, got %q", expected, dc.loadedSchema.Plural)
	}
	return nil
}

func (dc *domainContext) theLoadedSchemaShouldHaveEmoji(expected string) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("no schema loaded")
	}
	if dc.loadedSchema.Emoji != expected {
		return fmt.Errorf("expected emoji %q, got %q", expected, dc.loadedSchema.Emoji)
	}
	return nil
}

func (dc *domainContext) theLoadedPropertyAtIndexShouldBe(index int, expectedName string) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("loaded schema is nil")
	}
	if index >= len(dc.loadedSchema.Properties) {
		return fmt.Errorf("index %d out of range (have %d properties)", index, len(dc.loadedSchema.Properties))
	}
	got := dc.loadedSchema.Properties[index].Name
	if got != expectedName {
		return fmt.Errorf("property at index %d = %q, want %q", index, got, expectedName)
	}
	return nil
}

func initSharedSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a shared property file "([^"]*)" with type "([^"]*)" and emoji "([^"]*)"$`, dc.aSharedPropertyFileWithTypeAndEmoji)
	ctx.Step(`^a shared property file "([^"]*)" with type "([^"]*)" and options "([^"]*)"$`, dc.aSharedPropertyFileWithTypeAndOptions)
	ctx.Step(`^a shared property file "([^"]*)" with type "([^"]*)"$`, dc.aSharedPropertyFileWithType)
	ctx.Step(`^a shared property file "([^"]*)" with type "([^"]*)" and name override "([^"]*)"$`, dc.aSharedPropertyFileWithNameOverride)
	ctx.Step(`^an empty shared properties directory$`, dc.anEmptySharedPropertiesDirectory)
	ctx.Step(`^a non-YAML file "([^"]*)" in the properties directory$`, dc.aNonYAMLFileInPropertiesDirectory)
	ctx.Step(`^a shared properties file with described properties$`, dc.aSharedPropertiesFileWithDescribedProperties)
	ctx.Step(`^I load shared properties$`, dc.iLoadSharedProperties)
	ctx.Step(`^shared properties should contain (\d+) entries?$`, dc.sharedPropertiesShouldContainNEntries)
	ctx.Step(`^shared property "([^"]*)" should have type "([^"]*)"$`, dc.sharedPropertyShouldHaveType)
	ctx.Step(`^shared properties should have no errors$`, dc.sharedPropertiesShouldHaveNoErrors)
	ctx.Step(`^shared properties should have errors$`, dc.sharedPropertiesShouldHaveErrors)
	ctx.Step(`^a type schema "([^"]*)" with use "([^"]*)"$`, dc.aTypeSchemaWithUse)
	ctx.Step(`^a type schema "([^"]*)" with use "([^"]*)" and pin (\d+)$`, dc.aTypeSchemaWithUseAndPin)
	ctx.Step(`^a type schema "([^"]*)" with use "([^"]*)" and emoji "([^"]*)"$`, dc.aTypeSchemaWithUseAndEmoji)
	ctx.Step(`^a type schema "([^"]*)" with use "([^"]*)" and disallowed type override$`, dc.aTypeSchemaWithUseAndDisallowedTypeOverride)
	ctx.Step(`^a type schema "([^"]*)" with local property "([^"]*)"$`, dc.aTypeSchemaWithLocalProperty)
	ctx.Step(`^a type schema "([^"]*)" with duplicate use "([^"]*)"$`, dc.aTypeSchemaWithDuplicateUse)
	ctx.Step(`^a type schema "([^"]*)" with both use and name on same entry$`, dc.aTypeSchemaWithBothUseAndNameOnSameEntry)
	ctx.Step(`^I load type "([^"]*)"$`, dc.iLoadType)
	ctx.Step(`^the loaded type should have (\d+) propert(?:y|ies)$`, dc.theLoadedTypeShouldHaveNProperties)
	ctx.Step(`^the loaded property "([^"]*)" should have type "([^"]*)"$`, dc.theLoadedPropertyShouldHaveType)
	ctx.Step(`^the loaded property "([^"]*)" should have emoji "([^"]*)"$`, dc.theLoadedPropertyShouldHaveEmoji)
	ctx.Step(`^the loaded property "([^"]*)" should have pin (\d+)$`, dc.theLoadedPropertyShouldHavePin)
	ctx.Step(`^a type schema "([^"]*)" with use "([^"]*)" and description "([^"]*)"$`, dc.aTypeSchemaWithUseAndDescription)
	ctx.Step(`^a type schema "([^"]*)" with mixed use and name properties$`, dc.aTypeSchemaWithMixedUseAndNameProperties)
	ctx.Step(`^the loaded schema should have emoji "([^"]*)"$`, dc.theLoadedSchemaShouldHaveEmoji)
	ctx.Step(`^the loaded schema plural should be "([^"]*)"$`, dc.theLoadedSchemaPluralShouldBe)
	ctx.Step(`^the loaded property at index (\d+) should be "([^"]*)"$`, dc.theLoadedPropertyAtIndexShouldBe)
}

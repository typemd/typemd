package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

// ── Shared properties steps ──────────────────────────────────────────────

func (dc *domainContext) aSharedPropertiesFileWithDateAndSelectProperties(prop1, prop2 string) {
	content := fmt.Sprintf(`properties:
  - name: %s
    type: date
    emoji: 📅
  - name: %s
    type: select
    options:
      - value: high
      - value: medium
      - value: low
`, prop1, prop2)
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(content), 0644)
}

func (dc *domainContext) anEmptySharedPropertiesFile() {
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(""), 0644)
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

func (dc *domainContext) aSharedPropertiesFileWithDuplicateProperties(name string) {
	content := fmt.Sprintf(`properties:
  - name: %s
    type: date
  - name: %s
    type: string
`, name, name)
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(content), 0644)
}

func (dc *domainContext) aSharedPropertiesFileWithAnInvalidPropertyType() {
	content := `properties:
  - name: bad_prop
    type: invalid
`
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(content), 0644)
}

func (dc *domainContext) aSharedPropertiesFileWithAPropertyNamedName() {
	content := `properties:
  - name: name
    type: string
`
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(content), 0644)
}

func (dc *domainContext) aSharedPropertiesFileWithASelectPropertyMissingOptions() {
	content := `properties:
  - name: status
    type: select
`
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(content), 0644)
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
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aTypeSchemaWithUseAndPin(typeName, useName string, pin int) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
    pin: %d
`, typeName, useName, pin)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aTypeSchemaWithUseAndEmoji(typeName, useName, emoji string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
    emoji: %s
`, typeName, useName, emoji)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aTypeSchemaWithUseAndDisallowedTypeOverride(typeName, useName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
    type: string
`, typeName, useName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aTypeSchemaWithLocalProperty(typeName, propName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - name: %s
    type: string
`, typeName, propName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aTypeSchemaWithDuplicateUse(typeName, useName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
  - use: %s
`, typeName, useName, useName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aTypeSchemaWithBothUseAndNameOnSameEntry(typeName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: due_date
    name: my_date
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
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
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aSharedPropertiesFileWithDescribedProperties() {
	content := `properties:
  - name: due_date
    type: date
    emoji: 📅
    description: "A date something is due"
  - name: priority
    type: select
    description: "How important this is"
    options:
      - value: high
      - value: medium
      - value: low
`
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(content), 0644)
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
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
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
	ctx.Step(`^a shared properties file with "([^"]*)" date and "([^"]*)" select properties$`, dc.aSharedPropertiesFileWithDateAndSelectProperties)
	ctx.Step(`^an empty shared properties file$`, dc.anEmptySharedPropertiesFile)
	ctx.Step(`^I load shared properties$`, dc.iLoadSharedProperties)
	ctx.Step(`^shared properties should contain (\d+) entries?$`, dc.sharedPropertiesShouldContainNEntries)
	ctx.Step(`^shared property "([^"]*)" should have type "([^"]*)"$`, dc.sharedPropertyShouldHaveType)
	ctx.Step(`^a shared properties file with duplicate "([^"]*)" properties$`, dc.aSharedPropertiesFileWithDuplicateProperties)
	ctx.Step(`^a shared properties file with an invalid property type$`, dc.aSharedPropertiesFileWithAnInvalidPropertyType)
	ctx.Step(`^a shared properties file with a property named "([^"]*)"$`, dc.aSharedPropertiesFileWithAPropertyNamedName)
	ctx.Step(`^a shared properties file with a select property missing options$`, dc.aSharedPropertiesFileWithASelectPropertyMissingOptions)
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
	ctx.Step(`^a shared properties file with described properties$`, dc.aSharedPropertiesFileWithDescribedProperties)
	ctx.Step(`^a type schema "([^"]*)" with mixed use and name properties$`, dc.aTypeSchemaWithMixedUseAndNameProperties)
	ctx.Step(`^the loaded schema should have emoji "([^"]*)"$`, dc.theLoadedSchemaShouldHaveEmoji)
	ctx.Step(`^the loaded property at index (\d+) should be "([^"]*)"$`, dc.theLoadedPropertyAtIndexShouldBe)
}

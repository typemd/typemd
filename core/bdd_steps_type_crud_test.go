package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
	"gopkg.in/yaml.v3"
)

// ── Type CRUD step state ────────────────────────────────────────────────────

type typeCrudContext struct {
	dc             *domainContext
	schema         *TypeSchema
	yamlOutput     []byte
	roundTripSchema *TypeSchema
	objectCount    int
}

func newTypeCrudContext(dc *domainContext) *typeCrudContext {
	return &typeCrudContext{dc: dc}
}

// ── Given steps ─────────────────────────────────────────────────────────────

func (tc *typeCrudContext) aTypeSchemaWithPluralAndEmoji(name, plural, emoji string) {
	tc.schema = &TypeSchema{
		Name:   name,
		Plural: plural,
		Emoji:  emoji,
	}
}

func (tc *typeCrudContext) aTypeSchemaWithNoExtraFields(name string) {
	tc.schema = &TypeSchema{
		Name: name,
	}
}

func (tc *typeCrudContext) theSchemaHasAStringProperty(propName string) {
	tc.schema.Properties = append(tc.schema.Properties, Property{
		Name: propName,
		Type: "string",
	})
}

func (tc *typeCrudContext) theSchemaHasANumberPropertyWithPinAndEmoji(propName string, pin int, emoji string) {
	tc.schema.Properties = append(tc.schema.Properties, Property{
		Name:  propName,
		Type:  "number",
		Pin:   pin,
		Emoji: emoji,
	})
}

func (tc *typeCrudContext) theSchemaHasARelationPropertyTargeting(propName, target string) {
	tc.schema.Properties = append(tc.schema.Properties, Property{
		Name:   propName,
		Type:   "relation",
		Target: target,
	})
}

func (tc *typeCrudContext) theSchemaHasASelectPropertyWithOptions(propName, optionsCSV string) {
	opts := []Option{}
	for _, v := range strings.Split(optionsCSV, ",") {
		opts = append(opts, Option{Value: strings.TrimSpace(v)})
	}
	tc.schema.Properties = append(tc.schema.Properties, Property{
		Name:    propName,
		Type:    "select",
		Options: opts,
	})
}

func (tc *typeCrudContext) theSchemaHasANameTemplate(tmpl string) {
	tc.schema.NameTemplate = tmpl
}

func (tc *typeCrudContext) aTypeSchemaFileExistsOnDisk(name string) {
	data := fmt.Sprintf("name: %s\nproperties: []\n", name)
	path := filepath.Join(tc.dc.vault.TypesDir(), name+".yaml")
	os.WriteFile(path, []byte(data), 0644)
}

func (tc *typeCrudContext) iAddANumberPropertyToTheSchema(propName string) {
	tc.schema.Properties = append(tc.schema.Properties, Property{
		Name: propName,
		Type: "number",
	})
}

// ── When steps ──────────────────────────────────────────────────────────────

func (tc *typeCrudContext) iSerializeTheTypeSchema() error {
	data, err := MarshalTypeSchema(tc.schema)
	tc.yamlOutput = data
	tc.dc.lastErr = err
	return nil
}

func (tc *typeCrudContext) iDeserializeTheYAMLOutputBackToATypeSchema() error {
	var raw struct {
		Name       string     `yaml:"name"`
		Plural     string     `yaml:"plural,omitempty"`
		Emoji      string     `yaml:"emoji,omitempty"`
		Unique     bool       `yaml:"unique,omitempty"`
		Properties []Property `yaml:"properties"`
	}
	if err := yaml.Unmarshal(tc.yamlOutput, &raw); err != nil {
		return err
	}
	schema := &TypeSchema{
		Name:       raw.Name,
		Plural:     raw.Plural,
		Emoji:      raw.Emoji,
		Unique:     raw.Unique,
		Properties: raw.Properties,
	}
	// Extract NameTemplate like GetSchema does
	filtered := schema.Properties[:0]
	for _, prop := range schema.Properties {
		if prop.Name == NameProperty {
			schema.NameTemplate = prop.Template
			continue
		}
		filtered = append(filtered, prop)
	}
	schema.Properties = filtered
	tc.roundTripSchema = schema
	return nil
}

func (tc *typeCrudContext) iDeleteSchema(name string) {
	tc.dc.lastErr = tc.dc.vault.repo.DeleteSchema(name)
}

func (tc *typeCrudContext) iSaveTheTypeSchema() {
	tc.dc.lastErr = tc.dc.vault.SaveType(tc.schema)
}

func (tc *typeCrudContext) iDeleteType(name string) {
	tc.dc.lastErr = tc.dc.vault.DeleteType(name)
}

func (tc *typeCrudContext) iCountObjectsOfType(typeName string) {
	count, err := tc.dc.vault.CountObjectsByType(typeName)
	tc.objectCount = count
	tc.dc.lastErr = err
}

// ── Then steps ──────────────────────────────────────────────────────────────

func (tc *typeCrudContext) theYAMLOutputShouldContain(substr string) error {
	if !strings.Contains(string(tc.yamlOutput), substr) {
		return fmt.Errorf("expected YAML to contain %q, got:\n%s", substr, string(tc.yamlOutput))
	}
	return nil
}

func (tc *typeCrudContext) theYAMLOutputShouldNotContain(substr string) error {
	if strings.Contains(string(tc.yamlOutput), substr) {
		return fmt.Errorf("expected YAML NOT to contain %q, got:\n%s", substr, string(tc.yamlOutput))
	}
	return nil
}

func (tc *typeCrudContext) theRoundTripSchemaNameShouldBe(expected string) error {
	if tc.roundTripSchema.Name != expected {
		return fmt.Errorf("expected name %q, got %q", expected, tc.roundTripSchema.Name)
	}
	return nil
}

func (tc *typeCrudContext) theRoundTripSchemaShouldHaveNProperties(n int) error {
	if len(tc.roundTripSchema.Properties) != n {
		return fmt.Errorf("expected %d properties, got %d", n, len(tc.roundTripSchema.Properties))
	}
	return nil
}

func (tc *typeCrudContext) theRoundTripSchemaPropertyShouldHavePin(propName string, pin int) error {
	for _, p := range tc.roundTripSchema.Properties {
		if p.Name == propName {
			if p.Pin != pin {
				return fmt.Errorf("expected property %q pin %d, got %d", propName, pin, p.Pin)
			}
			return nil
		}
	}
	return fmt.Errorf("property %q not found", propName)
}

func (tc *typeCrudContext) theTypeSchemaFileShouldNotExistOnDisk(name string) error {
	path := filepath.Join(tc.dc.vault.TypesDir(), name+".yaml")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("expected file %s to not exist", path)
	}
	return nil
}

func (tc *typeCrudContext) theTypeSchemaFileShouldExistOnDisk(name string) error {
	path := filepath.Join(tc.dc.vault.TypesDir(), name+".yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("expected file %s to exist", path)
	}
	return nil
}

func (tc *typeCrudContext) loadingTypeShouldReturnASchemaWithNProperties(name string, n int) error {
	schema, err := tc.dc.vault.LoadType(name)
	if err != nil {
		return fmt.Errorf("LoadType(%q) failed: %v", name, err)
	}
	if len(schema.Properties) != n {
		return fmt.Errorf("expected %d properties, got %d", n, len(schema.Properties))
	}
	return nil
}

func (tc *typeCrudContext) theErrorMessageShouldContain(substr string) error {
	if tc.dc.lastErr == nil {
		return fmt.Errorf("expected an error, got nil")
	}
	if !strings.Contains(tc.dc.lastErr.Error(), substr) {
		return fmt.Errorf("expected error to contain %q, got %q", substr, tc.dc.lastErr.Error())
	}
	return nil
}

func (tc *typeCrudContext) theCountShouldBe(expected int) error {
	if tc.objectCount != expected {
		return fmt.Errorf("expected count %d, got %d", expected, tc.objectCount)
	}
	return nil
}

// ── Init ────────────────────────────────────────────────────────────────────

func initTypeCrudSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	tc := newTypeCrudContext(dc)

	// Given
	ctx.Step(`^a type schema "([^"]*)" with plural "([^"]*)" and emoji "([^"]*)"$`, tc.aTypeSchemaWithPluralAndEmoji)
	ctx.Step(`^a type schema "([^"]*)" with no extra fields$`, tc.aTypeSchemaWithNoExtraFields)
	ctx.Step(`^the schema has a "([^"]*)" string property$`, tc.theSchemaHasAStringProperty)
	ctx.Step(`^the schema has a "([^"]*)" number property with pin (\d+) and emoji "([^"]*)"$`, tc.theSchemaHasANumberPropertyWithPinAndEmoji)
	ctx.Step(`^the schema has a "([^"]*)" relation property targeting "([^"]*)"$`, tc.theSchemaHasARelationPropertyTargeting)
	ctx.Step(`^the schema has a "([^"]*)" select property with options "([^"]*)"$`, tc.theSchemaHasASelectPropertyWithOptions)
	ctx.Step(`^the schema has a name template "([^"]*)"$`, tc.theSchemaHasANameTemplate)
	ctx.Step(`^a type schema file "([^"]*)" exists on disk$`, tc.aTypeSchemaFileExistsOnDisk)
	ctx.Step(`^I add a "([^"]*)" number property to the schema$`, tc.iAddANumberPropertyToTheSchema)

	// When
	ctx.Step(`^I serialize the type schema$`, tc.iSerializeTheTypeSchema)
	ctx.Step(`^I deserialize the YAML output back to a TypeSchema$`, tc.iDeserializeTheYAMLOutputBackToATypeSchema)
	ctx.Step(`^I delete schema "([^"]*)"$`, tc.iDeleteSchema)
	ctx.Step(`^I save the type schema$`, tc.iSaveTheTypeSchema)
	ctx.Step(`^I delete type "([^"]*)"$`, tc.iDeleteType)
	ctx.Step(`^I count objects of type "([^"]*)"$`, tc.iCountObjectsOfType)

	// Then
	ctx.Step(`^the YAML output should contain "([^"]*)"$`, tc.theYAMLOutputShouldContain)
	ctx.Step(`^the YAML output should not contain "([^"]*)"$`, tc.theYAMLOutputShouldNotContain)
	ctx.Step(`^the round-trip schema name should be "([^"]*)"$`, tc.theRoundTripSchemaNameShouldBe)
	ctx.Step(`^the round-trip schema should have (\d+) properties$`, tc.theRoundTripSchemaShouldHaveNProperties)
	ctx.Step(`^the round-trip schema property "([^"]*)" should have pin (\d+)$`, tc.theRoundTripSchemaPropertyShouldHavePin)
	ctx.Step(`^the type schema file "([^"]*)" should not exist on disk$`, tc.theTypeSchemaFileShouldNotExistOnDisk)
	ctx.Step(`^the type schema file "([^"]*)" should exist on disk$`, tc.theTypeSchemaFileShouldExistOnDisk)
	ctx.Step(`^loading type "([^"]*)" should return a schema with (\d+) propert(?:y|ies)$`, tc.loadingTypeShouldReturnASchemaWithNProperties)
	ctx.Step(`^the error message should contain "([^"]*)"$`, tc.theErrorMessageShouldContain)
	ctx.Step(`^the count should be (\d+)$`, tc.theCountShouldBe)
}

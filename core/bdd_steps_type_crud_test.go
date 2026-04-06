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
	dc              *domainContext
	schema          *TypeSchema
	yamlOutput      []byte
	roundTripSchema *TypeSchema
	objectCount     int
	validationErrs  []error
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
	dirPath := filepath.Join(tc.dc.vault.TypesDir(), name)
	os.MkdirAll(dirPath, 0755)
	os.WriteFile(filepath.Join(dirPath, "schema.yaml"), []byte(data), 0644)
}

func (tc *typeCrudContext) iAddANumberPropertyToTheSchema(propName string) {
	tc.schema.Properties = append(tc.schema.Properties, Property{
		Name: propName,
		Type: "number",
	})
}

func (tc *typeCrudContext) aCustomTypeSchemaWithEmoji(typeName, emoji string) {
	data := fmt.Sprintf("name: %s\nemoji: %s\nproperties: []\n", typeName, emoji)
	mustWriteTypeSchema(tc.dc.vault, typeName, []byte(data))
}

func (tc *typeCrudContext) aCustomTagTypeSchemaWithoutUniqueField() {
	data := "name: tag\nemoji: \"🏷️\"\nproperties:\n  - name: color\n    type: string\n  - name: icon\n    type: string\n"
	mustWriteTypeSchema(tc.dc.vault, "tag", []byte(data))
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
		Version    string     `yaml:"version,omitempty"`
		Properties []Property `yaml:"properties"`
	}
	if err := yaml.Unmarshal(tc.yamlOutput, &raw); err != nil {
		return err
	}
	version := raw.Version
	if version == "" {
		version = DefaultSchemaVersion
	}
	schema := &TypeSchema{
		Name:       raw.Name,
		Plural:     raw.Plural,
		Emoji:      raw.Emoji,
		Unique:     raw.Unique,
		Version:    version,
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
	dirPath := filepath.Join(tc.dc.vault.TypesDir(), name, "schema.yaml")
	if _, err := os.Stat(dirPath); err == nil {
		return fmt.Errorf("expected type schema %s/schema.yaml to not exist", name)
	}
	return nil
}

func (tc *typeCrudContext) theTypeSchemaFileShouldExistOnDisk(name string) error {
	dirPath := filepath.Join(tc.dc.vault.TypesDir(), name, "schema.yaml")
	if _, err := os.Stat(dirPath); err == nil {
		return nil
	}
	return fmt.Errorf("expected type schema %s/schema.yaml to exist", name)
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

// ── Version steps ───────────────────────────────────────────────────────────

func (tc *typeCrudContext) theSchemaHasVersion(version string) {
	tc.schema.Version = version
}

func (tc *typeCrudContext) iValidateTheTypeSchema() {
	tc.validationErrs = ValidateSchema(tc.schema)
}

func (tc *typeCrudContext) theRoundTripSchemaVersionShouldBe(expected string) error {
	if tc.roundTripSchema.Version != expected {
		return fmt.Errorf("expected version %q, got %q", expected, tc.roundTripSchema.Version)
	}
	return nil
}

func (tc *typeCrudContext) noSchemaValidationErrorsShouldOccur() error {
	if len(tc.validationErrs) != 0 {
		return fmt.Errorf("expected no validation errors, got %v", tc.validationErrs)
	}
	return nil
}

func (tc *typeCrudContext) aSchemaValidationErrorShouldMention(substr string) error {
	if len(tc.validationErrs) == 0 {
		return fmt.Errorf("expected validation errors, got none")
	}
	for _, err := range tc.validationErrs {
		if strings.Contains(err.Error(), substr) {
			return nil
		}
	}
	return fmt.Errorf("expected error mentioning %q, got %v", substr, tc.validationErrs)
}

func (tc *typeCrudContext) comparingVersionWithShouldReturn(a, b string, expected int) error {
	result := CompareVersions(a, b)
	if result != expected {
		return fmt.Errorf("CompareVersions(%q, %q) = %d, want %d", a, b, result, expected)
	}
	return nil
}

// ── Domain Event steps ─────────────────────────────────────────────────────

func (tc *typeCrudContext) iSubscribeToDomainEvents() {
	tc.dc.capturedEvents = nil
	tc.dc.vault.Events.Subscribe(func(e DomainEvent) {
		tc.dc.capturedEvents = append(tc.dc.capturedEvents, e)
	})
}

func (tc *typeCrudContext) anEventShouldHaveBeenEmitted(eventName string) error {
	for _, e := range tc.dc.capturedEvents {
		if e.eventName() == eventName {
			return nil
		}
	}
	return fmt.Errorf("expected event %q, got %v", eventName, tc.dc.capturedEvents)
}

func (tc *typeCrudContext) noDomainEventsShouldHaveBeenEmitted() error {
	if len(tc.dc.capturedEvents) > 0 {
		names := make([]string, len(tc.dc.capturedEvents))
		for i, e := range tc.dc.capturedEvents {
			names[i] = e.eventName()
		}
		return fmt.Errorf("expected no events, got %v", names)
	}
	return nil
}

func (tc *typeCrudContext) theTypeSavedEventSchemaNameShouldBe(expected string) error {
	for _, e := range tc.dc.capturedEvents {
		if ts, ok := e.(TypeSaved); ok {
			if ts.Schema.Name != expected {
				return fmt.Errorf("expected TypeSaved schema name %q, got %q", expected, ts.Schema.Name)
			}
			return nil
		}
	}
	return fmt.Errorf("no TypeSaved event found")
}

func (tc *typeCrudContext) theTypeDeletedEventNameShouldBe(expected string) error {
	for _, e := range tc.dc.capturedEvents {
		if td, ok := e.(TypeDeleted); ok {
			if td.Name != expected {
				return fmt.Errorf("expected TypeDeleted name %q, got %q", expected, td.Name)
			}
			return nil
		}
	}
	return fmt.Errorf("no TypeDeleted event found")
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
	ctx.Step(`^a custom "([^"]*)" type schema with emoji "([^"]*)"$`, tc.aCustomTypeSchemaWithEmoji)
	ctx.Step(`^a custom tag type schema without unique field$`, tc.aCustomTagTypeSchemaWithoutUniqueField)

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

	// Domain Event steps
	ctx.Step(`^I subscribe to domain events$`, tc.iSubscribeToDomainEvents)
	ctx.Step(`^a "([^"]*)" event should have been emitted$`, tc.anEventShouldHaveBeenEmitted)
	ctx.Step(`^no domain events should have been emitted$`, tc.noDomainEventsShouldHaveBeenEmitted)
	ctx.Step(`^the TypeSaved event schema name should be "([^"]*)"$`, tc.theTypeSavedEventSchemaNameShouldBe)
	ctx.Step(`^the TypeDeleted event name should be "([^"]*)"$`, tc.theTypeDeletedEventNameShouldBe)

	// Version steps
	ctx.Step(`^the schema has version "([^"]*)"$`, tc.theSchemaHasVersion)
	ctx.Step(`^I validate the type schema$`, tc.iValidateTheTypeSchema)
	ctx.Step(`^the round-trip schema version should be "([^"]*)"$`, tc.theRoundTripSchemaVersionShouldBe)
	ctx.Step(`^no schema validation errors should occur$`, tc.noSchemaValidationErrorsShouldOccur)
	ctx.Step(`^a schema validation error should mention "([^"]*)"$`, tc.aSchemaValidationErrorShouldMention)
	ctx.Step(`^comparing version "([^"]*)" with "([^"]*)" should return (-?\d+)$`, tc.comparingVersionWithShouldReturn)
}

package core

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
	"gopkg.in/yaml.v3"
)

// ── Type Schema Version step state ──────────────────────────────────────────

type versionContext struct {
	dc              *domainContext
	schema          *TypeSchema
	yamlOutput      []byte
	roundTripSchema *TypeSchema
	validationErrs  []error
}

func newVersionContext(dc *domainContext) *versionContext {
	return &versionContext{dc: dc}
}

// ── Given steps ─────────────────────────────────────────────────────────────

func (vc *versionContext) aVersionedTypeSchemaWithVersion(name string, version int) {
	vc.schema = &TypeSchema{
		Name:    name,
		Version: version,
	}
}

// ── When steps ──────────────────────────────────────────────────────────────

func (vc *versionContext) iSerializeTheVersionedSchema() error {
	data, err := MarshalTypeSchema(vc.schema)
	vc.yamlOutput = data
	vc.dc.lastErr = err
	return nil
}

func (vc *versionContext) iDeserializeTheVersionedYAMLOutput() error {
	var raw struct {
		Name       string     `yaml:"name"`
		Plural     string     `yaml:"plural,omitempty"`
		Emoji      string     `yaml:"emoji,omitempty"`
		Unique     bool       `yaml:"unique,omitempty"`
		Version    int        `yaml:"version,omitempty"`
		Properties []Property `yaml:"properties"`
	}
	if err := yaml.Unmarshal(vc.yamlOutput, &raw); err != nil {
		return err
	}
	vc.roundTripSchema = &TypeSchema{
		Name:       raw.Name,
		Plural:     raw.Plural,
		Emoji:      raw.Emoji,
		Unique:     raw.Unique,
		Version:    raw.Version,
		Properties: raw.Properties,
	}
	return nil
}

func (vc *versionContext) iValidateTheVersionedSchema() {
	vc.validationErrs = ValidateSchema(vc.schema)
}

// ── Then steps ──────────────────────────────────────────────────────────────

func (vc *versionContext) theVersionedYAMLOutputShouldContain(substr string) error {
	if !strings.Contains(string(vc.yamlOutput), substr) {
		return fmt.Errorf("expected YAML to contain %q, got:\n%s", substr, string(vc.yamlOutput))
	}
	return nil
}

func (vc *versionContext) theVersionedYAMLOutputShouldNotContain(substr string) error {
	if strings.Contains(string(vc.yamlOutput), substr) {
		return fmt.Errorf("expected YAML NOT to contain %q, got:\n%s", substr, string(vc.yamlOutput))
	}
	return nil
}

func (vc *versionContext) theRoundTripSchemaVersionShouldBe(expected int) error {
	if vc.roundTripSchema.Version != expected {
		return fmt.Errorf("expected version %d, got %d", expected, vc.roundTripSchema.Version)
	}
	return nil
}

func (vc *versionContext) noSchemaValidationErrorsShouldOccur() error {
	if len(vc.validationErrs) != 0 {
		return fmt.Errorf("expected no validation errors, got %v", vc.validationErrs)
	}
	return nil
}

func (vc *versionContext) aSchemaValidationErrorShouldMention(substr string) error {
	if len(vc.validationErrs) == 0 {
		return fmt.Errorf("expected validation errors, got none")
	}
	for _, err := range vc.validationErrs {
		if strings.Contains(err.Error(), substr) {
			return nil
		}
	}
	return fmt.Errorf("expected error mentioning %q, got %v", substr, vc.validationErrs)
}

// ── Init ────────────────────────────────────────────────────────────────────

func initVersionSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	vc := newVersionContext(dc)

	// Given
	ctx.Step(`^a versioned type schema "([^"]*)" with version (-?\d+)$`, vc.aVersionedTypeSchemaWithVersion)

	// When
	ctx.Step(`^I serialize the versioned schema$`, vc.iSerializeTheVersionedSchema)
	ctx.Step(`^I deserialize the versioned YAML output$`, vc.iDeserializeTheVersionedYAMLOutput)
	ctx.Step(`^I validate the versioned schema$`, vc.iValidateTheVersionedSchema)

	// Then
	ctx.Step(`^the versioned YAML output should contain "([^"]*)"$`, vc.theVersionedYAMLOutputShouldContain)
	ctx.Step(`^the versioned YAML output should not contain "([^"]*)"$`, vc.theVersionedYAMLOutputShouldNotContain)
	ctx.Step(`^the round-trip schema version should be (-?\d+)$`, vc.theRoundTripSchemaVersionShouldBe)
	ctx.Step(`^no schema validation errors should occur$`, vc.noSchemaValidationErrorsShouldOccur)
	ctx.Step(`^a schema validation error should mention "([^"]*)"$`, vc.aSchemaValidationErrorShouldMention)
}

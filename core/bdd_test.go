package core

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

type frontmatterContext struct {
	props   map[string]any
	body    string
	output  []byte
	parsed  map[string]any
	pBody   string
	rawMD   string
	lastErr error
}

func (fc *frontmatterContext) propertiesWithSetTo(key, value string) {
	if fc.props == nil {
		fc.props = make(map[string]any)
	}
	fc.props[key] = value
}

func (fc *frontmatterContext) emptyProperties() {
	fc.props = make(map[string]any)
}

func (fc *frontmatterContext) iWriteFrontmatterWithNoBody() error {
	data, err := writeFrontmatter(fc.props, "", nil)
	fc.output = data
	fc.lastErr = err
	return err
}

func (fc *frontmatterContext) iWriteFrontmatterWithBody(body string) error {
	data, err := writeFrontmatter(fc.props, body, nil)
	fc.output = data
	fc.lastErr = err
	return err
}

func (fc *frontmatterContext) theOutputShouldStartWith(prefix string) error {
	if !strings.HasPrefix(string(fc.output), prefix) {
		return fmt.Errorf("expected output to start with %q, got %q", prefix, string(fc.output)[:20])
	}
	return nil
}

func (fc *frontmatterContext) theOutputShouldContain(substr string) error {
	if !strings.Contains(string(fc.output), substr) {
		return fmt.Errorf("expected output to contain %q, got:\n%s", substr, string(fc.output))
	}
	return nil
}

func (fc *frontmatterContext) theOutputShouldEqual(expected string) error {
	// Unescape newlines from Gherkin
	expected = strings.ReplaceAll(expected, `\n`, "\n")
	if string(fc.output) != expected {
		return fmt.Errorf("expected output %q, got %q", expected, string(fc.output))
	}
	return nil
}

func (fc *frontmatterContext) markdownContent(content *godog.DocString) {
	fc.rawMD = content.Content
}

func (fc *frontmatterContext) iParseTheFrontmatter() error {
	props, body, err := parseFrontmatter([]byte(fc.rawMD))
	fc.parsed = props
	fc.pBody = body
	fc.lastErr = err
	return err
}

func (fc *frontmatterContext) theParsedPropertyShouldBe(key, expected string) error {
	val, ok := fc.parsed[key]
	if !ok {
		return fmt.Errorf("property %q not found in parsed frontmatter", key)
	}
	got := fmt.Sprintf("%v", val)
	if got != expected {
		return fmt.Errorf("expected property %q to be %q, got %q", key, expected, got)
	}
	return nil
}

func (fc *frontmatterContext) theParsedBodyShouldBe(expected string) error {
	trimmed := strings.TrimSpace(fc.pBody)
	if trimmed != expected {
		return fmt.Errorf("expected body %q, got %q", expected, trimmed)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	fc := &frontmatterContext{}

	ctx.Step(`^properties with "([^"]*)" set to "([^"]*)"$`, fc.propertiesWithSetTo)
	ctx.Step(`^empty properties$`, fc.emptyProperties)
	ctx.Step(`^I write frontmatter with no body$`, fc.iWriteFrontmatterWithNoBody)
	ctx.Step(`^I write frontmatter with body "([^"]*)"$`, fc.iWriteFrontmatterWithBody)
	ctx.Step(`^the output should start with "([^"]*)"$`, fc.theOutputShouldStartWith)
	ctx.Step(`^the output should contain "([^"]*)"$`, fc.theOutputShouldContain)
	ctx.Step(`^the output should equal "([^"]*)"$`, fc.theOutputShouldEqual)
	ctx.Step(`^markdown content:$`, fc.markdownContent)
	ctx.Step(`^I parse the frontmatter$`, fc.iParseTheFrontmatter)
	ctx.Step(`^the parsed property "([^"]*)" should be "([^"]*)"$`, fc.theParsedPropertyShouldBe)
	ctx.Step(`^the parsed body should be "([^"]*)"$`, fc.theParsedBodyShouldBe)

	// Domain steps
	dc := newDomainContext()

	// Cleanup after each scenario
	ctx.After(func(hookCtx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if dc.vault != nil {
			dc.vault.Close()
		}
		if dc.rootDir != "" {
			os.RemoveAll(dc.rootDir)
		}
		return hookCtx, nil
	})

	initCommonSteps(ctx, dc)
	initVaultSteps(ctx, dc)
	initObjectSteps(ctx, dc)
	initRelationSteps(ctx, dc)
	initQuerySteps(ctx, dc)
	initValidateSteps(ctx, dc)
	initWikiLinkSteps(ctx, dc)
	initResolveSteps(ctx, dc)
	initPropertyTypeSteps(ctx, dc)
	initPropertyEmojiSteps(ctx, dc)
	initPinnedSteps(ctx, dc)
	initFilteringSteps(ctx, dc)
	initNameSteps(ctx, dc)
	initSharedSteps(ctx, dc)
	initSystemSteps(ctx, dc)
	initTagSteps(ctx, dc)
	initPluralSteps(ctx, dc)
	initNameTemplateSteps(ctx, dc)
	initUniqueSteps(ctx, dc)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

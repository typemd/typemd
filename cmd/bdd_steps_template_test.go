package cmd

import (
	"fmt"
	"os"

	"github.com/cucumber/godog"
	"github.com/typemd/typemd/core"
)

func initTemplateSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	// Setup steps
	ctx.Step(`^a template "([^"]*)" with body "([^"]*)"$`, cc.aTemplateWithBody)
	ctx.Step(`^a template "([^"]*)" with property "([^"]*)" value "([^"]*)" and body "([^"]*)"$`, cc.aTemplateWithPropertyAndBody)

	// List steps
	ctx.Step(`^I run template list$`, cc.iRunTemplateList)
	ctx.Step(`^I run template list with type "([^"]*)"$`, cc.iRunTemplateListWithType)
	ctx.Step(`^I run template list with json$`, cc.iRunTemplateListWithJSON)

	// Show steps
	ctx.Step(`^I run template show "([^"]*)"$`, cc.iRunTemplateShow)

	// Create steps
	ctx.Step(`^I run template create "([^"]*)"$`, cc.iRunTemplateCreate)
	ctx.Step(`^the template file "([^"]*)" should exist$`, cc.theTemplateFileShouldExist)

	// Delete steps
	ctx.Step(`^I run template delete "([^"]*)" with force$`, cc.iRunTemplateDeleteWithForce)
	ctx.Step(`^the template file "([^"]*)" should not exist$`, cc.theTemplateFileShouldNotExist)
}

func (cc *cmdContext) aTemplateWithBody(templateID, body string) error {
	typeName, name, err := parseTemplateArg(templateID)
	if err != nil {
		return err
	}
	tmpl := &core.Template{
		Name:       name,
		Properties: map[string]any{},
		Body:       body + "\n",
	}
	return cc.vault.SaveTemplate(typeName, name, tmpl)
}

func (cc *cmdContext) aTemplateWithPropertyAndBody(templateID, propKey, propValue, body string) error {
	typeName, name, err := parseTemplateArg(templateID)
	if err != nil {
		return err
	}
	tmpl := &core.Template{
		Name:       name,
		Properties: map[string]any{propKey: propValue},
		Body:       body + "\n",
	}
	return cc.vault.SaveTemplate(typeName, name, tmpl)
}

func (cc *cmdContext) iRunTemplateList() {
	cc.runCmd("template", "list")
}

func (cc *cmdContext) iRunTemplateListWithType(typeName string) {
	cc.runCmd("template", "list", typeName)
}

func (cc *cmdContext) iRunTemplateListWithJSON() {
	cc.runCmd("template", "list", "--json")
}

func (cc *cmdContext) iRunTemplateShow(templateID string) {
	cc.runCmd("template", "show", templateID)
}

func (cc *cmdContext) iRunTemplateCreate(templateID string) {
	// Override editor to no-op during tests
	oldEditor := openEditorFunc
	openEditorFunc = func(filePath string) error { return nil }
	defer func() { openEditorFunc = oldEditor }()

	cc.runCmd("template", "create", templateID)
}

func (cc *cmdContext) theTemplateFileShouldExist(templateID string) error {
	typeName, name, err := parseTemplateArg(templateID)
	if err != nil {
		return err
	}
	path := cc.vault.TemplatePath(typeName, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("expected template file %s to exist", path)
	}
	return nil
}

func (cc *cmdContext) iRunTemplateDeleteWithForce(templateID string) {
	cc.runCmd("template", "delete", templateID, "--force")
}

func (cc *cmdContext) theTemplateFileShouldNotExist(templateID string) error {
	typeName, name, err := parseTemplateArg(templateID)
	if err != nil {
		return err
	}
	path := cc.vault.TemplatePath(typeName, name)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("expected template file %s to not exist", path)
	}
	return nil
}

package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/cucumber/godog"
)

// ── Template steps ───────────────────────────────────────────────────────────

func (dc *domainContext) writeTemplateFile(typeName, templateName, content string) {
	os.MkdirAll(dc.vault.TypeTemplatesDir(typeName), 0755)
	path := dc.vault.TemplatePath(typeName, templateName)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		panic(fmt.Sprintf("write template: %v", err))
	}
}

func (dc *domainContext) aTemplateForTypeWithBody(templateName, typeName, body string) {
	dc.writeTemplateFile(typeName, templateName, body+"\n")
}

func (dc *domainContext) aTemplateForTypeWithFrontmatterAndBody(templateName, typeName, frontmatter, body string) {
	var content string
	if frontmatter != "" && body != "" {
		content = fmt.Sprintf("---\n%s\n---\n\n%s\n", frontmatter, body)
	} else if frontmatter != "" {
		content = fmt.Sprintf("---\n%s\n---\n", frontmatter)
	} else {
		content = body + "\n"
	}
	dc.writeTemplateFile(typeName, templateName, content)
}

func (dc *domainContext) anEmptyTemplatesDirectoryForType(typeName string) {
	dir := dc.vault.TypeTemplatesDir(typeName)
	os.MkdirAll(dir, 0755)
}

func (dc *domainContext) iListTemplatesForType(typeName string) {
	names, err := dc.vault.ListTemplates(typeName)
	dc.lastErr = err
	if err == nil {
		dc.templateNames = names
	}
}

func (dc *domainContext) theTemplateListShouldContain(expected string) error {
	expectedNames := strings.Split(expected, ", ")
	if len(dc.templateNames) != len(expectedNames) {
		return fmt.Errorf("template list has %d items, want %d: %v", len(dc.templateNames), len(expectedNames), dc.templateNames)
	}
	for _, name := range expectedNames {
		found := false
		for _, got := range dc.templateNames {
			if got == strings.TrimSpace(name) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("template list missing %q: %v", name, dc.templateNames)
		}
	}
	return nil
}

func (dc *domainContext) theTemplateListShouldBeEmpty() error {
	if len(dc.templateNames) != 0 {
		return fmt.Errorf("expected empty template list, got %v", dc.templateNames)
	}
	return nil
}

func (dc *domainContext) iLoadTemplateForType(templateName, typeName string) {
	tmpl, err := dc.vault.LoadTemplate(typeName, templateName)
	dc.lastErr = err
	dc.loadedTemplate = tmpl
}

func (dc *domainContext) theTemplatePropertyShouldBe(key, expected string) error {
	if dc.loadedTemplate == nil {
		return fmt.Errorf("no template loaded")
	}
	val := fmt.Sprintf("%v", dc.loadedTemplate.Properties[key])
	if val != expected {
		return fmt.Errorf("template property %q = %q, want %q", key, val, expected)
	}
	return nil
}

func (dc *domainContext) theTemplateBodyShouldBe(expected string) error {
	if dc.loadedTemplate == nil {
		return fmt.Errorf("no template loaded")
	}
	body := strings.TrimRight(dc.loadedTemplate.Body, "\n")
	if body != expected {
		return fmt.Errorf("template body = %q, want %q", body, expected)
	}
	return nil
}

func (dc *domainContext) theTemplateLoadShouldFail() error {
	if dc.lastErr == nil {
		return fmt.Errorf("expected template load to fail, got nil error")
	}
	return nil
}

func (dc *domainContext) iCreateAObjectNamedWithTemplate(typeName, name, templateName string) {
	obj, err := dc.vault.NewObject(typeName, name, templateName)
	dc.lastErr = err
	if err == nil {
		dc.objects[name] = obj
		dc.currentObject = obj
	}
}

func (dc *domainContext) theObjectBodyShouldBe(expected string) error {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	body := strings.TrimRight(got.Body, "\n")
	if body != expected {
		return fmt.Errorf("object body = %q, want %q", body, expected)
	}
	return nil
}

func (dc *domainContext) theObjectCreationShouldFail() error {
	if dc.lastErr == nil {
		return fmt.Errorf("expected object creation to fail, got nil error")
	}
	return nil
}

func initTemplateSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a template "([^"]*)" for type "([^"]*)" with body "([^"]*)"$`, dc.aTemplateForTypeWithBody)
	ctx.Step(`^a template "([^"]*)" for type "([^"]*)" with frontmatter "([^"]*)" and body "([^"]*)"$`, dc.aTemplateForTypeWithFrontmatterAndBody)
	ctx.Step(`^an empty templates directory for type "([^"]*)"$`, dc.anEmptyTemplatesDirectoryForType)
	ctx.Step(`^I list templates for type "([^"]*)"$`, dc.iListTemplatesForType)
	ctx.Step(`^the template list should contain "([^"]*)"$`, dc.theTemplateListShouldContain)
	ctx.Step(`^the template list should be empty$`, dc.theTemplateListShouldBeEmpty)
	ctx.Step(`^I load template "([^"]*)" for type "([^"]*)"$`, dc.iLoadTemplateForType)
	ctx.Step(`^the template property "([^"]*)" should be "([^"]*)"$`, dc.theTemplatePropertyShouldBe)
	ctx.Step(`^the template body should be "([^"]*)"$`, dc.theTemplateBodyShouldBe)
	ctx.Step(`^the template load should fail$`, dc.theTemplateLoadShouldFail)
	ctx.Step(`^I create a "([^"]*)" object named "([^"]*)" with template "([^"]*)"$`, dc.iCreateAObjectNamedWithTemplate)
	ctx.Step(`^the object body should be "([^"]*)"$`, dc.theObjectBodyShouldBe)
	ctx.Step(`^the object creation should fail$`, dc.theObjectCreationShouldFail)
}

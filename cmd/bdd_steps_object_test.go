package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

func initObjectSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	// Object create steps
	ctx.Step(`^I run object create "([^"]*)" "([^"]*)"$`, cc.iRunObjectCreate)
	ctx.Step(`^I run object create "([^"]*)"$`, cc.iRunObjectCreateNameOnly)
	ctx.Step(`^I run object create with no arguments$`, cc.iRunObjectCreateNoArgs)
	ctx.Step(`^I run object create "([^"]*)" "([^"]*)" with template "([^"]*)"$`, cc.iRunObjectCreateWithTemplate)
	ctx.Step(`^a vault with no default type$`, cc.aVaultWithNoDefaultType)
	ctx.Step(`^a vault with a "([^"]*)" template for "([^"]*)"$`, cc.aVaultWithTemplate)

	// Object show steps
	ctx.Step(`^a vault with a book "([^"]*)"$`, cc.aVaultWithABook)
	ctx.Step(`^a vault with a book "([^"]*)" with local property "([^"]*)" = "([^"]*)"$`, cc.aVaultWithABookWithLocalProperty)
	ctx.Step(`^I run object show for the created book$`, cc.iRunObjectShowForCreatedBook)
	ctx.Step(`^I run object show "([^"]*)"$`, cc.iRunObjectShow)

	// Object list steps
	ctx.Step(`^a vault with (\d+) books$`, cc.aVaultWithNBooks)
	ctx.Step(`^I run object list$`, cc.iRunObjectList)
	ctx.Step(`^I run object list with json$`, cc.iRunObjectListJSON)
	ctx.Step(`^the output should have (\d+) lines containing "([^"]*)"$`, cc.theOutputShouldHaveNLinesContaining)
}

// --- Object Create ---

func (cc *cmdContext) iRunObjectCreate(typeName, name string) {
	cc.runCmdAndTrack("object", "create", typeName, name)
}

func (cc *cmdContext) iRunObjectCreateNameOnly(name string) {
	cc.runCmdAndTrack("object", "create", name)
}

func (cc *cmdContext) iRunObjectCreateNoArgs() {
	cc.runCmd("object", "create")
}

func (cc *cmdContext) iRunObjectCreateWithTemplate(typeName, name, tmpl string) {
	cc.runCmdAndTrack("object", "create", typeName, name, "-t", tmpl)
}

func (cc *cmdContext) aVaultWithNoDefaultType() error {
	if err := cc.setupVault(); err != nil {
		return err
	}
	// Remove the config file so there's no default type.
	configPath := filepath.Join(cc.vaultDir, ".typemd", "config.yaml")
	return os.Remove(configPath)
}

// aVaultWithTemplate adds a template file to an already-initialized vault.
// It assumes "a vault is ready" was called first.
func (cc *cmdContext) aVaultWithTemplate(tmplName, typeName string) error {
	tmplDir := filepath.Join(cc.vaultDir, "templates", typeName)
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		return err
	}
	content := fmt.Sprintf(`---
author: "Template Author"
---

This is the %s template body.
`, tmplName)
	return os.WriteFile(filepath.Join(tmplDir, tmplName+".md"), []byte(content), 0644)
}

// --- Object Show ---

func (cc *cmdContext) aVaultWithABook(name string) error {
	if err := cc.setupVault(); err != nil {
		return err
	}
	obj, err := cc.vault.NewObject("book", name, "")
	if err != nil {
		return err
	}
	cc.createdObjectIDs = append(cc.createdObjectIDs, string(obj.ID))
	return nil
}

func (cc *cmdContext) aVaultWithABookWithLocalProperty(name, propKey, propValue string) error {
	if err := cc.aVaultWithABook(name); err != nil {
		return err
	}
	id := cc.createdObjectIDs[len(cc.createdObjectIDs)-1]
	obj, err := cc.vault.GetObject(id)
	if err != nil {
		return err
	}
	obj.Properties[propKey] = propValue
	return cc.vault.SaveObject(obj)
}

func (cc *cmdContext) iRunObjectShowForCreatedBook() {
	if len(cc.createdObjectIDs) == 0 {
		cc.lastErr = fmt.Errorf("no created objects to show")
		return
	}
	id := cc.createdObjectIDs[len(cc.createdObjectIDs)-1]
	cc.runCmd("object", "show", id)
}

func (cc *cmdContext) iRunObjectShow(objectID string) {
	cc.runCmd("object", "show", objectID)
}

// --- Object List ---

func (cc *cmdContext) aVaultWithNBooks(n int) error {
	if err := cc.setupVault(); err != nil {
		return err
	}
	for i := 1; i <= n; i++ {
		name := fmt.Sprintf("book-%d", i)
		if _, err := cc.vault.NewObject("book", name, ""); err != nil {
			return err
		}
	}
	return nil
}

func (cc *cmdContext) iRunObjectList() {
	cc.runCmd("object", "list")
}

func (cc *cmdContext) iRunObjectListJSON() {
	cc.runCmd("object", "list", "--json")
}

func (cc *cmdContext) theOutputShouldHaveNLinesContaining(n int, substr string) error {
	lines := strings.Split(strings.TrimSpace(cc.stdout), "\n")
	count := 0
	for _, line := range lines {
		if strings.Contains(line, substr) {
			count++
		}
	}
	if count != n {
		return fmt.Errorf("expected %d lines containing %q, got %d in output:\n%s", n, substr, count, cc.stdout)
	}
	return nil
}


package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

func (dc *domainContext) aConfigFileWithContent(content *godog.DocString) {
	path := filepath.Join(dc.vault.Dir(), configFileName)
	os.WriteFile(path, []byte(content.Content), 0644)
}

func (dc *domainContext) theDefaultTypeShouldBe(expected string) error {
	got := dc.vault.DefaultType()
	if got != expected {
		return fmt.Errorf("expected default type %q, got %q", expected, got)
	}
	return nil
}

func initVaultConfigSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a config file with content:$`, dc.aConfigFileWithContent)
	ctx.Step(`^the default type should be "([^"]*)"$`, dc.theDefaultTypeShouldBe)
}

package cmd

import (
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

func initFormatSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	ctx.Step(`^I run format$`, cc.iRunFormat)
	ctx.Step(`^I run format with type "([^"]*)"$`, cc.iRunFormatWithType)
	ctx.Step(`^I run format with dry-run$`, cc.iRunFormatDryRun)
	ctx.Step(`^an object with out-of-order frontmatter$`, cc.anObjectWithOutOfOrderFrontmatter)
	ctx.Step(`^the vault formatting is stable$`, cc.theVaultFormattingIsStable)
}

func (cc *cmdContext) iRunFormat() {
	cc.runCmd("format")
}

func (cc *cmdContext) iRunFormatWithType(typeName string) {
	cc.runCmd("format", "--type", typeName)
}

func (cc *cmdContext) iRunFormatDryRun() {
	cc.runCmd("format", "--dry-run")
}

func (cc *cmdContext) theVaultFormattingIsStable() {
	cc.runCmd("format")
	cc.runCmd("format")
}

func (cc *cmdContext) anObjectWithOutOfOrderFrontmatter() error {
	dir := filepath.Join(cc.vaultDir, "objects", "book")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	// Write frontmatter with author before name (out of canonical order).
	content := "---\nauthor: Alice\nname: Test Book\ncreated_at: \"2025-01-01T00:00:00Z\"\nupdated_at: \"2025-01-01T00:00:00Z\"\n---\n"
	return os.WriteFile(filepath.Join(dir, "test-book.md"), []byte(content), 0644)
}

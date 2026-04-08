package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/typemd/typemd/core"
)

func initImportSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	ctx.Step(`^a source directory with markdown files$`, cc.aSourceDirectoryWithMarkdownFiles)
	ctx.Step(`^I run import scan on the source directory$`, cc.iRunImportScanOnSourceDir)
	ctx.Step(`^I run import scan on a nonexistent path$`, cc.iRunImportScanNonexistent)
	ctx.Step(`^a plan file with a page object$`, cc.aPlanFileWithPageObject)
	ctx.Step(`^I run import execute with the plan file$`, cc.iRunImportExecuteWithPlanFile)
	ctx.Step(`^I run import execute "([^"]*)"$`, cc.iRunImportExecute)
	ctx.Step(`^a classifications file with a page$`, cc.aClassificationsFileWithPage)
	ctx.Step(`^I run import plan with the classifications file$`, cc.iRunImportPlanWithClassificationsFile)
}

func (cc *cmdContext) aSourceDirectoryWithMarkdownFiles() error {
	srcDir := filepath.Join(cc.vaultDir, "_import_src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		return err
	}
	content := "---\ntitle: Test Note\ntags: [go, cli]\n---\n\nHello world.\n"
	return os.WriteFile(filepath.Join(srcDir, "test-note.md"), []byte(content), 0644)
}

func (cc *cmdContext) iRunImportScanOnSourceDir() {
	srcDir := filepath.Join(cc.vaultDir, "_import_src")
	cc.runCmd("import", "scan", srcDir)
}

func (cc *cmdContext) iRunImportScanNonexistent() {
	cc.runCmd("import", "scan", filepath.Join(cc.vaultDir, "does-not-exist"))
}

func (cc *cmdContext) aPlanFileWithPageObject() error {
	plan := core.ImportPlan{
		Objects: []core.ObjectPlan{
			{
				SourcePath: "test.md",
				TypeName:   "page",
				Name:       "Test Page",
				Conflict:   "none",
			},
		},
		Order: []int{0},
	}
	data, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(cc.vaultDir, "plan.json"), data, 0644)
}

func (cc *cmdContext) iRunImportExecuteWithPlanFile() {
	cc.runCmd("import", "execute", filepath.Join(cc.vaultDir, "plan.json"))
}

func (cc *cmdContext) iRunImportExecute(file string) {
	cc.runCmd("import", "execute", file)
}

func (cc *cmdContext) aClassificationsFileWithPage() error {
	classifications := []core.ObjectPlan{
		{
			SourcePath: "note.md",
			TypeName:   "page",
			Name:       "My Note",
			Conflict:   "none",
		},
	}
	data, err := json.MarshalIndent(classifications, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(cc.vaultDir, "classifications.json"), data, 0644)
}

func (cc *cmdContext) iRunImportPlanWithClassificationsFile() {
	cc.runCmd("import", "plan", filepath.Join(cc.vaultDir, "classifications.json"))
}

package cmd

import (
	"github.com/cucumber/godog"
)

func initFormatSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	ctx.Step(`^I run format$`, cc.iRunFormat)
	ctx.Step(`^I run format with type "([^"]*)"$`, cc.iRunFormatWithType)
	ctx.Step(`^I run format with dry-run$`, cc.iRunFormatDryRun)
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

package cmd

import (
	"github.com/cucumber/godog"
)

func initInstructionsSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	ctx.Step(`^I run instructions with no arguments$`, cc.iRunInstructionsNoArgs)
	ctx.Step(`^I run instructions with json flag$`, cc.iRunInstructionsJSON)
	ctx.Step(`^I run instructions "([^"]*)"$`, cc.iRunInstructions)
	ctx.Step(`^I run instructions "([^"]*)" with skill flag$`, cc.iRunInstructionsSkill)
}

func (cc *cmdContext) iRunInstructionsNoArgs() {
	cc.runCmd("instructions")
}

func (cc *cmdContext) iRunInstructionsJSON() {
	cc.runCmd("instructions", "--json")
}

func (cc *cmdContext) iRunInstructions(name string) {
	cc.runCmd("instructions", name)
}

func (cc *cmdContext) iRunInstructionsSkill(name string) {
	cc.runCmd("instructions", name, "--skill")
}

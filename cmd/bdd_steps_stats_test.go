package cmd

import (
	"github.com/cucumber/godog"
)

func initStatsSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	// "a book object X exists" and "the output should start with" are registered in initCommonSteps
	ctx.Step(`^I run stats$`, cc.iRunStats)
	ctx.Step(`^I run stats with type "([^"]*)"$`, cc.iRunStatsWithType)
	ctx.Step(`^I run stats with json$`, cc.iRunStatsWithJSON)
}

func (cc *cmdContext) iRunStats() {
	cc.runCmd("stats")
}

func (cc *cmdContext) iRunStatsWithType(typeName string) {
	cc.runCmd("stats", "--type", typeName)
}

func (cc *cmdContext) iRunStatsWithJSON() {
	cc.runCmd("stats", "--json")
}

package cmd

import (
	"github.com/cucumber/godog"
)

func initSearchSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	ctx.Step(`^I run search "([^"]*)"$`, cc.iRunSearch)
	ctx.Step(`^I run search "([^"]*)" with json$`, cc.iRunSearchJSON)
	ctx.Step(`^I run search with no arguments$`, cc.iRunSearchNoArgs)
}

func (cc *cmdContext) iRunSearch(keyword string) {
	cc.runCmd("search", keyword)
}

func (cc *cmdContext) iRunSearchJSON(keyword string) {
	cc.runCmd("search", keyword, "--json")
}

func (cc *cmdContext) iRunSearchNoArgs() {
	cc.runCmd("search")
}

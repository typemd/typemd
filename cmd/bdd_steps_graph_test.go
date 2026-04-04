package cmd

import (
	"github.com/cucumber/godog"
)

func initGraphSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	ctx.Step(`^I run graph$`, cc.iRunGraph)
	ctx.Step(`^I run graph with type "([^"]*)"$`, cc.iRunGraphWithType)
	ctx.Step(`^I run graph with no-relations$`, cc.iRunGraphWithNoRelations)
	ctx.Step(`^I run graph with no-wikilinks$`, cc.iRunGraphWithNoWikiLinks)
}

func (cc *cmdContext) iRunGraph() {
	cc.runCmd("graph")
}

func (cc *cmdContext) iRunGraphWithType(typeName string) {
	cc.runCmd("graph", "--type", typeName)
}

func (cc *cmdContext) iRunGraphWithNoRelations() {
	cc.runCmd("graph", "--no-relations")
}

func (cc *cmdContext) iRunGraphWithNoWikiLinks() {
	cc.runCmd("graph", "--no-wikilinks")
}

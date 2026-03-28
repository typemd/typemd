package cmd

import (
	"strings"

	"github.com/cucumber/godog"
)

func initDebugSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	ctx.Step(`^I run command "([^"]*)"$`, func(command string) {
		args := strings.Fields(command)
		cc.runCmd(args...)
	})
}

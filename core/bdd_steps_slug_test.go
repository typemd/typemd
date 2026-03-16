package core

import (
	"fmt"

	"github.com/cucumber/godog"
)

type slugContext struct {
	result string
}

func (sc *slugContext) iSlugify(input string) {
	sc.result = Slugify(input)
}

func (sc *slugContext) theSlugShouldBe(expected string) error {
	if sc.result != expected {
		return fmt.Errorf("expected slug %q, got %q", expected, sc.result)
	}
	return nil
}

func initSlugSteps(ctx *godog.ScenarioContext, _ *domainContext) {
	sc := &slugContext{}
	ctx.Step(`^I slugify "([^"]*)"$`, sc.iSlugify)
	ctx.Step(`^the slug should be "([^"]*)"$`, sc.theSlugShouldBe)
}

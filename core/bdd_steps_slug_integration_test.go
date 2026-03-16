package core

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

func (dc *domainContext) theCurrentObjectIDShouldContain(substr string) error {
	if dc.currentObject == nil {
		return fmt.Errorf("no current object")
	}
	if !strings.Contains(dc.currentObject.ID, substr) {
		return fmt.Errorf("expected object ID %q to contain %q", dc.currentObject.ID, substr)
	}
	return nil
}

func (dc *domainContext) theCurrentObjectNamePropertyShouldBe(expected string) error {
	if dc.currentObject == nil {
		return fmt.Errorf("no current object")
	}
	name, ok := dc.currentObject.Properties[NameProperty]
	if !ok {
		return fmt.Errorf("name property not found")
	}
	got := fmt.Sprintf("%v", name)
	if got != expected {
		return fmt.Errorf("expected name property %q, got %q", expected, got)
	}
	return nil
}

func initSlugIntegrationSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^the current object ID should contain "([^"]*)"$`, dc.theCurrentObjectIDShouldContain)
	ctx.Step(`^the current object name property should be "([^"]*)"$`, dc.theCurrentObjectNamePropertyShouldBe)
}

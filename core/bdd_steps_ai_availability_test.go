package core

import (
	"fmt"

	"github.com/cucumber/godog"
)

func (dc *domainContext) aiServiceShouldBeAvailable() error {
	if dc.vault.AIService() == nil {
		return fmt.Errorf("expected AI service to be available, but it is nil")
	}
	return nil
}

func (dc *domainContext) aiServiceShouldNotBeAvailable() error {
	if dc.vault.AIService() != nil {
		return fmt.Errorf("expected AI service to be nil, but it is available")
	}
	return nil
}

func initAIAvailabilitySteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^AI service should be available$`, dc.aiServiceShouldBeAvailable)
	ctx.Step(`^AI service should not be available$`, dc.aiServiceShouldNotBeAvailable)
}

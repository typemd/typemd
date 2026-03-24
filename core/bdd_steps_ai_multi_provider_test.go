package core

import (
	"fmt"

	"github.com/cucumber/godog"
)

func (dc *domainContext) theAIDefaultProviderShouldBe(expected string) error {
	got := dc.vault.Config().AI.Default
	if got != expected {
		return fmt.Errorf("expected AI default provider %q, got %q", expected, got)
	}
	return nil
}

func (dc *domainContext) theAIProvidersMapShouldBeEmpty() error {
	providers := dc.vault.Config().AI.Providers
	if len(providers) != 0 {
		return fmt.Errorf("expected AI providers to be empty, got %d entries", len(providers))
	}
	return nil
}

func (dc *domainContext) theAIProviderShouldHaveType(name, expected string) error {
	p, ok := dc.vault.Config().AI.Providers[name]
	if !ok {
		return fmt.Errorf("provider %q not found in providers map", name)
	}
	if p.Type != expected {
		return fmt.Errorf("expected provider %q type %q, got %q", name, expected, p.Type)
	}
	return nil
}

func (dc *domainContext) theAIProviderShouldHaveModel(name, expected string) error {
	p, ok := dc.vault.Config().AI.Providers[name]
	if !ok {
		return fmt.Errorf("provider %q not found in providers map", name)
	}
	if p.Model != expected {
		return fmt.Errorf("expected provider %q model %q, got %q", name, expected, p.Model)
	}
	return nil
}

func (dc *domainContext) theAIProviderShouldHaveBaseURL(name, expected string) error {
	p, ok := dc.vault.Config().AI.Providers[name]
	if !ok {
		return fmt.Errorf("provider %q not found in providers map", name)
	}
	if p.BaseURL != expected {
		return fmt.Errorf("expected provider %q base_url %q, got %q", name, expected, p.BaseURL)
	}
	return nil
}

func (dc *domainContext) theAIProviderShouldHaveAPIKey(name, expected string) error {
	p, ok := dc.vault.Config().AI.Providers[name]
	if !ok {
		return fmt.Errorf("provider %q not found in providers map", name)
	}
	if p.APIKey != expected {
		return fmt.Errorf("expected provider %q api_key %q, got %q", name, expected, p.APIKey)
	}
	return nil
}

func initAIMultiProviderSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^the AI default provider should be "([^"]*)"$`, dc.theAIDefaultProviderShouldBe)
	ctx.Step(`^the AI providers map should be empty$`, dc.theAIProvidersMapShouldBeEmpty)
	ctx.Step(`^the AI provider "([^"]*)" should have type "([^"]*)"$`, dc.theAIProviderShouldHaveType)
	ctx.Step(`^the AI provider "([^"]*)" should have model "([^"]*)"$`, dc.theAIProviderShouldHaveModel)
	ctx.Step(`^the AI provider "([^"]*)" should have base_url "([^"]*)"$`, dc.theAIProviderShouldHaveBaseURL)
	ctx.Step(`^the AI provider "([^"]*)" should have api_key "([^"]*)"$`, dc.theAIProviderShouldHaveAPIKey)
}

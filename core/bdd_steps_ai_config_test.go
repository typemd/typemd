package core

import (
	"fmt"

	"github.com/cucumber/godog"
)

func (dc *domainContext) aiShouldBeEnabled() error {
	cfg := dc.vault.Config()
	if !cfg.AI.Enabled {
		return fmt.Errorf("expected AI to be enabled, but it is not")
	}
	return nil
}

func (dc *domainContext) aiShouldNotBeEnabled() error {
	cfg := dc.vault.Config()
	if cfg.AI.Enabled {
		return fmt.Errorf("expected AI to not be enabled, but it is")
	}
	return nil
}

func (dc *domainContext) theAIDescribePromptShouldBe(expected string) error {
	got := dc.vault.Config().AI.Prompts.Describe
	if got != expected {
		return fmt.Errorf("expected AI describe prompt %q, got %q", expected, got)
	}
	return nil
}

func (dc *domainContext) theAITagPromptShouldBe(expected string) error {
	got := dc.vault.Config().AI.Prompts.Tag
	if got != expected {
		return fmt.Errorf("expected AI tag prompt %q, got %q", expected, got)
	}
	return nil
}

func (dc *domainContext) theAIExplorePromptShouldBe(expected string) error {
	got := dc.vault.Config().AI.Prompts.Explore
	if got != expected {
		return fmt.Errorf("expected AI explore prompt %q, got %q", expected, got)
	}
	return nil
}

func (dc *domainContext) theAIExploreSampleCountShouldBe(expected int) error {
	got := dc.vault.Config().AI.Explore.SampleCount
	if got != expected {
		return fmt.Errorf("expected AI explore sample_count %d, got %d", expected, got)
	}
	return nil
}

func (dc *domainContext) theAIExploreBodyTruncateShouldBe(expected int) error {
	got := dc.vault.Config().AI.Explore.BodyTruncate
	if got != expected {
		return fmt.Errorf("expected AI explore body_truncate %d, got %d", expected, got)
	}
	return nil
}

func initAIConfigSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^AI should be enabled$`, dc.aiShouldBeEnabled)
	ctx.Step(`^AI should not be enabled$`, dc.aiShouldNotBeEnabled)
	ctx.Step(`^the AI describe prompt should be "([^"]*)"$`, dc.theAIDescribePromptShouldBe)
	ctx.Step(`^the AI tag prompt should be "([^"]*)"$`, dc.theAITagPromptShouldBe)
	ctx.Step(`^the AI explore prompt should be "([^"]*)"$`, dc.theAIExplorePromptShouldBe)
	ctx.Step(`^the AI explore sample count should be (\d+)$`, dc.theAIExploreSampleCountShouldBe)
	ctx.Step(`^the AI explore body truncate should be (\d+)$`, dc.theAIExploreBodyTruncateShouldBe)
}

package core

import (
	"fmt"

	"github.com/cucumber/godog"
)

func (dc *domainContext) everyConfigKeyShouldHaveADescription() error {
	infos := dc.vault.ConfigKeysInfo()
	for _, info := range infos {
		if info.Description == "" {
			return fmt.Errorf("config key %q has no description", info.Key)
		}
	}
	return nil
}

func (dc *domainContext) theConfigKeyInfoShouldHaveValue(key, expected string) error {
	infos := dc.vault.ConfigKeysInfo()
	for _, info := range infos {
		if info.Key == key {
			if info.Value != expected {
				return fmt.Errorf("expected config key %q value %q, got %q", key, expected, info.Value)
			}
			return nil
		}
	}
	return fmt.Errorf("config key %q not found in ConfigKeysInfo", key)
}

func (dc *domainContext) theConfigKeyInfoShouldHaveNonEmptyDefault(key string) error {
	infos := dc.vault.ConfigKeysInfo()
	for _, info := range infos {
		if info.Key == key {
			if info.Default == "" {
				return fmt.Errorf("config key %q has empty default", key)
			}
			return nil
		}
	}
	return fmt.Errorf("config key %q not found in ConfigKeysInfo", key)
}

func initConfigKeyInfoSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^every config key should have a description$`, dc.everyConfigKeyShouldHaveADescription)
	ctx.Step(`^the config key info for "([^"]*)" should have value "([^"]*)"$`, dc.theConfigKeyInfoShouldHaveValue)
	ctx.Step(`^the config key info for "([^"]*)" should have a non-empty default$`, dc.theConfigKeyInfoShouldHaveNonEmptyDefault)
}

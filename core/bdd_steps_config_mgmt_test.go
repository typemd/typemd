package core

import (
	"fmt"
	"os"
	"sort"

	"github.com/cucumber/godog"
)

func (dc *domainContext) iSetConfigTo(key, value string) {
	dc.lastErr = dc.vault.SetConfigValue(key, value)
}

func (dc *domainContext) theConfigValueShouldBe(key, expected string) error {
	val, err := dc.vault.GetConfigValue(key)
	if err != nil {
		return fmt.Errorf("config key %q: %w", key, err)
	}
	if val != expected {
		return fmt.Errorf("expected config %q = %q, got %q", key, expected, val)
	}
	return nil
}

func (dc *domainContext) theConfigKeyShouldNotBeKnown(key string) error {
	_, err := dc.vault.GetConfigValue(key)
	if err == nil {
		return fmt.Errorf("expected config key %q to return error, but it didn't", key)
	}
	return nil
}

func (dc *domainContext) theKnownConfigKeysShouldInclude(key string) error {
	keys := ConfigKeys()
	sort.Strings(keys)
	for _, k := range keys {
		if k == key {
			return nil
		}
	}
	return fmt.Errorf("expected known keys to include %q, got %v", key, keys)
}

func (dc *domainContext) theConfigFileShouldExist() error {
	path := dc.vault.ConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist at %s", path)
	}
	return nil
}

func initConfigMgmtSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^I set config "([^"]*)" to "([^"]*)"$`, dc.iSetConfigTo)
	ctx.Step(`^the config value "([^"]*)" should be "([^"]*)"$`, dc.theConfigValueShouldBe)
	ctx.Step(`^the config key "([^"]*)" should not be known$`, dc.theConfigKeyShouldNotBeKnown)
	ctx.Step(`^the known config keys should include "([^"]*)"$`, dc.theKnownConfigKeysShouldInclude)
	ctx.Step(`^the config file should exist$`, dc.theConfigFileShouldExist)
}

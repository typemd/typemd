package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
	"gopkg.in/yaml.v3"
)

// ── Starter type steps ──────────────────────────────────────────────────────

func (dc *domainContext) iListAvailableStarterTypes() {
	dc.starterTypes = StarterTypes()
}

func (dc *domainContext) iShouldGetNStarterTypes(expected int) error {
	if len(dc.starterTypes) != expected {
		return fmt.Errorf("expected %d starter types, got %d", expected, len(dc.starterTypes))
	}
	return nil
}

func (dc *domainContext) theStarterTypesShouldInclude(names string) error {
	nameList := strings.Split(names, ", ")
	have := make(map[string]bool)
	for _, st := range dc.starterTypes {
		have[st.Name] = true
	}
	for _, name := range nameList {
		name = strings.Trim(name, "\"")
		if !have[name] {
			return fmt.Errorf("starter type %q not found, have: %v", name, have)
		}
	}
	return nil
}

func (dc *domainContext) eachStarterTypeShouldHaveNameEmojiDescription() error {
	for _, st := range dc.starterTypes {
		if st.Name == "" {
			return fmt.Errorf("starter type has empty name")
		}
		if st.Emoji == "" {
			return fmt.Errorf("starter type %q has empty emoji", st.Name)
		}
		if st.Description == "" {
			return fmt.Errorf("starter type %q has empty description", st.Name)
		}
	}
	return nil
}

func (dc *domainContext) eachStarterTypeYAMLShouldParseAsValidTypeSchema() error {
	for _, st := range dc.starterTypes {
		var schema TypeSchema
		if err := yaml.Unmarshal(st.YAML, &schema); err != nil {
			return fmt.Errorf("starter type %q has invalid YAML: %v", st.Name, err)
		}
		if errs := ValidateSchema(&schema); len(errs) > 0 {
			return fmt.Errorf("starter type %q failed validation: %v", st.Name, errs)
		}
	}
	return nil
}

func (dc *domainContext) iWriteStarterTypesToTheVault(names string) {
	var nameList []string
	if names != "" {
		nameList = strings.Split(names, ",")
	}
	dc.lastErr = dc.vault.WriteStarterTypes(nameList)
}

func (dc *domainContext) iWriteAllStarterTypesToTheVault() {
	starters := StarterTypes()
	names := make([]string, len(starters))
	for i, st := range starters {
		names[i] = st.Name
	}
	dc.lastErr = dc.vault.WriteStarterTypes(names)
}

func (dc *domainContext) theFileShouldExist(relPath string) error {
	fullPath := filepath.Join(dc.rootDir, relPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("expected file %s to exist", relPath)
	}
	return nil
}

func (dc *domainContext) theFileShouldNotExist(relPath string) error {
	fullPath := filepath.Join(dc.rootDir, relPath)
	if _, err := os.Stat(fullPath); err == nil {
		return fmt.Errorf("expected file %s to not exist", relPath)
	}
	return nil
}

func (dc *domainContext) iShouldBeAbleToLoadType(name string) error {
	_, err := dc.vault.LoadType(name)
	if err != nil {
		return fmt.Errorf("failed to load type %q: %v", name, err)
	}
	return nil
}

func initStarterSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^I list available starter types$`, dc.iListAvailableStarterTypes)
	ctx.Step(`^I should get (\d+) starter types$`, dc.iShouldGetNStarterTypes)
	ctx.Step(`^the starter types should include (.+)$`, dc.theStarterTypesShouldInclude)
	ctx.Step(`^each starter type should have a name, emoji, and description$`, dc.eachStarterTypeShouldHaveNameEmojiDescription)
	ctx.Step(`^each starter type YAML should parse as a valid TypeSchema$`, dc.eachStarterTypeYAMLShouldParseAsValidTypeSchema)
	ctx.Step(`^I write starter types "([^"]*)" to the vault$`, dc.iWriteStarterTypesToTheVault)
	ctx.Step(`^I write all starter types to the vault$`, dc.iWriteAllStarterTypesToTheVault)
	ctx.Step(`^the file "([^"]*)" should exist$`, dc.theFileShouldExist)
	ctx.Step(`^the file "([^"]*)" should not exist$`, dc.theFileShouldNotExist)
	ctx.Step(`^I should be able to load type "([^"]*)"$`, dc.iShouldBeAbleToLoadType)
}

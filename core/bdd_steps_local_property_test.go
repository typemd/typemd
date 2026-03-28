package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

func (dc *domainContext) theDisplayPropertyShouldBeLocal(key string) error {
	for _, dp := range dc.displayProps {
		if dp.Key == key {
			if !dp.IsLocal {
				return fmt.Errorf("display property %q has IsLocal=false, expected true", key)
			}
			return nil
		}
	}
	return fmt.Errorf("display property %q not found", key)
}

func (dc *domainContext) theDisplayPropertyShouldNotBeLocal(key string) error {
	for _, dp := range dc.displayProps {
		if dp.Key == key {
			if dp.IsLocal {
				return fmt.Errorf("display property %q has IsLocal=true, expected false", key)
			}
			return nil
		}
	}
	return fmt.Errorf("display property %q not found", key)
}

func (dc *domainContext) noDisplayPropertyShouldBeLocal() error {
	for _, dp := range dc.displayProps {
		if dp.IsLocal {
			return fmt.Errorf("display property %q has IsLocal=true, expected none to be local", dp.Key)
		}
	}
	return nil
}

func (dc *domainContext) theTypeSchemaIsRemoved(typeName string) {
	dirPath := filepath.Join(dc.vault.TypesDir(), typeName)
	os.RemoveAll(dirPath)
	dc.vault.InvalidateSchemaCache()
}

func initLocalPropertySteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a "([^"]*)" object named "([^"]*)" exists with extra property "([^"]*)" set to "([^"]*)"$`, dc.aObjectNamedExistsWithPropertySetTo)
	ctx.Step(`^the type schema "([^"]*)" is removed$`, dc.theTypeSchemaIsRemoved)
	ctx.Step(`^the display property "([^"]*)" should be local$`, dc.theDisplayPropertyShouldBeLocal)
	ctx.Step(`^the display property "([^"]*)" should not be local$`, dc.theDisplayPropertyShouldNotBeLocal)
	ctx.Step(`^no display property should be local$`, dc.noDisplayPropertyShouldBeLocal)
}

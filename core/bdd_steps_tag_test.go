package core

import (
	"fmt"

	"github.com/cucumber/godog"
)

// ── Tag / tags property steps ───────────────────────────────────────────

func (dc *domainContext) theObjectShouldHavePropertyWithNilValue(propName string) error {
	val, ok := dc.currentObject.Properties[propName]
	if !ok {
		return fmt.Errorf("property %q not found", propName)
	}
	if val != nil {
		return fmt.Errorf("expected %q to be nil, got %v", propName, val)
	}
	return nil
}

func (dc *domainContext) iSetTagsOnTheObjectToATagReference() {
	tagObj, err := dc.vault.NewObject("tag", "test-tag", "")
	if err != nil {
		panic(fmt.Sprintf("create tag object failed: %v", err))
	}
	dc.currentObject.Properties[TagsProperty] = []any{tagObj.ID}
	if err := dc.vault.SaveObject(dc.currentObject); err != nil {
		panic(fmt.Sprintf("SaveObject failed: %v", err))
	}
}

// ── Tag uniqueness steps ────────────────────────────────────────────────

func (dc *domainContext) aRawDuplicateTagNamedExists(name string) {
	dc.aRawDuplicateObjectOfTypeNamedExists(TagTypeName, name)
}

func initTagSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	// Tag / tags property steps
	ctx.Step(`^the object should have property "([^"]*)" with nil value$`, dc.theObjectShouldHavePropertyWithNilValue)
	ctx.Step(`^I set tags on the object to a tag reference$`, dc.iSetTagsOnTheObjectToATagReference)

	// Tag uniqueness steps
	ctx.Step(`^a raw duplicate tag named "([^"]*)" exists$`, dc.aRawDuplicateTagNamedExists)
}

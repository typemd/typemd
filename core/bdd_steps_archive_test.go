package core

import (
	"fmt"

	"github.com/cucumber/godog"
)

// ── Object archive steps ────────────────────────────────────────────────

func (dc *domainContext) iArchiveTheObject() {
	dc.lastErr = dc.vault.SetArchived(dc.currentObject.ID, true)
	if dc.lastErr == nil {
		obj, err := dc.vault.GetObject(dc.currentObject.ID)
		if err == nil {
			dc.currentObject = obj
		}
	}
}

func (dc *domainContext) iUnarchiveTheObject() {
	dc.lastErr = dc.vault.SetArchived(dc.currentObject.ID, false)
	if dc.lastErr == nil {
		obj, err := dc.vault.GetObject(dc.currentObject.ID)
		if err == nil {
			dc.currentObject = obj
		}
	}
}

func (dc *domainContext) iArchiveTheNamedObject(name string) {
	for _, obj := range dc.objects {
		if obj.GetName() == name {
			dc.lastErr = dc.vault.SetArchived(obj.ID, true)
			return
		}
	}
	dc.lastErr = fmt.Errorf("object named %q not found", name)
}

func (dc *domainContext) theObjectShouldBeArchived() error {
	obj, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	if !obj.IsArchived() {
		return fmt.Errorf("expected object to be archived, but it is not")
	}
	return nil
}

func (dc *domainContext) theObjectShouldNotBeArchived() error {
	obj, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	if obj.IsArchived() {
		return fmt.Errorf("expected object to not be archived, but it is")
	}
	return nil
}

func (dc *domainContext) iQueryObjectsWithFilterIncludingArchived(filter string) {
	results, err := dc.vault.QueryObjects(parseTestFilter(filter), QueryIncludeArchived())
	dc.lastErr = err
	dc.queryResults = results
}

func (dc *domainContext) iGetTheObjectByID() {
	obj, err := dc.vault.GetObject(dc.currentObject.ID)
	dc.lastErr = err
	if err == nil {
		dc.currentObject = obj
	}
}

func initArchiveSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^I archive the object$`, dc.iArchiveTheObject)
	ctx.Step(`^I unarchive the object$`, dc.iUnarchiveTheObject)
	ctx.Step(`^I archive the "([^"]*)" object$`, dc.iArchiveTheNamedObject)
	ctx.Step(`^the object should be archived$`, dc.theObjectShouldBeArchived)
	ctx.Step(`^the object should not be archived$`, dc.theObjectShouldNotBeArchived)
	ctx.Step(`^I query objects with filter "([^"]*)" including archived$`, dc.iQueryObjectsWithFilterIncludingArchived)
	ctx.Step(`^I get the object by ID$`, dc.iGetTheObjectByID)
}

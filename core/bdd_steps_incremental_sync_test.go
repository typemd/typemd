package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

type incrementalSyncContext struct {
	dc          *domainContext
	objectPath  string
	objectID    string
	syncResult  *ReconcileResult
}

func newIncrementalSyncContext(dc *domainContext) *incrementalSyncContext {
	return &incrementalSyncContext{dc: dc}
}

func (isc *incrementalSyncContext) aTypeSchemaExists(typeName string) {
	data := fmt.Sprintf("name: %s\nproperties: []\n", typeName)
	dir := filepath.Join(isc.dc.vault.TypesDir(), typeName)
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "schema.yaml"), []byte(data), 0644)
}

func (isc *incrementalSyncContext) iCreateAnObjectNamed(typeName, name string) {
	obj, err := isc.dc.vault.Objects.Create(typeName, name, "")
	if err != nil {
		panic(fmt.Sprintf("create object failed: %v", err))
	}
	isc.objectID = obj.ID
	isc.objectPath = isc.dc.vault.ObjectPath(obj.Type, obj.Filename)
}

func (isc *incrementalSyncContext) iSyncFilesForTheCreatedObject() {
	events, result, err := isc.dc.vault.ReconcileFiles([]string{isc.objectPath})
	if err == nil {
		err = isc.dc.vault.Project(events)
	}
	isc.syncResult = result
	isc.dc.lastErr = err
}

func (isc *incrementalSyncContext) aFullSyncIsPerformed() {
	events, _, err := isc.dc.vault.Reconcile()
	if err != nil {
		panic(fmt.Sprintf("full reconcile failed: %v", err))
	}
	if err := isc.dc.vault.Project(events); err != nil {
		panic(fmt.Sprintf("project failed: %v", err))
	}
}

func (isc *incrementalSyncContext) theObjectFileIsDeletedFromDisk() {
	os.Remove(isc.objectPath)
}

func (isc *incrementalSyncContext) iSyncFilesForTheDeletedObjectPath() {
	events, result, err := isc.dc.vault.ReconcileFiles([]string{isc.objectPath})
	if err == nil {
		err = isc.dc.vault.Project(events)
	}
	isc.syncResult = result
	isc.dc.lastErr = err
}

func (isc *incrementalSyncContext) theObjectShouldBeSearchableBy(keyword string) error {
	results, err := isc.dc.vault.SearchObjects(keyword)
	if err != nil {
		return fmt.Errorf("search error: %w", err)
	}
	if len(results) == 0 {
		return fmt.Errorf("expected search results for %q, got none", keyword)
	}
	return nil
}

func (isc *incrementalSyncContext) theObjectShouldNotBeInTheIndex() error {
	results, err := isc.dc.vault.QueryObjects(nil)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}
	for _, obj := range results {
		if obj.ID == isc.objectID {
			return fmt.Errorf("expected object %s to be removed from index", isc.objectID)
		}
	}
	return nil
}

func (isc *incrementalSyncContext) iSyncWithEmptyPaths() {
	events, result, err := isc.dc.vault.ReconcileFiles(nil)
	if err == nil {
		err = isc.dc.vault.Project(events)
	}
	isc.syncResult = result
	isc.dc.lastErr = err
}

func initIncrementalSyncSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	isc := newIncrementalSyncContext(dc)

	ctx.Step(`^a "([^"]*)" type schema exists$`, isc.aTypeSchemaExists)
	ctx.Step(`^I create an object "([^"]*)" named "([^"]*)"$`, isc.iCreateAnObjectNamed)
	ctx.Step(`^I sync files for the created object$`, isc.iSyncFilesForTheCreatedObject)
	ctx.Step(`^a full sync is performed$`, isc.aFullSyncIsPerformed)
	ctx.Step(`^the object file is deleted from disk$`, isc.theObjectFileIsDeletedFromDisk)
	ctx.Step(`^I sync files for the deleted object path$`, isc.iSyncFilesForTheDeletedObjectPath)
	ctx.Step(`^the object should be searchable by "([^"]*)"$`, isc.theObjectShouldBeSearchableBy)
	ctx.Step(`^the object should not be in the index$`, isc.theObjectShouldNotBeInTheIndex)
	ctx.Step(`^I sync with empty paths$`, isc.iSyncWithEmptyPaths)
}

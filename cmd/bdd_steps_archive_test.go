package cmd

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/typemd/typemd/core"
)

func initArchiveSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	ctx.Step(`^the object is archived$`, cc.theObjectIsArchived)
	ctx.Step(`^the "([^"]*)" object is archived$`, cc.theNamedObjectIsArchived)
	ctx.Step(`^I run tmd object archive on the object$`, cc.iRunObjectArchive)
	ctx.Step(`^I run tmd object unarchive on the object$`, cc.iRunObjectUnarchive)
	ctx.Step(`^the object should be archived$`, cc.theObjectShouldBeArchivedCmd)
	ctx.Step(`^the object should not be archived$`, cc.theObjectShouldNotBeArchivedCmd)
	ctx.Step(`^I run object list with include-archived$`, cc.iRunObjectListIncludeArchived)
}

func (cc *cmdContext) theObjectIsArchived() error {
	if len(cc.createdObjectIDs) == 0 {
		return fmt.Errorf("no created objects to archive")
	}
	id := cc.createdObjectIDs[len(cc.createdObjectIDs)-1]
	return cc.vault.SetArchived(id, true)
}

func (cc *cmdContext) theNamedObjectIsArchived(name string) error {
	objects, err := cc.vault.QueryObjects(nil, core.QueryIncludeArchived())
	if err != nil {
		return fmt.Errorf("query objects: %w", err)
	}
	for _, obj := range objects {
		if obj.GetName() == name {
			return cc.vault.SetArchived(obj.ID, true)
		}
	}
	return fmt.Errorf("object named %q not found", name)
}

func (cc *cmdContext) iRunObjectArchive() {
	if len(cc.createdObjectIDs) == 0 {
		cc.lastErr = fmt.Errorf("no created objects to archive")
		return
	}
	id := cc.createdObjectIDs[len(cc.createdObjectIDs)-1]
	cc.runCmd("object", "archive", id)
}

func (cc *cmdContext) iRunObjectUnarchive() {
	if len(cc.createdObjectIDs) == 0 {
		cc.lastErr = fmt.Errorf("no created objects to unarchive")
		return
	}
	id := cc.createdObjectIDs[len(cc.createdObjectIDs)-1]
	cc.runCmd("object", "unarchive", id)
}

func (cc *cmdContext) theObjectShouldBeArchivedCmd() error {
	if len(cc.createdObjectIDs) == 0 {
		return fmt.Errorf("no created objects to check")
	}
	id := cc.createdObjectIDs[len(cc.createdObjectIDs)-1]

	vault, err := openVault(cc.vaultDir)
	if err != nil {
		return fmt.Errorf("openVault: %w", err)
	}
	defer vault.Close()

	obj, err := vault.GetObject(id)
	if err != nil {
		return fmt.Errorf("GetObject: %w", err)
	}
	if !obj.IsArchived() {
		return fmt.Errorf("expected object %s to be archived, but it is not", id)
	}
	return nil
}

func (cc *cmdContext) theObjectShouldNotBeArchivedCmd() error {
	if len(cc.createdObjectIDs) == 0 {
		return fmt.Errorf("no created objects to check")
	}
	id := cc.createdObjectIDs[len(cc.createdObjectIDs)-1]

	vault, err := openVault(cc.vaultDir)
	if err != nil {
		return fmt.Errorf("openVault: %w", err)
	}
	defer vault.Close()

	obj, err := vault.GetObject(id)
	if err != nil {
		return fmt.Errorf("GetObject: %w", err)
	}
	if obj.IsArchived() {
		return fmt.Errorf("expected object %s to not be archived, but it is", id)
	}
	return nil
}

func (cc *cmdContext) iRunObjectListIncludeArchived() {
	cc.runCmd("object", "list", "--include-archived")
}

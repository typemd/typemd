package cmd

import (
	"fmt"

	"github.com/cucumber/godog"
)

func initLockSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	ctx.Step(`^a vault with a "([^"]*)" object named "([^"]*)"$`, cc.aVaultWithATypedObject)
	ctx.Step(`^the object is locked$`, cc.theObjectIsLocked)
	ctx.Step(`^I run tmd object lock on the object$`, cc.iRunObjectLock)
	ctx.Step(`^I run tmd object unlock on the object$`, cc.iRunObjectUnlock)
	ctx.Step(`^the object should be locked$`, cc.theObjectShouldBeLocked)
	ctx.Step(`^the object should not be locked$`, cc.theObjectShouldNotBeLocked)
}

func (cc *cmdContext) aVaultWithATypedObject(typeName, name string) error {
	if err := cc.setupVault(); err != nil {
		return err
	}
	obj, err := cc.vault.NewObject(typeName, name, "")
	if err != nil {
		return err
	}
	cc.createdObjectIDs = append(cc.createdObjectIDs, string(obj.ID))
	return nil
}

func (cc *cmdContext) theObjectIsLocked() error {
	if len(cc.createdObjectIDs) == 0 {
		return fmt.Errorf("no created objects to lock")
	}
	id := cc.createdObjectIDs[len(cc.createdObjectIDs)-1]
	return cc.vault.SetLocked(id, true)
}

func (cc *cmdContext) iRunObjectLock() {
	if len(cc.createdObjectIDs) == 0 {
		cc.lastErr = fmt.Errorf("no created objects to lock")
		return
	}
	id := cc.createdObjectIDs[len(cc.createdObjectIDs)-1]
	cc.runCmd("object", "lock", id)
}

func (cc *cmdContext) iRunObjectUnlock() {
	if len(cc.createdObjectIDs) == 0 {
		cc.lastErr = fmt.Errorf("no created objects to unlock")
		return
	}
	id := cc.createdObjectIDs[len(cc.createdObjectIDs)-1]
	cc.runCmd("object", "unlock", id)
}

func (cc *cmdContext) theObjectShouldBeLocked() error {
	if len(cc.createdObjectIDs) == 0 {
		return fmt.Errorf("no created objects to check")
	}
	id := cc.createdObjectIDs[len(cc.createdObjectIDs)-1]

	vault, err := openVault(cc.vaultDir, false)
	if err != nil {
		return fmt.Errorf("openVault: %w", err)
	}
	defer vault.Close()

	obj, err := vault.GetObject(id)
	if err != nil {
		return fmt.Errorf("GetObject: %w", err)
	}
	if !obj.IsLocked() {
		return fmt.Errorf("expected object %s to be locked, but it is not", id)
	}
	return nil
}

func (cc *cmdContext) theObjectShouldNotBeLocked() error {
	if len(cc.createdObjectIDs) == 0 {
		return fmt.Errorf("no created objects to check")
	}
	id := cc.createdObjectIDs[len(cc.createdObjectIDs)-1]

	vault, err := openVault(cc.vaultDir, false)
	if err != nil {
		return fmt.Errorf("openVault: %w", err)
	}
	defer vault.Close()

	obj, err := vault.GetObject(id)
	if err != nil {
		return fmt.Errorf("GetObject: %w", err)
	}
	if obj.IsLocked() {
		return fmt.Errorf("expected object %s to not be locked, but it is", id)
	}
	return nil
}

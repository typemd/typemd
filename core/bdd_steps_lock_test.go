package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/cucumber/godog"
)

// ── Object lock steps ────────────────────────────────────────────────────

func (dc *domainContext) iLockTheObject() {
	dc.lastErr = dc.vault.SetLocked(dc.currentObject.ID, true)
	if dc.lastErr == nil {
		// Refresh in-memory object so subsequent steps see the locked state
		obj, err := dc.vault.GetObject(dc.currentObject.ID)
		if err == nil {
			dc.currentObject = obj
		}
	}
}

func (dc *domainContext) iUnlockTheObject() {
	dc.lastErr = dc.vault.SetLocked(dc.currentObject.ID, false)
	if dc.lastErr == nil {
		// Refresh in-memory object so subsequent steps see the unlocked state
		obj, err := dc.vault.GetObject(dc.currentObject.ID)
		if err == nil {
			dc.currentObject = obj
		}
	}
}

func (dc *domainContext) theObjectShouldBeLocked() error {
	obj, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	if !obj.IsLocked() {
		return fmt.Errorf("expected object to be locked, but it is not")
	}
	return nil
}

func (dc *domainContext) theObjectShouldNotBeLocked() error {
	obj, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	if obj.IsLocked() {
		return fmt.Errorf("expected object to not be locked, but it is")
	}
	return nil
}

func (dc *domainContext) theObjectFrontmatterShouldContain(text string) error {
	data, err := os.ReadFile(dc.vault.ObjectPath(dc.currentObject.Type, dc.currentObject.Filename))
	if err != nil {
		return fmt.Errorf("ReadFile error: %v", err)
	}
	if !strings.Contains(string(data), text) {
		return fmt.Errorf("frontmatter does not contain %q:\n%s", text, string(data))
	}
	return nil
}

func (dc *domainContext) theObjectFrontmatterShouldNotContain(text string) error {
	data, err := os.ReadFile(dc.vault.ObjectPath(dc.currentObject.Type, dc.currentObject.Filename))
	if err != nil {
		return fmt.Errorf("ReadFile error: %v", err)
	}
	if strings.Contains(string(data), text) {
		return fmt.Errorf("frontmatter should not contain %q:\n%s", text, string(data))
	}
	return nil
}

func (dc *domainContext) iLockTheSourceObject() {
	// Lock the book object (source) for relation lock testing
	for _, obj := range dc.objects {
		if obj.Type == "book" {
			dc.lastErr = dc.vault.SetLocked(obj.ID, true)
			if dc.lastErr == nil {
				refreshed, err := dc.vault.GetObject(obj.ID)
				if err == nil {
					dc.currentObject = refreshed
				}
			}
			return
		}
	}
}

func (dc *domainContext) iSetMultiplePropertiesOnTheObject(table *godog.Table) {
	props := make(map[string]any)
	for i, row := range table.Rows {
		if i == 0 {
			continue // skip header
		}
		props[row.Cells[0].Value] = row.Cells[1].Value
	}
	dc.lastErr = dc.vault.SetPropertyMultiple(dc.currentObject.ID, props)
}

func initLockSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^I lock the object$`, dc.iLockTheObject)
	ctx.Step(`^I unlock the object$`, dc.iUnlockTheObject)
	ctx.Step(`^the object should be locked$`, dc.theObjectShouldBeLocked)
	ctx.Step(`^the object should not be locked$`, dc.theObjectShouldNotBeLocked)
	ctx.Step(`^the object frontmatter should contain "([^"]*)"$`, dc.theObjectFrontmatterShouldContain)
	ctx.Step(`^the object frontmatter should not contain "([^"]*)"$`, dc.theObjectFrontmatterShouldNotContain)
	ctx.Step(`^I lock the source object$`, dc.iLockTheSourceObject)
	ctx.Step(`^I set multiple properties on the object:$`, dc.iSetMultiplePropertiesOnTheObject)
}

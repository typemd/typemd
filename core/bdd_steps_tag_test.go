package core

import (
	"fmt"
	"path/filepath"

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
	tagObj, err := dc.vault.NewObject("tag", "test-tag")
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

// ── Tag resolution steps ────────────────────────────────────────────────

func (dc *domainContext) aBookObjectExistsWithTagReferenceByID(bookName string) {
	book, err := dc.vault.NewObject("book", bookName)
	if err != nil {
		panic(fmt.Sprintf("create book: %v", err))
	}
	// Use the current tag object's full ID (expects a prior "go" tag via Background)
	tagObj := dc.objects["go"]
	if tagObj == nil {
		panic("tag object \"go\" not found — ensure a prior step creates it")
	}
	book.Properties[TagsProperty] = []any{tagObj.ID}
	if err := dc.vault.SaveObject(book); err != nil {
		panic(fmt.Sprintf("SaveObject failed: %v", err))
	}
	dc.objects[bookName] = book
	dc.currentObject = book
}

func (dc *domainContext) aBookObjectExistsWithTagReferenceByName(bookName, tagName string) {
	book, err := dc.vault.NewObject("book", bookName)
	if err != nil {
		panic(fmt.Sprintf("create book: %v", err))
	}
	book.Properties[TagsProperty] = []any{TagTypeName + "/" + tagName}
	if err := dc.vault.SaveObject(book); err != nil {
		panic(fmt.Sprintf("SaveObject failed: %v", err))
	}
	dc.objects[bookName] = book
	dc.currentObject = book
}

func (dc *domainContext) theBookShouldHaveATagRelationToTheTag() error {
	book := dc.currentObject
	rels, err := dc.vault.ListRelations(book.ID)
	if err != nil {
		return fmt.Errorf("list relations: %v", err)
	}
	for _, r := range rels {
		if r.Name == TagsProperty && r.FromID == book.ID {
			return nil
		}
	}
	return fmt.Errorf("no tag relation found for %s, got %v", book.ID, rels)
}

func (dc *domainContext) aTagObjectNamedShouldExistOnDisk(name string) error {
	pattern := filepath.Join(dc.vault.ObjectDir(TagTypeName), name+"-*.md")
	matches, _ := filepath.Glob(pattern)
	if len(matches) == 0 {
		return fmt.Errorf("expected tag object %q on disk, found none", name)
	}
	return nil
}

func initTagSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	// Tag / tags property steps
	ctx.Step(`^the object should have property "([^"]*)" with nil value$`, dc.theObjectShouldHavePropertyWithNilValue)
	ctx.Step(`^I set tags on the object to a tag reference$`, dc.iSetTagsOnTheObjectToATagReference)

	// Tag uniqueness steps
	ctx.Step(`^a raw duplicate tag named "([^"]*)" exists$`, dc.aRawDuplicateTagNamedExists)

	// Tag resolution steps
	ctx.Step(`^a "book" object named "([^"]*)" exists with tag reference by ID$`, dc.aBookObjectExistsWithTagReferenceByID)
	ctx.Step(`^a "book" object named "([^"]*)" exists with tag reference by name "([^"]*)"$`, dc.aBookObjectExistsWithTagReferenceByName)
	ctx.Step(`^the book should have a tag relation to the tag$`, dc.theBookShouldHaveATagRelationToTheTag)
	ctx.Step(`^a tag object named "([^"]*)" should exist on disk$`, dc.aTagObjectNamedShouldExistOnDisk)
}

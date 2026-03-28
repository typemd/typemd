package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

type relationSyncContext struct {
	dc         *domainContext
	syncResult *ReconcileResult
}

func newRelationSyncContext(dc *domainContext) *relationSyncContext {
	return &relationSyncContext{dc: dc}
}

// ── Setup steps ─────────────────────────────────────────────────────

func (rsc *relationSyncContext) aVaultIsReadyWithRelationSchemasForSync() {
	rsc.dc.aVaultIsInitialized()

	// Write relation-aware schemas BEFORE Open so they're used from the start
	bookSchema := []byte(`name: book
properties:
  - name: title
    type: string
  - name: author
    type: relation
    target: person
    bidirectional: true
    inverse: books
  - name: editor
    type: relation
    target: person
`)
	dir := filepath.Join(rsc.dc.vault.TypesDir(), "book")
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "schema.yaml"), bookSchema, 0644)

	personSchema := []byte(`name: person
properties:
  - name: books
    type: relation
    target: book
    multiple: true
    bidirectional: true
    inverse: author
`)
	dir = filepath.Join(rsc.dc.vault.TypesDir(), "person")
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "schema.yaml"), personSchema, 0644)

	if err := rsc.dc.vault.Open(); err != nil {
		panic(fmt.Sprintf("vault open failed: %v", err))
	}
}

func (rsc *relationSyncContext) aBookExistsWithAuthorNameReference(bookName, authorRef string) {
	book, err := rsc.dc.vault.NewObject("book", bookName, "")
	if err != nil {
		panic(fmt.Sprintf("create book: %v", err))
	}
	book.Properties["author"] = authorRef
	if err := rsc.dc.vault.SaveObject(book); err != nil {
		panic(fmt.Sprintf("save book: %v", err))
	}
	rsc.dc.objects[bookName] = book
	rsc.dc.currentObject = book
}

func (rsc *relationSyncContext) aBookExistsLinkedToThePersonVia(bookName, relation string) {
	book, err := rsc.dc.vault.NewObject("book", bookName, "")
	if err != nil {
		panic(fmt.Sprintf("create book: %v", err))
	}
	// Find the person object (last created person)
	var personID string
	for _, obj := range rsc.dc.objects {
		if obj.Type == "person" {
			personID = obj.ID
		}
	}
	if personID == "" {
		panic("no person object found")
	}
	book.Properties[relation] = personID
	if err := rsc.dc.vault.SaveObject(book); err != nil {
		panic(fmt.Sprintf("save book: %v", err))
	}
	rsc.dc.objects[bookName] = book
	rsc.dc.currentObject = book
}

func (rsc *relationSyncContext) aBookExistsWithAuthor(bookName, authorRef string) {
	book, err := rsc.dc.vault.NewObject("book", bookName, "")
	if err != nil {
		panic(fmt.Sprintf("create book: %v", err))
	}
	book.Properties["author"] = authorRef
	if err := rsc.dc.vault.SaveObject(book); err != nil {
		panic(fmt.Sprintf("save book: %v", err))
	}
	rsc.dc.objects[bookName] = book
	rsc.dc.currentObject = book
}

func (rsc *relationSyncContext) aBookExistsWithAuthorAndEditor(bookName, authorRef, editorRef string) {
	book, err := rsc.dc.vault.NewObject("book", bookName, "")
	if err != nil {
		panic(fmt.Sprintf("create book: %v", err))
	}
	book.Properties["author"] = authorRef
	book.Properties["editor"] = editorRef
	if err := rsc.dc.vault.SaveObject(book); err != nil {
		panic(fmt.Sprintf("save book: %v", err))
	}
	rsc.dc.objects[bookName] = book
	rsc.dc.currentObject = book
}

func (rsc *relationSyncContext) anotherPersonObjectNamedExists(name string) {
	obj, err := rsc.dc.vault.NewObject("person", name, "")
	if err != nil {
		panic(fmt.Sprintf("create person: %v", err))
	}
	// Store with a different key to avoid overwriting
	rsc.dc.objects[name+"-dup"] = obj
}

func (rsc *relationSyncContext) aPersonExistsWithBooksReferencesToBothBooks(personName string) {
	person, err := rsc.dc.vault.NewObject("person", personName, "")
	if err != nil {
		panic(fmt.Sprintf("create person: %v", err))
	}
	// Find all book objects
	var bookIDs []any
	for slug, obj := range rsc.dc.objects {
		if obj.Type == "book" && slug != personName {
			bookIDs = append(bookIDs, obj.ID)
		}
	}
	person.Properties["books"] = bookIDs
	if err := rsc.dc.vault.SaveObject(person); err != nil {
		panic(fmt.Sprintf("save person: %v", err))
	}
	rsc.dc.objects[personName] = person
	rsc.dc.currentObject = person
}

// ── Action steps ────────────────────────────────────────────────────

func (rsc *relationSyncContext) iSyncTheIndex() {
	events, result, err := rsc.dc.vault.Reconcile()
	if err == nil {
		err = rsc.dc.vault.Project(events)
	}
	rsc.syncResult = result
	rsc.dc.lastErr = err
}

// ── Assertion steps ─────────────────────────────────────────────────

func (rsc *relationSyncContext) theBookFileShouldHaveAuthorExpandedToThePersonsFullID() error {
	book := rsc.dc.currentObject
	freshBook, err := rsc.dc.vault.GetObject(book.ID)
	if err != nil {
		return fmt.Errorf("get book: %v", err)
	}
	authorVal, ok := freshBook.Properties["author"].(string)
	if !ok {
		return fmt.Errorf("author property not a string: %T", freshBook.Properties["author"])
	}
	// Find the person
	for _, obj := range rsc.dc.objects {
		if obj.Type == "person" {
			if authorVal == obj.ID {
				return nil
			}
		}
	}
	return fmt.Errorf("author %q does not match any person object", authorVal)
}

func (rsc *relationSyncContext) theBookFileShouldStillHaveAuthor(expected string) error {
	book := rsc.dc.currentObject
	freshBook, err := rsc.dc.vault.GetObject(book.ID)
	if err != nil {
		return fmt.Errorf("get book: %v", err)
	}
	authorVal, _ := freshBook.Properties["author"].(string)
	if authorVal != expected {
		return fmt.Errorf("expected author %q, got %q", expected, authorVal)
	}
	return nil
}

func (rsc *relationSyncContext) theBookFileShouldHaveAuthorReferencingThePersonsFullID() error {
	return rsc.theBookFileShouldHaveAuthorExpandedToThePersonsFullID()
}

func (rsc *relationSyncContext) theBookFileShouldHaveAuthorExpandedToPersonFullID(personName string) error {
	book := rsc.dc.currentObject
	freshBook, err := rsc.dc.vault.GetObject(book.ID)
	if err != nil {
		return fmt.Errorf("get book: %v", err)
	}
	authorVal, ok := freshBook.Properties["author"].(string)
	if !ok {
		return fmt.Errorf("author not a string: %T", freshBook.Properties["author"])
	}
	person := rsc.dc.objects[personName]
	if person == nil {
		return fmt.Errorf("person %q not found", personName)
	}
	if authorVal != person.ID {
		return fmt.Errorf("expected author %q, got %q", person.ID, authorVal)
	}
	return nil
}

func (rsc *relationSyncContext) theBookFileShouldHaveEditorExpandedToPersonFullID(personName string) error {
	book := rsc.dc.currentObject
	freshBook, err := rsc.dc.vault.GetObject(book.ID)
	if err != nil {
		return fmt.Errorf("get book: %v", err)
	}
	editorVal, ok := freshBook.Properties["editor"].(string)
	if !ok {
		return fmt.Errorf("editor not a string: %T", freshBook.Properties["editor"])
	}
	person := rsc.dc.objects[personName]
	if person == nil {
		return fmt.Errorf("person %q not found", personName)
	}
	if editorVal != person.ID {
		return fmt.Errorf("expected editor %q, got %q", person.ID, editorVal)
	}
	return nil
}

func (rsc *relationSyncContext) theSyncResultShouldHaveNExpansions(expected int) error {
	if rsc.syncResult == nil {
		return fmt.Errorf("no sync result")
	}
	if rsc.syncResult.Expanded != expected {
		return fmt.Errorf("expected %d expansions, got %d", expected, rsc.syncResult.Expanded)
	}
	return nil
}

func (rsc *relationSyncContext) theSyncResultShouldHaveNUnresolvedReferences(expected int) error {
	if rsc.syncResult == nil {
		return fmt.Errorf("no sync result")
	}
	if len(rsc.syncResult.Unresolved) != expected {
		return fmt.Errorf("expected %d unresolved, got %d", expected, len(rsc.syncResult.Unresolved))
	}
	return nil
}

func initRelationSyncSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	rsc := newRelationSyncContext(dc)

	// Background
	ctx.Step(`^a vault is ready with relation sync schemas$`, rsc.aVaultIsReadyWithRelationSchemasForSync)

	// Setup
	ctx.Step(`^a "book" object named "([^"]*)" exists with author name reference "([^"]*)"$`, rsc.aBookExistsWithAuthorNameReference)
	ctx.Step(`^a "book" object named "([^"]*)" exists linked to the person via "([^"]*)"$`, rsc.aBookExistsLinkedToThePersonVia)
	ctx.Step(`^a "book" object named "([^"]*)" exists with author "([^"]*)"$`, rsc.aBookExistsWithAuthor)
	ctx.Step(`^a "book" object named "([^"]*)" exists with author "([^"]*)" and editor "([^"]*)"$`, rsc.aBookExistsWithAuthorAndEditor)
	ctx.Step(`^another "person" object named "([^"]*)" exists$`, rsc.anotherPersonObjectNamedExists)
	ctx.Step(`^a "person" object named "([^"]*)" exists with books references to both books$`, rsc.aPersonExistsWithBooksReferencesToBothBooks)

	// Actions
	ctx.Step(`^I run a full relation sync$`, rsc.iSyncTheIndex)

	// Assertions
	ctx.Step(`^the book file should have author expanded to the person's full ID$`, rsc.theBookFileShouldHaveAuthorExpandedToThePersonsFullID)
	ctx.Step(`^the book file should still have author "([^"]*)"$`, rsc.theBookFileShouldStillHaveAuthor)
	ctx.Step(`^the book file should have author referencing the person's full ID$`, rsc.theBookFileShouldHaveAuthorReferencingThePersonsFullID)
	ctx.Step(`^the book file should have author expanded to "([^"]*)"'s full ID$`, rsc.theBookFileShouldHaveAuthorExpandedToPersonFullID)
	ctx.Step(`^the book file should have editor expanded to "([^"]*)"'s full ID$`, rsc.theBookFileShouldHaveEditorExpandedToPersonFullID)
	ctx.Step(`^the sync result should have (\d+) expansions$`, rsc.theSyncResultShouldHaveNExpansions)
	ctx.Step(`^the sync result should have (\d+) unresolved references?$`, rsc.theSyncResultShouldHaveNUnresolvedReferences)
}

package core

import (
	"fmt"

	"github.com/cucumber/godog"
)

// ── Shared relation sync setup steps (used by validate.feature) ─────────

func (dc *domainContext) aVaultIsReadyWithRelationSyncSchemas() {
	dc.aVaultIsInitialized()

	mustWriteTypeSchema(dc.vault, "book", []byte(`name: book
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
`))

	mustWriteTypeSchema(dc.vault, "person", []byte(`name: person
properties:
  - name: books
    type: relation
    target: book
    multiple: true
    bidirectional: true
    inverse: author
`))

	if err := dc.vault.Open(); err != nil {
		panic(fmt.Sprintf("vault open failed: %v", err))
	}
}

func (dc *domainContext) aBookExistsWithAuthorNameReference(bookName, authorRef string) {
	book, err := dc.vault.NewObject("book", bookName, "")
	if err != nil {
		panic(fmt.Sprintf("create book: %v", err))
	}
	book.Properties["author"] = authorRef
	if err := dc.vault.SaveObject(book); err != nil {
		panic(fmt.Sprintf("save book: %v", err))
	}
	dc.objects[bookName] = book
	dc.currentObject = book
}

func (dc *domainContext) aBookExistsLinkedToThePersonVia(bookName, relation string) {
	book, err := dc.vault.NewObject("book", bookName, "")
	if err != nil {
		panic(fmt.Sprintf("create book: %v", err))
	}
	var personID string
	for _, obj := range dc.objects {
		if obj.Type == "person" {
			personID = obj.ID
		}
	}
	if personID == "" {
		panic("no person object found")
	}
	book.Properties[relation] = personID
	if err := dc.vault.SaveObject(book); err != nil {
		panic(fmt.Sprintf("save book: %v", err))
	}
	dc.objects[bookName] = book
	dc.currentObject = book
}

func initRelationSyncSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a vault is ready with relation sync schemas$`, dc.aVaultIsReadyWithRelationSyncSchemas)
	ctx.Step(`^a "book" object named "([^"]*)" exists with author name reference "([^"]*)"$`, dc.aBookExistsWithAuthorNameReference)
	ctx.Step(`^a "book" object named "([^"]*)" exists linked to the person via "([^"]*)"$`, dc.aBookExistsLinkedToThePersonVia)
}

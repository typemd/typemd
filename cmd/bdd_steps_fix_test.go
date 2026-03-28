package cmd

import (
	"fmt"
	"os"

	"github.com/cucumber/godog"
)

func initFixSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	ctx.Step(`^a "note" object with a shorthand wiki-link exists$`, cc.aNoteObjectWithShorthandWikiLink)
	ctx.Step(`^an object with an ambiguous shorthand wiki-link exists$`, cc.anObjectWithAmbiguousShorthandWikiLink)
	ctx.Step(`^I run fix wikilinks$`, cc.iRunFixWikilinks)
}

func (cc *cmdContext) aNoteObjectWithShorthandWikiLink() error {
	// Create a note type schema
	noteSchema := `name: note
properties:
  - name: title
    type: string
`
	if err := os.MkdirAll(cc.vault.TypesDir()+"/note", 0755); err != nil {
		return err
	}
	if err := os.WriteFile(cc.vault.TypesDir()+"/note/schema.yaml", []byte(noteSchema), 0644); err != nil {
		return err
	}

	// Create two note objects
	target, err := cc.vault.NewObject("note", "target-note", "")
	if err != nil {
		return fmt.Errorf("create target: %w", err)
	}

	source, err := cc.vault.NewObject("note", "source-note", "")
	if err != nil {
		return fmt.Errorf("create source: %w", err)
	}

	// Write shorthand wiki-link in source body
	body := "---\nname: source-note\n---\n\nSee [[target-note]].\n"
	_ = target // just to reference it
	return os.WriteFile(cc.vault.ObjectPath(source.Type, source.Filename), []byte(body), 0644)
}

func (cc *cmdContext) anObjectWithAmbiguousShorthandWikiLink() error {
	// Create a book type schema (already exists from setupVault)
	// Create two books with the same slug prefix
	_, err := cc.vault.NewObject("book", "golang", "")
	if err != nil {
		return fmt.Errorf("create golang 1: %w", err)
	}

	_, err = cc.vault.NewObject("book", "golang", "")
	if err != nil {
		return fmt.Errorf("create golang 2: %w", err)
	}

	// Create a note that links to "book/golang"
	noteSchema := `name: note
properties:
  - name: title
    type: string
`
	os.MkdirAll(cc.vault.TypesDir()+"/note", 0755)
	os.WriteFile(cc.vault.TypesDir()+"/note/schema.yaml", []byte(noteSchema), 0644)

	source, err := cc.vault.NewObject("note", "my-note", "")
	if err != nil {
		return fmt.Errorf("create note: %w", err)
	}

	body := "---\nname: my-note\n---\n\nSee [[book/golang]].\n"
	return os.WriteFile(cc.vault.ObjectPath(source.Type, source.Filename), []byte(body), 0644)
}

func (cc *cmdContext) iRunFixWikilinks() {
	cc.runCmd("fix", "wikilinks")
}

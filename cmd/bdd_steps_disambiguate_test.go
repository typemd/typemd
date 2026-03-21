package cmd

import (
	"fmt"
	"os"

	"github.com/cucumber/godog"
)

func initDisambiguateSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	ctx.Step(`^a vault with two books sharing prefix "([^"]*)"$`, cc.aVaultWithTwoBooksSharing)
	ctx.Step(`^a vault with a person "([^"]*)"$`, cc.aVaultWithAPerson)
	ctx.Step(`^a link from the first book to the person via "([^"]*)"$`, cc.aLinkFromFirstBookToPersonVia)
	ctx.Step(`^I run object show "([^"]*)" in non-interactive mode$`, cc.iRunObjectShowNonInteractive)
	ctx.Step(`^I run object show "([^"]*)" selecting candidate (\d+)$`, cc.iRunObjectShowSelectingCandidate)
	ctx.Step(`^I run object show "([^"]*)" and cancel the picker$`, cc.iRunObjectShowAndCancel)
	ctx.Step(`^I run relation link with ambiguous from-id "([^"]*)" selecting candidate (\d+)$`, cc.iRunRelationLinkAmbiguousFromID)
	ctx.Step(`^I run relation unlink with ambiguous from-id "([^"]*)" selecting candidate (\d+)$`, cc.iRunRelationUnlinkAmbiguousFromID)
}

func (cc *cmdContext) aVaultWithTwoBooksSharing(prefix string) error {
	if err := cc.setupVault(); err != nil {
		return err
	}

	// Add author relation and person type so link/unlink scenarios work.
	typesDir := cc.vault.TypesDir()
	if err := os.WriteFile(typesDir+"/book.yaml", []byte(`name: book
emoji: "\U0001F4DA"
plural: books
properties:
  - name: author
    type: relation
    target: person
  - name: rating
    type: number
`), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(typesDir+"/person.yaml", []byte(`name: person
emoji: "\U0001F464"
plural: people
properties:
  - name: role
    type: string
`), 0644); err != nil {
		return err
	}
	cc.vault.InvalidateSchemaCache()
	if _, err := cc.vault.SyncIndex(); err != nil {
		return err
	}

	obj1, err := cc.vault.NewObject("book", prefix+"-first-edition", "")
	if err != nil {
		return fmt.Errorf("create first book: %w", err)
	}
	cc.createdObjectIDs = append(cc.createdObjectIDs, string(obj1.ID))

	obj2, err := cc.vault.NewObject("book", prefix+"-second-edition", "")
	if err != nil {
		return fmt.Errorf("create second book: %w", err)
	}
	cc.createdObjectIDs = append(cc.createdObjectIDs, string(obj2.ID))
	return nil
}

func (cc *cmdContext) aVaultWithAPerson(name string) error {
	obj, err := cc.vault.NewObject("person", name, "")
	if err != nil {
		return fmt.Errorf("create person: %w", err)
	}
	cc.createdObjectIDs = append(cc.createdObjectIDs, string(obj.ID))
	return nil
}

func (cc *cmdContext) aLinkFromFirstBookToPersonVia(relation string) error {
	if len(cc.createdObjectIDs) < 3 {
		return fmt.Errorf("need at least 3 created objects (2 books + 1 person), got %d", len(cc.createdObjectIDs))
	}
	bookID := cc.createdObjectIDs[0]
	personID := cc.createdObjectIDs[2] // person is the 3rd created object
	return cc.vault.LinkObjects(bookID, relation, personID)
}

func (cc *cmdContext) iRunObjectShowNonInteractive(prefix string) {
	cc.runCmd("object", "show", prefix)
}

// withInteractiveDisambiguation overrides disambiguateFunc and isInteractiveFunc
// for the duration of fn, then restores the originals.
func withInteractiveDisambiguation(selectFn func([]disambiguateItem) (string, error), fn func()) {
	origDisambiguate := disambiguateFunc
	origIsInteractive := isInteractiveFunc
	disambiguateFunc = selectFn
	isInteractiveFunc = func() bool { return true }
	defer func() {
		disambiguateFunc = origDisambiguate
		isInteractiveFunc = origIsInteractive
	}()
	fn()
}

func selectCandidate(candidate int) func([]disambiguateItem) (string, error) {
	return func(items []disambiguateItem) (string, error) {
		idx := candidate - 1
		if idx < 0 || idx >= len(items) {
			return "", fmt.Errorf("candidate %d out of range (have %d items)", candidate, len(items))
		}
		return items[idx].id, nil
	}
}

func cancelPicker([]disambiguateItem) (string, error) {
	return "", nil
}

func (cc *cmdContext) iRunObjectShowSelectingCandidate(prefix string, candidate int) {
	withInteractiveDisambiguation(selectCandidate(candidate), func() {
		cc.runCmd("object", "show", prefix)
	})
}

func (cc *cmdContext) iRunObjectShowAndCancel(prefix string) {
	withInteractiveDisambiguation(cancelPicker, func() {
		cc.runCmd("object", "show", prefix)
	})
}

func (cc *cmdContext) iRunRelationLinkAmbiguousFromID(prefix string, candidate int) {
	if len(cc.createdObjectIDs) < 3 {
		cc.lastErr = fmt.Errorf("need at least 3 objects (2 books + 1 person)")
		return
	}
	personID := cc.createdObjectIDs[2]
	withInteractiveDisambiguation(selectCandidate(candidate), func() {
		cc.runCmd("relation", "link", prefix, "author", personID)
	})
}

func (cc *cmdContext) iRunRelationUnlinkAmbiguousFromID(prefix string, candidate int) {
	if len(cc.createdObjectIDs) < 3 {
		cc.lastErr = fmt.Errorf("need at least 3 objects (2 books + 1 person)")
		return
	}
	personID := cc.createdObjectIDs[2]
	withInteractiveDisambiguation(selectCandidate(candidate), func() {
		cc.runCmd("relation", "unlink", prefix, "author", personID)
	})
}

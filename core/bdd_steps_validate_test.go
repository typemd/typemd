package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

// ── Validate steps ──────────────────────────────────────────────────────────

func (dc *domainContext) aTypeSchemaWithAStringProperty(typeName, propName string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: %s\n    type: string\n", typeName, propName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithASelectPropertyMissingOptions(typeName string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: status\n    type: select\n", typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) iValidateAllSchemas() {
	dc.schemaErrors = ValidateAllSchemas(dc.vault)
}

func (dc *domainContext) schemaShouldHaveNoErrors(typeName string) error {
	if errs, ok := dc.schemaErrors[typeName]; ok && len(errs) > 0 {
		return fmt.Errorf("expected no errors for %q, got %v", typeName, errs)
	}
	return nil
}

func (dc *domainContext) schemaShouldHaveErrors(typeName string) error {
	errs, ok := dc.schemaErrors[typeName]
	if !ok || len(errs) == 0 {
		return fmt.Errorf("expected errors for %q, got none", typeName)
	}
	return nil
}

func (dc *domainContext) anOrphanedRelationExists(fromID, toID string) {
	dc.vault.db.Exec("INSERT INTO relations (name, from_id, to_id) VALUES (?, ?, ?)",
		"author", fromID, toID)
	dc.vault.db.Exec("INSERT INTO objects (id, type, filename, properties, body) VALUES (?, ?, ?, ?, ?)",
		fromID, "book", "test-book", "{}", "")
}

func (dc *domainContext) iValidateRelations() {
	dc.relationErrors = ValidateRelations(dc.vault)
}

func (dc *domainContext) thereShouldBeNRelationErrors(expected int) error {
	if len(dc.relationErrors) != expected {
		return fmt.Errorf("relation errors = %d, want %d: %v", len(dc.relationErrors), expected, dc.relationErrors)
	}
	return nil
}

func (dc *domainContext) twoLinkedNotesExist() {
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "note.yaml"),
		[]byte("name: note\nproperties:\n  - name: title\n    type: string\n"), 0644)

	noteA, _ := dc.vault.NewObject("note", "alpha", "")
	noteB, _ := dc.vault.NewObject("note", "beta", "")
	dc.objects["alpha"] = noteA
	dc.objects["beta"] = noteB

	body := fmt.Sprintf("---\ntitle: Alpha\n---\n\nSee [[%s]].\n", noteB.ID)
	os.WriteFile(dc.vault.ObjectPath(noteA.Type, noteA.Filename), []byte(body), 0644)
	dc.vault.SyncIndex()
}

func (dc *domainContext) aNoteWithABrokenWikiLinkExists() {
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "note.yaml"),
		[]byte("name: note\nproperties:\n  - name: title\n    type: string\n"), 0644)

	note, _ := dc.vault.NewObject("note", "alpha", "")
	dc.objects["alpha"] = note

	body := "---\ntitle: Alpha\n---\n\nSee [[note/nonexistent-01jjjjjjjjjjjjjjjjjjjjjjjj]].\n"
	os.WriteFile(dc.vault.ObjectPath(note.Type, note.Filename), []byte(body), 0644)
	dc.vault.SyncIndex()
}

func (dc *domainContext) iValidateWikiLinks() {
	dc.wikiLinkErrors = ValidateWikiLinks(dc.vault)
}

func (dc *domainContext) thereShouldBeNoWikiLinkErrors() error {
	if len(dc.wikiLinkErrors) != 0 {
		return fmt.Errorf("expected no wiki-link errors, got %v", dc.wikiLinkErrors)
	}
	return nil
}

func (dc *domainContext) thereShouldBeNWikiLinkErrors(expected int) error {
	if len(dc.wikiLinkErrors) != expected {
		return fmt.Errorf("wiki-link errors = %d, want %d", len(dc.wikiLinkErrors), expected)
	}
	return nil
}

func (dc *domainContext) theErrorShouldMention(substr string) error {
	for _, err := range dc.wikiLinkErrors {
		if strings.Contains(err.Error(), substr) {
			return nil
		}
	}
	return fmt.Errorf("no error mentions %q", substr)
}

func initValidateSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a type schema "([^"]*)" with a "([^"]*)" string property$`, dc.aTypeSchemaWithAStringProperty)
	ctx.Step(`^a type schema "([^"]*)" with a select property missing options$`, dc.aTypeSchemaWithASelectPropertyMissingOptions)
	ctx.Step(`^I validate all schemas$`, dc.iValidateAllSchemas)
	ctx.Step(`^schema "([^"]*)" should have no errors$`, dc.schemaShouldHaveNoErrors)
	ctx.Step(`^schema "([^"]*)" should have errors$`, dc.schemaShouldHaveErrors)
	ctx.Step(`^an orphaned relation from "([^"]*)" to "([^"]*)" exists$`, dc.anOrphanedRelationExists)
	ctx.Step(`^I validate relations$`, dc.iValidateRelations)
	ctx.Step(`^there should be (\d+) relation errors?$`, dc.thereShouldBeNRelationErrors)
	ctx.Step(`^two linked notes exist$`, dc.twoLinkedNotesExist)
	ctx.Step(`^a note with a broken wiki-link exists$`, dc.aNoteWithABrokenWikiLinkExists)
	ctx.Step(`^I validate wiki-links$`, dc.iValidateWikiLinks)
	ctx.Step(`^there should be no wiki-link errors$`, dc.thereShouldBeNoWikiLinkErrors)
	ctx.Step(`^there should be (\d+) wiki-link errors?$`, dc.thereShouldBeNWikiLinkErrors)
	ctx.Step(`^the error should mention "([^"]*)"$`, dc.theErrorShouldMention)
}

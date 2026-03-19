package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

type schemaCacheContext struct {
	dc *domainContext
}

func newSchemaCacheContext(dc *domainContext) *schemaCacheContext {
	return &schemaCacheContext{dc: dc}
}

func (sc *schemaCacheContext) aTypeSchemaFileWithEmoji(name, emoji string) {
	data := fmt.Sprintf("name: %s\nemoji: \"%s\"\nproperties: []\n", name, emoji)
	dir := filepath.Join(sc.dc.vault.TypesDir(), name)
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "schema.yaml"), []byte(data), 0644)
}

func (sc *schemaCacheContext) iSaveTypeWithEmoji(name, emoji string) {
	schema := &TypeSchema{
		Name:  name,
		Emoji: emoji,
	}
	sc.dc.lastErr = sc.dc.vault.SaveType(schema)
}

func (sc *schemaCacheContext) theTypeSchemaFileIsChangedToEmojiOnDisk(name, emoji string) {
	data := fmt.Sprintf("name: %s\nemoji: \"%s\"\nproperties: []\n", name, emoji)
	dir := filepath.Join(sc.dc.vault.TypesDir(), name)
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "schema.yaml"), []byte(data), 0644)
}

func (sc *schemaCacheContext) iInvalidateTheSchemaCache() {
	sc.dc.vault.InvalidateSchemaCache()
}

func (sc *schemaCacheContext) theLoadedSchemaShouldHaveEmoji(expected string) error {
	if sc.dc.loadedSchema == nil {
		return fmt.Errorf("no schema loaded")
	}
	if sc.dc.loadedSchema.Emoji != expected {
		return fmt.Errorf("expected emoji %q, got %q", expected, sc.dc.loadedSchema.Emoji)
	}
	return nil
}

func initSchemaCacheSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	sc := newSchemaCacheContext(dc)

	ctx.Step(`^a type schema file "([^"]*)" with emoji "([^"]*)"$`, sc.aTypeSchemaFileWithEmoji)
	ctx.Step(`^I save type "([^"]*)" with emoji "([^"]*)"$`, sc.iSaveTypeWithEmoji)
	ctx.Step(`^the type schema file "([^"]*)" is changed to emoji "([^"]*)" on disk$`, sc.theTypeSchemaFileIsChangedToEmojiOnDisk)
	ctx.Step(`^I invalidate the schema cache$`, sc.iInvalidateTheSchemaCache)
	ctx.Step(`^the loaded schema should have emoji "([^"]*)"$`, sc.theLoadedSchemaShouldHaveEmoji)
}

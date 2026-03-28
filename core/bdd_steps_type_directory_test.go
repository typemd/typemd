package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

// ── Type directory step state ────────────────────────────────────────────────

type typeDirectoryContext struct {
	dc *domainContext
}

func newTypeDirectoryContext(dc *domainContext) *typeDirectoryContext {
	return &typeDirectoryContext{dc: dc}
}

// ── Given steps ─────────────────────────────────────────────────────────────

func (td *typeDirectoryContext) aTypeSchemaDirectoryWithContent(name string, content *godog.DocString) {
	dirPath := filepath.Join(td.dc.vault.TypesDir(), name)
	os.MkdirAll(dirPath, 0755)
	schemaPath := filepath.Join(dirPath, "schema.yaml")
	os.WriteFile(schemaPath, []byte(content.Content), 0644)
}

// ── Then steps ──────────────────────────────────────────────────────────────

func (td *typeDirectoryContext) theTypeSchemaDirectoryShouldExist(name string) error {
	dirPath := filepath.Join(td.dc.vault.TypesDir(), name)
	schemaPath := filepath.Join(dirPath, "schema.yaml")
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		return fmt.Errorf("expected type directory %s/schema.yaml to exist", name)
	}
	return nil
}

func (td *typeDirectoryContext) theTypeSchemaDirectoryShouldNotExist(name string) error {
	dirPath := filepath.Join(td.dc.vault.TypesDir(), name)
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		return fmt.Errorf("expected type directory %s to not exist", name)
	}
	return nil
}

func (td *typeDirectoryContext) theLoadedSchemaShouldHaveNProperties(n int) error {
	if td.dc.loadedSchema == nil {
		return fmt.Errorf("no schema loaded")
	}
	if len(td.dc.loadedSchema.Properties) != n {
		return fmt.Errorf("expected %d properties, got %d", n, len(td.dc.loadedSchema.Properties))
	}
	return nil
}

// ── Init ────────────────────────────────────────────────────────────────────

func initTypeDirectorySteps(ctx *godog.ScenarioContext, dc *domainContext) {
	td := newTypeDirectoryContext(dc)

	// Given
	ctx.Step(`^a type schema directory "([^"]*)" with schema content:$`, td.aTypeSchemaDirectoryWithContent)

	// Then
	ctx.Step(`^the type schema directory "([^"]*)" should exist$`, td.theTypeSchemaDirectoryShouldExist)
	ctx.Step(`^the type schema directory "([^"]*)" should not exist$`, td.theTypeSchemaDirectoryShouldNotExist)
	ctx.Step(`^the loaded schema should have (\d+) propert(?:y|ies)$`, td.theLoadedSchemaShouldHaveNProperties)
}

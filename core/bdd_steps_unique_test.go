package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

// ── Unique constraint steps ──────────────────────────────────────────────

func (dc *domainContext) aTypeSchemaWithUniqueConstraint(typeName string) {
	schema := fmt.Sprintf("name: %s\nunique: true\nproperties:\n  - name: role\n    type: string\n", typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) iLoadTheTypeSchema(typeName string) {
	schema, err := dc.vault.LoadType(typeName)
	dc.lastErr = err
	if err == nil {
		dc.loadedSchema = schema
	}
}

func (dc *domainContext) theLoadedSchemaShouldHaveUniqueTrue() error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("no schema loaded")
	}
	if !dc.loadedSchema.Unique {
		return fmt.Errorf("expected Unique to be true, got false")
	}
	return nil
}

func (dc *domainContext) theLoadedSchemaShouldHaveUniqueFalse() error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("no schema loaded")
	}
	if dc.loadedSchema.Unique {
		return fmt.Errorf("expected Unique to be false, got true")
	}
	return nil
}

func (dc *domainContext) aRawDuplicateObjectOfTypeNamedExists(typeName, name string) {
	// Create a raw object file bypassing the uniqueness check
	ulid := mustULID()
	filename := name + "-" + ulid
	objPath := dc.vault.ObjectPath(typeName, filename)
	os.MkdirAll(filepath.Dir(objPath), 0755)
	content := fmt.Sprintf("---\nname: %s\ncreated_at: 2026-01-01T00:00:00+08:00\nupdated_at: 2026-01-01T00:00:00+08:00\n---\n", name)
	os.WriteFile(objPath, []byte(content), 0644)
	// Also insert into DB
	propsJSON := fmt.Sprintf(`{"name":"%s"}`, name)
	dc.vault.db.Exec(
		"INSERT INTO objects (id, type, filename, properties, body) VALUES (?, ?, ?, ?, ?)",
		typeName+"/"+filename, typeName, filename, propsJSON, "",
	)
}

func (dc *domainContext) iValidateNameUniqueness() {
	dc.nameUniquenessErrors = ValidateNameUniqueness(dc.vault)
}

func (dc *domainContext) thereShouldBeNoNameUniquenessErrors() error {
	if len(dc.nameUniquenessErrors) > 0 {
		return fmt.Errorf("expected no name uniqueness errors, got %v", dc.nameUniquenessErrors)
	}
	return nil
}

func (dc *domainContext) thereShouldBeNameUniquenessErrors() error {
	if len(dc.nameUniquenessErrors) == 0 {
		return fmt.Errorf("expected name uniqueness errors, got none")
	}
	return nil
}

func initUniqueSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a type schema "([^"]*)" with unique constraint$`, dc.aTypeSchemaWithUniqueConstraint)
	ctx.Step(`^I load the type schema "([^"]*)"$`, dc.iLoadTheTypeSchema)
	ctx.Step(`^the loaded schema should have unique true$`, dc.theLoadedSchemaShouldHaveUniqueTrue)
	ctx.Step(`^the loaded schema should have unique false$`, dc.theLoadedSchemaShouldHaveUniqueFalse)
	ctx.Step(`^a raw duplicate object of type "([^"]*)" named "([^"]*)" exists$`, dc.aRawDuplicateObjectOfTypeNamedExists)
	ctx.Step(`^I validate name uniqueness$`, dc.iValidateNameUniqueness)
	ctx.Step(`^there should be no name uniqueness errors$`, dc.thereShouldBeNoNameUniquenessErrors)
	ctx.Step(`^there should be name uniqueness errors$`, dc.thereShouldBeNameUniquenessErrors)
}

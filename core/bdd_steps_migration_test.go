package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

// ── Migration steps ────────────────────────────────────────────────────────

func (dc *domainContext) createTypesAt(relPath string) {
	dir := filepath.Join(dc.rootDir, relPath)
	os.MkdirAll(filepath.Join(dir, "book"), 0755)
	os.WriteFile(filepath.Join(dir, "book", "schema.yaml"), []byte("name: book\n"), 0644)
}

func (dc *domainContext) typesExistAtOldLocation(relPath string) {
	dc.createTypesAt(relPath)
	// Remove empty new-location types/ (created by init) to simulate old layout
	newTypes := filepath.Join(dc.rootDir, "types")
	entries, _ := os.ReadDir(newTypes)
	if len(entries) == 0 {
		os.Remove(newTypes)
	}
}

func (dc *domainContext) typesExistAtNewLocation(relPath string) {
	dc.createTypesAt(relPath)
}

func (dc *domainContext) createPropertiesAt(relPath string) {
	path := filepath.Join(dc.rootDir, relPath)
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, []byte("- name: rating\n  type: number\n"), 0644)
}

func (dc *domainContext) propertiesFileExistsAtOldLocation(relPath string) {
	dc.createPropertiesAt(relPath)
	// Remove empty new-location properties/ (created by init) to simulate old layout
	newPropsDir := filepath.Join(dc.rootDir, "properties")
	entries, _ := os.ReadDir(newPropsDir)
	if len(entries) == 0 {
		os.Remove(newPropsDir)
	}
}

func (dc *domainContext) propertiesFileExistsAtNewLocation(relPath string) {
	dc.createPropertiesAt(relPath)
}

func (dc *domainContext) typesShouldExistAt(relPath string) error {
	dir := filepath.Join(dc.rootDir, relPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("expected types directory at %s to exist", relPath)
	}
	schema := filepath.Join(dir, "book", "schema.yaml")
	if _, err := os.Stat(schema); os.IsNotExist(err) {
		return fmt.Errorf("expected schema file at %s to exist", schema)
	}
	return nil
}

func (dc *domainContext) typesShouldNotExistAt(relPath string) error {
	dir := filepath.Join(dc.rootDir, relPath)
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("expected types directory at %s to NOT exist", relPath)
	}
	return nil
}

func (dc *domainContext) propertiesFileShouldExistAt(relPath string) error {
	path := filepath.Join(dc.rootDir, relPath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("expected properties file at %s to exist", relPath)
	}
	return nil
}

func (dc *domainContext) propertiesFileShouldNotExistAt(relPath string) error {
	path := filepath.Join(dc.rootDir, relPath)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("expected properties file at %s to NOT exist", relPath)
	}
	return nil
}

func (dc *domainContext) theLastErrorShouldMention(substr string) error {
	if dc.lastErr == nil {
		return fmt.Errorf("expected an error, got nil")
	}
	if !strings.Contains(dc.lastErr.Error(), substr) {
		return fmt.Errorf("expected error to contain %q, got %q", substr, dc.lastErr.Error())
	}
	return nil
}

// ── Per-property migration steps ──────────────────────────────────────────

func (dc *domainContext) aLegacyPropertiesFileWithTwoProperties(prop1, prop2 string) {
	dir := filepath.Join(dc.rootDir, "properties")
	os.MkdirAll(dir, 0755)
	content := fmt.Sprintf(`properties:
  - name: %s
    type: date
    emoji: 📅
  - name: %s
    type: select
    options:
      - value: high
      - value: low
`, prop1, prop2)
	os.WriteFile(filepath.Join(dir, "properties.yaml"), []byte(content), 0644)
}

func (dc *domainContext) anEmptyLegacyPropertiesFile() {
	dir := filepath.Join(dc.rootDir, "properties")
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "properties.yaml"), []byte("properties:\n"), 0644)
}

func (dc *domainContext) aPerPropertyFileExistsInPropertiesDirectory(filename string) {
	dir := filepath.Join(dc.rootDir, "properties")
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, filename), []byte("type: number\n"), 0644)
}

func (dc *domainContext) perPropertyFileShouldExistInPropertiesDirectory(filename string) error {
	path := filepath.Join(dc.rootDir, "properties", filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("expected per-property file %q to exist", filename)
	}
	return nil
}

func (dc *domainContext) legacyFileShouldNotExistInPropertiesDirectory(filename string) error {
	path := filepath.Join(dc.rootDir, "properties", filename)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("expected legacy file %q to NOT exist", filename)
	}
	return nil
}

func initMigrationSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^types exist at the old location "([^"]*)"$`, dc.typesExistAtOldLocation)
	ctx.Step(`^types exist at the new location "([^"]*)"$`, dc.typesExistAtNewLocation)
	ctx.Step(`^a properties file exists at the old location "([^"]*)"$`, dc.propertiesFileExistsAtOldLocation)
	ctx.Step(`^a properties file exists at the new location "([^"]*)"$`, dc.propertiesFileExistsAtNewLocation)
	ctx.Step(`^types should exist at "([^"]*)"$`, dc.typesShouldExistAt)
	ctx.Step(`^types should not exist at "([^"]*)"$`, dc.typesShouldNotExistAt)
	ctx.Step(`^a properties file should exist at "([^"]*)"$`, dc.propertiesFileShouldExistAt)
	ctx.Step(`^a properties file should not exist at "([^"]*)"$`, dc.propertiesFileShouldNotExistAt)
	ctx.Step(`^a properties file should still exist at "([^"]*)"$`, dc.propertiesFileShouldExistAt)
	ctx.Step(`^the last error should mention "([^"]*)"$`, dc.theLastErrorShouldMention)
	ctx.Step(`^a legacy properties file with "([^"]*)" and "([^"]*)" properties$`, dc.aLegacyPropertiesFileWithTwoProperties)
	ctx.Step(`^an empty legacy properties file$`, dc.anEmptyLegacyPropertiesFile)
	ctx.Step(`^a per-property file "([^"]*)" exists in properties directory$`, dc.aPerPropertyFileExistsInPropertiesDirectory)
	ctx.Step(`^per-property file "([^"]*)" should exist in properties directory$`, dc.perPropertyFileShouldExistInPropertiesDirectory)
	ctx.Step(`^legacy "([^"]*)" should not exist in properties directory$`, dc.legacyFileShouldNotExistInPropertiesDirectory)
}

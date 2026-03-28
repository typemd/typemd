package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

// ── Property filtering steps ────────────────────────────────────────────────

func (dc *domainContext) aTypeSchemaWithProperties(typeName, propList string) {
	props := strings.Split(propList, ",")
	var yamlProps string
	for _, p := range props {
		yamlProps += fmt.Sprintf("  - name: %s\n    type: string\n", strings.TrimSpace(p))
	}
	schema := fmt.Sprintf("name: %s\nproperties:\n%s", typeName, yamlProps)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func (dc *domainContext) aRawObjectFileWithProperties(relPath string, table *godog.Table) {
	var yamlContent string
	for _, row := range table.Rows[1:] { // skip header
		yamlContent += fmt.Sprintf("%s: %s\n", row.Cells[0].Value, row.Cells[1].Value)
	}
	content := fmt.Sprintf("---\n%s---\nSome body content.\n", yamlContent)

	fullPath := filepath.Join(dc.vault.ObjectsDir(), relPath)
	os.MkdirAll(filepath.Dir(fullPath), 0755)
	os.WriteFile(fullPath, []byte(content), 0644)
}

func (dc *domainContext) getIndexedProperties(objectID string) (map[string]any, error) {
	var propsJSON string
	err := dc.vault.db.QueryRow("SELECT properties FROM objects WHERE id = ?", objectID).Scan(&propsJSON)
	if err != nil {
		return nil, fmt.Errorf("query properties for %s: %v", objectID, err)
	}
	var props map[string]any
	if err := json.Unmarshal([]byte(propsJSON), &props); err != nil {
		return nil, fmt.Errorf("unmarshal properties for %s: %v", objectID, err)
	}
	return props, nil
}

func (dc *domainContext) theIndexedPropertiesForShouldContain(objectID, key string) error {
	props, err := dc.getIndexedProperties(objectID)
	if err != nil {
		return err
	}
	if _, ok := props[key]; !ok {
		return fmt.Errorf("indexed properties for %s do not contain %q, got: %v", objectID, key, props)
	}
	return nil
}

func (dc *domainContext) theIndexedPropertiesForShouldNotContain(objectID, key string) error {
	props, err := dc.getIndexedProperties(objectID)
	if err != nil {
		return err
	}
	if _, ok := props[key]; ok {
		return fmt.Errorf("indexed properties for %s should not contain %q, got: %v", objectID, key, props)
	}
	return nil
}

func (dc *domainContext) theIndexedPropertiesForShouldBeEmpty(objectID string) error {
	props, err := dc.getIndexedProperties(objectID)
	if err != nil {
		return err
	}
	if len(props) != 0 {
		return fmt.Errorf("expected empty properties for %s, got: %v", objectID, props)
	}
	return nil
}

func (dc *domainContext) theFileShouldStillContainInFrontmatter(relPath, expected string) error {
	fullPath := filepath.Join(dc.vault.ObjectsDir(), relPath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("read file %s: %v", relPath, err)
	}
	if !strings.Contains(string(data), expected) {
		return fmt.Errorf("file %s does not contain %q", relPath, expected)
	}
	return nil
}

func initFilteringSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a type schema "([^"]*)" with properties "([^"]*)"$`, dc.aTypeSchemaWithProperties)
	ctx.Step(`^a raw object file "([^"]*)" with properties:$`, dc.aRawObjectFileWithProperties)
	ctx.Step(`^the indexed properties for "([^"]*)" should contain "([^"]*)"$`, dc.theIndexedPropertiesForShouldContain)
	ctx.Step(`^the indexed properties for "([^"]*)" should not contain "([^"]*)"$`, dc.theIndexedPropertiesForShouldNotContain)
	ctx.Step(`^the indexed properties for "([^"]*)" should be empty$`, dc.theIndexedPropertiesForShouldBeEmpty)
	ctx.Step(`^the file "([^"]*)" should still contain "([^"]*)" in frontmatter$`, dc.theFileShouldStillContainInFrontmatter)
}

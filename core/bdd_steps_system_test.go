package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

// ── System property steps ────────────────────────────────────────────────

func (dc *domainContext) theSystemPropertyRegistryShouldContain(nameList string) error {
	expected := strings.Split(nameList, ", ")
	for i, s := range expected {
		expected[i] = strings.TrimSpace(s)
	}
	got := SystemPropertyNames()
	if len(got) != len(expected) {
		return fmt.Errorf("registry has %d entries, want %d: %v", len(got), len(expected), got)
	}
	for i, name := range expected {
		if got[i] != name {
			return fmt.Errorf("registry[%d] = %q, want %q", i, got[i], name)
		}
	}
	return nil
}

func (dc *domainContext) shouldBeASystemProperty(name string) error {
	if !IsSystemProperty(name) {
		return fmt.Errorf("%q should be a system property", name)
	}
	return nil
}

func (dc *domainContext) shouldNotBeASystemProperty(name string) error {
	if IsSystemProperty(name) {
		return fmt.Errorf("%q should not be a system property", name)
	}
	return nil
}

func (dc *domainContext) aTypeSchemaWithASystemProperty(typeName, propName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - name: %s
    type: datetime
`, typeName, propName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aSharedPropertiesFileWithASystemProperty(propName string) {
	content := fmt.Sprintf(`properties:
  - name: %s
    type: datetime
`, propName)
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(content), 0644)
}

func (dc *domainContext) theObjectShouldHaveATimestamp(propName string) error {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	val, ok := got.Properties[propName]
	if !ok || val == nil || val == "" {
		return fmt.Errorf("expected %q to be set, got %v", propName, val)
	}
	return nil
}

func (dc *domainContext) theObjectTimestampShouldNotHaveChanged(propName string) error {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	val := fmt.Sprintf("%v", got.Properties[propName])
	if dc.createdAtSnapshot == "" {
		return fmt.Errorf("no snapshot for %q", propName)
	}
	if val != dc.createdAtSnapshot {
		return fmt.Errorf("%q changed: was %q, now %q", propName, dc.createdAtSnapshot, val)
	}
	return nil
}

func (dc *domainContext) theObjectTimestampShouldBeRecent(propName string) error {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	val, ok := got.Properties[propName]
	if !ok || val == nil || val == "" {
		return fmt.Errorf("expected %q to be set", propName)
	}
	s := fmt.Sprintf("%v", val)
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("%q value %q is not valid RFC 3339: %v", propName, s, err)
	}
	if time.Since(parsed) > 5*time.Second {
		return fmt.Errorf("%q value %q is not recent (older than 5s)", propName, s)
	}
	return nil
}

func (dc *domainContext) theFrontmatterShouldHaveSystemPropertiesBeforeSchemaProperties() error {
	pairs := [][2]string{{"name", "created_at"}, {"created_at", "updated_at"}, {"updated_at", "title"}}
	for _, p := range pairs {
		if err := dc.theFrontmatterShouldHaveBefore(p[0], p[1]); err != nil {
			return err
		}
	}
	return nil
}

func (dc *domainContext) theIndexedPropertiesForTheObjectShouldContain(propName string) error {
	var propsJSON string
	err := dc.vault.db.QueryRow("SELECT properties FROM objects WHERE id = ?", dc.currentObject.ID).Scan(&propsJSON)
	if err != nil {
		return fmt.Errorf("query error: %v", err)
	}
	var props map[string]any
	if err := json.Unmarshal([]byte(propsJSON), &props); err != nil {
		return fmt.Errorf("unmarshal error: %v", err)
	}
	if _, ok := props[propName]; !ok {
		return fmt.Errorf("indexed properties do not contain %q: %v", propName, props)
	}
	return nil
}

func (dc *domainContext) createRawObjectFile(prefix, frontmatter string) {
	typeName := "book"
	filename := prefix + mustULID()
	objPath := dc.vault.ObjectPath(typeName, filename)
	os.MkdirAll(filepath.Dir(objPath), 0755)
	os.WriteFile(objPath, []byte("---\n"+frontmatter+"---\n"), 0644)
	dc.currentObject = &Object{
		ID:       typeName + "/" + filename,
		Type:     typeName,
		Filename: filename,
	}
}

func (dc *domainContext) aRawObjectFileWithoutTimestampsExists() {
	dc.createRawObjectFile("legacy-book-", "name: legacy-book\ntitle: Legacy\n")
}

func (dc *domainContext) rawObjectFileShouldNotContain(propName string) error {
	data, err := os.ReadFile(dc.vault.ObjectPath(dc.currentObject.Type, dc.currentObject.Filename))
	if err != nil {
		return fmt.Errorf("ReadFile error: %v", err)
	}
	content := string(data)
	if strings.Contains(content, propName+":") {
		return fmt.Errorf("%s was added to existing object:\n%s", propName, content)
	}
	return nil
}

func (dc *domainContext) theRawObjectFileShouldNotHaveTimestampsAdded() error {
	if err := dc.rawObjectFileShouldNotContain("created_at"); err != nil {
		return err
	}
	return dc.rawObjectFileShouldNotContain("updated_at")
}

func (dc *domainContext) theObjectShouldNotHaveProperty(propName string) error {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	if _, ok := got.Properties[propName]; ok {
		return fmt.Errorf("expected object to not have property %q, but it does", propName)
	}
	return nil
}

func (dc *domainContext) aRawObjectFileWithDescriptionExists() {
	dc.createRawObjectFile("desc-raw-book-", "name: desc-raw-book\ndescription: A raw book with description\ntitle: Raw Book\n")
}

func (dc *domainContext) theRawObjectFileShouldNotHaveDescriptionAdded() error {
	return dc.rawObjectFileShouldNotContain("description")
}

func (dc *domainContext) theFrontmatterShouldHaveBefore(first, second string) error {
	data, err := os.ReadFile(dc.vault.ObjectPath(dc.currentObject.Type, dc.currentObject.Filename))
	if err != nil {
		return fmt.Errorf("ReadFile error: %v", err)
	}
	content := string(data)
	firstIdx := strings.Index(content, first+":")
	secondIdx := strings.Index(content, second+":")
	if firstIdx == -1 {
		return fmt.Errorf("%q not found in frontmatter:\n%s", first, content)
	}
	if secondIdx == -1 {
		return fmt.Errorf("%q not found in frontmatter:\n%s", second, content)
	}
	if firstIdx > secondIdx {
		return fmt.Errorf("%q should come before %q in frontmatter", first, second)
	}
	return nil
}

func (dc *domainContext) shouldBeAnImmutableSystemProperty(name string) error {
	if !IsImmutableSystemProperty(name) {
		return fmt.Errorf("%q should be an immutable system property", name)
	}
	return nil
}

func (dc *domainContext) shouldNotBeAnImmutableSystemProperty(name string) error {
	if IsImmutableSystemProperty(name) {
		return fmt.Errorf("%q should not be an immutable system property", name)
	}
	return nil
}

func initSystemSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^the system property registry should contain "([^"]*)"$`, dc.theSystemPropertyRegistryShouldContain)
	ctx.Step(`^"([^"]*)" should be a system property$`, dc.shouldBeASystemProperty)
	ctx.Step(`^"([^"]*)" should not be a system property$`, dc.shouldNotBeASystemProperty)
	ctx.Step(`^a type schema "([^"]*)" with a system property "([^"]*)"$`, dc.aTypeSchemaWithASystemProperty)
	ctx.Step(`^a shared properties file with a system property "([^"]*)"$`, dc.aSharedPropertiesFileWithASystemProperty)
	ctx.Step(`^the object should have an? "([^"]*)" timestamp$`, dc.theObjectShouldHaveATimestamp)
	ctx.Step(`^the object "([^"]*)" should not have changed$`, dc.theObjectTimestampShouldNotHaveChanged)
	ctx.Step(`^the object "([^"]*)" should be recent$`, dc.theObjectTimestampShouldBeRecent)
	ctx.Step(`^the frontmatter should have system properties before schema properties$`, dc.theFrontmatterShouldHaveSystemPropertiesBeforeSchemaProperties)
	ctx.Step(`^the indexed properties for the object should contain "([^"]*)"$`, dc.theIndexedPropertiesForTheObjectShouldContain)
	ctx.Step(`^a raw object file without timestamps exists$`, dc.aRawObjectFileWithoutTimestampsExists)
	ctx.Step(`^the raw object file should not have timestamps added$`, dc.theRawObjectFileShouldNotHaveTimestampsAdded)
	ctx.Step(`^the object should not have property "([^"]*)"$`, dc.theObjectShouldNotHaveProperty)
	ctx.Step(`^a raw object file with description exists$`, dc.aRawObjectFileWithDescriptionExists)
	ctx.Step(`^the raw object file should not have description added$`, dc.theRawObjectFileShouldNotHaveDescriptionAdded)
	ctx.Step(`^the frontmatter should have "([^"]*)" before "([^"]*)"$`, dc.theFrontmatterShouldHaveBefore)
	ctx.Step(`^"([^"]*)" should be an immutable system property$`, dc.shouldBeAnImmutableSystemProperty)
	ctx.Step(`^"([^"]*)" should not be an immutable system property$`, dc.shouldNotBeAnImmutableSystemProperty)
}

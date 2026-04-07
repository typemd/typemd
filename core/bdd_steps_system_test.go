package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

// ── System property steps ────────────────────────────────────────────────

func assertRegistryContains(got []string, nameList, label string) error {
	expected := strings.Split(nameList, ", ")
	for i, s := range expected {
		expected[i] = strings.TrimSpace(s)
	}
	if len(got) != len(expected) {
		return fmt.Errorf("%s has %d entries, want %d: %v", label, len(got), len(expected), got)
	}
	for i, name := range expected {
		if got[i] != name {
			return fmt.Errorf("%s[%d] = %q, want %q", label, i, got[i], name)
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

func (dc *domainContext) aTypeSchemaWithASystemProperty(typeName, propName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - name: %s
    type: datetime
`, typeName, propName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(content))
}

func (dc *domainContext) aSharedPropertiesFileWithASystemProperty(propName string) {
	dir := dc.vault.SharedPropertiesDir()
	os.MkdirAll(dir, 0755)
	content := "type: datetime\n"
	os.WriteFile(filepath.Join(dir, propName+".yaml"), []byte(content), 0644)
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

func (dc *domainContext) shouldNotBeAnImmutableSystemProperty(name string) error {
	if IsImmutableSystemProperty(name) {
		return fmt.Errorf("%q should not be an immutable system property", name)
	}
	return nil
}

func (dc *domainContext) theStoredSystemPropertyRegistryShouldContain(nameList string) error {
	return assertRegistryContains(StoredPropertyNames(), nameList, "stored registry")
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

func (dc *domainContext) aRawObjectFileWithAComputedPropertyExists() {
	dc.createRawObjectFile("computed-book-", "name: computed-book\nobject_type: book\ntitle: Computed Test\n")
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

func (dc *domainContext) iSaveTheRawObject() {
	obj, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		dc.lastErr = err
		return
	}
	dc.currentObject = obj
	dc.lastErr = dc.vault.SaveObject(obj)
}

func (dc *domainContext) iGetPropertyOnTheObject(name string) {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		dc.lastErr = err
		return
	}
	dc.gotPropertyValue, dc.gotPropertyExists = got.GetProperty(name)
}

func (dc *domainContext) thePropertyValueShouldBe(expected string) error {
	val := fmt.Sprintf("%v", dc.gotPropertyValue)
	if val != expected {
		return fmt.Errorf("property value = %q, want %q", val, expected)
	}
	return nil
}

func (dc *domainContext) thePropertyShouldExist() error {
	if !dc.gotPropertyExists {
		return fmt.Errorf("expected property to exist, but it does not")
	}
	return nil
}

func (dc *domainContext) thePropertyShouldNotExist() error {
	if dc.gotPropertyExists {
		return fmt.Errorf("expected property to not exist, but it does (value: %v)", dc.gotPropertyValue)
	}
	return nil
}

func initSystemSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^"([^"]*)" should be a system property$`, dc.shouldBeASystemProperty)
	ctx.Step(`^a type schema "([^"]*)" with a system property "([^"]*)"$`, dc.aTypeSchemaWithASystemProperty)
	ctx.Step(`^a shared properties file with a system property "([^"]*)"$`, dc.aSharedPropertiesFileWithASystemProperty)
	ctx.Step(`^the object should have an? "([^"]*)" timestamp$`, dc.theObjectShouldHaveATimestamp)
	ctx.Step(`^the object "([^"]*)" should not have changed$`, dc.theObjectTimestampShouldNotHaveChanged)
	ctx.Step(`^the object "([^"]*)" should be recent$`, dc.theObjectTimestampShouldBeRecent)
	ctx.Step(`^the object should not have property "([^"]*)"$`, dc.theObjectShouldNotHaveProperty)
	ctx.Step(`^the frontmatter should have "([^"]*)" before "([^"]*)"$`, dc.theFrontmatterShouldHaveBefore)
	ctx.Step(`^the stored system property registry should contain "([^"]*)"$`, dc.theStoredSystemPropertyRegistryShouldContain)
	ctx.Step(`^a raw object file with a computed property exists$`, dc.aRawObjectFileWithAComputedPropertyExists)
	ctx.Step(`^the raw object file should not contain "([^"]*)"$`, dc.rawObjectFileShouldNotContain)
	ctx.Step(`^"([^"]*)" should not be an immutable system property$`, dc.shouldNotBeAnImmutableSystemProperty)
	ctx.Step(`^I save the raw object$`, dc.iSaveTheRawObject)
	ctx.Step(`^I get property "([^"]*)" on the object$`, dc.iGetPropertyOnTheObject)
	ctx.Step(`^the property value should be "([^"]*)"$`, dc.thePropertyValueShouldBe)
	ctx.Step(`^the property should exist$`, dc.thePropertyShouldExist)
	ctx.Step(`^the property should not exist$`, dc.thePropertyShouldNotExist)
}

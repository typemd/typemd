package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

// ── Object steps ────────────────────────────────────────────────────────────

func (dc *domainContext) aVaultIsReady() {
	dc.aVaultIsInitialized()
	// Write common test type schemas before Open (no longer built-in defaults)
	writeCommonTestTypeSchemas(dc.vault)
	if err := dc.vault.Open(); err != nil {
		panic(fmt.Sprintf("vault open failed: %v", err))
	}
}

// writeCommonTestTypeSchemas writes book, person, and note type schemas
// into the vault's types directory. These were previously built-in defaults
// but are now only available as user-defined types.
func writeCommonTestTypeSchemas(v *Vault) {
	mustWrite := func(path string, data []byte) {
		if err := os.WriteFile(path, data, 0644); err != nil {
			panic(fmt.Sprintf("writeCommonTestTypeSchemas: %v", err))
		}
	}
	mustWrite(filepath.Join(v.TypesDir(), "book.yaml"), []byte(`name: book
emoji: "📚"
properties:
  - name: title
    type: string
    emoji: "📖"
  - name: status
    type: select
    emoji: "📋"
    options:
      - value: to-read
      - value: reading
      - value: done
  - name: rating
    type: number
    emoji: "⭐"
`))
	mustWrite(filepath.Join(v.TypesDir(), "person.yaml"), []byte(`name: person
emoji: "👤"
properties:
  - name: role
    type: string
    emoji: "💼"
`))
	mustWrite(filepath.Join(v.TypesDir(), "note.yaml"), []byte(`name: note
emoji: "📝"
properties:
  - name: title
    type: string
    emoji: "🏷️"
`))
}

func (dc *domainContext) iCreateAObjectNamed(typeName, name string) {
	obj, err := dc.vault.NewObject(typeName, name, "")
	dc.lastErr = err
	if err == nil {
		dc.objects[name] = obj
		dc.currentObject = obj
	}
}

func (dc *domainContext) iCreateAnotherObjectNamed(typeName, name string) {
	dc.prevObject = dc.currentObject
	obj, err := dc.vault.NewObject(typeName, name, "")
	dc.lastErr = err
	if err == nil {
		dc.objects[name+"_2"] = obj
		dc.currentObject = obj
	}
}

func (dc *domainContext) theObjectFilenameShouldStartWith(prefix string) error {
	for _, obj := range dc.objects {
		if strings.HasPrefix(obj.Filename, prefix) {
			return nil
		}
	}
	return fmt.Errorf("no object filename starts with %q", prefix)
}

func (dc *domainContext) theObjectFilenameShouldHaveACharacterULIDSuffix(length int) error {
	for _, obj := range dc.objects {
		parts := strings.SplitN(obj.Filename, "-", 2)
		if len(parts) < 2 {
			continue
		}
		// ULID is the last 26 chars
		if len(obj.Filename) >= length {
			ulidPart := obj.Filename[len(obj.Filename)-length:]
			if len(ulidPart) == length {
				return nil
			}
		}
	}
	return fmt.Errorf("no object has a %d-character ULID suffix", length)
}

func (dc *domainContext) theObjectTypeShouldBe(expected string) error {
	for _, obj := range dc.objects {
		if obj.Type == expected {
			return nil
		}
	}
	return fmt.Errorf("no object has type %q", expected)
}

func (dc *domainContext) theObjectFileShouldExistOnDisk() error {
	for _, obj := range dc.objects {
		path := dc.vault.ObjectPath(obj.Type, obj.Filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("expected file %s to exist", path)
		}
	}
	return nil
}

func (dc *domainContext) theTwoObjectsShouldHaveDifferentIDs() error {
	if dc.prevObject == nil {
		return fmt.Errorf("no previous object to compare")
	}
	for _, obj := range dc.objects {
		if obj != dc.prevObject && obj.ID == dc.prevObject.ID {
			return fmt.Errorf("expected different IDs, both got %q", obj.ID)
		}
	}
	return nil
}

func (dc *domainContext) aObjectNamedExists(typeName, name string) {
	obj, err := dc.vault.NewObject(typeName, name, "")
	if err != nil {
		panic(fmt.Sprintf("create object %s/%s failed: %v", typeName, name, err))
	}
	dc.objects[name] = obj
	dc.prevObject = dc.currentObject
	dc.currentObject = obj
	// Snapshot created_at for system property tests
	if val, ok := obj.Properties["created_at"]; ok {
		dc.createdAtSnapshot = fmt.Sprintf("%v", val)
	}
}

func (dc *domainContext) iGetTheObjectByItsID() {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	dc.lastErr = err
	if err == nil {
		dc.retrieved = got
	}
}

func (dc *domainContext) theRetrievedObjectShouldMatchTheCreatedOne() error {
	if dc.retrieved == nil {
		return fmt.Errorf("no retrieved object found")
	}
	if dc.retrieved.ID != dc.currentObject.ID {
		return fmt.Errorf("retrieved ID = %q, want %q", dc.retrieved.ID, dc.currentObject.ID)
	}
	return nil
}

func (dc *domainContext) iSetPropertyToOnTheObject(key, value string) {
	dc.lastErr = dc.vault.SetProperty(dc.currentObject.ID, key, value)
}

func (dc *domainContext) theObjectPropertyShouldBe(key, expected string) error {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	val := fmt.Sprintf("%v", got.Properties[key])
	if val != expected {
		return fmt.Errorf("property %q = %q, want %q", key, val, expected)
	}
	return nil
}

func (dc *domainContext) iUpdateTheObjectBodyTo(body string) {
	dc.currentObject.Body = body
}

func (dc *domainContext) iUpdateTheObjectTitleTo(title string) {
	dc.currentObject.Properties["title"] = title
}

func (dc *domainContext) iSaveTheObject() {
	dc.lastErr = dc.vault.SaveObject(dc.currentObject)
}

func (dc *domainContext) theObjectFileShouldContain(expected string) error {
	data, err := os.ReadFile(dc.vault.ObjectPath(dc.currentObject.Type, dc.currentObject.Filename))
	if err != nil {
		return fmt.Errorf("ReadFile error: %v", err)
	}
	if !strings.Contains(string(data), expected) {
		return fmt.Errorf("file does not contain %q", expected)
	}
	return nil
}

func (dc *domainContext) gettingTheObjectByIDShouldReturnBody(expected string) error {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	if got.Body != expected {
		return fmt.Errorf("body = %q, want %q", got.Body, expected)
	}
	return nil
}

// ── Object with property/body setup steps ───────────────────────────────────

func (dc *domainContext) aObjectNamedExistsWithPropertySetTo(typeName, name, prop, value string) {
	dc.aObjectNamedExists(typeName, name)
	if err := dc.vault.SetProperty(dc.currentObject.ID, prop, value); err != nil {
		panic(fmt.Sprintf("SetProperty failed: %v", err))
	}
}

func (dc *domainContext) aObjectNamedExistsWithBody(typeName, name, body string) {
	dc.aObjectNamedExists(typeName, name)
	dc.currentObject.Body = body
	if err := dc.vault.SaveObject(dc.currentObject); err != nil {
		panic(fmt.Sprintf("saveObjectFile failed: %v", err))
	}
}

func initObjectSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a vault is ready$`, dc.aVaultIsReady)
	ctx.Step(`^I create a "([^"]*)" object named "([^"]*)"$`, dc.iCreateAObjectNamed)
	ctx.Step(`^I create another "([^"]*)" object named "([^"]*)"$`, dc.iCreateAnotherObjectNamed)
	ctx.Step(`^the object filename should start with "([^"]*)"$`, dc.theObjectFilenameShouldStartWith)
	ctx.Step(`^the object filename should have a (\d+)-character ULID suffix$`, dc.theObjectFilenameShouldHaveACharacterULIDSuffix)
	ctx.Step(`^the object type should be "([^"]*)"$`, dc.theObjectTypeShouldBe)
	ctx.Step(`^the object file should exist on disk$`, dc.theObjectFileShouldExistOnDisk)
	ctx.Step(`^the two objects should have different IDs$`, dc.theTwoObjectsShouldHaveDifferentIDs)
	ctx.Step(`^a "([^"]*)" object named "([^"]*)" exists$`, dc.aObjectNamedExists)
	ctx.Step(`^I get the object by its ID$`, dc.iGetTheObjectByItsID)
	ctx.Step(`^the retrieved object should match the created one$`, dc.theRetrievedObjectShouldMatchTheCreatedOne)
	ctx.Step(`^I set property "([^"]*)" to "([^"]*)" on the object$`, dc.iSetPropertyToOnTheObject)
	ctx.Step(`^the object property "([^"]*)" should be "([^"]*)"$`, dc.theObjectPropertyShouldBe)
	ctx.Step(`^I update the object body to "([^"]*)"$`, dc.iUpdateTheObjectBodyTo)
	ctx.Step(`^I update the object title to "([^"]*)"$`, dc.iUpdateTheObjectTitleTo)
	ctx.Step(`^I save the object$`, dc.iSaveTheObject)
	ctx.Step(`^the object file should contain "([^"]*)"$`, dc.theObjectFileShouldContain)
	ctx.Step(`^getting the object by ID should return body "([^"]*)"$`, dc.gettingTheObjectByIDShouldReturnBody)
	ctx.Step(`^a "([^"]*)" object named "([^"]*)" exists with property "([^"]*)" set to "([^"]*)"$`, dc.aObjectNamedExistsWithPropertySetTo)
	ctx.Step(`^a "([^"]*)" object named "([^"]*)" exists with body "([^"]*)"$`, dc.aObjectNamedExistsWithBody)
}

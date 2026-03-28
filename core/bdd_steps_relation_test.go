package core

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// ── Relation steps ──────────────────────────────────────────────────────────

func (dc *domainContext) aVaultIsReadyWithRelationSchemas() {
	dc.aVaultIsReady()

	bookSchema := []byte(`name: book
properties:
  - name: title
    type: string
  - name: author
    type: relation
    target: person
    bidirectional: true
    inverse: books
`)
	mustWriteTypeSchema(dc.vault, "book", bookSchema)

	personSchema := []byte(`name: person
properties:
  - name: name
    type: string
  - name: books
    type: relation
    target: book
    multiple: true
    bidirectional: true
    inverse: author
`)
	mustWriteTypeSchema(dc.vault, "person", personSchema)
}

func (dc *domainContext) iLinkToVia(sourceName, targetName, relation string) {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		dc.lastErr = fmt.Errorf("object %q or %q not found", sourceName, targetName)
		return
	}
	dc.lastErr = dc.vault.LinkObjects(source.ID, relation, target.ID)
}

func (dc *domainContext) iLinkTheFirstBookToTheSecondBookVia(relation string) {
	if dc.prevObject == nil || dc.currentObject == nil {
		dc.lastErr = fmt.Errorf("need at least 2 objects")
		return
	}
	dc.lastErr = dc.vault.LinkObjects(dc.prevObject.ID, relation, dc.currentObject.ID)
}

func (dc *domainContext) thePropertyOfShouldReference(prop, ownerName, targetName string) error {
	owner := dc.objects[ownerName]
	target := dc.objects[targetName]
	if owner == nil || target == nil {
		return fmt.Errorf("object %q or %q not found", ownerName, targetName)
	}
	obj, err := dc.vault.GetObject(owner.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	if obj.Properties[prop] != target.ID {
		return fmt.Errorf("%s.%s = %v, want %q", ownerName, prop, obj.Properties[prop], target.ID)
	}
	return nil
}

func (dc *domainContext) thePropertyOfShouldContain(prop, ownerName, targetName string) error {
	owner := dc.objects[ownerName]
	target := dc.objects[targetName]
	if owner == nil || target == nil {
		return fmt.Errorf("object %q or %q not found", ownerName, targetName)
	}
	obj, err := dc.vault.GetObject(owner.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	items, ok := obj.Properties[prop].([]any)
	if !ok {
		return fmt.Errorf("%s.%s type = %T, want []any", ownerName, prop, obj.Properties[prop])
	}
	for _, item := range items {
		if item == target.ID {
			return nil
		}
	}
	return fmt.Errorf("%s.%s does not contain %q", ownerName, prop, target.ID)
}

func (dc *domainContext) thePropertyOfShouldBeEmpty(prop, ownerName string) error {
	owner := dc.objects[ownerName]
	if owner == nil {
		return fmt.Errorf("object %q not found", ownerName)
	}
	obj, err := dc.vault.GetObject(owner.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	val := obj.Properties[prop]
	if val != nil {
		return fmt.Errorf("%s.%s = %v, want nil", ownerName, prop, val)
	}
	return nil
}

func (dc *domainContext) unlinkObjects(sourceName, targetName, relation string, both bool) {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		dc.lastErr = fmt.Errorf("object %q or %q not found", sourceName, targetName)
		return
	}
	dc.lastErr = dc.vault.UnlinkObjects(source.ID, relation, target.ID, both)
}

func (dc *domainContext) iUnlinkFromViaWithBothFlag(sourceName, targetName, relation string) {
	dc.unlinkObjects(sourceName, targetName, relation, true)
}

func (dc *domainContext) iUnlinkFromViaWithoutBothFlag(sourceName, targetName, relation string) {
	dc.unlinkObjects(sourceName, targetName, relation, false)
}

func (dc *domainContext) listingRelationsForShouldReturnNEntries(name string, expected int) error {
	obj := dc.objects[name]
	if obj == nil {
		return fmt.Errorf("object %q not found", name)
	}
	rels, err := dc.vault.ListRelations(obj.ID)
	if err != nil {
		return fmt.Errorf("ListRelations error: %v", err)
	}
	if len(rels) != expected {
		return fmt.Errorf("relations count = %d, want %d", len(rels), expected)
	}
	return nil
}

func (dc *domainContext) aTypeSchemaWithARelationPropertyMissingTarget(typeName string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: ref\n    type: relation\n", typeName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func (dc *domainContext) iLinkToANonExistentObjectVia(sourceName, relation string) {
	source := dc.objects[sourceName]
	if source == nil {
		dc.lastErr = fmt.Errorf("object %q not found", sourceName)
		return
	}
	dc.lastErr = dc.vault.LinkObjects(source.ID, relation, "person/nonexistent-01jjjjjjjjjjjjjjjjjjjjjjjj")
}

func (dc *domainContext) thePropertyOfShouldHaveNItems(prop, ownerName string, expected int) error {
	owner := dc.objects[ownerName]
	if owner == nil {
		return fmt.Errorf("object %q not found", ownerName)
	}
	obj, err := dc.vault.GetObject(owner.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	items, ok := obj.Properties[prop].([]any)
	if !ok {
		return fmt.Errorf("%s.%s type = %T, want []any", ownerName, prop, obj.Properties[prop])
	}
	if len(items) != expected {
		return fmt.Errorf("%s.%s has %d items, want %d", ownerName, prop, len(items), expected)
	}
	return nil
}

func (dc *domainContext) aTypeSchemaWithBidirectionalRelationToMissingInverse(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: reviewer
    type: relation
    target: person
    bidirectional: true
    inverse: reviewed_articles
`, typeName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func (dc *domainContext) iBuildDisplayPropertiesFor(name string) {
	obj := dc.objects[name]
	if obj == nil {
		dc.lastErr = fmt.Errorf("object %q not found", name)
		return
	}
	// Re-read the object to get latest properties
	freshObj, err := dc.vault.GetObject(obj.ID)
	if err != nil {
		dc.lastErr = err
		return
	}
	props, err := dc.vault.BuildDisplayProperties(freshObj)
	dc.lastErr = err
	dc.displayProps = props
}

func (dc *domainContext) theDisplayPropertiesShouldContainAReverseRelationWithIndicator(indicator string) error {
	if dc.displayProps == nil {
		return fmt.Errorf("no display properties built")
	}
	for _, dp := range dc.displayProps {
		if dp.IsReverse {
			formatted := dp.Format()
			if strings.Contains(formatted, indicator) {
				return nil
			}
		}
	}
	return fmt.Errorf("no reverse relation with indicator %q found in display properties", indicator)
}

func initRelationSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a vault is ready with relation schemas$`, dc.aVaultIsReadyWithRelationSchemas)
	ctx.Step(`^I link "([^"]*)" to "([^"]*)" via "([^"]*)"$`, dc.iLinkToVia)
	ctx.Step(`^I link the first book to the second book via "([^"]*)"$`, dc.iLinkTheFirstBookToTheSecondBookVia)
	ctx.Step(`^the "([^"]*)" property of "([^"]*)" should reference "([^"]*)"$`, dc.thePropertyOfShouldReference)
	ctx.Step(`^the "([^"]*)" property of "([^"]*)" should contain "([^"]*)"$`, dc.thePropertyOfShouldContain)
	ctx.Step(`^the "([^"]*)" property of "([^"]*)" should be empty$`, dc.thePropertyOfShouldBeEmpty)
	ctx.Step(`^the "([^"]*)" property of "([^"]*)" should have (\d+) items$`, dc.thePropertyOfShouldHaveNItems)
	ctx.Step(`^I unlink "([^"]*)" from "([^"]*)" via "([^"]*)" with both flag$`, dc.iUnlinkFromViaWithBothFlag)
	ctx.Step(`^I unlink "([^"]*)" from "([^"]*)" via "([^"]*)" without both flag$`, dc.iUnlinkFromViaWithoutBothFlag)
	ctx.Step(`^listing relations for "([^"]*)" should return (\d+) entries$`, dc.listingRelationsForShouldReturnNEntries)
	ctx.Step(`^a type schema "([^"]*)" with a relation property missing target$`, dc.aTypeSchemaWithARelationPropertyMissingTarget)
	ctx.Step(`^I link "([^"]*)" to a non-existent object via "([^"]*)"$`, dc.iLinkToANonExistentObjectVia)
	ctx.Step(`^a type schema "([^"]*)" with bidirectional relation to missing inverse$`, dc.aTypeSchemaWithBidirectionalRelationToMissingInverse)
	ctx.Step(`^I build display properties for "([^"]*)"$`, dc.iBuildDisplayPropertiesFor)
	ctx.Step(`^the display properties should contain a reverse relation with indicator "([^"]*)"$`, dc.theDisplayPropertiesShouldContainAReverseRelationWithIndicator)
}

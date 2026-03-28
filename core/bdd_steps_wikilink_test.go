package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/cucumber/godog"
)

// ── Wiki-link steps ─────────────────────────────────────────────────────────

func (dc *domainContext) aVaultIsReadyWithNoteSchemas() {
	dc.aVaultIsReady()
	mustWriteTypeSchema(dc.vault, "book",
		[]byte("name: book\nproperties:\n  - name: title\n    type: string\n"))
	mustWriteTypeSchema(dc.vault, "person",
		[]byte("name: person\nproperties:\n  - name: name\n    type: string\n"))
	mustWriteTypeSchema(dc.vault, "note",
		[]byte("name: note\nproperties:\n  - name: title\n    type: string\n"))
}

func (dc *domainContext) bodyContainsAWikiLinkTo(sourceName, targetName string) {
	source := dc.objects[sourceName]
	if source == nil {
		panic(fmt.Sprintf("source object %q not found", sourceName))
	}
	// If target is a known object slug, use its ID; otherwise treat as raw ID
	targetID := targetName
	if target, ok := dc.objects[targetName]; ok {
		targetID = target.ID
	}
	body := fmt.Sprintf("---\ntitle: %s\n---\n\nSee [[%s]].\n", sourceName, targetID)
	os.WriteFile(dc.vault.ObjectPath(source.Type, source.Filename), []byte(body), 0644)
}

func (dc *domainContext) bodyContainsAWikiLinkToWithDisplayText(sourceName, targetName, displayText string) {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		panic(fmt.Sprintf("object %q or %q not found", sourceName, targetName))
	}
	body := fmt.Sprintf("---\ntitle: %s\n---\n\nBy [[%s|%s]].\n", sourceName, target.ID, displayText)
	os.WriteFile(dc.vault.ObjectPath(source.Type, source.Filename), []byte(body), 0644)
}

func (dc *domainContext) iSyncTheIndex() {
	events, _, err := dc.vault.Reconcile()
	if err == nil {
		err = dc.vault.Project(events)
	}
	dc.lastErr = err
}

func (dc *domainContext) shouldHaveNWikiLinks(name string, expected int) error {
	obj := dc.objects[name]
	if obj == nil {
		return fmt.Errorf("object %q not found", name)
	}
	links, err := dc.vault.ListWikiLinks(obj.ID)
	if err != nil {
		return fmt.Errorf("ListWikiLinks error: %v", err)
	}
	dc.wikiLinks = links
	if len(links) != expected {
		return fmt.Errorf("wiki-links = %d, want %d", len(links), expected)
	}
	return nil
}

func (dc *domainContext) theWikiLinkTargetShouldBe(targetName string) error {
	target := dc.objects[targetName]
	if target == nil {
		return fmt.Errorf("object %q not found", targetName)
	}
	if len(dc.wikiLinks) == 0 {
		return fmt.Errorf("no wiki-links to check")
	}
	if dc.wikiLinks[0].ToID != target.ID {
		return fmt.Errorf("wiki-link ToID = %q, want %q", dc.wikiLinks[0].ToID, target.ID)
	}
	return nil
}

func (dc *domainContext) shouldHaveNBacklinksFrom(targetName string, expected int, sourceName string) error {
	target := dc.objects[targetName]
	if target == nil {
		return fmt.Errorf("object %q not found", targetName)
	}
	backlinks, err := dc.vault.ListBacklinks(target.ID)
	if err != nil {
		return fmt.Errorf("ListBacklinks error: %v", err)
	}
	if len(backlinks) != expected {
		return fmt.Errorf("backlinks = %d, want %d", len(backlinks), expected)
	}
	source := dc.objects[sourceName]
	if source == nil {
		return fmt.Errorf("source object %q not found", sourceName)
	}
	found := false
	for _, bl := range backlinks {
		if bl.FromID == source.ID {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("no backlink from %q found", sourceName)
	}
	return nil
}

func (dc *domainContext) theWikiLinkShouldHaveAnEmptyResolvedID() error {
	if len(dc.wikiLinks) == 0 {
		return fmt.Errorf("no wiki-links to check")
	}
	if dc.wikiLinks[0].ToID != "" {
		return fmt.Errorf("ToID = %q, want empty", dc.wikiLinks[0].ToID)
	}
	return nil
}

func (dc *domainContext) iChangeWikiLinkTo(sourceName, targetName string) {
	dc.bodyContainsAWikiLinkTo(sourceName, targetName)
}

func (dc *domainContext) wikiLinkShouldPointTo(sourceName, targetName string) error {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		return fmt.Errorf("object %q or %q not found", sourceName, targetName)
	}
	links, err := dc.vault.ListWikiLinks(source.ID)
	if err != nil {
		return fmt.Errorf("ListWikiLinks error: %v", err)
	}
	if len(links) != 1 {
		return fmt.Errorf("wiki-links = %d, want 1", len(links))
	}
	if links[0].ToID != target.ID {
		return fmt.Errorf("wiki-link ToID = %q, want %q", links[0].ToID, target.ID)
	}
	return nil
}

func (dc *domainContext) shouldHaveNBacklinks(targetName string, expected int) error {
	target := dc.objects[targetName]
	if target == nil {
		return fmt.Errorf("object %q not found", targetName)
	}
	backlinks, err := dc.vault.ListBacklinks(target.ID)
	if err != nil {
		return fmt.Errorf("ListBacklinks error: %v", err)
	}
	if len(backlinks) != expected {
		return fmt.Errorf("backlinks = %d, want %d", len(backlinks), expected)
	}
	return nil
}

func (dc *domainContext) theWikiLinkDisplayTextShouldBe(expected string) error {
	if len(dc.wikiLinks) == 0 {
		return fmt.Errorf("no wiki-links to check")
	}
	if dc.wikiLinks[0].DisplayText != expected {
		return fmt.Errorf("DisplayText = %q, want %q", dc.wikiLinks[0].DisplayText, expected)
	}
	return nil
}

func (dc *domainContext) bodyContainsDuplicateWikiLinksTo(sourceName, targetName string) {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		panic(fmt.Sprintf("object %q or %q not found", sourceName, targetName))
	}
	body := fmt.Sprintf("---\ntitle: %s\n---\n\nFirst [[%s]] and second [[%s]].\n", sourceName, target.ID, target.ID)
	os.WriteFile(dc.vault.ObjectPath(source.Type, source.Filename), []byte(body), 0644)
}

func (dc *domainContext) iDeleteTheObjectFromDisk(name string) {
	obj := dc.objects[name]
	if obj == nil {
		panic(fmt.Sprintf("object %q not found", name))
	}
	os.Remove(dc.vault.ObjectPath(obj.Type, obj.Filename))
}

func (dc *domainContext) shouldHaveADisplayPropertyFrom(targetName, propKey, sourceName string) error {
	target := dc.objects[targetName]
	source := dc.objects[sourceName]
	if target == nil || source == nil {
		return fmt.Errorf("object %q or %q not found", targetName, sourceName)
	}
	target, err := dc.vault.GetObject(target.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	props, err := dc.vault.BuildDisplayProperties(target)
	if err != nil {
		return fmt.Errorf("BuildDisplayProperties error: %v", err)
	}
	for _, p := range props {
		if p.Key == propKey && p.IsBacklink && p.FromID == source.ID {
			return nil
		}
	}
	return fmt.Errorf("no %q display property from %q found in %+v", propKey, sourceName, props)
}

func (dc *domainContext) iRenderTheBodyOf(name string) {
	obj := dc.objects[name]
	if obj == nil {
		panic(fmt.Sprintf("object %q not found", name))
	}
	// Re-read the object from disk to get the latest body
	obj, err := dc.vault.GetObject(obj.ID)
	if err != nil {
		panic(fmt.Sprintf("GetObject %q failed: %v", name, err))
	}
	dc.renderedBody = RenderWikiLinks(obj.Body)
}

func (dc *domainContext) theRenderedBodyShouldContain(expected string) error {
	if !strings.Contains(dc.renderedBody, expected) {
		return fmt.Errorf("rendered body %q does not contain %q", dc.renderedBody, expected)
	}
	return nil
}

func (dc *domainContext) theRenderedBodyShouldNotContain(unexpected string) error {
	if strings.Contains(dc.renderedBody, unexpected) {
		return fmt.Errorf("rendered body %q should not contain %q", dc.renderedBody, unexpected)
	}
	return nil
}

// ── Shorthand wiki-link steps ───────────────────────────────────────────────

func (dc *domainContext) anotherObjectNamedExists(typeName, name string) {
	obj, err := dc.vault.NewObject(typeName, name, "")
	if err != nil {
		panic(fmt.Sprintf("create another object %s/%s failed: %v", typeName, name, err))
	}
	dc.objects[name+"_2"] = obj
}

func (dc *domainContext) bodyContainsAShorthandWikiLink(sourceName, target string) {
	source := dc.objects[sourceName]
	if source == nil {
		panic(fmt.Sprintf("source object %q not found", sourceName))
	}
	body := fmt.Sprintf("---\nname: %s\n---\n\nSee [[%s]].\n", sourceName, target)
	os.WriteFile(dc.vault.ObjectPath(source.Type, source.Filename), []byte(body), 0644)
}

func (dc *domainContext) bodyContainsAShorthandWikiLinkWithDisplayText(sourceName, target, displayText string) {
	source := dc.objects[sourceName]
	if source == nil {
		panic(fmt.Sprintf("source object %q not found", sourceName))
	}
	body := fmt.Sprintf("---\nname: %s\n---\n\nSee [[%s|%s]].\n", sourceName, target, displayText)
	os.WriteFile(dc.vault.ObjectPath(source.Type, source.Filename), []byte(body), 0644)
}

func (dc *domainContext) theWikiLinkShouldResolveTo(targetName string) error {
	target := dc.objects[targetName]
	if target == nil {
		return fmt.Errorf("object %q not found", targetName)
	}
	if len(dc.wikiLinks) == 0 {
		return fmt.Errorf("no wiki-links to check")
	}
	if dc.wikiLinks[0].ToID != target.ID {
		return fmt.Errorf("wiki-link ToID = %q, want %q", dc.wikiLinks[0].ToID, target.ID)
	}
	return nil
}

func (dc *domainContext) bodyOnDiskShouldContainTheFullIDOf(sourceName, targetName string) error {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		return fmt.Errorf("object %q or %q not found", sourceName, targetName)
	}
	data, err := os.ReadFile(dc.vault.ObjectPath(source.Type, source.Filename))
	if err != nil {
		return fmt.Errorf("read file: %v", err)
	}
	expected := "[[" + target.ID + "]]"
	if !strings.Contains(string(data), expected) {
		return fmt.Errorf("body does not contain %q, got:\n%s", expected, string(data))
	}
	return nil
}

func (dc *domainContext) bodyOnDiskShouldContainTheFullIDOfWithDisplayText(sourceName, targetName, displayText string) error {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		return fmt.Errorf("object %q or %q not found", sourceName, targetName)
	}
	data, err := os.ReadFile(dc.vault.ObjectPath(source.Type, source.Filename))
	if err != nil {
		return fmt.Errorf("read file: %v", err)
	}
	expected := "[[" + target.ID + "|" + displayText + "]]"
	if !strings.Contains(string(data), expected) {
		return fmt.Errorf("body does not contain %q, got:\n%s", expected, string(data))
	}
	return nil
}

func (dc *domainContext) bodyOnDiskShouldContain(sourceName, expected string) error {
	source := dc.objects[sourceName]
	if source == nil {
		return fmt.Errorf("object %q not found", sourceName)
	}
	data, err := os.ReadFile(dc.vault.ObjectPath(source.Type, source.Filename))
	if err != nil {
		return fmt.Errorf("read file: %v", err)
	}
	if !strings.Contains(string(data), expected) {
		return fmt.Errorf("body does not contain %q, got:\n%s", expected, string(data))
	}
	return nil
}

func (dc *domainContext) iSyncTheIndexAndCaptureTheResult() {
	events, result, err := dc.vault.Reconcile()
	if err == nil {
		err = dc.vault.Project(events)
	}
	dc.lastErr = err
	dc.syncResult = result
}

func (dc *domainContext) theSyncResultShouldHaveNWikiLinksExpanded(expected int) error {
	if dc.syncResult == nil {
		return fmt.Errorf("no sync result captured")
	}
	if dc.syncResult.WikiLinksExpanded != expected {
		return fmt.Errorf("WikiLinksExpanded = %d, want %d", dc.syncResult.WikiLinksExpanded, expected)
	}
	return nil
}

func (dc *domainContext) theSyncResultShouldHaveNUnresolvedWikiLinks(expected int) error {
	if dc.syncResult == nil {
		return fmt.Errorf("no sync result captured")
	}
	if len(dc.syncResult.UnresolvedWikiLinks) != expected {
		return fmt.Errorf("UnresolvedWikiLinks = %d, want %d", len(dc.syncResult.UnresolvedWikiLinks), expected)
	}
	return nil
}

func (dc *domainContext) theUnresolvedWikiLinkReasonShouldBe(expected string) error {
	if dc.syncResult == nil {
		return fmt.Errorf("no sync result captured")
	}
	if len(dc.syncResult.UnresolvedWikiLinks) == 0 {
		return fmt.Errorf("no unresolved wiki-links")
	}
	if dc.syncResult.UnresolvedWikiLinks[0].Reason != expected {
		return fmt.Errorf("reason = %q, want %q", dc.syncResult.UnresolvedWikiLinks[0].Reason, expected)
	}
	return nil
}

func initWikiLinkSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a vault is ready with note schemas$`, dc.aVaultIsReadyWithNoteSchemas)
	ctx.Step(`^"([^"]*)" body contains a wiki-link to "([^"]*)"$`, dc.bodyContainsAWikiLinkTo)
	ctx.Step(`^"([^"]*)" body contains a wiki-link to "([^"]*)" with display text "([^"]*)"$`, dc.bodyContainsAWikiLinkToWithDisplayText)
	ctx.Step(`^I sync the index$`, dc.iSyncTheIndex)
	ctx.Step(`^"([^"]*)" should have (\d+) wiki-links?$`, dc.shouldHaveNWikiLinks)
	ctx.Step(`^the wiki-link target should be "([^"]*)"$`, dc.theWikiLinkTargetShouldBe)
	ctx.Step(`^"([^"]*)" should have (\d+) backlinks? from "([^"]*)"$`, dc.shouldHaveNBacklinksFrom)
	ctx.Step(`^the wiki-link should have an empty resolved ID$`, dc.theWikiLinkShouldHaveAnEmptyResolvedID)
	ctx.Step(`^I change "([^"]*)" wiki-link to "([^"]*)"$`, dc.iChangeWikiLinkTo)
	ctx.Step(`^"([^"]*)" wiki-link should point to "([^"]*)"$`, dc.wikiLinkShouldPointTo)
	ctx.Step(`^"([^"]*)" should have (\d+) backlinks$`, dc.shouldHaveNBacklinks)
	ctx.Step(`^the wiki-link display text should be "([^"]*)"$`, dc.theWikiLinkDisplayTextShouldBe)
	ctx.Step(`^"([^"]*)" body contains duplicate wiki-links to "([^"]*)"$`, dc.bodyContainsDuplicateWikiLinksTo)
	ctx.Step(`^I delete the object "([^"]*)" from disk$`, dc.iDeleteTheObjectFromDisk)
	ctx.Step(`^"([^"]*)" should have a "([^"]*)" display property from "([^"]*)"$`, dc.shouldHaveADisplayPropertyFrom)
	ctx.Step(`^I render the body of "([^"]*)"$`, dc.iRenderTheBodyOf)
	ctx.Step(`^the rendered body should contain "([^"]*)"$`, dc.theRenderedBodyShouldContain)
	ctx.Step(`^the rendered body should not contain "([^"]*)"$`, dc.theRenderedBodyShouldNotContain)

	// Shorthand wiki-link steps
	ctx.Step(`^another "([^"]*)" object named "([^"]*)" exists$`, dc.anotherObjectNamedExists)
	ctx.Step(`^"([^"]*)" body contains a shorthand wiki-link "([^"]*)"$`, dc.bodyContainsAShorthandWikiLink)
	ctx.Step(`^"([^"]*)" body contains a shorthand wiki-link "([^"]*)" with display text "([^"]*)"$`, dc.bodyContainsAShorthandWikiLinkWithDisplayText)
	ctx.Step(`^the wiki-link should resolve to "([^"]*)"$`, dc.theWikiLinkShouldResolveTo)
	ctx.Step(`^"([^"]*)" body on disk should contain the full ID of "([^"]*)"$`, dc.bodyOnDiskShouldContainTheFullIDOf)
	ctx.Step(`^"([^"]*)" body on disk should contain the full ID of "([^"]*)" with display text "([^"]*)"$`, dc.bodyOnDiskShouldContainTheFullIDOfWithDisplayText)
	ctx.Step(`^"([^"]*)" body on disk should contain "([^"]*)"$`, dc.bodyOnDiskShouldContain)
	ctx.Step(`^I sync the index and capture the result$`, dc.iSyncTheIndexAndCaptureTheResult)
	ctx.Step(`^the sync result should have (\d+) wiki-links? expanded$`, dc.theSyncResultShouldHaveNWikiLinksExpanded)
	ctx.Step(`^the sync result should have (\d+) unresolved wiki-links?$`, dc.theSyncResultShouldHaveNUnresolvedWikiLinks)
	ctx.Step(`^the unresolved wiki-link reason should be "([^"]*)"$`, dc.theUnresolvedWikiLinkReasonShouldBe)
}

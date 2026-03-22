## 1. Core: Shorthand Wiki-Link Resolution

- [x] 1.1 Write BDD scenarios for shorthand wiki-link resolution (type-qualified, same-type, full ID unchanged, ambiguous, not found)
- [x] 1.2 Implement BDD step definitions for shorthand resolution scenarios
- [x] 1.3 Add `resolveWikiLinkTarget()` function to classify and resolve wiki-link targets using nameIndex
- [x] 1.4 Refactor `syncWikiLinks()` to use `resolveWikiLinkTarget()` instead of plain diskIDs lookup
- [x] 1.5 Add unit tests for `resolveWikiLinkTarget()` edge cases (empty target, target with display text unaffected, slug vs name property matching)

## 2. Core: Write-Back Expanded Wiki-Links

- [x] 2.1 Write BDD scenarios for wiki-link write-back (shorthand expanded in body, display text preserved, unresolved not modified, already full ID not modified)
- [x] 2.2 Implement BDD step definitions for write-back scenarios
- [x] 2.3 Add body replacement logic to expand shorthand targets in-place (`expandWikiLinksInBody()`)
- [x] 2.4 Integrate write-back into `syncWikiLinksAndTags()` — track modified bodies and save via `repo.Save()`
- [x] 2.5 Add unit tests for `expandWikiLinksInBody()` edge cases (multiple links in one body, same target twice, link inside code block)

## 3. Core: SyncResult Extension

- [x] 3.1 Write BDD scenarios for sync result reporting (expanded count, unresolved wiki-links with reason/matches)
- [x] 3.2 Implement BDD step definitions for sync result scenarios
- [x] 3.3 Add `WikiLinksExpanded` and `UnresolvedWikiLinks` fields to `SyncResult`
- [x] 3.4 Populate sync result fields during wiki-link resolution
- [x] 3.5 Add toast notification integration for unresolved wiki-links in TUI

## 4. Core: Validate Wiki-Links with Shorthand Support

- [x] 4.1 Write BDD scenarios for validation (ambiguous shorthand reported with suggestions, resolvable shorthand passes)
- [x] 4.2 Implement BDD step definitions for validation scenarios
- [x] 4.3 Extend `ValidateWikiLinks()` to build name index and resolve shorthand targets
- [x] 4.4 Add ambiguity reporting with matching full IDs as suggestions
- [x] 4.5 Add unit tests for validation edge cases

## 5. CLI: tmd fix wikilinks Command

- [x] 5.1 Write BDD scenarios for `tmd fix wikilinks` (expand shorthand, report summary, ambiguous reported, no changes needed)
- [x] 5.2 Implement BDD step definitions for fix wikilinks scenarios
- [x] 5.3 Add `FixWikiLinks()` method to Vault that walks objects, resolves, and writes back
- [x] 5.4 Add `tmd fix wikilinks` Cobra command wired to `FixWikiLinks()`
- [x] 5.5 Add unit tests for fix command output formatting

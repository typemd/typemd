## Why

Wiki-links currently require the full object ID including ULID (e.g., `[[book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz]]`), which is cumbersome to type and hard to remember. Supporting shorthand forms (`[[type/name]]` and `[[name]]`) makes wiki-links practical for human authoring while maintaining precise resolution at sync time.

## What Changes

- Extend wiki-link resolution in `syncWikiLinks()` to support three formats: full ID (`[[type/name-ulid]]`), type-qualified name (`[[type/name]]`), and same-type shorthand (`[[name]]`)
- Write-back resolved full IDs to source files after successful resolution (like relation name expansion)
- Extend `ValidateWikiLinks()` to detect ambiguous shorthand links and suggest full IDs
- Add `tmd fix wikilinks` subcommand to batch-expand all shorthand wiki-links to full IDs
- Report unresolved/ambiguous wiki-links via toast notifications during sync (reusing existing SyncResult patterns)

## Capabilities

### New Capabilities

- `shorthand-wikilink-resolution`: Resolution of shorthand wiki-link syntax (`[[type/name]]` and `[[name]]`) during sync, including write-back and ambiguity handling
- `fix-wikilinks-command`: CLI command `tmd fix wikilinks` to expand all shorthand wiki-links to full IDs

### Modified Capabilities

- `wiki-links`: Validation now detects ambiguous shorthand links and suggests full IDs; broken link reporting includes shorthand context

## Impact

- `core/projector.go` — `syncWikiLinks()` gains name-based resolution and write-back logic
- `core/validate.go` — `ValidateWikiLinks()` extended for shorthand + ambiguity reporting
- `core/name_resolve.go` — Reuse `buildNameIndex()` and `resolveByName()` (already exist for relations)
- `core/sync.go` — `SyncResult` extended with wiki-link expansion/unresolved counts
- `cmd/` — New `tmd fix wikilinks` command
- No DB schema changes — `to_id` continues to store the resolved full ID

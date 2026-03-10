## Why

Object display titles are currently derived from the filename via `DisplayName()` (which strips the ULID suffix). This couples the display name to filesystem constraints — no spaces, no casing flexibility, and no way to rename without changing the file. Introducing `name` as a built-in, required property decouples display from storage and unblocks #186 (name templates).

## What Changes

- Add `name` as an implicit, required system property on all objects
- `name` is always the first key in YAML frontmatter
- On object creation (`NewObject`), `name` is auto-populated from the slug
- TUI and CLI display `name` instead of `DisplayName()` (filename-derived)
- `name` is a reserved word — type schemas must not define a property named `name` (validation error)
- `sync` writes `name` into existing objects that lack it (migration via `DisplayName()` fallback)
- Provide a centralized `GetName()` method to prepare for future #186 name template support

## Capabilities

### New Capabilities
- `system-property`: System-level `name` property — definition, storage, migration, display, and validation rules

### Modified Capabilities
- `type-schema`: Add validation rule rejecting `name` as a user-defined property name
- `tui-layout`: Title panel and list view use `name` instead of `DisplayName()`

## Impact

- **core/**: `Object` struct gains `GetName()`, `NewObject` sets `name`, `Sync` migrates existing objects, `TypeSchema` validation rejects `name` property, `writeFrontmatter` pins `name` first
- **tui/**: List view and detail title panel read `name` instead of `DisplayName()`
- **cmd/**: `new` command passes through (no explicit change needed, handled by core)
- **mcp/**: No change expected (uses core APIs)
- **Breaking**: None — `name` is additive; `DisplayName()` remains available as fallback/utility

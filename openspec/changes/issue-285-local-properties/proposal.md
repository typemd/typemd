## Why

When an object's frontmatter contains properties not defined in the type schema, these "local properties" appear mixed in with schema-defined properties — making it impossible for users to tell which properties are formally part of the type and which are file-local additions. A visual separation improves clarity and prevents confusion about what is schema-managed vs. ad-hoc.

## What Changes

- Add `IsLocal bool` field to `DisplayProperty` in core, marking properties that exist in the object but not in the type schema or system property registry
- TUI property panel: render a visual separator (`── Local Properties ──`) between schema-defined and local properties; local properties are read-only (cursor skips them)
- CLI `tmd object show`: display local properties in a separate section with a "Local Properties" header
- `BuildDisplayProperties()` sets `IsLocal: true` on properties not in schema

## Capabilities

### New Capabilities

- `local-property-display`: Visual identification and separation of properties that exist in object frontmatter but are not defined in the type schema, across core, TUI, and CLI

### Modified Capabilities

- `property-editing`: Local properties are displayed but marked non-editable (read-only) in the TUI property editor — cursor navigation skips them

## Impact

- `core/display.go` — add `IsLocal` field to `DisplayProperty`
- `core/query_service.go` — `BuildDisplayProperties()` marks non-schema, non-system properties as `IsLocal: true`
- `tui/prop_editor.go` — render separator before first local property; mark local items as non-editable
- `cmd/show.go` — separate "Local Properties" section in `tmd object show` output
- No database or index changes (local properties remain excluded from SQLite index per #174)

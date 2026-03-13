## Why

Type names are always singular (e.g., `book`, `person`). When displaying collection headers in the TUI or CLI output, the singular form reads awkwardly (e.g., "▼ 📚 book (3)"). Adding an optional `plural` field to type schemas allows UIs to display grammatically correct plural forms for type group headers.

## What Changes

- Add an optional `plural` field to `TypeSchema` struct and YAML schema
- Add a `PluralName()` method that returns `Plural` if set, otherwise falls back to `Name` (no automatic `s` suffix — this avoids issues with non-English type names like Chinese)
- Update built-in `tag` type to include `plural: "tags"`
- Update TUI group headers to display plural form (e.g., "▼ 📚 books (3)")
- Update CLI `tmd type show` to display the plural name when set

## Capabilities

### New Capabilities

- `plural-display-name`: Type schemas support an optional `plural` field for grammatically correct collection display names

### Modified Capabilities

- `type-schema`: TypeSchema struct gains a new optional `plural` field with `PluralName()` accessor
- `tui-object-list`: TUI group headers use plural form instead of singular type name

## Impact

- `core/type_schema.go` — `TypeSchema` struct, `defaultTypes`, validation, `LoadType`
- `tui/list.go` — group header rendering, `buildGroups` function
- `tui/app.go` — `typeGroup` struct gains `Plural` field
- Existing type schema YAML files are unaffected (field is optional)

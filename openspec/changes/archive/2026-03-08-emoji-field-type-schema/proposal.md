## Why

Types are currently identified only by text names across CLI output, TUI, and other interfaces. Adding an optional emoji field to type schemas enables visual identification, making it faster to scan and distinguish types in compact UI contexts (e.g., TUI lists, search results, MCP tool output).

## What Changes

- Add an optional `emoji` field to the `TypeSchema` struct
- Support `emoji` in type schema YAML files (`.typemd/types/*.yaml`)
- Include `emoji` in built-in default types (book → 📚, person → 👤, note → 📝)
- Expose `emoji` via `tmd type show` and `tmd type list` CLI output
- No uniqueness constraint — different types may use the same emoji

## Capabilities

### New Capabilities

- `type-emoji`: Optional emoji field on type schemas for visual identification in UI contexts

### Modified Capabilities

_(none — no existing specs are affected)_

## Impact

- **core/type_schema.go** — `TypeSchema` struct gains `Emoji` field; default types updated; validation unchanged (emoji is optional, no constraints)
- **cmd/type_show.go** — Display emoji in type detail view
- **cmd/type_list.go** — Display emoji alongside type names
- **examples/** — Update example type YAML files to include emoji
- No breaking changes; existing YAML files without `emoji` continue to work

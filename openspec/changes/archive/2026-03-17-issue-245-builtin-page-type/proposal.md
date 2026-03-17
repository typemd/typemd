## Why

Users who skip starter types during `tmd init` have no type available for general content. The only built-in type is `tag`, which backs the system `tags` property and is not suitable for writing. A general-purpose content type should exist in every vault by default, giving users an immediate starting point for free-form writing without requiring any setup.

## What Changes

- Add `page` as a second built-in type in `defaultTypes` (alongside `tag`) with emoji 📄, plural "pages", unique false, no custom properties
- Add `PageTypeName = "page"` constant for consistent reference
- Protect `page` from deletion (same mechanism as `tag`)
- Update `tmd init` to fallback to `page` as `cli.default_type` when no starter types are selected
- `page` can be overridden by custom `.typemd/types/page.yaml` (existing mechanism, no new code)

## Capabilities

### New Capabilities

- `builtin-type`: Built-in type definitions (tag, page), deletion protection, and listing behavior

### Modified Capabilities

- `type-schema`: Update "Only tag is a built-in type" requirement to include `page`; add requirement for `page` built-in defaults
- `starter-type-templates`: Update `tmd init` default type resolution to fallback to `page`

## Impact

- `core/type_schema.go`: Add page to `defaultTypes` map, add `PageTypeName` constant
- `core/system_property.go`: Add `PageTypeName` constant (alongside `TagTypeName`)
- `cmd/init.go`: Update `resolveDefaultType()` fallback logic
- Existing tests for `defaultTypes` behavior (deletion, listing) may need updating to account for the new built-in type
- Documentation: type schema docs, README

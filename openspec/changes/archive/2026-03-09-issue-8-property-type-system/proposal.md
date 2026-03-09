## Why

typemd's type schema currently supports only 4 property types (string, number, enum, relation). Users need richer data modeling — dates for tracking timelines, URLs for linking resources, checkboxes for boolean states, and multi-select for tagging.

## What Changes

- Add 5 new property types: `date`, `datetime`, `url`, `checkbox`, `multi_select`
- **BREAKING**: Rename `enum` type to `select` for consistency with common knowledge tools
- **BREAKING**: Replace `values` array with `options` object array (`{value, label}`) for select/multi_select
- Add migration support in `tmd migrate` for enum → select conversion
- Update validation to enforce type-specific rules (date format, URL scheme, options membership)

## Capabilities

### New Capabilities
- `predefined-property`: Defines the complete set of supported property types, their validation rules, and schema format. Covers string, number, date, datetime, url, checkbox, select, multi_select, and relation.

### Modified Capabilities
- `type-schema`: The emoji field spec is unchanged, but the Property struct gains new type values and the `options` field replaces `values` for select types.

## Impact

- **core/type_schema.go**: Property struct changes (new types, Options field, remove Values)
- **core/validate.go**: New validation rules per type (date format, URL scheme, checkbox bool, options membership)
- **core/object.go**: Property parsing may need type coercion (YAML auto-parses dates/bools)
- **core/migrate.go**: enum → select migration logic
- **core/display.go**: Type-aware formatting (dates, checkboxes, URLs)
- **cmd/**: No new commands, but existing commands benefit from richer types
- **examples/**: Update example vault schemas to use new types
- **BDD features**: New scenarios for property type validation and migration

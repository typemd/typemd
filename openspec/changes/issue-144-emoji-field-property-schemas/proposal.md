## Why

Properties in type schemas are identified only by text names. In compact UI contexts (TUI properties panel, table columns), property names consume significant horizontal space. An optional emoji field on property definitions would provide a compact visual identifier, enabling UIs to use emojis as space-efficient labels while maintaining readability.

## What Changes

- Add an optional `emoji` field to the `Property` struct in type schemas
- Validate emoji uniqueness within a type's properties scope
- Validate emoji uniqueness within global/predefined properties scope (separate from type scope)
- Reject duplicate emojis during schema validation (`ValidateSchema`)
- Update default type schemas to include example emojis

## Capabilities

### New Capabilities
- `property-emoji`: Optional emoji field on property definitions with per-type uniqueness validation

### Modified Capabilities
- `type-schema`: Property struct gains an `emoji` field; `ValidateSchema` adds duplicate emoji checking

## Impact

- `core/type_schema.go` — `Property` struct, `ValidateSchema()` function
- `.typemd/types/*.yaml` — type schema files gain optional `emoji` field
- `examples/` — example vault schemas updated with sample emojis
- No breaking changes — `emoji` is optional with zero-value default

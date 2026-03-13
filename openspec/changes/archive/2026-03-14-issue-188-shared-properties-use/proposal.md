## Why

Each type schema defines all its properties independently. Common properties like `due_date` or `priority` must be redefined in every type that uses them, leading to duplication and inconsistency. A shared properties system lets users define once, reference everywhere.

## What Changes

- Add `.typemd/properties.yaml` file format for defining shared property definitions
- Add `use: <name>` syntax in type schemas to reference shared properties (only `pin` and `emoji` overrides allowed)
- Extend `LoadType()` to resolve `use` entries by merging the full shared property definition with local overrides
- Add validation rules for shared properties:
  - Shared properties file: duplicate property names → error
  - Type schema defines a `name` property that conflicts with a shared property name → error
  - `use` references a non-existent shared property → error
  - `use` entry includes fields other than `pin` and `emoji` → error
  - Resolved properties have duplicate names (including across `use` and `name` entries) → error

## Capabilities

### New Capabilities
- `shared-properties`: Define shared property definitions in `.typemd/properties.yaml` and reference them from type schemas via `use` keyword with `pin`/`emoji` overrides

### Modified Capabilities
- `type-schema`: `LoadType()` must resolve `use` entries; `ValidateSchema()` must check `use` validation rules

## Impact

- `core/type_schema.go` — Property struct gains `Use` field; LoadType() gains resolution logic; ValidateSchema() gains new rules
- `core/vault.go` — New method to load shared properties file
- `.typemd/properties.yaml` — New file format
- Downstream code (TUI, MCP, CLI) should be unaffected — they receive fully resolved TypeSchema

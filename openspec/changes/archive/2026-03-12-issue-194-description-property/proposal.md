## Why

Every object needs a brief summary for list displays, search results, MCP responses, and Web UI previews. Currently there is no standard place for this — users must rely on the markdown body or type-specific properties. A `description` system property provides a consistent, discoverable field across all object types.

## What Changes

- Add `description` as a stored system property (type: `text`, single-line string)
- `description` is optional — default is empty string or absent
- Frontmatter ordering: `name` → `description` → `created_at` → `updated_at` → schema properties
- Schema validation rejects `description` in type schemas and shared properties (reserved)
- `description` is indexed in SQLite for full-text search (already covered by JSON properties in FTS5)
- Legacy objects without `description` continue to work without modification

## Capabilities

### New Capabilities

_(none — this extends existing system property infrastructure)_

### Modified Capabilities

- `system-property-registry`: Registry expands from 3 to 4 entries; `description` inserted between `name` and `created_at`; `SystemPropertyNames()` returns updated order; `IsSystemProperty("description")` returns true
- `system-properties`: New requirement for `description` behavior — optional text field, no auto-population, no migration during sync, frontmatter ordering updated

## Impact

- **core/system_property.go** — Add `DescriptionProperty` constant and registry entry
- **core/type_schema.go** — `OrderedPropKeys()` automatically handles new registry order (no code change expected)
- **core/type_schema.go** — `ValidateSchema()` automatically rejects `description` via `IsSystemProperty()` (no code change expected)
- **core/shared_properties.go** — `ValidateSharedProperties()` automatically rejects `description` via `IsSystemProperty()` (no code change expected)
- **core/features/system_property.feature** — Update BDD scenarios for new registry order and `description` behavior
- **core/system_property_test.go** — Update unit tests for new registry order and `description` validation
- **No SQLite schema change** — properties stored as JSON, FTS5 already indexes the JSON blob
- **No migration needed** — `description` is optional, legacy objects work as-is

## Why

There is no built-in mechanism to track where knowledge came from. Users who ingest external content (books, articles, videos) have no standard way to record the original source, when it was processed, or link extracted entities back to it. Adding `source` as the third built-in type (alongside `tag` and `page`) establishes a shared vocabulary for provenance tracking.

## What Changes

- Add `source` as a built-in type with emoji 📥, plural "sources", unique false
- Define three custom properties: `url` (text, 🔗), `author` (text, ✍️), `ingested_at` (date, 📅)
- Register `SourceTypeName = "source"` constant alongside existing `TagTypeName` and `PageTypeName`
- Like other built-in types: exists without YAML file, cannot be deleted, can be overridden by custom `types/source/schema.yaml`

## Capabilities

### New Capabilities

- `builtin-source-type`: Registration and behavior of the built-in `source` type — default schema, deletion protection, override support

### Modified Capabilities

*(none — no existing spec-level requirements change)*

## Impact

- `core/system_property.go` — new constant
- `core/type_schema.go` — new entry in `defaultTypes` map
- `core/features/type_crud.feature` — new BDD scenarios
- `core/type_crud_test.go` — new unit tests
- Documentation (CLAUDE.md data model section) — mention `source` alongside `tag` and `page`

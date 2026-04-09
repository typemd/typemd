## Why

Objects can only be referenced by their canonical `name` in wiki-links and search. Users who know an object by an alternative name (a translation, abbreviation, or former title) have no way to link to it without looking up its exact name. Adding `aliases` as a stored system property lets objects be discovered and linked by any of their known names.

## What Changes

- Add `aliases` as a stored system property of type `list[text]` — an optional string array written to frontmatter after `tags`
- Register `aliases` in the system property registry as a user-authored stored property
- Update frontmatter serialization order: `name`, `description`, `created_at`, `updated_at`, `tags`, `aliases`, then schema properties
- Update wiki-link name resolution to match against aliases in addition to the canonical name
- Update the name index (`buildNameIndex`) to include aliases as additional lookup keys
- Reject `aliases` as a property name in type schemas and shared properties (reserved system property)

## Capabilities

### New Capabilities

- `object-aliases`: `aliases` system property — storage, frontmatter serialization, schema validation guard, and wiki-link resolution by alias

### Modified Capabilities

- `wiki-links`: Wiki-link name resolution now also matches against object aliases (existing shorthand resolution behavior extended)

## Impact

- **`core/system_property.go`** — add `aliases` to system property registry
- **`core/object.go`** — frontmatter marshal/unmarshal for `aliases` field
- **`core/reconciler.go` / `core/name_resolve.go`** — alias lookup in name index and wiki-link resolution
- **`core/sqlite_index.go`** — index aliases as additional name keys
- **`core/schema_validate.go`** — reject `aliases` in type schemas and shared properties
- No breaking changes to existing frontmatter (field is optional; absent = empty)

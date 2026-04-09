## Context

typemd objects are identified by their canonical `name` property and slug for wiki-link resolution. The name index (`buildNameIndex`) currently maps slugified name and filename slug → object ID. There is no mechanism for objects to declare alternative names (translations, abbreviations, former titles).

The system property registry (`system_property.go`) defines the fixed set of system properties. Adding a new stored property requires touching: registry, object struct, frontmatter marshal/unmarshal, index building, and validation.

## Goals / Non-Goals

**Goals:**
- Add `aliases` as a `list[text]` stored system property written to frontmatter after `tags`
- Index aliases alongside name/slug for wiki-link shorthand resolution
- Guard `aliases` from being declared in type schemas or shared properties
- Extend wiki-link name resolution to match against aliases

**Non-Goals:**
- Alias uniqueness enforcement across the vault (ambiguous aliases use existing "multiple matches" behavior)
- Search-by-alias in `tmd object list` or query filters (future enhancement)
- TUI alias editing widget (follows standard `list[text]` property editing when implemented)
- Alias sanitization / normalization (stored as-is, same as `name`)

## Decisions

### 1. Store aliases in frontmatter, not a separate file

**Decision:** `aliases` is stored as a YAML list in object frontmatter under the `aliases` key.

**Rationale:** Consistent with all other stored system properties (`tags`, `name`, etc.). Files remain the source of truth. No additional file format or migration needed. Aliases travel with the object on copy/move.

**Alternative considered:** Separate `.typemd/aliases.yaml` registry — rejected because it breaks the "files are source of truth" principle and adds a synchronization surface.

### 2. Position aliases after `tags` in frontmatter

**Decision:** Serialization order: `name`, `description`, `created_at`, `updated_at`, `tags`, `aliases`, then schema properties.

**Rationale:** Groups system properties together in a predictable order. `aliases` is less frequently set than `tags`, so placing it after `tags` keeps the most common properties near the top. Existing objects without `aliases` are unaffected (field is omitted when empty).

### 3. Aliases indexed as additional name-index keys

**Decision:** In `buildNameIndex`, each alias string is slugified and added as an additional lookup key mapping to the same object ID. The existing slug and name-property keys are preserved.

**Rationale:** Reuses the existing name resolution path (`resolveByName`) without a new data structure. Ambiguity handling (multiple matches) is already implemented — aliases that collide with names or other aliases use the same "warn, don't resolve" behavior.

**Alternative considered:** Separate alias index table in SQLite — overkill for the lookup frequency; the in-memory name index rebuilt on each reconcile is sufficient.

### 4. Aliases reserved in schema validation

**Decision:** Schema validation rejects any property named `aliases` in type schemas and shared properties, with the same error message used for other reserved system properties.

**Rationale:** Prevents user-defined properties from shadowing the system `aliases` field. Consistent with existing reserved property guard.

## Risks / Trade-offs

- **Alias collision with object names** → same behavior as duplicate names: both entries exist in the index, resolution picks the exact-name match preferentially; if both are aliases, ambiguous. No new error handling required.
- **Large alias lists** → index rebuild is O(objects × aliases). Acceptable given typical vault sizes; no lazy loading needed.
- **Frontmatter position change on re-save** → existing objects whose frontmatter has `aliases` out of order will be normalized on next save (existing behavior for all system properties).

## Migration Plan

No migration required. `aliases` is optional; existing objects without the field behave identically. The SQLite index is rebuilt on first sync after the update — aliases are indexed at that point.

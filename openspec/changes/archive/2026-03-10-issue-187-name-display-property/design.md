## Context

Object display names are currently derived from the filename via `DisplayName()` which strips the ULID suffix. This couples display to filesystem constraints (no spaces, no casing, immutable after creation). The `title` property exists in some type schemas as a regular data field but has no system-level role.

Key code paths:
- `core/object.go`: `DisplayName()` → `StripULID(Filename)`
- `tui/list.go:123`: list view uses `DisplayName()`
- `tui/detail.go:24,26`: title panel uses `DisplayName()`
- `core/object.go:52-94`: `writeFrontmatter()` with `OrderedPropKeys()`

## Goals / Non-Goals

**Goals:**
- Introduce `name` as a system property stored in frontmatter, decoupled from filename
- Provide `GetName()` as the single access point for display name
- Migrate existing objects seamlessly during sync
- Validate that type schemas don't redefine `name`

**Non-Goals:**
- Name templates (#186) — not implemented, only the `GetName()` interface is prepared
- Renaming files when `name` changes — filename remains stable
- Making `name` searchable/indexable beyond current property storage — future enhancement

## Decisions

### 1. `name` lives in the properties map, not as a struct field

**Decision**: Store `name` in `Object.Properties["name"]` rather than adding a dedicated `Object.Name` field.

**Rationale**: Keeping `name` in the properties map means it flows through existing frontmatter read/write paths naturally. No schema changes to the SQLite `properties` JSON column. When #186 adds name templates, `GetName()` can compute the value without changing storage.

**Alternative rejected**: Dedicated struct field — would require parallel storage, special handling in frontmatter parsing, and database schema changes.

### 2. `GetName()` method with fallback

**Decision**: `GetName()` reads `Properties["name"]`; if missing or empty, falls back to `DisplayName()`.

**Rationale**: Provides forward compatibility — callers use `GetName()` everywhere, and when #186 adds template logic, only `GetName()` needs to change. The fallback ensures no breakage for objects not yet migrated.

### 3. `name` pinned first in frontmatter via `OrderedPropKeys()`

**Decision**: Modify `OrderedPropKeys()` to always emit `name` as the first key, before schema-defined properties.

**Rationale**: Consistent, human-friendly frontmatter. Users always see `name` at the top when editing markdown files.

### 4. Reserved name validation in `ValidateTypeSchema()`

**Decision**: Add a check in `ValidateTypeSchema()` that rejects any property with `name: name`.

**Rationale**: Prevents confusion between user-defined `name` properties and the system `name`. Fail early at schema load time.

### 5. Migration via sync, not a separate command

**Decision**: During `Vault.Sync()`, detect objects missing `name` and add it from `DisplayName()`.

**Rationale**: Sync already reads and writes all objects. Adding migration here means zero extra user steps — just run `tmd sync` or open the TUI (which syncs on start).

## Risks / Trade-offs

- **[Risk] Existing `title` properties may confuse users** → Document that `name` is the display title; `title` remains a data property. Consider deprecating `title` in built-in schemas in a future issue.
- **[Risk] Large vaults may have slow first sync** → Migration is O(n) file writes, but this is a one-time cost and sync is already O(n).
- **[Risk] `name` collides with user data** → Mitigated by schema validation rejecting `name` as a property name. Existing frontmatter `name` keys in objects are fine — they become the system property.

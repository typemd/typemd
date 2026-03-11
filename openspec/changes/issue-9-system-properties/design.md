## Context

typemd currently has one system property (`name`) implemented via hardcoded constants and scattered special-case logic across `object.go`, `type_schema.go`, `shared_properties.go`, and `sync.go`. Adding `created_at` and `updated_at` (and eventually `created_by`/`modified_by` from #77) would multiply these special cases. A registry mechanism is needed to centralize system property management.

The existing `name` property pattern establishes the integration points:
- `NewObject()` auto-sets the value
- `ValidateSchema()` / `ValidateSharedProperties()` reject reserved names
- `SyncIndex()` preserves system properties during property filtering
- `OrderedPropKeys()` places system properties before schema-defined properties

## Goals / Non-Goals

**Goals:**
- Introduce a system property registry that centralizes all system property definitions
- Add `created_at` and `updated_at` timestamp properties to all objects
- Refactor existing `name` property to use the registry
- Ensure timestamps use RFC 3339 format with local timezone

**Non-Goals:**
- Migration of existing objects (deferred to #77 which handles git-based metadata)
- `created_by` / `modified_by` properties (separate issue #77)
- Indexing system properties as separate SQLite columns (stay in JSON `properties` column)
- Making system properties user-configurable

## Decisions

### 1. System property registry as a package-level slice

Define system properties as a `[]SystemProperty` slice in a new file `core/system_property.go`. Each entry declares the property name, type, and an optional description. Helper functions (`IsSystemProperty`, `SystemPropertyNames`) provide the lookup interface.

**Why a slice, not a map:** Slices preserve insertion order, which matters for frontmatter output ordering (name → created_at → updated_at → schema properties). A map would require a separate ordering mechanism.

**Why not an interface with auto-set callbacks:** The auto-set logic for each property differs significantly (`name` from slug, `created_at` from `time.Now()`, `updated_at` on every save). Encoding these as callbacks in the registry adds abstraction without reducing complexity. Instead, `NewObject` and `SaveObject` handle the setting logic directly, while the registry handles identification and validation.

**Alternatives considered:**
- Continue hardcoding each property: rejected because #77 will add 2 more, making 5 total scattered special cases
- Use a config file for system properties: rejected because these are hardcoded by typemd, not user-configurable

### 2. RFC 3339 with local timezone, stored as string

Timestamps are formatted via `time.Now().Format(time.RFC3339)` and stored as strings in YAML frontmatter. Go's `time.RFC3339` naturally includes the local timezone offset.

**Why string, not time.Time in the map:** YAML parsers auto-convert datetime strings to `time.Time`, but when we write back via `yaml.Marshal`, a `time.Time` value may serialize differently than the original string. Storing as string gives us control over the exact format. The existing `validateDatetime` function already accepts both `time.Time` and string.

**Alternatives considered:**
- UTC only: rejected for human readability in local-first tool
- Unix timestamp: rejected for same reason

### 3. `updated_at` only on SaveObject

`updated_at` is set in `saveObjectFile()` (the internal method called by both `SaveObject` and `SetProperty`). Manual edits to .md files outside of typemd will not update `updated_at`.

**Why not detect changes via SyncIndex:** File mtime is unreliable (git checkout, file copy, etc.) and would introduce false positives. Keeping it to typemd operations only is simple and correct.

### 4. Frontmatter ordering: system properties first, in registry order

`OrderedPropKeys` will iterate the system property registry to place all system properties first (in their defined order: name, created_at, updated_at), followed by schema-defined properties, then extras alphabetically.

### 5. Existing objects: graceful absence

Objects without `created_at`/`updated_at` continue to work. No migration during SyncIndex. The properties are simply absent from frontmatter and the properties map. Only `name` has migration logic (backfill from filename) because it predates this change.

## Risks / Trade-offs

- **[Risk] YAML auto-parsing of datetime strings** → The `gopkg.in/yaml.v3` parser may auto-convert RFC 3339 strings to `time.Time` on read. Mitigation: `parseFrontmatter` already returns `map[string]any`, and `validateDatetime` handles both types. On write, we'll ensure `created_at`/`updated_at` are always stored as strings.
- **[Risk] Refactoring `name` to use registry may break existing tests** → Mitigation: the registry is additive; `NameProperty` constant remains, `IsSystemProperty("name")` returns true. Existing code that checks `== NameProperty` continues to work; we only change validation and sync filtering to use the registry.
- **[Trade-off] No auto-update on manual file edits** → Acceptable for a local-first tool where typemd commands are the primary interface. Users who need mtime can check the filesystem directly.

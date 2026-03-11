## Why

Objects in typemd have no temporal metadata. Users cannot see when an object was created or last modified. This is foundational for any knowledge management tool — timestamps enable sorting, filtering, and understanding the evolution of a vault. The `name` system property exists as a single hardcoded special case, but adding more system properties requires a proper registry mechanism.

## What Changes

- Introduce a **system property registry** that centrally defines all system-managed properties (`name`, `created_at`, `updated_at`)
- Refactor existing `name` property handling to use the registry instead of hardcoded checks
- Add `created_at` (datetime) — set once on object creation, never modified
- Add `updated_at` (datetime) — set on creation, updated on every `SaveObject` call
- Timestamps use RFC 3339 format with local timezone (e.g., `2026-03-07T12:00:00+08:00`)
- Existing objects without timestamps gracefully remain empty (migration deferred to #77)
- Type schemas and shared properties cannot define properties that conflict with system property names

## Capabilities

### New Capabilities

- `system-property-registry`: Centralized registry for system properties with helpers like `IsSystemProperty()` and `SystemPropertyNames()`, replacing scattered hardcoded checks

### Modified Capabilities

- `system-property`: Add `created_at` and `updated_at` requirements alongside existing `name` requirements. Update frontmatter ordering and sync behavior to handle all system properties via the registry.

## Impact

- **core/object.go** — `NewObject` sets timestamps; `SaveObject` updates `updated_at`; new registry types and helpers
- **core/type_schema.go** — `OrderedPropKeys` uses registry for ordering; validation uses `IsSystemProperty()`
- **core/shared_properties.go** — validation uses `IsSystemProperty()`
- **core/sync.go** — `SyncIndex` preserves all system properties via registry
- **No breaking changes** — existing objects without timestamps continue to work; `name` behavior is unchanged

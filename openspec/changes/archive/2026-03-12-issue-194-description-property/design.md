## Context

typemd objects have three system properties: `name`, `created_at`, `updated_at`. These are managed by a central registry in `core/system_property.go`. The registry drives validation (rejecting reserved names in schemas), frontmatter ordering, and index sync. Adding a new system property follows this established pattern exactly.

## Goals / Non-Goals

**Goals:**

- Add `description` to the system property registry as a `text` type property
- Maintain frontmatter ordering: `name` → `description` → `created_at` → `updated_at`
- Ensure `description` is automatically reserved (rejected in type schemas and shared properties)
- Full backward compatibility with existing objects

**Non-Goals:**

- Auto-populating `description` on object creation (it's user-authored, unlike timestamps)
- Adding `description` to existing objects during sync/migration
- Adding a `GetDescription()` convenience method (not needed yet; can be added later)
- Full-text search changes (FTS5 already indexes the JSON properties blob)

## Decisions

### 1. Registry-only change

**Decision:** Add one entry to `systemProperties` slice in `system_property.go`.

**Rationale:** The existing infrastructure (`IsSystemProperty`, `SystemPropertyNames`, `OrderedPropKeys`, `ValidateSchema`, `ValidateSharedProperties`, `SyncIndex`) all derive behavior from the registry. No changes needed in these functions — they already iterate the registry dynamically.

**Alternative considered:** Adding special-case handling in each function. Rejected because the registry pattern was specifically designed to avoid this.

### 2. Insert position: between name and created_at

**Decision:** `description` is index 1 in the registry (after `name` at index 0, before `created_at`).

**Rationale:** Follows the identity-first ordering principle from issue #199: identity properties (`name`, `description`) group together, then timestamps. This positions us well for the future frontmatter reordering.

### 3. No migration during sync

**Decision:** `SyncIndex` will not add `description` to existing object files.

**Rationale:** Consistent with `created_at`/`updated_at` behavior — sync preserves existing system properties but does not add missing ones (except `name`, which has special migration logic). `description` is optional and user-authored; auto-adding an empty value adds noise.

## Risks / Trade-offs

- **[Low] Existing tests hardcode registry order** → Update all tests that assert on `SystemPropertyNames()` order or registry length. Mitigation: search for hardcoded `["name", "created_at", "updated_at"]` patterns.
- **[Low] BDD scenarios reference 3-property registry** → Update scenario expectations. Mitigation: covered by test-first approach in tasks.

## Context

typemd objects are always editable — there is no mechanism to protect individual objects from accidental modifications. Users managing a knowledge base may want to "freeze" certain objects (e.g., reference entries, published notes) so they cannot be edited through CLI or TUI without explicitly unlocking first.

The system already has a system property registry (`system_property.go`) with immutable properties (`created_at`, `updated_at`) that are auto-managed. The `locked` property is different: it's user-authored (like `name` and `description`) but controls write access to the entire object.

## Goals / Non-Goals

**Goals:**
- Add `locked` as a boolean system property stored in object frontmatter
- Provide `tmd object lock <id>` and `tmd object unlock <id>` CLI commands
- Guard all write operations in `ObjectService` against locked objects
- Show visual lock indicator (🔒) in TUI and prevent editing locked objects
- Surface clear error messages when operations are blocked by lock

**Non-Goals:**
- Bulk lock/unlock commands
- Lock with password or encryption key
- Lock expiry or time-based locking
- Type-level locking (all objects of a type)
- Property-level locking (individual properties within an object)

## Decisions

### 1. `locked` as a stored system property

**Decision:** Register `locked` in the `systemProperties` registry as a boolean, user-authored (not immutable) system property.

**Rationale:** Using the existing system property mechanism means `locked` appears in frontmatter consistently, is recognized during sync, and follows established patterns. It's user-authored because users set it explicitly (like `name`), not auto-managed (like `updated_at`).

**Alternatives considered:**
- Separate `.lock` file per object — adds filesystem complexity, harder to discover
- Metadata in SQLite index only — violates "files are source of truth" principle
- Type schema flag — too coarse, locks all objects of a type

### 2. Guard at ObjectService level

**Decision:** Check `locked` in `ObjectService.Save()`, `SetProperty()`, and `SetPropertyMultiple()`. Return a typed `ErrObjectLocked` error.

**Rationale:** ObjectService is the command boundary in the CQRS architecture — all mutations flow through it. Guarding here catches CLI, TUI, and any future consumers. The lock/unlock commands bypass this guard by directly writing the `locked` property without going through the standard save path.

**Alternatives considered:**
- Guard at Object entity level — too low, would block lock/unlock operations themselves
- Guard at Vault facade level — misses direct ObjectService usage

### 3. Lock/unlock commands bypass the guard

**Decision:** `tmd object lock` and `tmd object unlock` write the `locked` property directly, bypassing the `ObjectService.Save()` guard.

**Rationale:** This is a controlled operation that changes the lock state itself. The commands use `ObjectService.SetLocked(id, bool)` — a dedicated method that skips the lock check.

### 4. Frontmatter ordering

**Decision:** Place `locked` after `tags` (last among system properties) in frontmatter output order.

**Rationale:** `locked` is a meta-property about the object's editability, not content. Placing it last among system properties keeps content-related properties (`name`, `description`) prominent. When absent or false, it's omitted from frontmatter (`omitempty`).

### 5. TUI behavior for locked objects

**Decision:** Show 🔒 icon next to the object name in the sidebar and properties panel. When a user attempts to enter edit mode (Tab to properties, Enter on a property), show a toast: "Object is locked. Unlock to edit." Do not enter prop editor mode at all.

**Rationale:** Preventing edit mode entry is simpler and clearer than allowing entry but blocking each individual edit. The toast provides immediate feedback on why editing is blocked.

## Risks / Trade-offs

- **[Risk] Users lock objects and forget** → Mitigation: clear 🔒 visual indicator in TUI sidebar and properties panel; unlock is a simple `tmd object unlock` command or direct frontmatter edit
- **[Risk] Lock check adds overhead to every save** → Mitigation: single boolean check, negligible performance impact
- **[Risk] Sync writes back to locked objects** → Mitigation: Projector sync (name resolution, wikilink expansion) is a system operation that MUST bypass the lock — it operates on file content, not user-initiated edits. The lock guard is only in ObjectService, not in Projector.

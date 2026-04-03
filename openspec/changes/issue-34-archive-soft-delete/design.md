## Context

typemd objects are stored as Markdown files. Deleting an object permanently removes the file. The only recovery path is git history. The `locked` system property (boolean, omitempty) is the closest precedent — it's a mutable boolean that controls object behavior without altering the file's location.

The current system has no concept of "hidden but not deleted" objects. All queries return all objects; filtering is user-driven via explicit `FilterRule`.

## Goals / Non-Goals

**Goals:**

- Add `archived` as a stored, mutable boolean system property (omitempty, like `locked`)
- Provide `tmd object archive` and `tmd object unarchive` CLI commands
- Exclude archived objects from default queries (QueryService-level)
- Allow opt-in inclusion via `--include-archived` flag on `tmd object list` / `tmd object query`
- Archived objects remain resolvable by wiki-links and direct `GetObject` calls

**Non-Goals:**

- Trash/recycle bin with auto-expiry
- Moving files to a separate archive folder (files stay in `objects/`)
- TUI archive management UI (e.g., dedicated archive panel)
- Bulk archive via glob patterns (single object only in this iteration)
- Write-guard behavior (archived ≠ locked — archived objects can still be edited)

## Decisions

### 1. System property, not a special folder

**Decision:** `archived` is a frontmatter boolean property, not a file-system operation.

**Rationale:** Matches the existing `locked` pattern. Files stay in `objects/`, preserving wiki-links and relations. No file-move logic needed. The reconciler and projector already handle system properties generically.

**Alternative considered:** Moving files to `objects/.archived/` — rejected because it breaks wiki-link resolution, requires path migration logic, and complicates the reconciler.

### 2. Default exclusion at QueryService level

**Decision:** Add `IncludeArchived bool` to `QueryOptions` (or equivalent parameter). When false (default), `QueryService.Query()` injects `{Property: "archived", Operator: "is_not", Value: "true"}` before delegating to the index.

**Rationale:** Centralizing the filter at QueryService ensures consistent behavior across CLI, TUI, MCP, and Web. Each consumer doesn't need to remember to filter archived objects.

**Alternative considered:** Filtering at each call site (CLI, TUI, etc.) — rejected because it's error-prone and violates DRY.

### 3. SetArchived mirrors SetLocked

**Decision:** Add `ObjectService.SetArchived(id, archived)` following the exact same pattern as `SetLocked`. When `archived` is false, delete the key from properties (omitempty). Emit `ObjectUpserted` event.

**Rationale:** Reuses the proven pattern. No lock guard bypass needed (unlike `SetLocked`, there's no mutation restriction on archived objects).

### 4. CLI commands under `tmd object`

**Decision:** Add `tmd object archive <id>` and `tmd object unarchive <id>` as subcommands of `objectCmd`, following the `lock.go` pattern.

**Rationale:** Consistent with existing object lifecycle commands (`tmd object lock/unlock`). Uses `resolveIDInteractive` for prefix matching.

### 5. No SQLite schema migration

**Decision:** Store `archived` in the existing `properties` JSON blob column, not as a dedicated column.

**Rationale:** The SQLite index already stores all properties as JSON. The reconciler automatically includes any system property from `SystemPropertyNames()`. JSON extraction (`json_extract(properties, '$.archived')`) is already supported by `FilterRuleToSQL`. Adding a dedicated column would require a migration mechanism that doesn't exist yet.

## Risks / Trade-offs

- **[Query performance]** JSON extraction for `archived` filtering on every query is slightly slower than a dedicated indexed column. → Mitigation: For typical vault sizes (< 10k objects), this is negligible. A dedicated column can be added later if profiling shows a bottleneck.
- **[Implicit filter surprise]** Users may be confused when archived objects "disappear" from queries. → Mitigation: `tmd object archive` prints a clear message. `tmd object list --include-archived` provides easy recovery. The TUI sidebar respects the same filter.
- **[Relation integrity]** Archiving an object with incoming relations may confuse users when related objects reference a "missing" object. → Mitigation: Archived objects are not deleted — `GetObject(id)` still returns them. Wiki-links and relations still resolve. Only default query results are affected.

## Context

The TUI `refreshData()` is triggered by `fileChangedMsg` from a file watcher whenever any file in `objects/` changes. Currently, every trigger runs a full `Projector.Sync()` — walking all object files, parsing frontmatter, upserting every object into SQLite, rebuilding the entire FTS index, then querying all objects back to rebuild the sidebar. This is O(N) for N total objects, regardless of how many files actually changed.

The watcher already detects which files changed but discards this information, sending an empty `fileChangedMsg{}`. The index already supports single-object `Upsert()` and `Remove()`. The building blocks for incremental sync exist but aren't wired together.

## Goals / Non-Goals

**Goals:**
- Incremental index sync: only process changed files, not the entire vault
- Incremental FTS: update only affected entries instead of full rebuild
- Schema cache at Vault level to avoid repeated YAML file reads during `buildGroups()`
- Watch `.typemd/types/` for external schema changes to invalidate cache
- Configurable debounce interval via `tui.debounce_ms` in `.typemd/config.yaml`
- Fallback to full sync on watcher overflow or initial startup

**Non-Goals:**
- Incremental sidebar rebuild (sidebar still does full `QueryObjects` + `buildGroups` after sync)
- Incremental wikilink sync (still full rebuild per sync)
- Incremental tag relation sync (still full rebuild per sync)
- Async/background sync (sync still blocks the TUI update loop)

## Decisions

### 1. `fileChangedMsg` carries file paths

`fileChangedMsg` will change from empty struct to carry `[]string` of changed file paths collected during the debounce window. The watcher collects all events (Write/Create/Remove/Rename) during debounce and deduplicates paths.

**Alternative considered:** Single file path (take last event) — rejected because debounce window can accumulate multiple file changes (e.g., `tmd object create` creates object + updates tag).

### 2. New `Projector.SyncFiles()` method for incremental sync

A new method `SyncFiles(paths []string) (*SyncResult, error)` will:
1. Classify each path as create/update vs. delete (check if file exists)
2. For existing files: `repo.Get(id)` → filter properties → `index.Upsert()`
3. For deleted files: `index.Remove(id)`
4. Run full wikilink and tag sync (reuse existing `Sync()` logic)
5. Skip full FTS rebuild (handled incrementally by `Upsert`)

The existing `Sync()` method remains unchanged for full sync scenarios.

**Path-to-ID conversion:** Object files follow `objects/<type>/<slug>.md` convention. The ID is `<type>/<slug>` (strip `objects/` prefix and `.md` suffix). This logic already exists in `LocalObjectRepository`.

### 3. Incremental FTS in `Upsert()`

`SQLiteObjectIndex.Upsert()` will update the FTS entry atomically alongside the main table update. SQLite FTS5 content-sync tables can use triggers, but since we use a manual `Rebuild()` approach currently, we'll add explicit FTS delete+insert within `Upsert()` and `Remove()`.

The `Rebuild()` method remains for full sync but is no longer called during incremental sync.

### 4. Schema cache at Vault level

Add a `schemaCache map[string]*TypeSchema` field to `Vault`. `LoadType()` checks cache first, falls through to `repo.GetSchema()` on miss. Invalidation:
- `SaveType()` / `DeleteType()` → invalidate that type's entry
- `MigrateSchemas()` → invalidate all
- External `.typemd/types/` file change (via watcher) → invalidate all (simpler than tracking which type file changed)

### 5. Watcher watches `.typemd/types/` too

The watcher will additionally monitor the `.typemd/types/` directory. Schema file changes produce a distinct message type (`schemaChangedMsg`) that triggers schema cache invalidation + full data refresh (since schema changes can affect property filtering for all objects of that type).

### 6. Debounce from config

`VaultConfig` gains a `TUI` sub-struct with `DebounceMs int`. The watcher reads this value at startup. Default is 200ms if not configured or zero.

### 7. Fallback to full sync

`refreshData()` uses incremental sync when `fileChangedMsg` carries paths. Falls back to full `Projector.Sync()` when:
- `fileChangedMsg` has empty paths (initial sync, watcher overflow)
- `SyncFiles()` returns an error (corrupted state recovery)

## Risks / Trade-offs

- **[Wikilink/tag full rebuild negates some gains]** → Acceptable trade-off for correctness. These operations query the index (fast) rather than walking the filesystem (slow). Can be optimized incrementally in a future issue.
- **[Schema cache staleness from external edits]** → Mitigated by watching `.typemd/types/` directory. Cache invalidates on any type file change.
- **[Path-to-ID conversion may fail for edge cases]** → Use existing repository path conventions. If conversion fails, fall back to full sync.
- **[FTS incremental update correctness]** → SQLite FTS5 content tables require manual sync. Must delete old FTS entry before inserting new one to avoid stale search results.

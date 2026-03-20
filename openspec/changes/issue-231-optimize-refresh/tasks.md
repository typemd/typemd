## 1. Core: Incremental FTS in Upsert/Remove

- [x] 1.1 Write unit tests for incremental FTS update in Upsert (search finds updated content)
- [x] 1.2 Write unit tests for FTS cleanup in Remove (search no longer finds removed object)
- [x] 1.3 Update `SQLiteObjectIndex.Upsert()` to delete+insert FTS entry atomically (NOT NEEDED — FTS triggers already handle this)
- [x] 1.4 Update `SQLiteObjectIndex.Remove()` to delete FTS entry before removing object row (NOT NEEDED — FTS triggers already handle this)

## 2. Core: Vault-level Schema Cache

- [x] 2.1 Write BDD scenarios for schema cache (first load, cache hit, invalidation on SaveType/DeleteType/MigrateSchemas)
- [x] 2.2 Implement BDD step definitions for schema cache scenarios
- [x] 2.3 Add `schemaCache map[string]*TypeSchema` to Vault and update `LoadType()` to check cache first
- [x] 2.4 Invalidate cache entry in `SaveType()` and `DeleteType()`
- [x] 2.5 Invalidate entire cache in `MigrateSchemas()`
- [x] 2.6 Add `InvalidateSchemaCache()` method for external change notifications

## 3. Core: VaultConfig TUI Namespace

- [x] 3.1 Write unit tests for `tui.debounce_ms` config key (load, get, set)
- [x] 3.2 Add `TUIConfig` sub-struct with `DebounceMs int` to `VaultConfig`
- [x] 3.3 Register `tui.debounce_ms` in `configKeyRegistry`

## 4. Core: Projector.SyncFiles

- [x] 4.1 Write BDD scenarios for `SyncFiles` (file created/updated, file deleted, fallback on error)
- [x] 4.2 Implement BDD step definitions for `SyncFiles` scenarios
- [x] 4.3 Add path-to-ID conversion helper (strip `objects/` prefix and `.md` suffix)
- [x] 4.4 Implement `Projector.SyncFiles(paths []string)` — incremental upsert/remove + full wikilink/tag sync
- [x] 4.5 Add unit tests for path-to-ID edge cases (nested dirs, non-.md files, invalid paths)

## 5. TUI: Watcher Enhancements

- [x] 5.1 Update `fileChangedMsg` to carry `[]string` of changed file paths
- [x] 5.2 Update watcher to collect and deduplicate file paths during debounce window
- [x] 5.3 Read `tui.debounce_ms` from vault config and use as debounce interval
- [x] 5.4 Add `.typemd/types/` directory monitoring with distinct `schemaChangedMsg`
- [x] 5.5 Write unit tests for path deduplication logic

## 6. TUI: Incremental Refresh Flow

- [x] 6.1 Update `refreshData()` to call `SyncFiles()` when paths are available, fallback to `Sync()` otherwise
- [x] 6.2 Handle `schemaChangedMsg` — invalidate schema cache and trigger full refresh
- [x] 6.3 Write unit tests for fallback behavior (empty paths, SyncFiles error)

## Why

`refreshData()` performs a full sync on every file change: walks all object files, re-parses frontmatter, upserts everything into SQLite, rebuilds FTS, then re-queries all objects and rebuilds sidebar groups. This becomes a performance bottleneck as the vault grows (100+ objects). The watcher already knows which files changed but discards that information.

## What Changes

- `fileChangedMsg` carries a list of changed file paths (collected during debounce window) instead of being an empty struct
- Incremental index sync: only `Get()` + `Upsert()` the changed files, with fallback to full sync on watcher overflow or initial startup
- Incremental FTS: update only the affected FTS entries instead of full `Rebuild()`
- Wikilinks and tags still use full rebuild after incremental object sync
- Sidebar still uses full rebuild (`QueryObjects` + `buildGroups`) after sync
- Vault-level schema cache with invalidation on `SaveType()` / `DeleteType()` / `MigrateSchemas()` and external `.typemd/types/` file changes
- Watcher additionally monitors `.typemd/types/` directory for external schema edits
- Debounce interval configurable via `tui.debounce_ms` in `.typemd/config.yaml` (default 200)

## Capabilities

### New Capabilities
- `incremental-sync`: Incremental file-to-index synchronization with fallback to full sync
- `schema-cache`: Vault-level type schema caching with invalidation strategies

### Modified Capabilities
- `vault-config`: Add `tui.debounce_ms` configuration key
- `tui-session-state`: `fileChangedMsg` now carries changed file paths

## Impact

- `core/projector.go` — new incremental sync method alongside existing `Sync()`
- `core/sqlite_object_index.go` — incremental FTS update in `Upsert()`
- `core/vault.go` or `core/type_schema.go` — schema cache layer
- `core/vault_config.go` — new `tui.debounce_ms` config key
- `tui/watcher.go` — collect file paths during debounce, watch `.typemd/types/`
- `tui/app.go` — pass file paths to incremental sync, handle schema invalidation

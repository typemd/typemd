## Why

The query pipeline uses a legacy `key=value` filter string format that only supports equality checks. Meanwhile, a structured `FilterRule` type with 20+ operators (contains, gt, before, etc.) and `FilterRuleToSQL()` already exist but are disconnected from the query path. After #263 integrated `FilterRuleToSQL` into the codebase, the legacy string format has no external consumers — `tmd query` is the only CLI entry point and can be removed. Unifying on `[]FilterRule` eliminates the parallel filter systems and lets `ViewConfig.Filter` flow naturally through queries.

## What Changes

- **BREAKING**: Remove `tmd query` command (`cmd/query.go`)
- Replace `filter string` parameter with `[]FilterRule` in `ObjectIndex.Query()`, `QueryService.Query()`, and `Vault.QueryObjects()`
- Remove `parseFilter()` and `filterCondition` type from `core/query.go`
- Migrate all internal callers to construct `[]FilterRule` instead of filter strings:
  - `stats.go` — `VaultStats()`, `TypeStats()`
  - `type_schema.go` — `CountObjectsByType()`
  - `object_service.go` — name uniqueness check
  - `validate.go` — validate all objects
  - `migrate.go` — migrate objects by type
  - `tui/view_mode.go` — view mode queries
  - `tui/app.go` — object list loading
  - `cmd/object_list.go` — `tmd object list`
- Update `SQLiteObjectIndex.Query()` to iterate `[]FilterRule` and call `FilterRuleToSQL()`, with special handling for `type` column (not in properties JSON)
- Update all BDD scenarios and unit tests

## Capabilities

### New Capabilities
- `structured-query-filter`: Structured `[]FilterRule`-based query pipeline replacing legacy `key=value` filter strings

### Modified Capabilities
- `vault-stats`: Internal callers change from filter strings to `[]FilterRule` (implementation change, query behavior unchanged)
- `tui-object-list`: Object list loading changes from filter strings to `[]FilterRule` (implementation change, list behavior unchanged)

## Impact

- **core/**: `object_index.go` (interface change), `query_service.go`, `query.go`, `sqlite_object_index.go`, `object_service.go`, `stats.go`, `type_schema.go`, `validate.go`, `migrate.go`
- **cmd/**: `query.go` (removed), `object_list.go` (updated)
- **tui/**: `app.go`, `view_mode.go` (updated)
- **Tests**: BDD step definitions and unit tests updated across `core/` and `tui/`
- **No data model changes**: SQLite schema, file format, and YAML configs unchanged

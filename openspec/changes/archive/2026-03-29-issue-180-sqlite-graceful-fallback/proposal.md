## Why

Key query operations (`QueryObjects()`, `SearchObjects()`) hard-fail when the SQLite index is unavailable (missing, corrupt, or connection error). This violates typemd's core architectural principle: **files are the source of truth, SQLite is optional acceleration**. Adding filesystem-based fallback paths ensures critical read operations degrade gracefully instead of failing entirely.

## What Changes

- Add fallback detection in `QueryService` that catches index errors and falls back to filesystem scanning via `ObjectRepository.Walk()`
- Implement in-memory filtering for `Query()` fallback (property matching against `Object` fields)
- Implement substring-based search for `Search()` fallback (matching against name, description, body)
- Emit warning log when operating in fallback mode
- Vault stats methods (`VaultStats()`, `TypeStats()`) also gain fallback paths

## Capabilities

### New Capabilities

- `query-fallback`: Filesystem-based fallback for `QueryObjects()` and `SearchObjects()` when SQLite index is unavailable, including in-memory filtering, substring search, and degraded-mode warnings

### Modified Capabilities

_(none — existing query behavior is preserved when index is available)_

## Impact

- **Code**: `core/query_service.go` (fallback logic), possibly new `core/query_fallback.go` for filesystem-based filtering/search helpers
- **APIs**: No API changes — same signatures, same return types; callers are unaware of fallback
- **Performance**: Fallback is O(n) file reads + in-memory filtering vs indexed queries; acceptable for small-to-medium vaults, documented as degraded
- **Dependencies**: No new dependencies

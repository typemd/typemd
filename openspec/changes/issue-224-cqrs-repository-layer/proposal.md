## Why

The `Vault` struct in `core/` is a god object with 50+ methods that mixes domain logic, use case orchestration, and infrastructure concerns (direct `os.*` file I/O and raw SQL queries). This makes it impossible to provide alternative storage backends for the upcoming Web UI (`try.typemd.io` via GitHub API, Wails desktop app via Go bindings) without duplicating the entire Vault. Extracting a Repository layer with DDD-aligned interfaces enables CQRS separation and multi-backend support.

## What Changes

- **Define `ObjectRepository` interface** — write-side repository for entity persistence and retrieval by known ID. Returns domain entities (`*Object`, `*TypeSchema`), not raw bytes. Path/storage conventions encapsulated in implementations.
- **Define `ObjectIndex` interface** — read-side service for search, query, and discovery. Returns lightweight `ObjectResult` records, not full entities. Also exposes projector write methods for index maintenance.
- **Extract `SQLiteObjectIndex`** — move all SQLite query/write operations out of Vault into a concrete `ObjectIndex` implementation.
- **Extract `LocalObjectRepository`** — move all filesystem I/O out of Vault into a concrete `ObjectRepository` implementation.
- **Extract `Projector`** — formalize `SyncIndex` as an independent component that reads from `ObjectRepository.Walk()` and writes to `ObjectIndex.Upsert()`.
- **Refactor `Vault` to Use Case layer** — Vault depends only on `ObjectRepository` and `ObjectIndex` interfaces, coordinating commands (dual-write to both) and queries (delegate to index).

## Capabilities

### New Capabilities

None — this is a pure internal refactoring. No new external behaviors are introduced.

### Modified Capabilities

None — all existing behaviors are preserved. The refactoring changes internal structure (how) without changing external requirements (what).

## Impact

- **core/** — Major refactoring of vault.go, object.go, sync.go, relation.go, query.go, wikilink.go, type_schema.go, template.go, shared_properties.go, list.go, migrate.go
- **cmd/** — Minimal changes; commands use Vault which maintains the same public API
- **tui/** — Minimal changes; TUI uses Vault which maintains the same public API
- **mcp/** — Minimal changes; MCP uses Vault which maintains the same public API
- **Testing** — All existing BDD scenarios and unit tests must continue to pass; new tests for repository implementations
- **No breaking changes** to Vault's public method signatures during this refactoring

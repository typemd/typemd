## Why

Deleting an object permanently removes the file with no recovery option (unless recovered via git). A soft delete mechanism using an `archived` frontmatter flag would hide objects from default views without destroying data, giving users a safe way to declutter their vault.

## What Changes

- Register `archived` as a stored system property (boolean, default false/absent)
- Add `tmd archive <id>` and `tmd unarchive <id>` CLI commands
- Update `QueryService.QueryObjects()` to exclude archived objects by default
- Add `--include-archived` flag to `tmd list` and `tmd query` commands
- Update SQLite index to include `archived` field for efficient filtering
- Archived objects remain resolvable by wiki-links and relations (archived ≠ deleted)

## Capabilities

### New Capabilities

- `object-archive`: Soft delete mechanism via `archived` system property, with CLI commands to archive/unarchive objects and query-level filtering to hide archived objects by default

### Modified Capabilities

- `object-format`: Adding `archived` as a new system property to the object frontmatter format
- `wiki-links`: Archived objects must still resolve in wiki-link lookups (no behavioral spec change needed — just confirming archived ≠ deleted)

## Impact

- **core/**: `Object` entity gains `archived` system property; `ObjectService` gets archive/unarchive commands; `QueryService` filters archived by default; `Reconciler` and `Projector` handle the new field; SQLite schema adds `archived` column
- **cmd/**: New `tmd archive` and `tmd unarchive` commands; `--include-archived` flag on `tmd list` and `tmd query`
- **mcp/**: Search results should respect archived filtering (existing `search` tool)
- **tui/**: Sidebar and views should exclude archived objects by default (via QueryService)

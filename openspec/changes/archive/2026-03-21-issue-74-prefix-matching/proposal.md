## Why

After ULID suffixes were introduced, relation references in YAML frontmatter require the full ULID-suffixed object ID (e.g., `author: person/john-doe-01jqr3k5mpbvn8e0f2g7h9txyz`). This makes manually editing frontmatter tedious and error-prone. Users should be able to write short prefixes like `author: person/john-doe` and have the system resolve and expand them automatically.

Additionally, non-tag relation properties written directly in frontmatter are not synced to the SQLite `relations` table by the Projector — only `ObjectService.Link()` creates relation index records. This means hand-edited relations are invisible to queries and backlink displays until a full re-link cycle.

## What Changes

- **Prefix resolution in Projector sync**: When the Projector encounters a relation value without a ULID suffix, it resolves the prefix via `GlobIDs` to find the matching object.
- **Auto-expand and write-back**: On unique match, the Projector expands the prefix to the full ID and writes the updated frontmatter back to the file.
- **General relation sync**: The Projector syncs all schema-defined relation properties (not just tags) to the SQLite `relations` table during both full and incremental sync.
- **SyncResult reporting**: `SyncResult` reports prefix expansions (resolved count) and unresolvable references (ambiguous, not found).
- **Validation integration**: `tmd validate` reports unresolvable relation prefixes as warnings.

## Capabilities

### New Capabilities
- `relation-prefix-resolution`: Prefix matching and auto-expansion for relation values in frontmatter during Projector sync.
- `relation-sync`: Syncing all schema-defined relation properties (not just tags) from frontmatter to the SQLite relations table during Projector sync.

### Modified Capabilities
- `object-relations`: Adding the requirement that relation values in frontmatter can use prefix form and are auto-expanded during sync.
- `incremental-sync`: Extending incremental sync to include relation property syncing alongside wikilinks and tags.

## Impact

- **core/projector.go** — New `syncRelations()` method with prefix resolution + write-back + relation indexing
- **core/sync.go** — Extended `SyncResult` and `syncContext` structs
- **core/query_service.go** — `Resolve()` may be reused or a lower-level resolver extracted
- **core/validate.go** — Enhanced `ValidateRelations` to detect prefix references
- **core/object_service.go** — No changes needed; `Link()`/`Unlink()` continue to work with full IDs
- **TUI file watcher** — No changes needed; watcher is single-shot so write-back during sync won't cause loops

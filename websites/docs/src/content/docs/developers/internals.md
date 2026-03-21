---
title: Internals
description: Internal mechanisms — caching, sync, query pipeline, and formatting.
sidebar:
  order: 4
---

This page documents TypeMD's internal mechanisms that are invisible to users but important for contributors to understand. These cover the schema cache, incremental sync, relation sync, name resolution, query pipeline, and property formatting.

## Schema cache

The Vault maintains an in-memory cache of type schemas to avoid repeated disk reads.

### Cache behavior

- **First load**: `LoadType("book")` reads from `ObjectRepository.GetSchema("book")` and caches the result.
- **Subsequent loads**: Returns the cached schema without reading from disk.

### Invalidation rules

| Trigger | Scope |
|---------|-------|
| `SaveType()` | Invalidates the specific type's cache entry |
| `DeleteType()` | Invalidates the specific type's cache entry |
| `MigrateSchemas()` | Invalidates the entire cache |
| External file change in `.typemd/types/` | Invalidates the entire cache |

External file changes are detected by the TUI's file watcher. When a schema file is modified outside the Vault API, the watcher emits a schema change message that triggers full cache invalidation and data refresh.

## Incremental sync

The TUI file watcher supports incremental sync to avoid full index rebuilds on every file change.

### Debounce and path collection

The watcher collects changed file paths during a configurable debounce window (default: 200ms, configurable via `tui.debounce_ms` in `.typemd/config.yaml`). Duplicate paths within the same window are deduplicated.

### SyncFiles flow

`Projector.SyncFiles(paths []string)` synchronizes only the specified files:

1. For each path that still exists on disk: read the object via `ObjectRepository.Get()`, filter properties against the type schema, and upsert into the index.
2. For each path that no longer exists: remove the object's entry from the index.
3. After incremental object sync: perform full wikilink sync, full tag relation sync, and relation sync for changed objects.
4. Upsert/Remove operations update the corresponding FTS entry atomically — no full `Rebuild()` needed.

### Fallback to full sync

The TUI falls back to full `Projector.Sync()` when:

- `fileChangedMsg` has an empty paths list
- `Projector.SyncFiles()` returns an error
- Initial startup (first data load)

### Schema file monitoring

The watcher also monitors `.typemd/types/`. Schema file changes produce a distinct message that triggers schema cache invalidation and full data refresh — incremental sync does not apply to schema changes.

## Relation sync

During sync, the Projector reads each object's frontmatter, identifies relation properties defined in the type schema, and inserts corresponding records into the SQLite `relations` table.

### Sync behavior by relation type

- **Single-value relation**: One record per relation (e.g., `author: person/john-doe-01abc...` → one `relations` row).
- **Multi-value relation**: One record per value (e.g., `books: [book/a-01abc..., book/b-01xyz...]` → two rows).
- **Non-existent targets**: Skipped — if the referenced object doesn't exist on disk, no relation record is inserted.
- **Non-relation properties**: Ignored — only properties with `type: relation` in the schema are processed.

### Full sync vs. incremental sync

- **Full sync** (`Projector.Sync()`): Deletes all non-tag relation records and rebuilds from frontmatter. Tag relations are managed separately by `syncTagRelations`.
- **Incremental sync** (`Projector.SyncFiles()`): Deletes relation records only for changed objects and rebuilds them from updated frontmatter. Relations for unchanged objects remain intact.

## Name resolution (relation prefix resolution)

During sync, relation values without a ULID suffix are treated as `type/name` references and resolved to full object IDs.

### Resolution rules

| Input | Behavior |
|-------|----------|
| `person/john-doe-01abc...` (has ULID) | Treated as full ID, no resolution needed |
| `person/john-doe` (no ULID, unique match) | Resolved to full ID, file is updated |
| `person/nobody` (no match) | Left unchanged, reported as unresolved in `SyncResult` |
| `person/john` (multiple matches) | Left unchanged, reported as ambiguous in `SyncResult` |

### Name index

The Projector builds a per-type name index from walked objects. Each object's slug and original name are indexed. Duplicate names within the same type produce ambiguous entries.

### Auto-expansion write-back

When a prefix is successfully resolved, the Projector writes the expanded full ID back to the object's frontmatter file. Multiple properties can be expanded in a single file write. Unresolvable prefixes are left unchanged.

### Multi-value resolution

For `multiple: true` relations, each value in the array is resolved independently. A mix of full IDs and prefixes is handled correctly — full IDs are kept as-is, prefixes are resolved individually.

### SyncResult reporting

`SyncResult` includes:

- `Expanded` — count of successfully resolved prefixes
- `Unresolved` — list of unresolved references with prefix and reason

### Shared resolution method

The name resolution method is shared between relation prefix resolution and future wiki-link shorthand resolution. It takes a per-type name index and resolves `type/name` references to full object IDs.

## Structured query filter

The query pipeline uses structured `[]FilterRule` parameters instead of raw filter strings.

### FilterRule structure

Each rule specifies a property name, operator, and value. An empty slice returns all objects (no filtering). Multiple rules are combined with AND logic.

### SQL mapping

| Property | SQL mapping |
|----------|------------|
| `type` | Direct column reference: `type = ?` |
| Any other property | JSON extraction: `json_extract(properties, '$.property') = ?` |

### TypeFilter convenience

`TypeFilter(typeName string) []FilterRule` returns a single-element filter for type-based queries — the most common query pattern.

## Display property formatting

`DisplayProperty` provides two formatting methods for rendering property values in the TUI and CLI.

### FormatValue()

Returns the formatted value without key prefix:

| Property type | Example input | Output |
|--------------|---------------|--------|
| `string` | `"Robert Martin"` | `Robert Martin` |
| `date` | `2024-01-15` | `2024-01-15` |
| `datetime` | `2024-01-15T10:30:00` | `2024-01-15T10:30:00` |
| `multi_select` | `["go", "cli"]` | `[go, cli]` |
| `checkbox` (true) | `true` | `✓` |
| `checkbox` (false) | `false` | `` (empty) |
| `relation` | `person/robert-martin-01abc` | `→ person/robert-martin` |
| `backlink` | `note/my-note-01abc` | `⟵ note/my-note` |
| `reverse relation` | `book/clean-code-01abc` | `← book/clean-code` |
| `nil` | `nil` | `` (empty) |

### Format()

Composes output as `key + ": " + FormatValue()`. Used for the property detail panel. Delegates to `FormatValue()` for the value portion.

### Usage in view mode

View mode table rows use `FormatValue()` for property columns and preview panels, ensuring consistent formatting across all display contexts.

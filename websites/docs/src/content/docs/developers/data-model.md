---
title: Data Model
description: How TypeMD indexes data and manages the SQLite acceleration layer.
sidebar:
  order: 1
---

TypeMD uses a dual-storage design: Markdown files are the source of truth, and a SQLite database provides fast indexing and search. This page covers the technical details of the indexing layer.

For user-facing file structure details, see [File Structure](/advanced/file-structure).

## Files as source of truth

The fundamental design principle is that **files are always the source of truth**. The SQLite index is an acceleration layer that can be rebuilt from files at any time. If the index is deleted or corrupted, no data is lost — running `tmd reindex` or opening the vault will rebuild it.

This design enables:

- Version control with Git (only Markdown and YAML files matter)
- Manual editing with any text editor
- Syncing via file-based tools (Dropbox, iCloud, Syncthing)
- Portability between machines without database migration

## SQLite index

The index is stored at `.typemd/index.db` and contains:

- **Object metadata** — type, filename, and all frontmatter properties
- **Wiki-link records** — extracted from Object body content, used for backlink tracking
- **Full-text search index** — powered by [FTS5](https://www.sqlite.org/fts5.html), covering filenames, properties, and body content

The index file is not meant to be edited directly. It is managed entirely by TypeMD.

## Auto-sync behavior

The index is automatically synced in these situations:

- **Vault open** — when opening a vault with an empty or missing database (e.g. after a fresh `git clone`), TypeMD walks all Object files and populates the index
- **CLI/TUI operations** — creating, saving, or deleting Objects updates the index immediately
- **Manual reindex** — run `tmd reindex` to rebuild the index after editing files outside of TypeMD

The sync process is handled by the **Projector** component, which walks all Object files via the repository and upserts entries into the index. See [Architecture](/developers/architecture) for how the Projector fits into the system.

## Querying architecture

TypeMD provides two query paths:

### Structured queries

The query pipeline uses structured `FilterRule` conditions to filter Objects. Each rule specifies a property, an operator, and a value. Queries run against the SQLite index for performance, returning lightweight `ObjectResult` projections rather than full Object entities.

Filter rules are used programmatically through `Vault.QueryObjects()` and in view configurations (`.typemd/types/<name>/views/<view>.yaml`).

### Full-text search

Use `tmd search` to search across filenames, properties, and body content. Powered by SQLite FTS5:

```bash
tmd search "concurrency"
tmd search "golang" --json
```

### TUI search

In the TUI, press `/` to enter search mode. Results are filtered in real-time. Press `Esc` to clear results and return to the full list.

## Vault configuration

An optional configuration file at `.typemd/config.yaml` provides vault-level settings. The file uses interface-layer namespacing:

```yaml
# .typemd/config.yaml
cli:
  default_type: page
```

Currently supported settings:

| Key | Description |
|-----|-------------|
| `cli.default_type` | Default object type for `tmd object create` when type argument is omitted |

The config file is loaded during vault open. If the file is missing or empty, all settings use their zero values (no error). Invalid YAML produces an error.

## Unique constraint enforcement

Type schemas can opt into name uniqueness by setting `unique: true`. When enabled, TypeMD prevents creating multiple Objects of the same type with identical `name` values.

```yaml
# .typemd/types/person.yaml
name: person
unique: true  # only one person per name
properties:
  - name: role
    type: string
```

Uniqueness is enforced at creation time by checking the index for existing Objects with the same type and name. It is also validated by `tmd type validate`. The built-in `tag` type has `unique: true` enabled by default.

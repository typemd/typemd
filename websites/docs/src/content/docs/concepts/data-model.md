---
title: Data Model
description: How TypeMD stores and indexes data.
sidebar:
  order: 6
---

## Storage

Objects are stored as Markdown files with YAML frontmatter under `objects/<type>/`. The full Object ID is `type/<slug>-<ulid>`, e.g. `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`.

When created via the CLI (`tmd object create`), a 26-character lowercase ULID is automatically appended to the slug for uniqueness. Objects created manually (without the CLI) do not require a ULID — this is backward compatible.

```
vault/
├── .typemd/
│   ├── types/              # type schema definitions (YAML)
│   │   ├── book.yaml
│   │   └── person.yaml
│   └── index.db            # SQLite index (auto-updated)
└── objects/
    ├── book/
    │   └── golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md
    └── person/
        └── alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz.md
```

## Indexing

TypeMD uses SQLite with FTS5 for indexing. The index is stored at `.typemd/index.db` and contains:

- Object metadata (type, filename, properties)
- Wiki-link records extracted from object body content (for backlink tracking)
- Full-text search index over filenames, properties, and body content

The index is automatically synced when opening a vault with an empty or missing database (e.g. after a fresh clone). It is also kept up-to-date when using the TUI or CLI commands. Use `tmd reindex` to rebuild after manual file edits outside of TypeMD.

## Querying

TypeMD provides two ways to find Objects: structured queries and full-text search.

### Structured queries

Use `tmd query` to filter Objects by properties. Conditions use `key=value` format, separated by spaces (AND logic).

```bash
tmd query "type=book"
tmd query "type=book status=reading"
tmd query "type=book" --json
```

### Full-text search

Use `tmd search` to search across filenames, properties, and body content. Powered by SQLite FTS5.

```bash
tmd search "concurrency"
tmd search "golang" --json
```

### TUI search

In the TUI, press `/` to enter search mode. Results are filtered in real-time. Press `Esc` to clear results and return to the full list.

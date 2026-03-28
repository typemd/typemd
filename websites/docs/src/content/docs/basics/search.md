---
title: Search
description: Full-text search across your vault.
sidebar:
  order: 4
---

Search lets you find Objects by matching text across filenames, properties, and body content. It is powered by SQLite FTS5, so results are fast even in large vaults.

## How it works

When you open a vault, TypeMD builds a full-text search index from all Object files. The index covers:

- **Filenames** — the Object's slug and ULID
- **Properties** — all frontmatter values (system and schema-defined)
- **Body content** — the Markdown body of each Object

The index is stored at `.typemd/index.db` and automatically stays in sync when you use the CLI or TUI. The index is automatically synced every time the vault is opened, even after manual file edits.

## CLI usage

```bash
tmd search "concurrency"
tmd search "golang" --json
```

| Option | Description |
|--------|-------------|
| `--json` | Output results in JSON format |

## TUI usage

In the TUI, press `/` to enter search mode. Results are filtered in real-time as you type. Press `Esc` to clear the search and return to the full list.

## See also

- [Queries](/basics/queries) — structured filtering by property values
- [tmd search](/tui/search) — CLI reference
- [Data Model](/developers/data-model#sqlite-index) — how the index works

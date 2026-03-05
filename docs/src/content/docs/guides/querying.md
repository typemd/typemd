---
title: Querying
description: Filter and search objects in your vault.
sidebar:
  order: 4
---

TypeMD provides two ways to find objects: structured queries and full-text search.

## Structured Queries

Use `tmd query` to filter objects by properties. Conditions use `key=value` format, separated by spaces (AND logic).

```bash
tmd query "type=book"
tmd query "type=book status=reading"
tmd query "type=book" --json
```

## Full-text Search

Use `tmd search` to search across filenames, properties, and body content. Powered by SQLite FTS5.

```bash
tmd search "concurrency"
tmd search "golang" --json
```

## TUI Search

In the TUI, press `/` to enter search mode. Results are filtered in real-time. Press `Esc` to clear results and return to the full list.

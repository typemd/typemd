---
title: --reindex
description: Rebuild the SQLite index and search database.
sidebar:
  order: 8
---

The `--reindex` flag forces a full sync of the `objects/` directory to the database, cleans up orphaned relations, and rebuilds the full-text search index. It can be combined with any command.

> **Note:** When opening a vault, TypeMD automatically syncs the index if it is empty or missing. You only need `--reindex` when files have been edited while the vault was not open.

```bash
# Reindex and launch TUI
tmd --reindex

# Reindex and start MCP server
tmd mcp --reindex

# Reindex and run a query
tmd query "type=book" --reindex
```

## Orphaned relation cleanup

When an object is deleted from disk, any relations pointing to or from that object become orphaned. During reindex, these dangling references are automatically detected and removed from the index. A warning is displayed listing the affected relations:

```
Warning: Found 2 orphaned relation(s):
  book/golang-in-action -> person/deleted-author (relation: "author")
  person/deleted-author -> book/golang-in-action (relation: "books")
Orphaned relations have been removed from the index.
```

> **Note:** This only cleans up the SQLite index. The frontmatter in `.md` files is not modified — use `tmd relation unlink --both` to properly remove relations before deleting objects.

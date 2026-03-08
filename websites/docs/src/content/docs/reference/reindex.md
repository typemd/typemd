---
title: tmd reindex
description: Rebuild the SQLite index and search database.
sidebar:
  order: 8
---

Scans the `objects/` directory, syncs all files to the database, cleans up orphaned relations, and rebuilds the full-text search index. Use after manually editing files outside of TypeMD.

> **Note:** When opening a vault, TypeMD automatically syncs the index if it is empty or missing. You only need `tmd reindex` when files have been edited while the vault was not open.

```bash
tmd reindex
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

---
title: tmd reindex
description: Rebuild the SQLite index and search database.
sidebar:
  order: 8
---

Scans the `objects/` directory, syncs all files to the database, and rebuilds the full-text search index. Use after manually editing files outside of TypeMD.

> **Note:** When opening a vault, TypeMD automatically syncs the index if it is empty or missing. You only need `tmd reindex` when files have been edited while the vault was not open.

```bash
tmd reindex
```

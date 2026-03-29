---
title: tmd search
description: Full-text search across your vault.
sidebar:
  order: 5
---

Full-text search across filenames, properties, and body content. Powered by SQLite FTS5, with automatic fallback to substring matching when the index is unavailable.

```bash
tmd search "concurrency"
tmd search "golang" --json
```

## Flags

| Flag | Description |
|------|-------------|
| `--json` | Output results as JSON |

## Example output

```
book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz
note/concurrency-patterns-01jqr4a2bcdef0123456789xyz
```

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

## Options

| Option | Description |
|--------|-------------|
| `--json` | Output results in JSON format |

---
title: tmd template list
description: List available templates.
sidebar:
  order: 12
---

Lists all available templates in the vault, optionally filtered by type.

```bash
tmd template list
tmd template list book
tmd template list --json
```

Example output:

```
book/review
book/summary
note/meeting
```

## Flags

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON array with `type` and `name` fields |

When no type argument is provided, lists templates across all types. When a type is specified, only templates for that type are shown. If no templates exist, the output is empty.

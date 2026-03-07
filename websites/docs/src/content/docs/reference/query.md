---
title: tmd query
description: Filter objects by properties.
sidebar:
  order: 4
---

Filter objects by properties. Conditions use `key=value` format, separated by spaces (AND logic).

```bash
tmd query "type=book"
tmd query "type=book status=reading"
tmd query "type=book" --json
```

## Options

| Option | Description |
|--------|-------------|
| `--json` | Output results in JSON format |

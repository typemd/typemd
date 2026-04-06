---
title: tmd stats
description: Vault-wide summary and per-type property aggregate statistics.
sidebar:
  order: 20
---

Displays statistics about your vault. Without flags, shows a vault-wide summary. With `--type`, shows per-property aggregate statistics for a specific type.

## Vault-wide summary

```bash
tmd stats
```

Shows each type's emoji, plural name, object count, last updated date, and a total across all types.

### Example output

```
📚 books       12   last updated 2026-03-15
👤 people       8   last updated 2026-03-14
💡 ideas        5   last updated 2026-03-10
🏷️  tags         3   last updated 2026-02-28
────────────────────────────
  Total         28
```

## Per-type property stats

```bash
tmd stats --type book
```

Shows aggregate statistics for each property of the specified type:

| Property type | Aggregations |
|---------------|-------------|
| `number` | sum, avg, min, max |
| `select` | value frequency counts |
| `checkbox` | true/false counts |
| `date` | earliest, latest |
| `relation` | link count |

### Example output

```
📚 book (12 objects)
────────────────────────────────────────

rating (number)
  filled: 12/12
  sum: 46.5  avg: 3.88  min: 2  max: 5

status (select)
  filled: 12/12
  reading         ████ 4
  done            ██████ 6
  to-read         ██ 2

favorite (checkbox)
  filled: 12/12
  true: 3  false: 9

published (date)
  filled: 12/12
  earliest: 2008-08-01  latest: 2024-11-15

author (relation)
  filled: 8/12
  links: 8
```

## JSON output

```bash
tmd stats --json
tmd stats --type book --json
```

Both modes support `--json` for machine-readable output.

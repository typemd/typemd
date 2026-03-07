---
title: Wiki-links
description: Link objects inline with wiki-style syntax and track backlinks.
sidebar:
  order: 7.5
---

Wiki-links let you reference other objects inline within markdown body content using `[[target]]` syntax. Backlinks are automatically computed so linked objects know who references them.

## Syntax

```markdown
See [[book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz]] for details.
```

With optional display text:

```markdown
See [[book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz|Clean Code]] for details.
```

The target must be a full Object ID (including the ULID suffix).

## How It Works

1. When the index is synced (`tmd reindex` or auto-sync), the indexer parses `[[...]]` patterns from each object's body
2. Targets are resolved against existing objects in the database
3. Wiki-link records are stored in the SQLite index for fast backlink lookups
4. Unresolvable targets are stored as broken links (detectable via `tmd validate`)

## Backlinks

Backlinks appear as a built-in `backlinks` property when viewing an object. If object A contains `[[B]]`, then B's properties will show:

```
backlinks: ⟵ A
```

Backlinks are visible in both `tmd show` and the TUI properties panel.

## Wiki-links vs Relations

| | Wiki-links | Relations |
|--|-----------|-----------|
| **Defined in** | Markdown body content | Type schema properties |
| **Syntax** | `[[target]]` | `tmd link` command |
| **Schema required** | No | Yes (`type: relation`) |
| **Bidirectional** | Via backlinks (read-only) | Configurable (`bidirectional: true`) |
| **Stored in** | `wikilinks` table | `relations` table + frontmatter |

## Validation

`tmd validate` includes a wiki-link validation phase that detects broken links — references to objects that do not exist.

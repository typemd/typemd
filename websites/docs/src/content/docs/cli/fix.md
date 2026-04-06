---
title: tmd fix
description: Fix common issues in your vault.
sidebar:
  order: 18
---

The `fix` command provides subcommands to automatically correct common issues in your vault.

## tmd fix wikilinks

Expands shorthand wiki-link targets to full Object IDs.

```bash
tmd fix wikilinks
```

### What it does

Walks all objects in the vault and replaces shorthand wiki-link targets with their resolved full IDs (`type/name-ulid`).

Shorthand formats that get expanded:

| Format | Example | Resolution |
|--------|---------|------------|
| `[[type/name]]` | `[[book/clean-code]]` | Resolved by name within the specified type |
| `[[name]]` | `[[clean-code]]` | Resolved by name within the source object's type |

Full IDs (`[[type/name-ulid]]`) are already complete and are not modified.

### Output

```
Expanded 3 wiki-link(s) to full IDs.
```

If shorthand links match multiple objects (ambiguous), they are reported but not expanded:

```
  note/my-note-01abc: ambiguous [[golang]] — matches: book/golang-intro-01def, book/golang-guide-01ghi
```

If all wiki-links are already full IDs:

```
All wiki-links are already full IDs. No changes needed.
```

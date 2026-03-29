---
title: tmd log
description: Show git commit history for a specific object.
sidebar:
  order: 11.6
---

Displays the git commit history for a specific object file. Wraps `git log --follow` to track history across renames.

## Usage

```bash
tmd log <object-id>
```

Supports prefix matching — you can omit the ULID suffix if the prefix uniquely identifies an object. If a prefix matches multiple objects, an interactive picker is shown.

### Example output

```
commit a1b2c3d (HEAD -> main)
Author: Alice <alice@example.com>
Date:   Sat Mar 29 2026

    update book properties

commit e4f5g6h
Author: Alice <alice@example.com>
Date:   Fri Mar 28 2026

    add clean-code book
```

## Compact output

```bash
tmd log --oneline book/clean-code
```

Shows each commit on a single line with abbreviated hash and subject.

## Edge cases

- **Vault not in a git repository** — exits with error: "vault is not inside a git repository"
- **Object with no commits** — displays: "no commits found for this object"
- **Renamed object** — `--follow` tracks history across file renames

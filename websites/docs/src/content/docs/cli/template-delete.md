---
title: tmd template delete
description: Delete a template file.
sidebar:
  order: 9.9
---

Deletes a template file from the vault.

```bash
tmd template delete book/review
tmd template delete book/review --force
```

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--force` | `-f` | Skip confirmation prompt |

In interactive terminals, the command prompts for confirmation before deleting. Use `--force` to skip the prompt (useful for scripting). The argument must be in `type/name` format.
